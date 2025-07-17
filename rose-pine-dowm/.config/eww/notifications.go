package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
)

type Notification struct {
	Summary string
	Body    string
	Id      uint32
	Icon    string
	Actions []string
	AppName string
	Sender  string
}

var (
	notifications []Notification
	notifMutex    sync.Mutex
	nextId        uint32 = 0
	conn          *dbus.Conn
)

func main() {
	var err error
	conn, err = dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reply, err := conn.RequestName("org.freedesktop.Notifications", dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		panic("Failed to get D-Bus name")
	}

	conn.Export(&NotificationServer{}, "/org/freedesktop/Notifications", "org.freedesktop.Notifications")
	setupLogger()

	pipePath := "/tmp/notifications.pipe"
	os.Remove(pipePath)
	err = os.MkdirAll("/tmp", 0755)
	if err != nil {
		panic(err)
	}
	if err := syscall.Mkfifo(pipePath, 0666); err != nil {
		panic(err)
	}

	go listenToPipe(pipePath)

	select {}
}

type NotificationServer struct{}

func (n *NotificationServer) Notify(appName string, replacesID uint32, appIcon, summary, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	log.Println(appIcon)
	log.Println("New notification from:", appName)
	id := replacesID
	if replacesID == 0 {
		nextId++
		id = nextId
	}

	sender := conn.Names()[0]
	addNotification(Notification{
		Summary: summary,
		Body:    body,
		Icon:    appIcon,
		Id:      id,
		Actions: actions,
		AppName: appName,
		Sender:  sender,
	})
	return id, nil
}

func addNotification(n Notification) {
	notifMutex.Lock()
	notifications = append([]Notification{n}, notifications...)
	notifMutex.Unlock()

	printState()

	go func(id uint32) {
		time.Sleep(10 * time.Second)
		notifMutex.Lock()
		for i, notif := range notifications {
			if notif.Id == id {
				notifications = append(notifications[:i], notifications[i+1:]...)
				break
			}
		}
		notifMutex.Unlock()
		printState()
	}(n.Id)
}

func printState() {
	str := ""
	for _, notif := range notifications {
		str += `
		(eventbox :onclick 'echo "`+strconv.Itoa(int(notif.Id))+`|close"> /tmp/notifications.pipe'
		(box :class 'notif'
             (box :orientation 'h' :space-evenly 'false'
             (image :hexpand "true" :image-width 80 :image-height 80 :path '` + notif.Icon + `')
             (box :orientation 'v' :hexpand 'true' :space-evenly "false"
             (label :wrap "true" :class "notif-summary" :halign "start" :text '` + notif.Summary + `')
             (label :wrap "true" :class "notif-body" :halign "start" :text '` + notif.Body + `')
             (box :orientation "h" :class "notif-actions" :space-evenly "true"          
        `
		for i := 0; i < len(notif.Actions); i += 2 {
			str += `(button :class "action" :onclick 'echo "` + strconv.Itoa(int(notif.Id)) + `|action|` + notif.Actions[i] + `" > /tmp/notifications.pipe' "` + notif.Actions[i+1] + `")`
		}
		str += ")))))"
	}
	str = strings.ReplaceAll(str, "\n", " ")
	box := `(box :orientation "v" :height 0 :class "notifications" `
	if len(notifications)==0{
		box += `:visible "false" :style "min-height: 0px;"`
	}
	str = box + str + ")"
	fmt.Println(str)
	log.Println(str)
}

func (n *NotificationServer) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"body", "body-markup", "actions"}, nil
}

func (n *NotificationServer) GetServerInformation() (string, string, string, string, *dbus.Error) {
	return "GoNotify", "sam", "1.0", "1.2", nil
}

func (n *NotificationServer) CloseNotification(id uint32) *dbus.Error {
	notifMutex.Lock()
	defer notifMutex.Unlock()
	for i, notif := range notifications {
		if notif.Id == id {
			notifications = append(notifications[:i], notifications[i+1:]...)
			break
		}
	}
	printState()
	return nil
}

func setupLogger() {
	f, err := os.OpenFile("/tmp/notifications.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}

func listenToPipe(pipePath string) {
	for {
		f, err := os.OpenFile(pipePath, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			log.Println("Pipe open error:", err)
			time.Sleep(time.Second)
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			handlePipeCommand(scanner.Text())
		}
		f.Close()
	}
}

func handlePipeCommand(cmd string) {
	parts := strings.Split(cmd, "|")
	if len(parts) < 2 {
		return
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	notifMutex.Lock()
	defer notifMutex.Unlock()

	var notif *Notification
	for i := range notifications {
		if notifications[i].Id == uint32(id) {
			notif = &notifications[i]
			break
		}
	}
	if notif == nil {
		return
	}

	switch parts[1] {
	case "action":
		if len(parts) == 3 {
			log.Printf("Action invoked: id=%d, action=%s, app=%s", id, parts[2], notif.AppName)
			
			// Send the ActionInvoked signal
			err := conn.Emit(
				"/org/freedesktop/Notifications",
				"org.freedesktop.Notifications.ActionInvoked",
				uint32(id),
				parts[2],
			)
			if err != nil {
				log.Println("Failed to emit ActionInvoked signal:", err)
			}

			// Also send NotificationClosed signal with reason=2 (dismissed by user)
			err = conn.Emit(
				"/org/freedesktop/Notifications",
				"org.freedesktop.Notifications.NotificationClosed",
				uint32(id),
				uint32(2),
			)
			if err != nil {
				log.Println("Failed to emit NotificationClosed signal:", err)
			}

			// Remove the notification
			for i, n := range notifications {
				if n.Id == uint32(id) {
					notifications = append(notifications[:i], notifications[i+1:]...)
					break
				}
			}
			printState()
		}

	case "close":
		err = conn.Emit(
			"/org/freedesktop/Notifications",
			"org.freedesktop.Notifications.NotificationClosed",
			uint32(id),
			uint32(2),
		)
		if err != nil {
			log.Println("Failed to emit NotificationClosed signal:", err)
		}
		for i, n := range notifications {
			if n.Id == uint32(id) {
				notifications = append(notifications[:i], notifications[i+1:]...)
				break
			}
		}
		printState()
	}
}

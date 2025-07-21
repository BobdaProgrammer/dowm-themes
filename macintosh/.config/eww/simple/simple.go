
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
    "image"
	"image/color"
    "image/png"
	"time"
	"path/filepath"
	"github.com/godbus/dbus/v5"
)

type Notification struct {
	Summary string
	Body    string
	Id      uint32
	Icon    string
	Actions []string
	AppName string
	Urgent bool
	Sender  string
}

var (
	notification Notification
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

func expandPath(path string) string {
    if strings.HasPrefix(path, "~/") {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            log.Printf("Error getting home directory: %v", err)
            return path // Return original path if we can't get home dir
        }
        return filepath.Join(homeDir, path[2:]) // Remove ~/ and join with home
    }
    return path
}
func extractImageData(hints map[string]dbus.Variant) string {
    imageData, exists := hints["image-data"]
    if !exists {
		imageData, exists = hints["image_data"]
		if !exists{
        	return ""
		}
    }
    
    // The image-data format is: (width, height, rowstride, has_alpha, bits_per_sample, channels, data)
    data, ok := imageData.Value().([]interface{})
    if !ok || len(data) < 7 {
        return ""
    }
    
    width := int(data[0].(int32))
    height := int(data[1].(int32))
    rowstride := int(data[2].(int32))
    hasAlpha := data[3].(bool)
    bitsPerSample := int(data[4].(int32))
    channels := int(data[5].(int32))
    imageBytes := data[6].([]byte)
    
    log.Printf("Image data: %dx%d, rowstride: %d, alpha: %v, bits: %d, channels: %d, data_len: %d", 
        width, height, rowstride, hasAlpha, bitsPerSample, channels, len(imageBytes))
    
    // Convert RGBA data to PNG and save
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            offset := y*rowstride + x*channels
            if offset+channels <= len(imageBytes) {
                img.Set(x, y, color.RGBA{
                    R: imageBytes[offset],
                    G: imageBytes[offset+1],
                    B: imageBytes[offset+2],
                    A: 255, // Default alpha
                })
                if hasAlpha && channels >= 4 {
                    img.SetRGBA(x, y, color.RGBA{
                        R: imageBytes[offset],
                        G: imageBytes[offset+1], 
                        B: imageBytes[offset+2],
                        A: imageBytes[offset+3],
                    })
                }
            }
        }
    }
    
    // Save to temp file
    tempFile := "/tmp/notification_img.png"
    file, err := os.Create(tempFile)
    if err != nil {
        log.Printf("Error creating temp file: %v", err)
        return ""
    }
    defer file.Close()
    
    err = png.Encode(file, img)
    if err != nil {
        log.Printf("Error encoding PNG: %v", err)
        return ""
    }
    
    return tempFile
}
func (n *NotificationServer) Notify(appName string, replacesID uint32, appIcon, summary, body string, actions []string, hints map[string]dbus.Variant, expireTimeout int32) (uint32, *dbus.Error) {
	//fmt.Println(hints)
	for key := range hints{
		fmt.Println(key)
	}
	fmt.Println(appIcon)   
	file := extractImageData(hints)
	if file == "" && appIcon!=""{
		appIcon = expandPath(appIcon)
	}else if file != ""{
		appIcon = file
	}
	log.Println("New notification from:", appName)
	id := replacesID
	if replacesID == 0 {
		nextId++
		id = nextId
	}
	urgent:=false

	if levelVar, exists := hints["urgency"]; exists {
		if level, ok := levelVar.Value().(byte); ok && level == 2 {
			urgent = true
		}
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
		Urgent:	 urgent,
	})
	return id, nil
}

func addNotification(n Notification) {
	notification = n
	printState()

	go func(id uint32) {
		oldId := n.Id
		time.Sleep(5 * time.Second)
		if notification.Id == oldId{
			clearState()
		}
	}(n.Id)
}

func clearState(){
	fmt.Println("")
	exec.Command("eww", "close", "notifications").Run()
	exec.Command("eww", "update", `NOTIFICATION=''`).Run()
}

func printState() {
	class := "notif"
	if notification.Urgent{
		class+=" urgent"
	}
	str := `
	(eventbox :onclick 'echo "`+strconv.Itoa(int(notification.Id))+`|close"> /tmp/notifications.pipe'
	(box :class "`+class+`"
	(box :orientation "h" :space-evenly "false"`
	image := `(box :hexpand "true" :class "image-container" (image :image-width 90 :image-height 90 :path "` + notification.Icon + `"))`
	notif:=`(box :orientation "v" :class "notification_text" :hexpand "true" :vexpand "true" :space-evenly "false"
	(label :wrap "true" :class "notif-summary" :halign "center" :vexpand "true" :valign "end" :limit-width 20 :text "` + notification.Summary + `")
	(label :wrap "true" :class "notif-body" :halign "center" :vexpand "true" :valign "center" :limit-width 40 :text "` + notification.Body + `")
	(box :orientation "h" :class "notif-actions" :vexpand "true" :valign "start" :space-evenly "true"          
	`
	if notification.Icon != ""{
		str+=image
	}
	str+=notif
	for i := 0; i < len(notification.Actions); i += 2 {
		str += `(button :class "action" :onclick 'echo "` + strconv.Itoa(int(notification.Id)) + `|action|` + notification.Actions[i] + `" > /tmp/notifications.pipe' "` + notification.Actions[i+1] + `")`
	}
	str += ")))))"
	str = strings.ReplaceAll(str, "\n", " ")
	
	// Execute the commands
	updateCmd := exec.Command("eww", "update", "NOTIFICATION="+str)
	err := updateCmd.Run()
	if err != nil {
		log.Printf("Error updating eww: %v", err)
	}
	
	openCmd := exec.Command("eww", "open", "notifications")
	err = openCmd.Run()
	if err != nil {
		log.Printf("Error opening eww window: %v", err)
	}
	
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
	clearState()
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



	switch parts[1] {
	case "action":
		if len(parts) == 3 {
			log.Printf("Action invoked: id=%d, action=%s, app=%s", id, parts[2], notification.AppName)
			
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
			clearState()
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
		clearState()
	}
}

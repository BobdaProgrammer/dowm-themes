
#!/bin/bash

# Check for correct number of arguments
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <number> <total>"
  exit 1
fi

number=$1
total=$2

# Avoid division by zero
if [ "$total" -eq 0 ]; then
  echo "Error: total cannot be zero"
  exit 1
fi

# Calculate percentage
percentage=$(awk "BEGIN { printf \"%.2f\", ($number / $total) * 100 }")
echo "$percentage"

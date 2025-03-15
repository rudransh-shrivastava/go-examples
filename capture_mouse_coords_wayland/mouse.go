package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type InputEvent struct {
	Time  TimeVal
	Type  uint16
	Code  uint16
	Value int32
}

type TimeVal struct {
	Sec  int64
	Usec int64
}

func main() {
	// Use `la /dev/input/by-id/*` to find mouse device
	// For me, the output was
	// 	lrwxrwxrwx 1 root root 9 Feb 19 22:37 usb-Logitech_USB_Receiver-if01-event-kbd -> ../event6
	// lrwxrwxrwx 1 root root 9 Feb 19 22:37 usb-Logitech_USB_Receiver-if01-event-mouse -> ../event5
	// lrwxrwxrwx 1 root root 9 Feb 19 22:37 usb-Logitech_USB_Receiver-if01-mouse -> ../mouse2
	// the mouse is event5

	devicePath := "/dev/input/event4"

	f, err := os.Open(devicePath)
	if err != nil {
		fmt.Printf("Failed to open device: %v\n", err)
		return
	}
	defer f.Close()

	event := &InputEvent{}
	for {
		err := binary.Read(f, binary.LittleEndian, event)
		if err != nil {
			fmt.Printf("Failed to read event: %v\n", err)
			return
		}

		// EV_REL (2) is for relative movement
		// REL_X (0) is for x-axis movement
		// REL_Y (1) is for y-axis movement
		if event.Type == 2 { // EV_REL
			switch event.Code {
			case 0: // REL_X
				fmt.Printf("X movement: %d\n", event.Value)
			case 1: // REL_Y
				fmt.Printf("Y movement: %d\n", event.Value)
			}
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// interface (default "any" to capture on all interfaces)
	iface := flag.String("i", "any", "Interface to capture packets from (e.g., eth0, wlan0, any)")
	flag.Parse()

	// Open the device for capturing
	handle, err := pcap.OpenLive(*iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening device %s: %v", *iface, err)
	}
	defer handle.Close()

	fmt.Printf("Listening on interface %s...\n", *iface)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		processPacket(packet)
	}
}

// processPacket prints out details about each captured packet.
func processPacket(packet gopacket.Packet) {
	timestamp := packet.Metadata().Timestamp.Format(time.RFC3339)
	fmt.Printf("Time: %s\n", timestamp)

	// Network layer (IPv4/IPv6, etc.)
	if netLayer := packet.NetworkLayer(); netLayer != nil {
		src, dst := netLayer.NetworkFlow().Endpoints()
		fmt.Printf("Network: %s -> %s\n", src, dst)
	}

	// Transport layer (TCP, UDP, etc.)
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		switch t := transportLayer.(type) {
		case *layers.TCP:
			fmt.Printf("TCP: %s -> %s\n", t.SrcPort, t.DstPort)
		case *layers.UDP:
			fmt.Printf("UDP: %s -> %s\n", t.SrcPort, t.DstPort)
		default:
			fmt.Printf("Other Transport: %s\n", transportLayer.LayerType())
		}
	} else {
		fmt.Println("No Transport Layer found")
	}

	// print payload ?
	// if appLayer := packet.ApplicationLayer(); appLayer != nil {
	// 	fmt.Printf("Payload: %s\n", appLayer.Payload())
	// }

	fmt.Println("------------------------------")
}

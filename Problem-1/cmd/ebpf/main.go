package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

const (
	// Default port to drop TCP packets
	defaultPort = 4040
)

func main() {
	// Remove memory lock limits to ensure the eBPF program can be loaded.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	var objs drop_packetsObjects
	if err := loadDrop_packetsObjects(&objs, nil); err != nil {
		log.Fatal("Loading eBPF objects:", err)
	}
	defer objs.Close()

	// Specify the network interface to attach the XDP program to.
	ifname := "lo" //  Change this to the desired network interface on your machine.
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatalf("Getting interface %s: %s", ifname, err)
	}

	// Attach the XDP program (drop_tcp_packets) to the specified network interface.
	xdpLink, err := link.AttachXDP(link.XDPOptions{
		Program:   objs.DropTcpPackets,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("Attaching XDP:", err)
	}
	defer xdpLink.Close()

	log.Printf("Dropping TCP packets on port %d on %s..", defaultPort, ifname)

	// Update the drop_port map with the port number to drop TCP packets.
	port := defaultPort
	if len(os.Args) > 1 {
		var portArg int
		if _, err := fmt.Sscanf(os.Args[1], "%d", &portArg); err == nil {
			port = portArg
		}
	}

	key := uint32(0)
	value := uint16(port)
	if err := objs.DropPort.Update(&key, &value, ebpf.UpdateAny); err != nil {
		log.Fatal("Updating drop_port map:", err)
	}

	log.Printf("Configured to drop TCP packets on port %d", port)

	// Set up a signal handler to handle program interruption (e.g., Ctrl+C).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Received signal, exiting..")
}

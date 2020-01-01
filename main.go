package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// defaultAddressTimeout defines the duration to wait
	// for new addresses.
	defaultAddressTimeout = time.Minute * 10

	// defaultNodeTimeout defines the timeout time waiting for
	// a response from a node.
	defaultNodeTimeout = time.Second * 3
)

var (
	manager  *Manager
	globalWg sync.WaitGroup
	cfg      *config
)

func main() {
	var err error
	cfg, err = loadConfig()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	manager, err = NewManager(filepath.Join(defaultHomeDir,
		activeNetParams.Name))
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	if !cfg.Check {
		manager.AddAddresses([]net.IP{net.ParseIP(cfg.Seeder)})

		globalWg.Add(1)
		go creep()

		dnsServer := NewDNSServer(cfg.Host, cfg.Nameserver, cfg.Listen)
		go dnsServer.Start()

		globalWg.Wait()
	}else {
		log.Printf("check the peer %v \n",cfg.CheckPeer)
		manager.AddAddresses([]net.IP{net.ParseIP(cfg.CheckPeer)})
		var wg sync.WaitGroup
		wg.Add(1)
		ip := net.ParseIP(cfg.CheckPeer)
		doCreep(ip,wg)
		wg.Wait()
	}
}

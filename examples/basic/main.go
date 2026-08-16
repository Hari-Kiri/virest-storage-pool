package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Hari-Kiri/virest/storagePool"
)

func main() {
	uri := "qemu:///system"
	if v := os.Getenv("VIREST_LIBVIRT_URI"); v != "" {
		uri = v
	}

	conn, err := storagePool.Connect(uri)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	pools, err := conn.List(0, 0)
	if err != nil {
		log.Fatalf("list: %v", err)
	}
	fmt.Printf("found %d storage pools on %s\n", len(pools), uri)
	n := len(pools)
	for i := 0; i < n; i++ {
		p := pools[i]
		fmt.Printf("- %s (%s) state=%v\n", p.Name, p.Uuid, p.State)
	}
}

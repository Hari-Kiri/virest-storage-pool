package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Hari-Kiri/virest/storagePool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	uri := "qemu:///system"
	if v := os.Getenv("VIREST_LIBVIRT_URI"); v != "" {
		uri = v
	}

	conn, err := storagePool.Connect(uri)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	pools, err := conn.List(0, 0)
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	fmt.Printf("found %d storage pools on %s\n", len(pools), uri)
	n := len(pools)
	for i := 0; i < n; i++ {
		p := pools[i]
		fmt.Printf("- %s (%s) state=%v\n", p.Name, p.Uuid, p.State)
	}
	return nil
}

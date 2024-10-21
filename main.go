package main

import (
	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/virest-storage-pool/handlers/storagePool"
)

func main() {
	goalMakeHandler.HandleRequest(storagePool.PoolDefine, "/pool-define")
	goalMakeHandler.HandleRequest(storagePool.PoolUndefine, "/pool-undefine")
	goalMakeHandler.Serve("Gerandong", 8000)
}

package main

import (
	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/virest-storage-pool/handlers/poolDefine"
)

func main() {
	goalMakeHandler.HandleRequest(poolDefine.PoolDefine, "/pool-define")
	goalMakeHandler.Serve("Gerandong", 8000)
}

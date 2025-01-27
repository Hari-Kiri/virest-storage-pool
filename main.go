package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/handlers/storagePool"
	"libvirt.org/go/libvirt"
)

func main() {
	// Read .env file
	readEnvFile, errorReadEnvFile := os.ReadFile(".env")
	if errorReadEnvFile != nil {
		temboLog.FatalLogging("failed read env file:", errorReadEnvFile)
	}

	// Set environment variables based on .env file
	rows := strings.Split(string(readEnvFile), "\n")
	for i := 0; i < len(rows); i++ {
		columns := strings.Split(rows[i], "=")
		os.Setenv(columns[0], columns[1])
	}

	// Convert environment variable which is hold port number to int
	portFromEnv, errorGetPortFromEnv := strconv.Atoi(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_PORT"))
	if errorGetPortFromEnv != nil {
		temboLog.FatalLogging("failed get port from env:", errorGetPortFromEnv)
	}

	temboLog.InfoLogging("registering a default event implementation based on the poll() system call...")
	errorEventRegisterDefaultImpl := libvirt.EventRegisterDefaultImpl()
	if errorEventRegisterDefaultImpl != nil {
		temboLog.FatalLogging("failed registers a default event implementation based on the poll() system call:", errorEventRegisterDefaultImpl)
	}
	temboLog.InfoLogging("registering a default event implementation based on the poll() system call, success!")

	// Make handler
	goalMakeHandler.HandleRequest(storagePool.Authenticate, "/storage-pool/authenticate")
	goalMakeHandler.HandleRequest(storagePool.FindStoragePoolSource, "/storage-pool/find-storage-pool-sources")
	goalMakeHandler.HandleRequest(storagePool.PoolList, "/storage-pool/list")
	goalMakeHandler.HandleRequest(storagePool.PoolInfo, "/storage-pool/info")
	goalMakeHandler.HandleRequest(storagePool.PoolDetail, "/storage-pool/detail")
	goalMakeHandler.HandleRequest(storagePool.GetUid, "/storage-pool/get-uid")
	goalMakeHandler.HandleRequest(storagePool.GetGid, "/storage-pool/get-gid")
	goalMakeHandler.HandleRequest(storagePool.PoolDefine, "/storage-pool/define")
	goalMakeHandler.HandleRequest(storagePool.PoolBuild, "/storage-pool/build")
	goalMakeHandler.HandleRequest(storagePool.PoolCreate, "/storage-pool/create")
	goalMakeHandler.HandleRequest(storagePool.PoolAutostart, "/storage-pool/autostart")
	goalMakeHandler.HandleRequest(storagePool.PoolDestroy, "/storage-pool/destroy")
	goalMakeHandler.HandleRequest(storagePool.PoolUndefine, "/storage-pool/undefine")
	goalMakeHandler.HandleRequest(storagePool.PoolDelete, "/storage-pool/delete")
	goalMakeHandler.HandleRequest(storagePool.PoolRefresh, "/storage-pool/refresh")
	goalMakeHandler.HandleRequest(storagePool.PoolCapabilities, "/storage-pool/capabilities")
	goalMakeHandler.HandleRequest(storagePool.PoolEvent, "/storage-pool/event")
	goalMakeHandler.Serve(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"), portFromEnv)
}

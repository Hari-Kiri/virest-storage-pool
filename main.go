package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/Hari-Kiri/goalMakeHandler"
	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/handlers/storagePool"
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

	// Make handler
	goalMakeHandler.HandleRequest(storagePool.PoolList, "/pool-list")
	goalMakeHandler.HandleRequest(storagePool.PoolDetail, "/pool-detail")
	goalMakeHandler.HandleRequest(storagePool.PoolDefine, "/pool-define")
	goalMakeHandler.HandleRequest(storagePool.PoolBuild, "/pool-build")
	goalMakeHandler.HandleRequest(storagePool.PoolCreate, "/pool-create")
	goalMakeHandler.HandleRequest(storagePool.PoolUndefine, "/pool-undefine")
	goalMakeHandler.Serve(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"), portFromEnv)
}

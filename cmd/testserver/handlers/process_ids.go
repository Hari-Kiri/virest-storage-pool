package handlers

import "os"

func processUID() int { return os.Getuid() }
func processGID() int { return os.Getgid() }

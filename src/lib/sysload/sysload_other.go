// +build !linux

package sysload

import (
	"log"
	"runtime"
	"time"
)

func getLoadAvgInfo() (load LoadAvgInfo) {
	load.Time = time.Now()
	log.Println("LoadAvg is not implemented in OS: " + runtime.GOOS + "\n")
	return
}

func getMemInfo() (load MemInfo) {
	load.Time = time.Now()
	log.Println("meminfo is not implemented in OS:" + runtime.GOOS + "\n")
	return
}

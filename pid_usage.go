package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func bytesToMiB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

func millisecondsToHHMMSS(ms int64) string {
	seconds := ms / 1000
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getProcessUsage(pid int32) {
	p, err := process.NewProcess(pid)
	if err != nil {
		log.Fatalf("Failed to create process: %v", err)
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		log.Fatalf("Failed to get CPU usage: %v", err)
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		log.Fatalf("Failed to get memory usage: %v", err)
	}

	memPercent, err := p.MemoryPercent()
	if err != nil {
		log.Fatalf("Failed to get memory percent: %v", err)
	}

	numThreads, err := p.NumThreads()
	if err != nil {
		log.Fatalf("Failed to get number of threads: %v", err)
	}

	createTime, err := p.CreateTime()
	if err != nil {
		log.Fatalf("Failed to get process create time: %v", err)
	}

	uptimeMs := time.Now().UnixMilli() - createTime
	uptimeStr := millisecondsToHHMMSS(uptimeMs)

	fmt.Printf("PID: %d\n", pid)
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent)
	fmt.Printf("Memory Usage: %.2f MiB\n", bytesToMiB(memInfo.RSS))
	fmt.Printf("Virtual Memory Usage: %.2f MiB\n", bytesToMiB(memInfo.VMS))
	fmt.Printf("Memory Percentage: %.2f%%\n", memPercent)
	fmt.Printf("Number of Threads: %d\n", numThreads)
	fmt.Printf("Process Uptime: %s\n", uptimeStr)
}

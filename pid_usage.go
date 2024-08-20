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

func getProcessUsage(pid int32) {
	p, err := process.NewProcess(pid)
	if err != nil {
		log.Fatalf("failed to create process: %v", err)
	}

	cpuPercent, err := p.CPUPercent()
	if err != nil {
		log.Fatalf("failed to get CPU usage: %v", err)
	}

	memInfo, err := p.MemoryInfo()
	if err != nil {
		log.Fatalf("failed to get memory usage: %v", err)
	}

	memPercent, err := p.MemoryPercent()
	if err != nil {
		log.Fatalf("failed to get memory percent: %v", err)
	}

	numThreads, err := p.NumThreads()
	if err != nil {
		log.Fatalf("failed to get number of threads: %v", err)
	}

	createTime, err := p.CreateTime()
	if err != nil {
		log.Fatalf("failed to get process create time: %v", err)
	}

	fmt.Printf("PID: %d\n", pid)
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent)
	fmt.Printf("Memory Usage: %.2f MiB\n", bytesToMiB(memInfo.RSS))
	fmt.Printf("Virtual Memory Usage: %.2f MiB\n", bytesToMiB(memInfo.VMS))
	fmt.Printf("Memory Percentage: %.2f%%\n", memPercent)
	fmt.Printf("Number of Threads: %d\n", numThreads)
	fmt.Printf("Process Uptime: %s\n", time.Since(time.Unix(createTime/1000, 0)).String())
}

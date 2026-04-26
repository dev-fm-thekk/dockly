package dockerapi

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/moby/moby/client"
)

type MemoryUsage struct {
	Used  int64
	Total int64
}

type CPUUsage struct {
	Percentage float64
}

type Metrics struct {
	Status           bool
	Memory           MemoryUsage
	CPU              CPUUsage
	ActiveContainers int
	Images           int
}

func getSystemStatus() bool {
	cli := GetClient()
	if cli == nil {
		return false
	}
	_, err := cli.Ping(context.Background(), client.PingOptions{})
	return err == nil
}

func FetchMetrics() Metrics {
	metrics := Metrics{
		Status:           false,
		Memory:           MemoryUsage{},
		CPU:              CPUUsage{},
		ActiveContainers: 0,
		Images:           0,
	}

	cli := GetClient()
	if cli == nil {
		return metrics
	}

	ctx := context.Background()

	if _, err := cli.Ping(ctx, client.PingOptions{}); err != nil {
		return metrics
	}
	metrics.Status = true

	info, err := cli.Info(ctx, client.InfoOptions{})
	if err == nil {
		metrics.ActiveContainers = info.Info.ContainersRunning
		metrics.Images = info.Info.Images
		metrics.Memory.Total = info.Info.MemTotal
	}

	containers, err := cli.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		return metrics
	}

	var totalMemUsed int64
	var totalCPUPercentage float64
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, ctr := range containers.Items {
		wg.Add(1)
		go func(containerID string) {
			defer wg.Done()
			stats, err := cli.ContainerStats(ctx, containerID, client.ContainerStatsOptions{Stream: false, IncludePreviousSample: true})
			if err != nil {
				return
			}
			defer stats.Body.Close()

			var v struct {
				MemoryStats struct {
					Usage int64 `json:"usage"`
					Stats struct {
						Cache int64 `json:"cache"`
					} `json:"stats"`
				} `json:"memory_stats"`
				CPUStats struct {
					CPUUsage struct {
						TotalUsage  uint64   `json:"total_usage"`
						PercpuUsage []uint64 `json:"percpu_usage"`
					} `json:"cpu_usage"`
					SystemUsage uint64 `json:"system_cpu_usage"`
				} `json:"cpu_stats"`
				PreCPUStats struct {
					CPUUsage struct {
						TotalUsage uint64 `json:"total_usage"`
					} `json:"cpu_usage"`
					SystemUsage uint64 `json:"system_cpu_usage"`
				} `json:"precpu_stats"`
			}

			if err := json.NewDecoder(stats.Body).Decode(&v); err != nil {
				return
			}

			// Calculate CPU percentage
			cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
			systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)

			var cpuPercent float64
			if systemDelta > 0.0 && cpuDelta > 0.0 {
				cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
			}

			memUsed := v.MemoryStats.Usage
			if v.MemoryStats.Stats.Cache > 0 {
				memUsed -= v.MemoryStats.Stats.Cache
			}

			mu.Lock()
			totalCPUPercentage += cpuPercent
			if memUsed > 0 {
				totalMemUsed += memUsed
			}
			mu.Unlock()
		}(ctr.ID)
	}

	wg.Wait()

	metrics.CPU.Percentage = totalCPUPercentage
	metrics.Memory.Used = totalMemUsed

	return metrics
}

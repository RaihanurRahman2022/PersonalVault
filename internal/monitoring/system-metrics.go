package monitoring

import (
	"context"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v4/cpu"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type SystemMetrics struct {
	// CPU Metrics
	cpuUsagePercent metric.Float64Gauge
	cpuLoad1        metric.Float64Gauge
	cpuLoad5        metric.Float64Gauge
	cpuLoad15       metric.Float64Gauge

	// Memory Metrics
	memoryTotal     metric.Int64Gauge
	memoryUsed      metric.Int64Gauge
	memoryFree      metric.Int64Gauge
	memoryAvailable metric.Int64Gauge
	memoryCached    metric.Int64Gauge
	memoryBuffers   metric.Int64Gauge
	memorySwapTotal metric.Int64Gauge
	memorySwapUsed  metric.Int64Gauge

	// Disk Metrics
	diskTotal        metric.Int64Gauge
	diskUsed         metric.Int64Gauge
	diskFree         metric.Int64Gauge
	diskUsagePercent metric.Float64Gauge
	diskReadBytes    metric.Int64Counter
	diskWriteBytes   metric.Int64Counter

	// Network Metrics
	networkBytesReceived metric.Int64Counter
	networkBytesSent     metric.Int64Counter
	networkPacketsRecv   metric.Int64Counter
	networkPacketsSent   metric.Int64Counter

	// Go Runtime Metrics
	goMemoryUsage    metric.Int64UpDownCounter
	goHeapSize       metric.Int64Gauge
	goStackSize      metric.Int64Gauge
	goGoroutineCount metric.Int64Gauge
	goGcDuration     metric.Float64Histogram
}

func NewSystemMetrics() *SystemMetrics {
	meter := otel.Meter("personal-vault-system")

	// CPU Metrics
	cpuUsagePercent, _ := meter.Float64Gauge("system_cpu_usage_percent", metric.WithDescription("CPU usages percentage"))

	cpuLoad1, _ := meter.Float64Gauge(
		"system_load_1m",
		metric.WithDescription("System load average 1 minute"),
	)

	cpuLoad5, _ := meter.Float64Gauge(
		"system_load_5m",
		metric.WithDescription("System load average 5 minutes"),
	)

	cpuLoad15, _ := meter.Float64Gauge(
		"system_load_15m",
		metric.WithDescription("System load average 15 minutes"),
	)

	// Memory Metrics
	memoryTotal, _ := meter.Int64Gauge(
		"system_memory_total_bytes",
		metric.WithDescription("Total system memory in bytes"),
	)

	memoryUsed, _ := meter.Int64Gauge(
		"system_memory_used_bytes",
		metric.WithDescription("Used system memory in bytes"),
	)

	memoryFree, _ := meter.Int64Gauge(
		"system_memory_free_bytes",
		metric.WithDescription("Free system memory in bytes"),
	)

	memoryAvailable, _ := meter.Int64Gauge(
		"system_memory_available_bytes",
		metric.WithDescription("Available system memory in bytes"),
	)

	memoryCached, _ := meter.Int64Gauge(
		"system_memory_cached_bytes",
		metric.WithDescription("Cached system memory in bytes"),
	)

	memoryBuffers, _ := meter.Int64Gauge(
		"system_memory_buffers_bytes",
		metric.WithDescription("Buffered system memory in bytes"),
	)

	memorySwapTotal, _ := meter.Int64Gauge(
		"system_memory_swap_total_bytes",
		metric.WithDescription("Total swap memory in bytes"),
	)

	memorySwapUsed, _ := meter.Int64Gauge(
		"system_memory_swap_used_bytes",
		metric.WithDescription("Used swap memory in bytes"),
	)

	// Disk Metrics
	diskTotal, _ := meter.Int64Gauge(
		"system_disk_total_bytes",
		metric.WithDescription("Total disk space in bytes"),
	)

	diskUsed, _ := meter.Int64Gauge(
		"system_disk_used_bytes",
		metric.WithDescription("Used disk space in bytes"),
	)

	diskFree, _ := meter.Int64Gauge(
		"system_disk_free_bytes",
		metric.WithDescription("Free disk space in bytes"),
	)

	diskUsagePercent, _ := meter.Float64Gauge(
		"system_disk_usage_percent",
		metric.WithDescription("Disk usage percentage"),
	)

	diskReadBytes, _ := meter.Int64Counter(
		"system_disk_read_bytes_total",
		metric.WithDescription("Total disk read bytes"),
	)

	diskWriteBytes, _ := meter.Int64Counter(
		"system_disk_write_bytes_total",
		metric.WithDescription("Total disk write bytes"),
	)

	// Network Metrics
	networkBytesReceived, _ := meter.Int64Counter(
		"system_network_receive_bytes_total",
		metric.WithDescription("Total network bytes received"),
	)

	networkBytesSent, _ := meter.Int64Counter(
		"system_network_transmit_bytes_total",
		metric.WithDescription("Total network bytes sent"),
	)

	networkPacketsRecv, _ := meter.Int64Counter(
		"system_network_receive_packets_total",
		metric.WithDescription("Total network packets received"),
	)

	networkPacketsSent, _ := meter.Int64Counter(
		"system_network_transmit_packets_total",
		metric.WithDescription("Total network packets sent"),
	)

	// Go Runtime Metrics
	goMemoryUsage, _ := meter.Int64UpDownCounter(
		"go_memory_usage_bytes",
		metric.WithDescription("Go runtime memory usage in bytes"),
	)

	goHeapSize, _ := meter.Int64Gauge(
		"go_heap_size_bytes",
		metric.WithDescription("Go heap size in bytes"),
	)

	goStackSize, _ := meter.Int64Gauge(
		"go_stack_size_bytes",
		metric.WithDescription("Go stack size in bytes"),
	)

	goGoroutineCount, _ := meter.Int64Gauge(
		"go_goroutines_total",
		metric.WithDescription("Number of goroutines"),
	)

	goGcDuration, _ := meter.Float64Histogram(
		"go_gc_duration_seconds",
		metric.WithDescription("Go GC duration in seconds"),
	)

	return &SystemMetrics{
		// CPU
		cpuUsagePercent: cpuUsagePercent,
		cpuLoad1:        cpuLoad1,
		cpuLoad5:        cpuLoad5,
		cpuLoad15:       cpuLoad15,

		// Memory
		memoryTotal:     memoryTotal,
		memoryUsed:      memoryUsed,
		memoryFree:      memoryFree,
		memoryAvailable: memoryAvailable,
		memoryCached:    memoryCached,
		memoryBuffers:   memoryBuffers,
		memorySwapTotal: memorySwapTotal,
		memorySwapUsed:  memorySwapUsed,

		// Disk
		diskTotal:        diskTotal,
		diskUsed:         diskUsed,
		diskFree:         diskFree,
		diskUsagePercent: diskUsagePercent,
		diskReadBytes:    diskReadBytes,
		diskWriteBytes:   diskWriteBytes,

		// Network
		networkBytesReceived: networkBytesReceived,
		networkBytesSent:     networkBytesSent,
		networkPacketsRecv:   networkPacketsRecv,
		networkPacketsSent:   networkPacketsSent,

		// Go Runtime
		goMemoryUsage:    goMemoryUsage,
		goHeapSize:       goHeapSize,
		goStackSize:      goStackSize,
		goGoroutineCount: goGoroutineCount,
		goGcDuration:     goGcDuration,
	}

}

func (sm *SystemMetrics) RecordSystemMetrics() {
	ctx := context.Background()

	// CPU Metrics
	if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
		sm.cpuUsagePercent.Record(ctx, cpuPercent[0])
	}

	if loadAvg, err := load.Avg(); err == nil {
		sm.cpuLoad1.Record(ctx, loadAvg.Load1)
		sm.cpuLoad5.Record(ctx, loadAvg.Load5)
		sm.cpuLoad15.Record(ctx, loadAvg.Load15)
	}

	// Memory Metrics
	if memInfo, err := mem.VirtualMemory(); err == nil {
		sm.memoryTotal.Record(ctx, int64(memInfo.Total))
		sm.memoryUsed.Record(ctx, int64(memInfo.Used))
		sm.memoryFree.Record(ctx, int64(memInfo.Free))
		sm.memoryAvailable.Record(ctx, int64(memInfo.Available))
		sm.memoryCached.Record(ctx, int64(memInfo.Cached))
		sm.memoryBuffers.Record(ctx, int64(memInfo.Buffers))
	}

	if swapInfo, err := mem.SwapMemory(); err == nil {
		sm.memorySwapTotal.Record(ctx, int64(swapInfo.Total))
		sm.memorySwapUsed.Record(ctx, int64(swapInfo.Used))
	}

	// Disk Metrics
	if diskUsage, err := disk.Usage("/"); err == nil {
		sm.diskTotal.Record(ctx, int64(diskUsage.Total))
		sm.diskUsed.Record(ctx, int64(diskUsage.Used))
		sm.diskFree.Record(ctx, int64(diskUsage.Free))
		sm.diskUsagePercent.Record(ctx, diskUsage.UsedPercent)
	}

	if diskIO, err := disk.IOCounters(); err == nil {
		for _, io := range diskIO {
			sm.diskReadBytes.Add(ctx, int64(io.ReadBytes))
			sm.diskWriteBytes.Add(ctx, int64(io.WriteBytes))
		}
	}

	// Network Metrics
	if netIO, err := net.IOCounters(true); err == nil {
		for _, io := range netIO {
			sm.networkBytesReceived.Add(ctx, int64(io.BytesRecv))
			sm.networkBytesSent.Add(ctx, int64(io.BytesSent))
			sm.networkPacketsRecv.Add(ctx, int64(io.PacketsRecv))
			sm.networkPacketsSent.Add(ctx, int64(io.PacketsSent))
		}
	}

	// Go Runtime Metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sm.goMemoryUsage.Add(ctx, int64(m.Alloc))
	sm.goHeapSize.Record(ctx, int64(m.HeapSys))
	sm.goStackSize.Record(ctx, int64(m.StackSys))
	sm.goGoroutineCount.Record(ctx, int64(runtime.NumGoroutine()))
	sm.goGcDuration.Record(ctx, float64(m.PauseTotalNs)/1e9)
}

func (sm *SystemMetrics) StartSystemMonitoring() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			sm.RecordSystemMetrics()
		}
	}()
}

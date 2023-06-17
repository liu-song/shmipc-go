package shmipc

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

type PrometheusMonitor struct {
	receiveSyncEventCount prometheus.Gauge
	sendSyncEventCount    prometheus.Gauge
	outFlowBytes          prometheus.Gauge
	inFlowBytes           prometheus.Gauge
	sendQueueCount        prometheus.Gauge
	receiveQueueCount     prometheus.Gauge

	allocShmErrorCount     prometheus.Gauge
	fallbackWriteCount     prometheus.Gauge
	fallbackReadCount      prometheus.Gauge
	eventConnErrorCount    prometheus.Gauge
	queueFullErrorCount    prometheus.Gauge
	activeStreamCount      prometheus.Gauge
	hotRestartSuccessCount prometheus.Gauge
	hotRestartErrorCount   prometheus.Gauge

	capacityOfShareMemory prometheus.Gauge
	allInUsedShareMemory  prometheus.Gauge
}

func NewPrometheusMonitor(addr, path string) *PrometheusMonitor {

	registry := prom.NewRegistry()

	http.Handle(path, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}))
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatal("Unable to start a promhttp server, err: " + err.Error())
		}
	}()

	receiveSyncEventCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "receive_sync_event_count",
		Help: "The SyncEvent count that session had received",
	})
	sendSyncEventCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "send_sync_event_count",
		Help: "The SyncEvent count that session had sent",
	})
	outFlowBytes := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "out_flow_bytes",
		Help: "The out flow in bytes that session had sent",
	})
	inFlowBytes := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "in_flow_bytes",
		Help: "The in flow in bytes that session had receive",
	})
	sendQueueCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "send_queue_count",
		Help: "The pending count of send queue",
	})
	receiveQueueCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "receive_queue_count",
		Help: "The pending count of receive queue",
	})
	allocShmErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "alloc_shm_error_count",
		Help: "The error count of allocating share memory",
	})
	fallbackWriteCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fallback_write_count",
		Help: "The count of the fallback data write to unix/tcp connection",
	})
	fallbackReadCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "fallback_read_count",
		Help: "The error count of receiving fallback data from unix/tcp connection every period",
	})
	eventConnErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "event_conn_error_count",
		Help: "The error count of unix/tcp connection which usually happened in that the peer's process exit(crashed or other reason)",
	})
	queueFullErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "queue_full_error_count",
		Help: "The error count due to the IO-Queue(SendQueue or ReceiveQueue) is full which usually happened in that the peer was busy",
	})
	activeStreamCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_stream_count",
		Help: "Current all active stream count",
	})
	hotRestartSuccessCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hot_restart_success_count",
		Help: "The successful count of hot restart",
	})
	hotRestartErrorCount := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hot_restart_error_count",
		Help: "The failed count of hot restart",
	})
	capacityOfShareMemory := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "capacity_of_share_memory_in_bytes",
		Help: "Capacity of all share memory",
	})
	allInUsedShareMemory := promauto.NewGauge(prometheus.GaugeOpts{
		Name: "all_in_used_share_memory_in_bytes",
		Help: "Current in-used share memory",
	})
	registry.MustRegister(
		receiveSyncEventCount,
		sendSyncEventCount,
		outFlowBytes,
		inFlowBytes,
		sendQueueCount,
		receiveQueueCount,

		allocShmErrorCount,
		fallbackWriteCount,
		fallbackReadCount,
		eventConnErrorCount,
		queueFullErrorCount,
		activeStreamCount,
		hotRestartSuccessCount,
		hotRestartErrorCount,

		capacityOfShareMemory,
		allInUsedShareMemory,
	)
	return &PrometheusMonitor{
		receiveSyncEventCount: receiveSyncEventCount,
		sendSyncEventCount:    sendSyncEventCount,
		outFlowBytes:          outFlowBytes,
		inFlowBytes:           inFlowBytes,
		sendQueueCount:        sendQueueCount,
		receiveQueueCount:     receiveQueueCount,

		allocShmErrorCount:     allocShmErrorCount,
		fallbackWriteCount:     fallbackWriteCount,
		fallbackReadCount:      fallbackReadCount,
		eventConnErrorCount:    eventConnErrorCount,
		queueFullErrorCount:    queueFullErrorCount,
		activeStreamCount:      activeStreamCount,
		hotRestartSuccessCount: hotRestartSuccessCount,
		hotRestartErrorCount:   hotRestartErrorCount,

		capacityOfShareMemory: capacityOfShareMemory,
		allInUsedShareMemory:  allInUsedShareMemory,
	}
}

func (m *PrometheusMonitor) OnEmitSessionMetrics(performanceMetrics PerformanceMetrics, stabilityMetrics StabilityMetrics, shareMemoryMetrics ShareMemoryMetrics, session *Session) {

	m.receiveSyncEventCount.Set(float64(performanceMetrics.ReceiveSyncEventCount))
	m.sendSyncEventCount.Set(float64(performanceMetrics.SendSyncEventCount))
	m.outFlowBytes.Set(float64(performanceMetrics.OutFlowBytes))
	m.inFlowBytes.Set(float64(performanceMetrics.InFlowBytes))
	m.sendQueueCount.Set(float64(performanceMetrics.SendQueueCount))
	m.receiveQueueCount.Set(float64(performanceMetrics.ReceiveQueueCount))

	m.allocShmErrorCount.Set(float64(stabilityMetrics.AllocShmErrorCount))
	m.fallbackWriteCount.Set(float64(stabilityMetrics.FallbackWriteCount))
	m.fallbackReadCount.Set(float64(stabilityMetrics.FallbackReadCount))
	m.eventConnErrorCount.Set(float64(stabilityMetrics.EventConnErrorCount))
	m.queueFullErrorCount.Set(float64(stabilityMetrics.QueueFullErrorCount))
	m.activeStreamCount.Set(float64(stabilityMetrics.ActiveStreamCount))
	m.hotRestartSuccessCount.Set(float64(stabilityMetrics.HotRestartSuccessCount))
	m.hotRestartErrorCount.Set(float64(stabilityMetrics.HotRestartErrorCount))

	m.capacityOfShareMemory.Set(float64(shareMemoryMetrics.CapacityOfShareMemoryInBytes))
	m.allInUsedShareMemory.Set(float64(shareMemoryMetrics.AllInUsedShareMemoryInBytes))
}

func (m *PrometheusMonitor) Flush() error {
	// Write metrics to a log file
	file, err := os.OpenFile("metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write metrics to the log file
	_, err = fmt.Fprintf(file, "Flushed metrics at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

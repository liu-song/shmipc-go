package metrics

import (
	"fmt"
	"github.com/cloudwego/shmipc-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"os"
	"time"
)

type PrometheusMonitor struct {
	// Some fields here
	receiveSyncEventCount  prometheus.Counter
	sendSyncEventCount     prometheus.Counter
	outFlowBytes           prometheus.Counter
	inFlowBytes            prometheus.Counter
	sendQueueCount         prometheus.Gauge
	receiveQueueCount      prometheus.Gauge
	allocShmErrorCount     prometheus.Counter
	fallbackWriteCount     prometheus.Counter
	fallbackReadCount      prometheus.Counter
	eventConnErrorCount    prometheus.Counter
	queueFullErrorCount    prometheus.Counter
	activeStreamCount      prometheus.Gauge
	hotRestartSuccessCount prometheus.Counter
	hotRestartErrorCount   prometheus.Counter
	capacityOfShareMemory  prometheus.Gauge
	allInUsedShareMemory   prometheus.Gauge
}

func NewPrometheusMonitor() *PrometheusMonitor {
	return &PrometheusMonitor{
		receiveSyncEventCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "receive_sync_event_count",
			Help: "The SyncEvent count that session had received",
		}),
		sendSyncEventCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "send_sync_event_count",
			Help: "The SyncEvent count that session had sent",
		}),
		outFlowBytes: promauto.NewCounter(prometheus.CounterOpts{
			Name: "out_flow_bytes",
			Help: "The out flow in bytes that session had sent",
		}),
		inFlowBytes: promauto.NewCounter(prometheus.CounterOpts{
			Name: "in_flow_bytes",
			Help: "The in flow in bytes that session had receive",
		}),
		sendQueueCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "send_queue_count",
			Help: "The pending count of send queue",
		}),
		receiveQueueCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "receive_queue_count",
			Help: "The pending count of receive queue",
		}),
		allocShmErrorCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "alloc_shm_error_count",
			Help: "The error count of allocating share memory",
		}),
		fallbackWriteCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "fallback_write_count",
			Help: "The count of the fallback data write to unix/tcp connection",
		}),
		fallbackReadCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "fallback_read_count",
			Help: "The error count of receiving fallback data from unix/tcp connection every period",
		}),
		eventConnErrorCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "event_conn_error_count",
			Help: "The error count of unix/tcp connection which usually happened in that the peer's process exit(crashed or other reason)",
		}),
		queueFullErrorCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "queue_full_error_count",
			Help: "The error count due to the IO-Queue(SendQueue or ReceiveQueue) is full which usually happened in that the peer was busy",
		}),
		activeStreamCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "active_stream_count",
			Help: "Current all active stream count",
		}),
		hotRestartSuccessCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "hot_restart_success_count",
			Help: "The successful count of hot restart",
		}),
		hotRestartErrorCount: promauto.NewCounter(prometheus.CounterOpts{
			Name: "hot_restart_error_count",
			Help: "The failed count of hot restart",
		}),
		capacityOfShareMemory: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "capacity_of_share_memory_in_bytes",
			Help: "Capacity of all share memory",
		}),
		allInUsedShareMemory: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "all_in_used_share_memory_in_bytes",
			Help: "Current in-used share memory",
		}),
	}
}

func (m *PrometheusMonitor) OnEmitSessionMetrics(performanceMetrics shmipc.PerformanceMetrics, stabilityMetrics shmipc.StabilityMetrics, shareMemoryMetrics shmipc.ShareMemoryMetrics, session *shmipc.Session) {
	m.receiveSyncEventCount.Add(float64(performanceMetrics.ReceiveSyncEventCount))
	m.sendSyncEventCount.Add(float64(performanceMetrics.SendSyncEventCount))
	m.outFlowBytes.Add(float64(performanceMetrics.OutFlowBytes))
	m.inFlowBytes.Add(float64(performanceMetrics.InFlowBytes))
	m.sendQueueCount.Set(float64(performanceMetrics.SendQueueCount))
	m.receiveQueueCount.Set(float64(performanceMetrics.ReceiveQueueCount))
	m.allocShmErrorCount.Add(float64(stabilityMetrics.AllocShmErrorCount))
	m.fallbackWriteCount.Add(float64(stabilityMetrics.FallbackWriteCount))
	m.fallbackReadCount.Add(float64(stabilityMetrics.FallbackReadCount))
	m.eventConnErrorCount.Add(float64(stabilityMetrics.EventConnErrorCount))
	m.queueFullErrorCount.Add(float64(stabilityMetrics.QueueFullErrorCount))
	m.activeStreamCount.Set(float64(stabilityMetrics.ActiveStreamCount))
	m.hotRestartSuccessCount.Add(float64(stabilityMetrics.HotRestartSuccessCount))
	m.hotRestartErrorCount.Add(float64(stabilityMetrics.HotRestartErrorCount))
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

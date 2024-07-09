package generator

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"time"
)

type Generator struct {
	receiverAddr    string
	ticker          *time.Ticker
	requestsTotal   prometheus.Counter
	requestDuration prometheus.Histogram
}

func (g *Generator) Start() {
	fmt.Println("Starting generator")
	go func() {
		for {
			select {
			case <-g.ticker.C:
				g.probe()

			}
		}
	}()
}

func (g *Generator) probe() {
	start := time.Now()
	resp, err := http.Get(g.receiverAddr) // "http://localhost:9966")
	total := time.Since(start)

	// TODO: add labels to the metrics, based on json
	g.requestsTotal.Inc()
	g.requestDuration.Observe(total.Seconds())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Status)
	}
	fmt.Println(total)
}

func StartGenerator() {
	g := &Generator{
		receiverAddr: "http://localhost:9966",
	}
	fmt.Printf("Starting generator %v", g)
}

func NewGenerator(receiverAddr string, interval time.Duration) *Generator {
	return &Generator{
		receiverAddr: receiverAddr,
		ticker:       time.NewTicker(interval),
		requestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "generator_requests_total",
			Help: "The total number of requests sent by the generator",
		}),
		requestDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "generator_request_duration_seconds",
			Help: "The total number of requests sent by the generator",
		}),
	}
}

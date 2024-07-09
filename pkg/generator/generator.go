package generator

import (
	"encoding/json"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"time"
)

type ReceiverResponse struct {
	ClusterName, Node, Ip, Zone, Region string
}

type Generator struct {
	receiverAddr    string
	ticker          *time.Ticker
	testName        string
	clusterName     string
	region          string
	zone            string
	node            string
	requestsTotal   prometheus.Counter
	requestDuration *prometheus.HistogramVec
}

func (g *Generator) Start() {
	go func() {
		for {
			select {
			case <-g.ticker.C:
				go g.probe()

			}
		}
	}()
}

func (g *Generator) probe() {
	var r ReceiverResponse

	start := time.Now()
	resp, err := http.Get(g.receiverAddr)
	total := time.Since(start)

	// TODO: add labels to the metrics, based on json
	g.requestsTotal.Inc()
	labels := prometheus.Labels{
		"test":     g.testName,
		"cluster":  g.clusterName,
		"g_region": g.region,
		"g_zone":   g.zone,
		"g_node":   g.node,
		"r_region": "-",
		"r_zone":   "-",
		"r_node":   "-",
	}
	if err != nil {
		labels["status"] = "error"
	} else {
		labels["status"] = strconv.Itoa(resp.StatusCode)
		body := json.NewDecoder(resp.Body)
		err := body.Decode(&r)
		if err != nil {
			labels["status"] = "error-invalid-response"
			labels["r_region"] = r.Region
			labels["r_zone"] = r.Zone
			labels["r_node"] = r.Node
		}
	}
	g.requestDuration.With(labels).Observe(total.Seconds())
}

func NewGenerator(receiverAddr string, testName string, clusterName string, region string, zone string, node string, interval time.Duration) *Generator {
	buckets := []float64{0.0005, 0.001, 0.005, 0.010, 0.020, 0.050, .1, .2, .75, 1, 2}
	return &Generator{
		receiverAddr: receiverAddr,
		ticker:       time.NewTicker(interval),
		testName:     testName,
		clusterName:  clusterName,
		region:       region,
		zone:         zone,
		node:         node,
		requestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "khcm_generator_requests_total",
			Help: "The total number of requests sent by the generator",
		}),
		requestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "khcm_generator_request_duration_seconds",
			Help:    "",
			Buckets: buckets,
		},
			[]string{"test", "cluster", "g_region", "g_zone", "g_node", "r_region", "r_zone", "r_node", "status"},
		),
	}
}

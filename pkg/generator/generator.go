package generator

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	maasClient "stash.grupa.onet/go/go-maas.git/client"
	"time"
	"github.com/tarfik/khttp-connectivity-monitor/pkg/utils"
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
    size            int64
	response_size   int64
	requestsTotal   prometheus.Counter
	requestDuration *prometheus.HistogramVec
	maas            *maasClient.MaasClient
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
	client := &http.Client{}
	req, err := http.NewRequest("POST", g.receiverAddr, bytes.NewBuffer(utils.RandStringBytes(g.size)))
	req.Header.Set("Content-Type", "application/octet-stream")
    req.Header.Add("X-Khttp-Response-Size", strconv.FormatInt(g.response_size, 10))
    resp, err := client.Do(req)

	total := time.Since(start)
	maas_path := ""

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
		maas_path = fmt.Sprintf(
			"%s.%s.generator-%s.receiver-%s",
			g.testName, g.clusterName,
			strings.ReplaceAll(g.node, ".", "_"),
			"unknown",
		)
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
		maas_path = fmt.Sprintf(
			"%s.%s.generator-%s.receiver-%s",
			g.testName, g.clusterName,
			strings.ReplaceAll(g.node, ".", "_"),
			strings.ReplaceAll(r.Node, ".", "_"),
		)
	}
	g.maas.Number(fmt.Sprintf("%s.respose_time", maas_path), total.Seconds())
	g.maas.Count(fmt.Sprintf("%s.status.%s", maas_path, labels["status"]))
	g.requestDuration.With(labels).Observe(total.Seconds())
}

func NewGenerator(receiverAddr string, testName string, clusterName string, region string, zone string, node string, interval time.Duration, size int64, maas *maasClient.MaasClient) *Generator {
	buckets := []float64{0.0005, 0.001, 0.005, 0.010, 0.020, 0.050, .1, .2, .75, 1, 2}
	return &Generator{
		receiverAddr: receiverAddr,
		ticker:       time.NewTicker(interval),
		testName:     testName,
		clusterName:  clusterName,
		region:       region,
		zone:         zone,
		node:         node,
		maas:         maas,
		size:         size,
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

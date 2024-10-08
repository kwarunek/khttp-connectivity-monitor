package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kwarunek/khttp-connectivity-monitor/pkg/generator"
	"github.com/kwarunek/khttp-connectivity-monitor/pkg/receiver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version string
)

const (
	RECEIVER  = "receiver"
	GENERATOR = "generator"
)

func main() {
	log.SetLevel(log.InfoLevel)
	// a config watcher
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("khttp")
	viper.SetDefault("port", 9966)
	viper.SetDefault("size", 1500)
	viper.SetDefault("response_size", 1500)
	viper.SetDefault("host", "localhost")
	viper.AutomaticEnv()
	viper.SetDefault("testName", "test")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	addr := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	mode := viper.GetString("mode")

	if mode == RECEIVER {
		receiver.SetupReceiver(
			addr,
			viper.GetString("region"),
			viper.GetString("zone"),
			viper.GetString("clusterName"),
			viper.GetString("node"),
			viper.GetInt64("response_size"),
		)
	} else {
		interval, err := time.ParseDuration(viper.GetString("generatorInterval"))
		if err != nil {
			log.Fatalf("Failed to parse interval: %v", err)
		}
		g := generator.NewGenerator(
			viper.GetString("probedAddr"),
			viper.GetString("testName"),
			viper.GetString("clusterName"),
			viper.GetString("region"),
			viper.GetString("zone"),
			viper.GetString("node"),
			interval,
			viper.GetInt64("size"),
		)
		g.Start()
	}
	log.Infof("Starting [version: %s] in mode: %s on %s", version, mode, addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

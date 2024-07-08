package main

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/kwarunek/khttp-connectivity-monitor/pkg/generator"
	"github.com/kwarunek/khttp-connectivity-monitor/pkg/receiver"

)

var (
	version       string
)

const (
    RECEIVER = "receiver"
    GENERATOR = "generator"
)

func main() {
	log.SetLevel(log.InfoLevel)
	addr := flag.String("addr", "0.0.0.0:9966", "listen address")
	flag.Parse()

	// a config watcher
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
        viper.SetEnvPrefix("khttp")
        viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
        mode := viper.GetString("mode")
        log.Infof("Starting [version: %s] in mode: %s on %s", version, viper.GetString("mode"), *addr)

        if mode == RECEIVER {
            clusterName := viper.GetString("clusterName")
            podIp := viper.GetString("podIp")
            vmIp := viper.GetString("vmIp")
            hostIp := viper.GetString("hostIp")
            receiver.ServeReceiver(addr, clusterName, podIp, vmIp, hostIp)
        } else {
            generator.StartGenerator()
        }
}

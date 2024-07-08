package receiver

import (
	"flag"
	"fmt"
	"net/http"
    	viper "github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)


func ServeReceiver(addr String, clusterName String, podIp string, vmIp String, HostIp String) {

    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        json = fmt.Sprintf(`{"clusterName": "%s", "podIp": "%s", "vmIp": "%s", "hostIp": "%s"}`, clusterName, podIp, vmIp, HostIp)
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, json)
    })
   log.Fatal(http.ListenAndServe(*addr, nil))
}

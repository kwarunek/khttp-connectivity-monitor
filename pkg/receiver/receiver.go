package receiver

import (
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
)


func ServeReceiver(addr *string, clusterName string, podIp string, vmIp string, HostIp string) {

    json := fmt.Sprintf(`{"clusterName": "%s", "podIp": "%s", "vmIp": "%s", "hostIp": "%s"}`, clusterName, podIp, vmIp, HostIp)
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, json)
    })
   log.Fatal(http.ListenAndServe(*addr, nil))
}

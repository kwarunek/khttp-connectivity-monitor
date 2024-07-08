package receiver

import (
	"fmt"
	"net/http"
	log "github.com/sirupsen/logrus"
)


func ServeReceiver(addr *string, testName string, region string, zone string, clusterName string, node string, ip string) {

    json := fmt.Sprintf(`{"testName": "%s", "clusterName": "%s", "node": "%s", "ip": "%s", "zone": "%s", "region": "%s"}`, testName, clusterName, node, ip, zone, region)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, json)
    })
   log.Fatal(http.ListenAndServe(*addr, nil))
}

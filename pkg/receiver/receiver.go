package receiver

import (
	"fmt"
	"net/http"
)

func SetupReceiver(addr string, region string, zone string, clusterName string, node string, ip string) {

	json := fmt.Sprintf(`{"clusterName": "%s", "node": "%s", "ip": "%s", "zone": "%s", "region": "%s"}`, clusterName, node, ip, zone, region)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, json)
	})
}

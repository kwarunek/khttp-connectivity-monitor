package receiver

import (
	"fmt"
	"net/http"
	"github.com/tarfik/khttp-connectivity-monitor/pkg/utils"
)

func SetupReceiver(addr string, region string, zone string, clusterName string, node string, reponse_size int64) {

	json := fmt.Sprintf(`{"clusterName": "%s", "node": "%s", "zone": "%s", "region": "%s", "response": "%s"}`, clusterName, node, zone, region, bytes.NewBuffer(RandStringBytes(reponse_size))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, json)
	})
}
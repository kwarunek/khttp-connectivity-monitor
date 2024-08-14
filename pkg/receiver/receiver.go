package receiver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/kwarunek/khttp-connectivity-monitor/pkg/utils"
)

func SetupReceiver(addr string, region string, zone string, clusterName string, node string, reponse_size int64) {

	json := fmt.Sprintf(`{"clusterName": "%s", "node": "%s", "zone": "%s", "region": "%s", "response": "RESPONSE"}`, clusterName, node, zone, region)
	rand_bytes_cnt := reponse_size - int64(len(json)) + int64(len("RESPONSE"))
	if rand_bytes_cnt > 0 {
		json = strings.Replace(json, "RESPONSE", string(utils.RandStringBytes(rand_bytes_cnt)), -1)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, json)
	})
}

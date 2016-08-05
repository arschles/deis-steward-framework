package brokerapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deis/steward/mode"
	"github.com/deis/steward/web"
	"github.com/gorilla/mux"
	"github.com/juju/loggo"
)

func provisioningHandler(logger loggo.Logger, provisioner mode.Provisioner, auth *web.BasicAuth) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		instanceID, ok := vars[instanceIDPathKey]
		if !ok {
			http.Error(w, "missing instance ID", http.StatusBadRequest)
			return
		}
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "authorization missing", http.StatusBadRequest)
			return
		}
		if username != auth.Username || password != auth.Password {
			http.Error(w, "wrong login credentials", http.StatusUnauthorized)
			return
		}
		provisionReq := new(mode.ProvisionRequest)
		if err := json.NewDecoder(r.Body).Decode(provisionReq); err != nil {
			logger.Debugf("error decoding provision request (%s)", err)
			http.Error(w, "error decoding provision request", http.StatusBadRequest)
			return
		}
		resp, err := provisioner.Provision(instanceID, provisionReq)
		if err != nil {
			logger.Debugf("error provisioning (%s)", err)
			http.Error(w, "error provisioning", http.StatusInternalServerError)
			return
		}
		respStr := fmt.Sprintf(`{"operation":"%s"}`, resp.Operation)
		w.WriteHeader(resp.Status)
		w.Write([]byte(respStr))
	})
}
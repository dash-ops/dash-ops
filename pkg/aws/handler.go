package aws

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

// Account config
type Account struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func accountsHandler(dashConfig dashYaml) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var accounts []Account

		for _, account := range dashConfig.AWS {
			c := Account{
				Name: account.Name,
				Key:  commons.UnderScoreString(account.Name),
			}

			accounts = append(accounts, c)
		}

		commons.RespondJSON(w, http.StatusOK, accounts)
	}
}

func permissionsHandler(permission permission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commons.RespondJSON(w, http.StatusOK, permission)
	}
}

func ec2InstancesHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := client.GetInstances()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, instances)
	}
}

func ec2InstanceStartHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data, err := client.StartInstance(vars["instanceId"])
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, data)
	}
}

func ec2InstanceStopHandler(client Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data, err := client.StopInstance(vars["instanceId"])
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, data)
	}
}

// MakeAWSInstanceHandlers Add aws module endpoints
func MakeAWSInstanceHandlers(r *mux.Router, fileConfig []byte) {
	dashConfig := loadConfig(fileConfig)

	r.HandleFunc("/aws/accounts", accountsHandler(dashConfig)).
		Methods("GET", "OPTIONS").
		Name("awsAccounts")

	for _, account := range dashConfig.AWS {
		client, err := NewClient(account)
		if err != nil {
			fmt.Println(err.Error())
		}

		accountRoute := r.PathPrefix("/aws/" + commons.UnderScoreString(account.Name)).Subrouter()

		accountRoute.HandleFunc("/permissions", permissionsHandler(account.Permission)).
			Methods("GET", "OPTIONS").
			Name("k8sPermissions")

		accountRoute.HandleFunc("/ec2/instances", ec2InstancesHandler(client)).
			Methods("GET", "OPTIONS").
			Name("ec2Instances")

		accountRoute.HandleFunc("/ec2/instance/start/{instanceId}", ec2InstanceStartHandler(client)).
			Methods("POST", "OPTIONS").
			Name("ec2InstanceStart")

		accountRoute.HandleFunc("/ec2/instance/stop/{instanceId}", ec2InstanceStopHandler(client)).
			Methods("POST", "OPTIONS").
			Name("ec2InstanceStop")
	}
}

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

func awsAccountsHandler(dashConfig dashYaml) http.HandlerFunc {
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

func awsPermissionsHandler(permission awsPermission) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		commons.RespondJSON(w, http.StatusOK, permission)
	}
}

func ec2InstancesHandler(awsClient AwsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := awsClient.GetInstances()
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, instances)
	}
}

func ec2InstanceStartHandler(awsClient AwsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data, err := awsClient.StartInstance(vars["instanceId"])
		if err != nil {
			commons.RespondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		commons.RespondJSON(w, http.StatusOK, data)
	}
}

func ec2InstanceStopHandler(awsClient AwsClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data, err := awsClient.StopInstance(vars["instanceId"])
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

	r.HandleFunc("/aws/accounts", awsAccountsHandler(dashConfig)).
		Methods("GET", "OPTIONS").
		Name("awsAccounts")

	for _, account := range dashConfig.AWS {
		awsClient, err := NewAwsClient(account)
		if err != nil {
			fmt.Println(err.Error())
		}

		accountRoute := r.PathPrefix("/aws/" + commons.UnderScoreString(account.Name)).Subrouter()

		accountRoute.HandleFunc("/permissions", awsPermissionsHandler(account.Permission)).
			Methods("GET", "OPTIONS").
			Name("k8sPermissions")

		accountRoute.HandleFunc("/ec2/instances", ec2InstancesHandler(awsClient)).
			Methods("GET", "OPTIONS").
			Name("ec2Instances")

		accountRoute.HandleFunc("/ec2/instance/start/{instanceId}", ec2InstanceStartHandler(awsClient)).
			Methods("POST", "OPTIONS").
			Name("ec2InstanceStart")

		accountRoute.HandleFunc("/ec2/instance/stop/{instanceId}", ec2InstanceStopHandler(awsClient)).
			Methods("POST", "OPTIONS").
			Name("ec2InstanceStop")
	}
}

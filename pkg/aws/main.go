package aws

import (
	"fmt"
	"net/http"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
)

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

	awsClient, err := NewAwsClient(dashConfig.AWS)
	if err != nil {
		fmt.Println(err.Error())
	}
	r.HandleFunc("/v1/ec2/instances", ec2InstancesHandler(awsClient)).
		Methods("GET", "OPTIONS").
		Name("ec2Instances")

	r.HandleFunc("/v1/ec2/instance/start/{instanceId}", ec2InstanceStartHandler(awsClient)).
		Methods("POST", "OPTIONS").
		Name("ec2InstanceStart")

	r.HandleFunc("/v1/ec2/instance/stop/{instanceId}", ec2InstanceStopHandler(awsClient)).
		Methods("POST", "OPTIONS").
		Name("ec2InstanceStop")
}

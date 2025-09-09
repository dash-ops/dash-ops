package aws

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dash-ops/dash-ops/pkg/commons"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) GetInstances() ([]Instance, error) {
	args := m.Called()
	return args.Get(0).([]Instance), args.Error(1)
}

func (m *mockClient) StartInstance(instanceID string) (InstanceOutput, error) {
	args := m.Called(instanceID)
	return args.Get(0).(InstanceOutput), args.Error(1)
}

func (m *mockClient) StopInstance(instanceID string) (InstanceOutput, error) {
	args := m.Called(instanceID)
	return args.Get(0).(InstanceOutput), args.Error(1)
}

func TestAccountsHandler(t *testing.T) {
	mockConfig := dashYaml{
		AWS: []config{
			{Name: "AWS PROD"},
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/accounts", nil)

	handler := accountsHandler(mockConfig)
	handler.ServeHTTP(rr, req)

	var accounts []Account
	json.NewDecoder(rr.Body).Decode(&accounts)

	assert.Equal(t, "AWS PROD", accounts[0].Name, "return name account")
	assert.Equal(t, "aws_prod", accounts[0].Key, "return key account")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestPermissionsHandler(t *testing.T) {
	mockPermission := permission{
		EC2: ec2Permissions{
			Start: []string{"bla"},
			Stop:  []string{"ble"},
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/permissions", nil)

	handler := permissionsHandler(mockPermission)
	handler.ServeHTTP(rr, req)

	var resultPermission permission
	json.NewDecoder(rr.Body).Decode(&resultPermission)

	assert.Equal(t, mockPermission, resultPermission, "return permissions aws plugin")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestEc2InstancesHandler(t *testing.T) {
	mockInstances := []Instance{
		{InstanceID: "111", Name: "Test"},
		{InstanceID: "222", Name: "Test2"},
	}

	client := new(mockClient)
	client.On("GetInstances").Return(mockInstances, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/instances", nil)

	handler := ec2InstancesHandler(client)
	handler.ServeHTTP(rr, req)

	var instances []Instance
	json.NewDecoder(rr.Body).Decode(&instances)

	assert.Equal(t, mockInstances[0].Name, instances[0].Name, "return instance name in position 0")
	assert.Equal(t, mockInstances[1].Name, instances[1].Name, "return instance name in position 1")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestEc2InstancesHandlerWithError(t *testing.T) {
	mockErr := errors.New("message error")

	clientWithError := new(mockClient)
	clientWithError.On("GetInstances").Return([]Instance{}, mockErr)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/instances", nil)

	handlerWithError := ec2InstancesHandler(clientWithError)
	handlerWithError.ServeHTTP(rr, req)

	var respError commons.ResponseError
	json.NewDecoder(rr.Body).Decode(&respError)

	expecteData := commons.ResponseError{Error: mockErr.Error()}
	assert.Equal(t, expecteData, respError, "return error message")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return status 500")
}

func TestEc2InstanceStartHandler(t *testing.T) {
	mockOutput := InstanceOutput{
		CurrentState:  "",
		PreviousState: "",
	}

	client := new(mockClient)
	client.On("StartInstance", "111").Return(mockOutput, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/instance/start/111", nil)

	router := mux.NewRouter()
	router.HandleFunc("/instance/start/{instanceId}", ec2InstanceStartHandler(client))
	router.ServeHTTP(rr, req)

	var output InstanceOutput
	json.NewDecoder(rr.Body).Decode(&output)

	assert.Equal(t, mockOutput.CurrentState, output.CurrentState, "return current state of the instance")
	assert.Equal(t, mockOutput.PreviousState, output.PreviousState, "return previous state of the instance")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestEc2InstanceStartHandlerWithError(t *testing.T) {
	mockErr := errors.New("message error")

	clientWithError := new(mockClient)
	clientWithError.On("StartInstance", "111").Return(InstanceOutput{}, mockErr)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/instance/start/111", nil)

	router := mux.NewRouter()
	router.HandleFunc("/instance/start/{instanceId}", ec2InstanceStartHandler(clientWithError))
	router.ServeHTTP(rr, req)

	var respError commons.ResponseError
	json.NewDecoder(rr.Body).Decode(&respError)

	expecteData := commons.ResponseError{Error: mockErr.Error()}
	assert.Equal(t, expecteData, respError, "return error message")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return status 500")
}

func TestEc2InstanceStopHandler(t *testing.T) {
	mockOutput := InstanceOutput{
		CurrentState:  "",
		PreviousState: "",
	}

	client := new(mockClient)
	client.On("StopInstance", "111").Return(mockOutput, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/instance/stop/111", nil)

	router := mux.NewRouter()
	router.HandleFunc("/instance/stop/{instanceId}", ec2InstanceStopHandler(client))
	router.ServeHTTP(rr, req)

	var output InstanceOutput
	json.NewDecoder(rr.Body).Decode(&output)

	assert.Equal(t, mockOutput.CurrentState, output.CurrentState, "return current state of the instance")
	assert.Equal(t, mockOutput.PreviousState, output.PreviousState, "return previous state of the instance")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestEc2InstanceStopHandlerWithError(t *testing.T) {
	mockErr := errors.New("message error")

	clientWithError := new(mockClient)
	clientWithError.On("StopInstance", "111").Return(InstanceOutput{}, mockErr)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/instance/stop/111", nil)

	router := mux.NewRouter()
	router.HandleFunc("/instance/stop/{instanceId}", ec2InstanceStopHandler(clientWithError))
	router.ServeHTTP(rr, req)

	var respError commons.ResponseError
	json.NewDecoder(rr.Body).Decode(&respError)

	expecteData := commons.ResponseError{Error: mockErr.Error()}
	assert.Equal(t, expecteData, respError, "return error message")
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return status 500")
}

func TestMakeAWSInstanceHandlers(t *testing.T) {
	fileConfig := []byte(`aws:
  - name: 'AWS Test'
    region: us-east-1
    accessKeyId: 1234
    secretAccessKey: 4321
    ec2Config:
      skipList:
        - "test"`)

	r := mux.NewRouter()
	MakeAWSInstanceHandlers(r, fileConfig)

	path, err := r.GetRoute("awsAccounts").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/aws/accounts", path)
	path, err = r.GetRoute("awsPermissions").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/aws/aws_test/permissions", path)
}

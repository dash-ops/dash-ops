package kubernetes

import (
	"context"
	"encoding/json"
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

func (m *mockClient) GetNodes() ([]Node, error) {
	args := m.Called()
	return args.Get(0).([]Node), args.Error(1)
}

func (m *mockClient) GetNamespaces() ([]Namespace, error) {
	args := m.Called()
	return args.Get(0).([]Namespace), args.Error(1)
}

func (m *mockClient) GetDeployments(filters deploymentFilter) ([]Deployment, error) {
	args := m.Called()
	return args.Get(0).([]Deployment), args.Error(1)
}

func (m *mockClient) Scale(name string, ns string, replicas int32) error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockClient) GetPods(filters podFilter) ([]Pod, error) {
	args := m.Called()
	return args.Get(0).([]Pod), args.Error(1)
}

func (m *mockClient) GetPodLogs(filters podFilter) ([]ContainerLog, error) {
	args := m.Called()
	return args.Get(0).([]ContainerLog), args.Error(1)
}

func (m *mockClient) RestartDeployment(name string, ns string) error {
	args := m.Called()
	return args.Error(0)
}

func TestClustersHandler(t *testing.T) {
	mockConfig := dashYaml{
		Kubernetes: []config{
			{Name: "Kube Prod"},
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/clusters", nil)

	handler := clustersHandler(mockConfig)
	handler.ServeHTTP(rr, req)

	var clusters []Cluster
	json.NewDecoder(rr.Body).Decode(&clusters)

	assert.Equal(t, "Kube Prod", clusters[0].Name, "return name clusters")
	assert.Equal(t, "default", clusters[0].Context, "return context clusters")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestPermissionsHandler(t *testing.T) {
	mockPermission := permission{
		Deployments: deploymentsPermissions{
			Namespaces: []string{"apps"},
			Restart:    []string{"bla"},
			Scale:      []string{"ble"},
		},
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/permissions", nil)

	handler := permissionsHandler(mockPermission)
	handler.ServeHTTP(rr, req)

	var resultPermission permission
	json.NewDecoder(rr.Body).Decode(&resultPermission)

	assert.Equal(t, mockPermission, resultPermission, "return permissions kubernetes plugin")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestNodesHandler(t *testing.T) {
	mockNodes := []Node{
		{Name: "Test"},
	}

	client := new(mockClient)
	client.On("GetNodes").Return(mockNodes, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/nodes", nil)

	handler := nodesHandler(client)
	handler.ServeHTTP(rr, req)

	var nodes []Node
	json.NewDecoder(rr.Body).Decode(&nodes)

	assert.Equal(t, mockNodes, nodes, "return nodes kubernetes")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestNamespacesHandler(t *testing.T) {
	mockNamespaces := []Namespace{
		{Name: "default"},
	}

	client := new(mockClient)
	client.On("GetNamespaces").Return(mockNamespaces, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/namespaces", nil)

	handler := namespacesHandler(client)
	handler.ServeHTTP(rr, req)

	var namespaces []Namespace
	json.NewDecoder(rr.Body).Decode(&namespaces)

	assert.Equal(t, mockNamespaces, namespaces, "return namespaces kubernetes")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestDeploymentsHandler(t *testing.T) {
	mockDeployments := []Deployment{
		{Name: "project01", Namespace: "default"},
	}

	client := new(mockClient)
	client.On("GetDeployments").Return(mockDeployments, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/deployments", nil)

	handler := deploymentsHandler(client)
	handler.ServeHTTP(rr, req)

	var deployments []Deployment
	json.NewDecoder(rr.Body).Decode(&deployments)

	assert.Equal(t, mockDeployments, deployments, "return deployments kubernetes")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestDeploymentRestartHandler(t *testing.T) {
	mockPermission := permission{
		Deployments: deploymentsPermissions{
			Namespaces: []string{"apps"},
			Restart:    []string{"bla"},
		},
	}

	mockUserData := commons.UserData{
		Groups: []string{"bla"},
	}

	client := new(mockClient)
	client.On("RestartDeployment").Return(nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/deployment/restart/apps/nginx", nil)
	ctx := context.WithValue(req.Context(), commons.UserDataKey, mockUserData)
	req = req.WithContext(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/deployment/restart/{namespace}/{name}", deploymentRestartHandler(client, mockPermission))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestDeploymentScaleHandler(t *testing.T) {
	mockPermission := permission{
		Deployments: deploymentsPermissions{
			Namespaces: []string{"apps"},
			Scale:      []string{"bla"},
		},
	}

	mockUserData := commons.UserData{
		Groups: []string{"bla"},
	}

	client := new(mockClient)
	client.On("Scale").Return(nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/deployment/scale/apps/nginx/3", nil)
	ctx := context.WithValue(req.Context(), commons.UserDataKey, mockUserData)
	req = req.WithContext(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/deployment/scale/{namespace}/{name}/{replicas}", deploymentScaleHandler(client, mockPermission))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestPodsHandler(t *testing.T) {
	mockPods := []Pod{
		{Name: "project01",
			Namespace: "default"},
	}

	client := new(mockClient)
	client.On("GetPods").Return(mockPods, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/pods", nil)

	handler := podsHandler(client)
	handler.ServeHTTP(rr, req)

	var pods []Pod
	json.NewDecoder(rr.Body).Decode(&pods)

	assert.Equal(t, mockPods, pods, "return pods kubernetes")
	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestPodLogsHandler(t *testing.T) {
	client := new(mockClient)
	client.On("GetPodLogs").Return([]ContainerLog{
		{
			Name: "Apps-001",
			Log:  "xpto",
		},
	}, nil)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/pod/nginx/logs?namespace=default", nil)

	router := mux.NewRouter()
	router.HandleFunc("/pod/{name}/logs", podLogsHandler(client)).
		Queries("namespace", "{namespace}")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "should return status 200")
}

func TestMakeKubernetesHandlers(t *testing.T) {
	fileConfig := []byte(`kubernetes:
  - name: 'Kube Test'
    kubeconfig: /root/.kube/config
    context: xpto`)

	r := mux.NewRouter()
	MakeKubernetesHandlers(r, fileConfig)

	path, err := r.GetRoute("k8sClusters").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/k8s/clusters", path)
	path, err = r.GetRoute("k8sPermissions").GetPathTemplate()
	assert.Nil(t, err)
	assert.Equal(t, "/k8s/xpto/permissions", path)
}

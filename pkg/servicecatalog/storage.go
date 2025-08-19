package servicecatalog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Storage interface for service catalog
type Storage interface {
	ListServices(filter ServiceFilter) ([]ServiceSummary, error)
	GetService(id string) (*Service, error)
	CreateService(req CreateServiceRequest) (*Service, error)
	UpdateService(id string, service *Service) error
	DeleteService(id string) error
}

// FileStorage implements Storage interface using file system
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(basePath string) *FileStorage {
	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		fmt.Printf("Warning: Could not create catalog directory: %v\n", err)
	}
	return &FileStorage{basePath: basePath}
}

// ListServices returns all services with optional filtering
func (fs *FileStorage) ListServices(filter ServiceFilter) ([]ServiceSummary, error) {
	var services []ServiceSummary

	// Read all .json files in the catalog directory
	files, err := ioutil.ReadDir(fs.basePath)
	if err != nil {
		return services, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		service, err := fs.readServiceFile(filepath.Join(fs.basePath, file.Name()))
		if err != nil {
			fmt.Printf("Warning: Could not read service file %s: %v\n", file.Name(), err)
			continue
		}

		// Apply filters
		if fs.matchesFilter(service, filter) {
			summary := ServiceSummary{
				ID:          service.ID,
				Name:        service.Name,
				DisplayName: service.DisplayName,
				Description: service.Description,
				Tier:        service.Tier,
				Team:        service.Team,
				Squad:       service.Squad,
				Tags:        service.Tags,
				Regions:     service.Regions,
				IngressType: service.IngressType,
				Status:      service.Status,
				UpdatedAt:   service.UpdatedAt,
			}
			services = append(services, summary)
		}
	}

	return services, nil
}

// GetService returns a specific service by ID
func (fs *FileStorage) GetService(id string) (*Service, error) {
	filename := filepath.Join(fs.basePath, id+".json")
	return fs.readServiceFile(filename)
}

// CreateService creates a new service
func (fs *FileStorage) CreateService(req CreateServiceRequest) (*Service, error) {
	service := &Service{
		ID:          uuid.New().String(),
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Tier:        req.Tier,
		Team:        req.Team,
		Squad:       req.Squad,
		Tags:        req.Tags,
		Regions:     []string{}, // Will be populated later
		IngressType: "internal", // Default
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}

	if req.DisplayName == "" {
		service.DisplayName = req.Name
	}

	return service, fs.writeServiceFile(service)
}

// UpdateService updates an existing service
func (fs *FileStorage) UpdateService(id string, service *Service) error {
	service.UpdatedAt = time.Now()
	return fs.writeServiceFile(service)
}

// DeleteService deletes a service
func (fs *FileStorage) DeleteService(id string) error {
	filename := filepath.Join(fs.basePath, id+".json")
	return os.Remove(filename)
}

// Helper methods

func (fs *FileStorage) readServiceFile(filename string) (*Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var service Service
	err = json.Unmarshal(data, &service)
	return &service, err
}

func (fs *FileStorage) writeServiceFile(service *Service) error {
	filename := filepath.Join(fs.basePath, service.ID+".json")
	data, err := json.MarshalIndent(service, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (fs *FileStorage) matchesFilter(service *Service, filter ServiceFilter) bool {
	// Filter by tier
	if filter.Tier != "" && filter.Tier != "all" && service.Tier != filter.Tier {
		return false
	}

	// Filter by team
	if filter.Team != "" && service.Team != filter.Team {
		return false
	}

	// Filter by status
	if filter.Status != "" && service.Status != filter.Status {
		return false
	}

	// Filter by search (name, description, tags)
	if filter.Search != "" {
		search := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(service.Name), search) &&
			!strings.Contains(strings.ToLower(service.DisplayName), search) &&
			!strings.Contains(strings.ToLower(service.Description), search) {
			// Check tags
			found := false
			for _, tag := range service.Tags {
				if strings.Contains(strings.ToLower(tag), search) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}

	return true
}

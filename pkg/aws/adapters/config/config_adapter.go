package config

import (
	"fmt"

	"gopkg.in/yaml.v2"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// ConfigAdapter handles AWS configuration parsing
type ConfigAdapter struct{}

// NewConfigAdapter creates a new config adapter
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// ParseAWSConfigFromFileConfig parses AWS config from file bytes
func (ca *ConfigAdapter) ParseAWSConfigFromFileConfig(fileConfig []byte) ([]awsModels.AWSAccount, error) {
	var config struct {
		AWS []struct {
			Name            string `yaml:"name"`
			Region          string `yaml:"region"`
			AccessKeyID     string `yaml:"accessKeyId"`
			SecretAccessKey string `yaml:"secretAccessKey"`
			Permission      struct {
				EC2 struct {
					Start []string `yaml:"start"`
					Stop  []string `yaml:"stop"`
					View  []string `yaml:"view"`
				} `yaml:"ec2"`
			} `yaml:"permission"`
			EC2Config struct {
				SkipList    []string `yaml:"skipList"`
				DefaultTags []struct {
					Key   string `yaml:"key"`
					Value string `yaml:"value"`
				} `yaml:"defaultTags"`
			} `yaml:"ec2Config"`
		} `yaml:"aws"`
	}

	if err := yaml.Unmarshal(fileConfig, &config); err != nil {
		return nil, fmt.Errorf("failed to parse AWS configuration: %w", err)
	}

	if len(config.AWS) == 0 {
		return nil, fmt.Errorf("no AWS configuration found")
	}

	var accounts []awsModels.AWSAccount
	for _, awsConfig := range config.AWS {
		// Generate key from name
		key := generateAccountKey(awsConfig.Name)

		// Convert default tags
		var defaultTags []awsModels.Tag
		for _, tag := range awsConfig.EC2Config.DefaultTags {
			defaultTags = append(defaultTags, awsModels.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}

		account := awsModels.AWSAccount{
			Name:            awsConfig.Name,
			Key:             key,
			Region:          awsConfig.Region,
			AccessKeyID:     awsConfig.AccessKeyID,
			SecretAccessKey: awsConfig.SecretAccessKey,
			Permissions: awsModels.AccountPermissions{
				EC2: awsModels.EC2Permissions{
					Start: awsConfig.Permission.EC2.Start,
					Stop:  awsConfig.Permission.EC2.Stop,
					View:  awsConfig.Permission.EC2.View,
				},
			},
			EC2Config: awsModels.EC2Config{
				SkipList:    awsConfig.EC2Config.SkipList,
				DefaultTags: defaultTags,
			},
			Status: awsModels.AccountStatusUnknown, // Will be validated later
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// generateAccountKey generates a normalized key from account name
func generateAccountKey(name string) string {
	// Simple implementation - in production, this might be more sophisticated
	key := name
	// Replace spaces with underscores and convert to lowercase
	for i, char := range key {
		if char == ' ' {
			key = key[:i] + "_" + key[i+1:]
		}
	}
	return key
}

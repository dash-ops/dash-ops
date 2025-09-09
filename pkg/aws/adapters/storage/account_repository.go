package storage

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v2"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// AccountRepositoryAdapter implements AccountRepository interface using file-based storage
type AccountRepositoryAdapter struct {
	fileConfig []byte
	accounts   map[string]*awsModels.AWSAccount
}

// NewAccountRepositoryAdapter creates a new account repository adapter
func NewAccountRepositoryAdapter(fileConfig []byte) (*AccountRepositoryAdapter, error) {
	adapter := &AccountRepositoryAdapter{
		fileConfig: fileConfig,
		accounts:   make(map[string]*awsModels.AWSAccount),
	}

	// Parse configuration and load accounts
	if err := adapter.loadAccounts(); err != nil {
		return nil, fmt.Errorf("failed to load accounts: %w", err)
	}

	return adapter, nil
}

// GetAccount gets a specific AWS account configuration
func (ara *AccountRepositoryAdapter) GetAccount(ctx context.Context, accountKey string) (*awsModels.AWSAccount, error) {
	if accountKey == "" {
		return nil, fmt.Errorf("account key is required")
	}

	account, exists := ara.accounts[accountKey]
	if !exists {
		return nil, fmt.Errorf("account %s not found", accountKey)
	}

	return account, nil
}

// ListAccounts lists all configured AWS accounts
func (ara *AccountRepositoryAdapter) ListAccounts(ctx context.Context) ([]awsModels.AWSAccount, error) {
	var accounts []awsModels.AWSAccount
	for _, account := range ara.accounts {
		accounts = append(accounts, *account)
	}

	return accounts, nil
}

// ValidateAccount validates account credentials and connectivity
func (ara *AccountRepositoryAdapter) ValidateAccount(ctx context.Context, account *awsModels.AWSAccount) error {
	if account == nil {
		return fmt.Errorf("account cannot be nil")
	}

	// Basic validation
	if err := account.Validate(); err != nil {
		return fmt.Errorf("account validation failed: %w", err)
	}

	// TODO: Implement actual AWS credential validation
	// This would involve creating an AWS client and making a simple API call
	// For now, we'll just mark it as valid if basic validation passes
	account.Status = awsModels.AccountStatusActive
	account.Error = ""

	return nil
}

// UpdateAccountStatus updates account status and error information
func (ara *AccountRepositoryAdapter) UpdateAccountStatus(ctx context.Context, accountKey string, status awsModels.AccountStatus, errorMsg string) error {
	account, exists := ara.accounts[accountKey]
	if !exists {
		return fmt.Errorf("account %s not found", accountKey)
	}

	account.Status = status
	account.Error = errorMsg

	return nil
}

// loadAccounts loads accounts from file configuration
func (ara *AccountRepositoryAdapter) loadAccounts() error {
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

	if err := yaml.Unmarshal(ara.fileConfig, &config); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

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

		account := &awsModels.AWSAccount{
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

		ara.accounts[key] = account
	}

	return nil
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

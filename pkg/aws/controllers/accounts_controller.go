package aws

import (
	"context"
	"fmt"

	awsModels "github.com/dash-ops/dash-ops/pkg/aws/models"
)

// AccountsController orchestrates account-related operations
type AccountsController struct {
	accounts []awsModels.AWSAccount
}

func NewAccountsController(accounts []awsModels.AWSAccount) *AccountsController {
	return &AccountsController{accounts: accounts}
}

func (c *AccountsController) ListAccounts(ctx context.Context) ([]awsModels.AWSAccount, error) {
	return c.accounts, nil
}

func (c *AccountsController) GetAccount(ctx context.Context, accountKey string) (*awsModels.AWSAccount, error) {
	if accountKey == "" {
		return nil, fmt.Errorf("account key is required")
	}

	for _, account := range c.accounts {
		if account.Key == accountKey {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("account %s not found", accountKey)
}

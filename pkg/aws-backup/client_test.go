package aws

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	mockConfig := config{
		Name:            "AWS Dev Account",
		Region:          "us-east-1",
		AccessKeyID:     "1234",
		SecretAccessKey: "4321",
	}

	client, err := NewClient(mockConfig)

	assert.Nil(t, err)
	assert.Equal(t, "aws.client", reflect.TypeOf(client).String())
}

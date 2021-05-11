package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	res, err := gameProvider.Message(context.TODO(), "A001", "hello")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
	assert.NotNil(t, res.Data)
	assert.Equal(t, "A001", res.Data["to"])
	assert.Equal(t, "hello", res.Data["message"])
}

func TestOnline(t *testing.T) {
	res, err := gameProvider.Online(context.TODO(), "A001")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
}

func TestOffline(t *testing.T) {
	res, err := gameProvider.Offline(context.TODO(), "A001")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
}

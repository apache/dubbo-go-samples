package integration

import (
	"context"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	res, err := gameProvider.Login(context.TODO(), "A001")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
	assert.NotNil(t, res.Data)
	assert.Equal(t, "A001", res.Data["to"])
	assert.Equal(t, 0, res.Data["score"])
}

func TestOnline(t *testing.T) {
	res, err := gameProvider.Score(context.TODO(), "A001", "1")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
}

func TestOffline(t *testing.T) {
	res, err := gameProvider.Rank(context.TODO(), "A001")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int32(0), res.Code)
}

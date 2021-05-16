package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
    res, err := gameProvider.Login(context.TODO(), "A001")
    assert.Nil(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, int32(0), res.Code)
    assert.NotNil(t, res.Data)
    assert.Equal(t, "A001", res.Data["to"])
    assert.Equal(t, 0, res.Data["score"])
}

func TestScore(t *testing.T) {
    res, err := gameProvider.Score(context.TODO(), "A001", "1")
    assert.Nil(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, int32(0), res.Code)
}

func TestRank(t *testing.T) {
    res, err := gameProvider.Rank(context.TODO(), "A001")
    assert.Nil(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, int32(0), res.Code)
}
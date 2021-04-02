package integration

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
    res, err := gateProvider.Send(context.TODO(), "A001", "hello")
    assert.Nil(t, err)
    assert.NotNil(t, res)
    assert.Equal(t, int32(0), res.Code)
    assert.NotNil(t, res.Data)
    assert.Equal(t, "A001", res.Data["to"])
    assert.Equal(t, "hello", res.Data["message"])
}
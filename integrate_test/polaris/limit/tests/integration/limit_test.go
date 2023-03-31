package integration

import (
	"context"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPolarisLimit(t *testing.T) {

	var successCount, failCount int64
	for i := 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		_, err := userProvider.GetUser(context.TODO(), &User{Name: "Alex03"})
		if err != nil {
			failCount++
		} else {
			successCount++
		}
	}
	//current limiting effect
	assert.Equal(t, true, failCount > 0)

}

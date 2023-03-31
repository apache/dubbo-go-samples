package integration

import (
	"context"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPolarisRegistry(t *testing.T) {
	user, err := userProvider.GetUser(context.TODO(), &User{Name: "Alex001"})
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.ID, "A001")
	assert.Equal(t, user.Name, "Alex Stocks")
	assert.Equal(t, user.Age, 18)

	user, err = userProviderWithCustomRegistryGroupAndVersion.GetUser(context.TODO(), &User{Name: "Alex001"})
	assert.Nil(t, err)
	assert.Equal(t, user.Name, "Alex Stocks from UserProviderWithCustomGroupAndVersion")
	assert.Equal(t, user.Age, 18)

}

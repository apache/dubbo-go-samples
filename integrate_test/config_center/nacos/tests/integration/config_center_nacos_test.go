package integration

import (
	"context"
	greet "github.com/apache/dubbo-go-samples/config_center/nacos/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSayHello(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}

	ctx := context.Background()

	reply, err := greeterProvider.Greet(ctx, req)

	assert.Nil(t, err)
	assert.Equal(t, "hello world", reply.Greeting)
}

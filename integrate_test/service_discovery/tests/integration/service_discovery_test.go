package integration

import (
	"context"
	"testing"

	greet "github.com/apache/dubbo-go-samples/service_discovery/service/proto"
	"github.com/stretchr/testify/assert"
)

func TestGreet(t *testing.T) {
	req := &greet.GreetRequest{Name: "hello world"}
	ctx := context.Background()
	reply, err := greetService.Greet(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", reply.Greeting)

}

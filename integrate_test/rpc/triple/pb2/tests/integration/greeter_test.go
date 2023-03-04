package integration

import (
	"context"
	"testing"
)

import (
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/triple/pb2/models"
)

func TestStream(t *testing.T) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, tripleConstant.TripleCtxKey("tri-req-id"), "triple-request-id-demo")

	req := models.HelloRequest{
		Name: "dubbo-go",
	}

	r, err := greeterProvider.SayHelloStream(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, r)

	for i := 0; i < 2; i++ {
		err = r.Send(&req)
		assert.Nil(t, err)
	}

	rspUser := &models.User{}
	err = r.RecvMsg(rspUser)
	assert.Nil(t, err)

	err = r.Send(&req)
	assert.Nil(t, err)

	rspUser2 := &models.User{}
	err = r.RecvMsg(rspUser2)
	assert.Nil(t, err)
}

func TestUnary(t *testing.T) {

	req := models.HelloRequest{
		Name: "dubbo-go",
	}
	user, err := greeterProvider.SayHello(context.Background(), &req)
	assert.Nil(t, err)
	assert.NotNil(t, user)

}

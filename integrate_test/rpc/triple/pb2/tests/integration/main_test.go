package integration

import (
	"testing"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/triple/pb2/api"
)

var greeterProvider = new(api.GreeterClientImpl)

func TestMain(m *testing.M) {

	config.SetConsumerService(greeterProvider)

	if err := config.Load(); err != nil {
		panic(err)
	}

}

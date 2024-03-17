package integration

import (
	"os"
	"testing"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/service_discovery/interface/proto"
)

var greetService greet.GreetService

func TestMain(m *testing.M) {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20022"),
	)
	if err != nil {
		panic(err)
	}

	greetService, err = greet.NewGreetService(cli)

	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

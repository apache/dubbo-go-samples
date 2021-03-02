package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(TigerService))
}

type TigerService struct {
}

func (t *TigerService) GetId() (int, error) {
	return 3, nil
}

func (t *TigerService) GetName() (string, error) {
	fmt.Println("I am a Tiger!")
	return "Tiger", nil
}

func (t *TigerService) Yell() (string, error) {
	fmt.Println("Tiger Tiger!")
	return "Tiger Tiger!", nil
}

func (t *TigerService) Reference() string {
	return "TigerService"
}

package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(CatService))
}

type CatService struct {
}

func (c *CatService) GetId() (int, error) {
	return 1, nil
}

func (c *CatService) GetName() (string, error) {
	fmt.Println("I am a Cat!")
	return "Cat", nil
}

func (c *CatService) Yell() (string, error) {
	fmt.Println("Meow Meow!")
	return "Meow Meow!", nil
}

func (c *CatService) Reference() string {
	return "CatService"
}

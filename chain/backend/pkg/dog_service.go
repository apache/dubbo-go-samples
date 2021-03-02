package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(DogService))
}

type DogService struct {
}

func (d *DogService) GetId() (int, error) {
	return 0, nil
}

func (d *DogService) GetName() (string, error) {
	fmt.Println("I am a Dog!")
	return "Dog", nil
}

func (d *DogService) Yell() (string, error) {
	fmt.Println("Woof Woof!")
	return "Woof Woof!", nil
}

func (d *DogService) Reference() string {
	return "DogService"
}

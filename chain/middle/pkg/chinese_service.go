package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

type DogService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (d *DogService) Reference() string {
	return "DogService"
}

type TigerService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (t *TigerService) Reference() string {
	return "TigerService"
}

func init() {
	dog := new(DogService)
	config.SetConsumerService(dog)
	tiger := new(TigerService)
	config.SetConsumerService(tiger)

	config.SetProviderService(&ChineseService{
		dog:   dog,
		tiger: tiger,
	})
}

type ChineseService struct {
	dog   *DogService
	tiger *TigerService
}

func (c *ChineseService) Have() (string, error) {
	name, _ := c.dog.GetName()
	return "I'm Chinese and I have a " + name, nil
}

func (c *ChineseService) Hear() (string, error) {
	name, _ := c.tiger.GetName()
	yell, _ := c.tiger.Yell()
	return fmt.Sprintf("I'm Chinese and I heard a %s yells like %s", name, yell), nil
}

func (c *ChineseService) Reference() string {
	return "ChineseService"
}

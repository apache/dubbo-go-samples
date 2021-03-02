package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

type CatService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (c *CatService) Reference() string {
	return "CatService"
}

type LionService struct {
	GetId   func() (int, error)
	GetName func() (string, error)
	Yell    func() (string, error)
}

func (l *LionService) Reference() string {
	return "LionService"
}

func init() {
	cat := new(CatService)
	config.SetConsumerService(cat)
	lion := new(LionService)
	config.SetConsumerService(lion)

	config.SetProviderService(&AmericanService{
		cat:  cat,
		lion: lion,
	})
}

type AmericanService struct {
	cat  *CatService
	lion *LionService
}

func (a *AmericanService) Have() (string, error) {
	name, _ := a.cat.GetName()
	return "I'm American and I have a " + name, nil
}

func (a *AmericanService) Hear() (string, error) {
	name, _ := a.lion.GetName()
	yell, _ := a.lion.Yell()
	return fmt.Sprintf("I'm American and I heard a %s yells like %s", name, yell), nil
}

func (a *AmericanService) Reference() string {
	return "AmericanService"
}

package pkg

import (
	"fmt"
	"github.com/apache/dubbo-go/config"
)

func init() {
	config.SetProviderService(new(LionService))
}

type LionService struct {
}

func (l *LionService) GetId() (int, error) {
	return 2, nil
}

func (l *LionService) GetName() (string, error) {
	fmt.Println("I am a Lion!")
	return "Lion", nil
}

func (l *LionService) Yell() (string, error) {
	fmt.Println("Lion Lion!")
	return "Lion Lion!", nil
}

func (l *LionService) Reference() string {
	return "LionService"
}

package pkg

type ChineseService struct {
	Have func() (string, error)
	Hear func() (string, error)
}

func (c *ChineseService) Reference() string {
	return "ChineseService"
}

type AmericanService struct {
	Have func() (string, error)
	Hear func() (string, error)
}

func (a *AmericanService) Reference() string {
	return "AmericanService"
}

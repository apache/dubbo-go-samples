package pojo

type Info struct {
	Name string `json:"name"`
	Score int	`json:"score"`
}

func (m Info) JavaClassName() string {
	return "org.apache.dubbo.pojo.Info"
}
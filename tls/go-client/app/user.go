package main

import (
	"context"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
	"time"
)

var userProvider = new(UserProvider)

func init() {
	config.SetConsumerService(userProvider)
	hessian.RegisterPOJO(&User{})
}

type User struct {
	Id		string
	Name	string
	Age		int32
	Time	time.Time
}

type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

func (User) JavaClassName() string {
	return "com.ikurento.user.User"
}

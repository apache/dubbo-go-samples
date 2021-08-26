package main

import (
	"context"
)

import (
	gxlog "github.com/dubbogo/gost/log"
)

type User struct {
	Id   string
	Name string
	Age  int32
}

func (u *User) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.User"
}

type UserProvider struct {
}

func (u *UserProvider) GetUser(ctx context.Context, usr *User) (*User, error) {
	gxlog.CInfo("req:%#v", usr)
	rsp := User{"12345", "" + usr.Name, 18}
	gxlog.CInfo("rsp:%#v", rsp)
	return &rsp, nil
}

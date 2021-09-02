package main

import (
	"context"
)

type User struct {
	Id   string
	Name string
	Age  int32
}

type UserProvider struct {
	GetUser func(ctx context.Context, req *User) (*User, error)
}

func (u *User) JavaClassName() string {
	return "com.apache.dubbo.sample.basic.User"
}

package pkg

import (
	"context"
)

type User struct {
	ID   string
	Name string
	Age  int32
}

type UserProvider struct {
	GetUser func(ctx context.Context, req string) (*User, error)
}

func (u *User) JavaClassName() string {
	return "org.apache.dubbo.User"
}

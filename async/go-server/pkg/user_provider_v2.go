package pkg

import (
	"context"
	"fmt"
	gxlog "github.com/dubbogo/gost/log"
)

type UserProviderV2 struct {
}

func (u *UserProviderV2) getUser(userID string) (*User, error) {
	if user, ok := userMap[userID]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userID)
}

func (u *UserProviderV2) SayHello(ctx context.Context, userID string) error {
	var (
		err  error
		user *User
	)

	user, err = u.getUser(userID)
	if err != nil {
		panic(err)
	}

	gxlog.CInfo("hello, %s", user.Name)
	return nil
}

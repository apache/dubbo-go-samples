package pkg

import (
	"context"
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"

	"github.com/dubbogo/gost/log"

	perrors "github.com/pkg/errors"
)

func init() {
	config.SetProviderService(new(UserProvider2))
}

type UserProvider2 struct {
}

func (u *UserProvider2) getUser(userId string) (*User, error) {
	if user, ok := UserMap[userId]; ok {
		return &user, nil
	}

	return nil, fmt.Errorf("invalid user id:%s", userId)
}

func (u *UserProvider2) GetUser(ctx context.Context, req []interface{}) (*User, error) {
	var (
		err  error
		user *User
	)

	gxlog.CInfo("req:%#v", req)
	user, err = u.getUser(req[0].(string))
	if err == nil {
		gxlog.CInfo("rsp:%#v", user)
	}
	return user, err
}

func (u *UserProvider2) GetUser0(id string, name string, age int) (User, error) {
	var err error

	gxlog.CInfo("id:%s, name:%s, age:%d", id, name, age)
	user, err := u.getUser(id)
	if err != nil {
		return User{}, err
	}
	if user.Name != name {
		return User{}, perrors.New("name is not " + user.Name)
	}
	if user.Age != age {
		return User{}, perrors.New(fmt.Sprintf("age is not %d", user.Age))
	}
	return *user, err
}

func (u *UserProvider2) GetUser3() error {
	return nil
}

func (u *UserProvider2) GetUsers(req []interface{}) ([]User, error) {
	var err error

	gxlog.CInfo("req:%s", req)
	t := req[0].(map[string]interface{})
	user, err := u.getUser(t["ID"].(string))
	if err != nil {
		return nil, err
	}
	gxlog.CInfo("user:%v", user)

	return []User{*user}, err
}

func (u *UserProvider2) GetUser1(req []interface{}) (*User, error) {
	err := perrors.New("test error")
	return nil, err
}

//func (u *UserProvider2) Reference() string {
//	return "UserProvider2"
//}

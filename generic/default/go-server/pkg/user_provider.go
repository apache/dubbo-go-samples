package pkg

import (
	"context"
	"strconv"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
)

type UserProvider struct {
}

func (u *UserProvider) GetUser1(_ context.Context, userID string) (*User, error) {
	logger.Infof("req:%#v", userID)
	rsp := User{userID, "Joe", 48, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser2(_ context.Context, userID string, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userID, name)
	rsp := User{userID, name, 48, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser3(_ context.Context, userCode int) (*User, error) {
	logger.Infof("req:%#v", userCode)
	rsp := User{strconv.Itoa(userCode), "Alex Stocks", 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetUser4(_ context.Context, userCode int, name string) (*User, error) {
	logger.Infof("req:%#v, %#v", userCode, name)
	rsp := User{strconv.Itoa(userCode), name, 18, time.Now()}
	logger.Infof("rsp:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) GetOneUser(_ context.Context) (*User, error) {
	return &User{
		ID:   "1000",
		Name: "xavierniu",
		Age:  24,
		Time: time.Now(),
	}, nil
}

func (u *UserProvider) GetUsers(_ context.Context, userIdList []string) (*UserResponse, error) {
	logger.Infof("req:%#v", userIdList)
	var users []*User
	for _, i := range userIdList {
		users = append(users, userMap[i])
	}
	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) GetUsersMap(_ context.Context, userIdList []string) (map[string]*User, error) {
	logger.Infof("req:%#v", userIdList)
	var users = make(map[string]*User)
	for _, i := range userIdList {
		users[i] = userMap[i]
	}
	return users, nil
}

func (u *UserProvider) QueryUser(_ context.Context, user *User) (*User, error) {
	logger.Infof("req1:%#v", user)
	rsp := User{user.ID, user.Name, user.Age, time.Now()}
	logger.Infof("rsp1:%#v", rsp)
	return &rsp, nil
}

func (u *UserProvider) QueryUsers(_ context.Context, users []*User) (*UserResponse, error) {
	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) QueryAll(_ context.Context) (*UserResponse, error) {
	users := []*User{
		{
			ID:   "001",
			Name: "Joe",
			Age:  18,
			Time: time.Now(),
		},
		{
			ID:   "002",
			Name: "Wen",
			Age:  20,
			Time: time.Now(),
		},
	}

	return &UserResponse{
		Users: users,
	}, nil
}

func (u *UserProvider) MethodMapper() map[string]string {
	return map[string]string{
		"QueryUser":  "queryUser",
		"QueryUsers": "queryUsers",
		"QueryAll":   "queryAll",
	}
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}

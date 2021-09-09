package pkg

type UserResponse struct {
	Users []*User
}

func (u *UserResponse) JavaClassName() string {
	return "org.apache.dubbo.UserResponse"
}

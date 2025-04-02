package model

import (
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *User) hashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return cast.ToString(p), err
}

func (u *User) VerifyPassword() bool {
	hash, err := u.hashPassword(u.Password)
	println(cast.ToString(hash), err)
	return false
}

func (u *User) VerifyLogin() bool {
	return false
}

func (u *User) AddUser() int {
	hash, err := u.hashPassword(u.Password)
	if u.VerifyLogin() {
		return http.StatusConflict
	}
	println(cast.ToString(hash), err)
	return http.StatusOK
}

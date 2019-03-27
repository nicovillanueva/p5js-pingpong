package server

import (
	"fmt"
	"math/rand"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func newUser(u, p string) *user {
	return &user{u, p}
}

func (u *user) save() (int, error) {
	i := rand.Int()
	fmt.Printf("saving user: %+v with id %d\n", u, i)
	return i, nil
}

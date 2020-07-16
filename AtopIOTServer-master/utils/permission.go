package utils

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
)

var once sync.Once
var enforcer *casbin.Enforcer

func GetEnforcer() *casbin.Enforcer {

	once.Do(func() {
		var err error

		enforcer, err = casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
		if err != nil {
			fmt.Println("create enforcer error", err)
			panic(err)
		}
	})
	return enforcer
}

// RequRequestPermission information about coming request
type RequestPermission struct {
	Method     string
	Name       string
	Role       string
	Collection string
}

//CheckPermission  check request permission
func CheckPermission(r *RequestPermission) bool {
	sub := r.Role
	obj := r.Collection
	act := r.Method
	fmt.Println("CheckPermission")
	ok, err := GetEnforcer().Enforce(sub, obj, act)
	fmt.Println("CheckPermission ok err", ok, err)
	if err != nil {
		//handle error
		return false
	}
	return ok
}

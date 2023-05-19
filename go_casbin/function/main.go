package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
	"strings"
)

//官网：https://casbin.org/zh/docs/function

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

// 法一,运用自带的函数
//
//	func main() {
//		e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
//		if err != nil {
//			log.Fatalf("NewEnforecer failed:%v\n", err)
//		}
//
//		check(e, "dajun", "user/dajun/1", "read")
//		check(e, "lizi", "user/lizi/2", "read")
//		check(e, "dajun", "user/lizi/1", "read")
//	}
func KeyMatch(key1, key2 string) bool {
	i := strings.Index(key2, "*")
	if i == -1 {
		return key1 == key2
	}

	if len(key1) > i {
		return key1[:i] == key2[:i]
	}

	return key1 == key2[:i]
}
func KeyMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(KeyMatch(name1, name2)), nil
}

// 法二：自定义函数
func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy_private.csv")
	e.AddFunction("my_func", KeyMatchFunc)
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	check(e, "dajun", "data/1", "read")
	check(e, "dajun", "data/2", "read")
	check(e, "dajun", "data/1", "write")
	check(e, "dajun", "mydata", "read")
}

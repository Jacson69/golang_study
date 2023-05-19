package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"log"
)

func check(e *casbin.Enforcer, sub, domain, obj, act string) {
	ok, _ := e.Enforce(sub, domain, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s in %s\n", sub, act, obj, domain)
	} else {
		fmt.Printf("%s CANNOT %s %s in %s\n", sub, act, obj, domain)
	}
}

func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		log.Fatalf("NewEnforecer failed:%v\n", err)
	}

	//check(e, "dajun", "data", "read")
	//check(e, "lizi", "data", "write")
	//check(e, "dajun", "data", "write")
	//check(e, "lizi", "data", "read")

	//check(e, "dajun", "prod.data", "read")
	//check(e, "dajun", "prod.data", "write")
	//check(e, "lizi", "dev.data", "read")
	//check(e, "lizi", "dev.data", "write")
	//check(e, "lizi", "prod.data", "write")

	//check(e, "dajun", "data", "read")
	//check(e, "dajun", "data", "write")
	//check(e, "lizi", "data", "read")
	//check(e, "lizi", "data", "write")

	check(e, "dajun", "tenant1", "data1", "read")
	check(e, "dajun", "tenant2", "data2", "read")
}

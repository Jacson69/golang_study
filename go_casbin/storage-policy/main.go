package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
)

func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}

func main() {
	// 两种写法
	//a, _ := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
	a, _ := gormadapter.NewAdapter("mysql", "root:123456@tcp(192.168.78.173:3306)/casbin", true)
	e, _ := casbin.NewEnforcer("./model.conf", a)

	check(e, "dajun", "data1", "read")
	check(e, "lizi", "data2", "write")
	check(e, "dajun", "data1", "write")
	check(e, "dajun", "data2", "read")
}

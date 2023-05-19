package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

// 字符串模型
func check(e *casbin.Enforcer, sub, obj, act string) {
	ok, _ := e.Enforce(sub, obj, act)
	if ok {
		fmt.Printf("%s CAN %s %s\n", sub, act, obj)
	} else {
		fmt.Printf("%s CANNOT %s %s\n", sub, act, obj)
	}
}
func main() {
	//写法一
	//m := model.NewModel()
	//m.AddDef("r", "r", "sub, obj, act")
	//m.AddDef("p", "p", "sub, obj, act")
	//m.AddDef("e", "e", "some(where (p.eft == allow))")
	//m.AddDef("m", "m", "r.sub == p.sub && r.obj == p.obj && r.act == p.act")
	//
	//a := fileadapter.NewAdapter("./policy.csv")
	//e, err := casbin.NewEnforcer(m, a)
	//if err != nil {
	//	log.Fatalf("NewEnforecer failed:%v\n", err)
	//}
	//
	//check(e, "dajun", "data1", "read")
	//check(e, "lizi", "data2", "write")
	//check(e, "dajun", "data1", "write")
	//check(e, "dajun", "data2", "read")

	//写法二
	text := `
  [request_definition]
  r = sub, obj, act
  
  [policy_definition]
  p = sub, obj, act
  
  [policy_effect]
  e = some(where (p.eft == allow))
  
  [matchers]
  m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
  `

	m, _ := model.NewModelFromString(text)
	a := fileadapter.NewAdapter("./policy.csv")
	e, _ := casbin.NewEnforcer(m, a)

	check(e, "dajun", "data1", "read")
	check(e, "lizi", "data2", "write")
	check(e, "dajun", "data1", "write")
	check(e, "dajun", "data2", "read")
}

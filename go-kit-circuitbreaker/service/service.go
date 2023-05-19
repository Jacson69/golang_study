package service

import "errors"

type Service interface {
	// Add calculate a+b
	Add(a, b int) int

	// Subtract calculate a-b
	Substract(a, b int) int

	// Multiply calculate a*b
	Multiply(a, b int) int

	// Divide calculate a/b
	Divide(a, b int) (int, error)

	// 服务注册新增的
	// HealthCheck check service health status
	HealthCheck() bool

	//登录
	Login(name, pwd string) (string, error)
}

// ArithmeticService implement Service interface
type ArithmeticService struct {
}

func (s ArithmeticService) Add(a, b int) int {
	return a + b
}

func (s ArithmeticService) Substract(a, b int) int {
	return a - b
}

func (s ArithmeticService) Multiply(a, b int) int {
	return a * b
}

func (s ArithmeticService) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("the dividend can not be zero!")
	}
	return a / b, nil
}

// ArithmeticService实现HealthCheck
// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true。
func (s ArithmeticService) HealthCheck() bool {
	return true
}

func (s ArithmeticService) Login(name, pwd string) (string, error) {
	if name == "name" && pwd == "pwd" {
		token, err := Sign(name, pwd)
		return token, err
	}

	return "", errors.New("Your name or password dismatch")
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service

package service

import (
	"context"
)

type Service interface {
	Context() context.Context // 获取上下文
	Start()                   // 开启服务: 服务创建之后运行
	Stop()                    // 停止服务：停止服务逻辑
	Close()                   // 停止后销毁时调度
}

// BaseService 服务基类
type BaseService struct {
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 关闭方法

	name string // 服务名
}

func NewService(parent context.Context, name string) *BaseService {
	server := new(BaseService)
	server.ctx, server.cancel = context.WithCancel(parent)
	server.name = name
	return server
}

// Service interface

func (s *BaseService) Context() context.Context {
	return s.ctx
}

func (s *BaseService) Start() {
	panic(" Service Sub Class need to realize Service interface func Start() error")
}

func (s *BaseService) Stop() {
	s.cancel()
}

func (s *BaseService) Close() {
}

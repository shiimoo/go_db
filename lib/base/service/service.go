package service

import (
	"context"
)

type Service interface {
	Context() context.Context // 获取上下文
	Start() error             // 开启服务: 服务创建之后运行
	Stop() error              // 停止服务：停止服务逻辑
	Close() error             // 销毁服务: 关闭服务之后执行
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

func (s *BaseService) Start() error {
	panic(" Service Sub Class need to realize Service interface func Start() error")
}

func (s *BaseService) Stop() error {
	panic(" Service Sub Class need to realize Service interface func Stop() error")
}

func (s *BaseService) Close() error {
	panic(" Service Sub Class need to realize Service interface func Close() error")
}

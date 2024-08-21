package wss

import (
	"github.com/shiimoo/godb/lib/base/service"
)

// WssServer webSocket服务
type WssServer struct {
	*service.BaseService

	// listener   *net.TCPListener // 监听器
	// key        string           // 标识符

	// linkMu sync.RWMutex      // 管理锁
	// links  map[uint]*TcpLink // 链接管理
}

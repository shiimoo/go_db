package tcp

import (
	"context"
	"net"
	"sync"

	"github.com/shiimoo/godb/lib/base/errors"
	"github.com/shiimoo/godb/lib/base/service"
	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/shiimoo/godb/lib/mlog"
)

// TcpServer tcp服务
type TcpServer struct {
	*service.BaseService

	listenAddr *net.TCPAddr     // 监听地址
	listener   *net.TCPListener // 监听器

	linkMu sync.RWMutex      // 管理锁
	links  map[uint]*TcpLink // 链接管理
}

func NewServer(parent context.Context, address string) (*TcpServer, error) {
	// @param address "0.0.0.0:8080"

	var server *TcpServer = nil
	listenAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return server, errors.NewErr(ErrTcpListerAddrError, address, err)
	}

	tcpListener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return server, errors.NewErr(ErrTcpListerError, err)
	}
	// create
	server = new(TcpServer)
	server.BaseService = service.NewService(parent, "TcpServer") // todo 名称规范待定
	server.listenAddr = listenAddr
	server.listener = tcpListener
	return server, nil
}

func (s *TcpServer) createLinkObj(baseLink *net.TCPConn) {
	linkObj := NewLink(s.Context(), baseLink, uint(snowflake.Gen()))
	s.addLinkObj(linkObj)
	linkObj.Start()
}

func (s *TcpServer) addLinkObj(linkObj *TcpLink) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	s.links[linkObj.Key()] = linkObj
}

// Service interface

func (s *TcpServer) Start() error {
	for {
		select {
		case <-s.Context().Done():
		default:
			// 监听链接
			tcpConn, err := s.listener.AcceptTCP()
			if err != nil {
				mlog.Warn("tcp", "acceptTCP", err.Error())
			} else {
				s.createLinkObj(tcpConn)
			}
		}
	}
}

// func (s *TcpServer) Stop() error {
// 	panic(" Service Sub Class need to realize Service interface func Stop() error")
// }

// func (s *TcpServer) Close() error {
// 	panic(" Service Sub Class need to realize Service interface func Close() error")
// }

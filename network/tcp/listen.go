package tcp

import (
	"context"
	"fmt"
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

	listener *net.TCPListener // 监听器
	key      string           // 标识符

	linkMu sync.RWMutex      // 管理锁
	links  map[uint]*TcpLink // 链接管理
}

func NewServer(parent context.Context, key, address string) (*TcpServer, error) {
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
	server.BaseService = service.NewService(parent, fmt.Sprintf("TcpServer_%s", key))
	server.listener = tcpListener
	server.key = key
	server.links = make(map[uint]*TcpLink)
	return server, nil
}

func (s *TcpServer) createLinkObj(baseLink *net.TCPConn) {
	linkObj := NewLink(s.Context(), s, baseLink, uint(snowflake.Gen()))
	s.addLinkObj(linkObj)
	linkObj.Start()
}

func (s *TcpServer) addLinkObj(linkObj *TcpLink) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	s.links[linkObj.Key()] = linkObj
}

func (s *TcpServer) _removeLinkObj(linkObj *TcpLink, err error) {
	if _, found := s.links[linkObj.Key()]; found {
		delete(s.links, linkObj.Key())
		fmt.Println("关闭链接 ", linkObj.Key(), err)
		linkObj.Stop()
	}
}
func (s *TcpServer) removeLinkObj(linkObj *TcpLink, err error) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	s._removeLinkObj(linkObj, err)
}

func (s *TcpServer) removeLinkObjById(key uint, err error) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	if linkObj, found := s.links[key]; found {
		delete(s.links, key)
		fmt.Println("关闭链接 ", key, err)
		linkObj.Stop()
	}
}

// Listen interface

func (s *TcpServer) Dispatch(bs []byte) {
	fmt.Println("todo 接受到的数据处理", len(bs), bs)
}

func (s *TcpServer) LinkCount() int {
	return len(s.links)
}

func (s *TcpServer) RemoveLink(id uint) {
	s.removeLinkObjById(id, nil)
}

func (s *TcpServer) SendData(id uint, data []byte) {
	if len(data) == 0 {
		return // 数据为空
	}
	s.linkMu.Lock()
	linkObj, found := s.links[id]
	s.linkMu.Unlock()
	if !found {
		return // 链接不存在
	}
	if err := linkObj.Write(data); err != nil {
		s.removeLinkObj(linkObj, err)
	}
}

// Service interface

func (s *TcpServer) Start() {
	go func() {
		for {
			select {
			case <-s.Context().Done():
				s.Close()
				return
			default:
				s.listener.Accept()
				// 监听链接
				tcpConn, err := s.listener.AcceptTCP()
				if err != nil {
					mlog.Warn("tcp", "acceptTCP", err.Error())
				} else {
					s.createLinkObj(tcpConn)
				}
			}
		}
	}()
}

func (s *TcpServer) Close() {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	for _, linkObj := range s.links {
		s._removeLinkObj(linkObj, nil)
	}
}

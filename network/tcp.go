package network

import (
	"context"
	"fmt"
	"net"

	"github.com/shiimoo/godb/lib/base/snowflake"
	"github.com/shiimoo/godb/lib/mlog"
)

type TcpLink struct {
	*baseLink

	// tcpServer *tcp.TcpServer // 归属的tcp服务 todo listen
}

func NewLink(parent context.Context, baseLink *net.TCPConn, id uint) *TcpLink {
	link := new(TcpLink)
	link.baseLink = newBaseLink(parent, baseLink)
	// link.tcpServer = tcpServer
	return link
}

func (l *TcpLink) Key() uint {
	return l.ID()
}

// TcpListenServer tcp服务
type TcpListenServer struct {
	*baseListenServer
}

func NewServer(parent context.Context, key, address string) (*TcpListenServer, error) {
	// @param address "0.0.0.0:8080"

	serverObj := new(TcpListenServer)
	base, err := newBaseListenServer(parent, address)
	if err != nil {
		return nil, err
	}
	// CREATE
	serverObj.baseListenServer = base
	return serverObj, nil
}
func (b *TcpListenServer) NetType() string {
	return "tcp"
}

func (s *TcpListenServer) createLinkObj(baseLink *net.TCPConn) {
	linkObj := network.NewLink(s.Context(), baseLink, uint(snowflake.Gen()))
	s.addLinkObj(linkObj)
	linkObj.Start()
}

func (s *TcpListenServer) addLinkObj(linkObj *network.TcpLink) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	s.links[linkObj.Key()] = linkObj
}

func (s *TcpListenServer) _removeLinkObj(linkObj *network.TcpLink, err error) {
	if _, found := s.links[linkObj.Key()]; found {
		delete(s.links, linkObj.Key())
		fmt.Println("关闭链接 ", linkObj.Key(), err)
		linkObj.Stop()
	}
}
func (s *TcpListenServer) removeLinkObj(linkObj *network.TcpLink, err error) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	s._removeLinkObj(linkObj, err)
}

func (s *TcpListenServer) removeLinkObjById(key uint, err error) {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	if linkObj, found := s.links[key]; found {
		delete(s.links, key)
		fmt.Println("关闭链接 ", key, err)
		linkObj.Stop()
	}
}

// Listen interface

func (s *TcpListenServer) Dispatch(bs []byte) {
	fmt.Println("todo 接受到的数据处理", len(bs), bs)
}

func (s *TcpListenServer) LinkCount() int {
	return len(s.links)
}

func (s *TcpListenServer) RemoveLink(id uint) {
	s.removeLinkObjById(id, nil)
}

func (s *TcpListenServer) SendData(id uint, data []byte) {
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

func (s *TcpListenServer) Start() {
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

func (s *TcpListenServer) Close() {
	s.linkMu.Lock()
	defer s.linkMu.Unlock()
	for _, linkObj := range s.links {
		s._removeLinkObj(linkObj, nil)
	}
}

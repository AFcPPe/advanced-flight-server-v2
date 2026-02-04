package handler

import (
	"github.com/panjf2000/gnet/v2"
)

// Handler 处理器接口
type Handler interface {
	Handle(conn gnet.Conn) error
}

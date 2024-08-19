package network

type Link interface {
	Key() any              // 标识格式待定
	Start() error          // 开启监听
	Close() error          // 关闭链接
	Read() ([]byte, error) // 读取数据
	Write([]byte) error    // 写入数据
}

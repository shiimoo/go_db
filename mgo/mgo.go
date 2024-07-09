package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Indexs mongo 索引结构重写(复刻结构primitive.E) {{字段名:1/-1}, {关键字:值}}
type Indexs bson.D

// 命令结果
type opResult struct {
	err     error // 错误信息
	results []any // 结果参数 todo 改成唯一参数，去掉[]
}

func newOpResult() *opResult {
	return &opResult{
		err:     nil,
		results: make([]any, 0),
	}
}

// 尾部追加
func (r *opResult) addResult(args ...any) {
	r.results = append(r.results, args...)
}

// 操作指令
type op struct {
	cmd          string         // 操作指令集
	args         []any          // 指令所需参数
	resultAccept chan *opResult // 结果接受channel
}

// 参数列表
func newOp(cmd string) *op {
	o := new(op)
	o.cmd = cmd
	o.args = make([]any, 0)
	o.resultAccept = make(chan *opResult)
	return o
}

// 追加参数
func (o *op) append(args ...any) {
	o.args = append(o.args, args...)
}

// 参数出栈
func (o *op) getResult() *opResult {
	return <-o.resultAccept
}

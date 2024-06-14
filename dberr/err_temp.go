package dberr

// ErrTemp 错误模板 : 暂定理解为错误类型
type TempErr string

func (t TempErr) Error() string {
	return string(t)
}

var (
	Text = TempErr("TempErr Text :%v")
)
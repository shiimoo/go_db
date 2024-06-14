package dberr

import "fmt"

type Err struct {
	temp TempErr // 模板
	vals []any   // 错误参数
}

func (e Err) Error() string {
	if e.vals == nil {
		return e.Format()
	}
	return fmt.Sprintf(e.Format(), e.vals...)
}

func (e Err) Format() string {
	return e.temp.Error()
}

// IsSameErr 是否体同一错误;
// 对比的错误需要满足DbErr, 任意一错误不满足，均会判定为不同(false);
// 该方法比较的是错误模板， 而非错误本身.
func IsSameErr(errA, errB error) bool {
	eA, ok := errA.(Err)
	if !ok {
		return false
	}
	eB, ok := errB.(Err)
	if !ok {
		return false
	}
	return eA.Format() == eB.Format()
}

// NewErr 创建error
func NewErr(temp TempErr, vals ...any) error {
	if vals == nil { // 没有参数时直接用模板
		return temp
	}
	return Err{temp, vals}
}

// 安全沙盒控制
package savectrl

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/shiimoo/godb/dberr"
)

func SaveBox(handler func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			errInfos := []string{fmt.Sprintf("%v", r)}
			stackList := strings.Split(string(debug.Stack()), "\n")[7:]
			err = dberr.NewErr(ErrSaveBoxPanic, strings.Join(append(errInfos, stackList...), "\n"))
		}
	}()
	return handler()
}

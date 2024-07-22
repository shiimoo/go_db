package mlog

// ***** todo 日志格式参数 *****
// Println 默认输出
func Println(msg string) {
	GetMgr().Println(msg)
}

func Debug(mod, label, msg string) {
	GetMgr().Debug(mod, label, msg)
}

func Debugf(mod, label, format string, datas ...Data) {
	GetMgr().Debugf(mod, label, format, datas...)
}

func Info(mod, label, msg string) {
	GetMgr().Info(mod, label, msg)
}

func Infof(mod, label, format string, datas ...Data) {
	GetMgr().Infof(mod, label, format, datas...)
}

func Warn(mod, label, msg string) {
	GetMgr().Warn(mod, label, msg)
}

func Warnf(mod, label, format string, datas ...Data) {
	GetMgr().Warnf(mod, label, format, datas...)
}

func Error(mod, label, msg string) {
	GetMgr().Error(mod, label, msg)
}

func Errorf(mod, label, format string, datas ...Data) {
	GetMgr().Errorf(mod, label, format, datas...)
}

func Fatal(mod, label, msg string) {
	GetMgr().Fatal(mod, label, msg)
}

func Fatalf(mod, label, format string, datas ...Data) {
	GetMgr().Fatalf(mod, label, format, datas...)
}

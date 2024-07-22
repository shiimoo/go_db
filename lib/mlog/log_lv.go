package mlog

import "fmt"

// ***** 日志级别 *****

type logLv string

func (l logLv) String() string {
	return fmt.Sprintf("%-5s", string(l))
}

const (
	LvDebug = logLv("DEBUG") // **调试**信息: 主要用于开发阶段，记录详细的调试信息。
	LvInfo  = logLv("INFO")  // **一般**信息: 记录程序运行过程中的重要里程碑或状态变化。
	LvWarn  = logLv("WARN")  // **警告**信息: 表明可能出现潜在错误的情况，但不影响程序当前运行。
	LvError = logLv("ERROR") // **错误**信息: 指出程序运行中发生的错误，但程序仍能继续运行。
	LvFatal = logLv("FATAL") // **致命**错误: 导致程序无法继续运行或必须立即停止的错误。
)

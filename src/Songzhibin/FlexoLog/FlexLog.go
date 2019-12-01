package FLexlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// 多功能日志库(输出至文件版)

//定义等级制
//	DEBUG -> 1
//	INFO  -> 2
//	WARN  -> 3
//	ERROR -> 4
//	CRITICAL -> 5

// 声明等级类
type Lever uint

const (
	DEBUG Lever = iota + 1
	INFO
	WARN
	ERROR
	CRITICAL
)

// 创建FilesLog结构体
type FilesLog struct {
	// Lever:显示等级 1-5 超过5默认为1
	Lever              Lever
	FilePath, FileName string
	Handle             *os.File
	MaxSize            int
}

// 创建构造FilesLog函数
func NewFileLog(Lever Lever, FilePath, FileName string, MaxSize int) (NewFileLogObj *FilesLog) {
	// 判断传入等级来指定日志的写入等级
	switch {
	case Lever >= DEBUG && Lever <= CRITICAL:
	default:
		Lever = DEBUG
	}
	// 创建构造FilesLog 返回指针
	NewFileLogObj = &FilesLog{
		Lever:    Lever,
		FilePath: FilePath,
		FileName: FileName,
		MaxSize:  MaxSize,
	}
	NewFileLogObj.Handle = NewFile(FilePath, FileName)
	return
}

// 创建构造FileHand函数
func NewFile(FilePath, FileName string) (NewFileObj *os.File) {
	OpenPath := filepath.Join(FilePath, FileName)
	NewFileObj, err := os.OpenFile(OpenPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(fmt.Sprintf("打开文件/创建文件出现异常:%s", err))
	}
	return
}

// 创建*FilesLog结构体对应方法
// Debug
func (f *FilesLog) Debug(format string, a ...interface{}) {
	if f.Lever <= 1 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		// 调用封装偏函数
		f.Encapsulation(format, "DEBUG", a...)
	}
	return
}

// Info
func (f *FilesLog) Info(format string, a ...interface{}) {
	if f.Lever <= 2 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		f.Encapsulation(format, "INFO", a...)
	}
	return
}

// Warn
func (f *FilesLog) Warn(format string, a ...interface{}) {
	if f.Lever <= 3 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		f.Encapsulation(format, "WARN", a...)
	}
	return
}

// Error
func (f *FilesLog) Error(format string, a ...interface{}) {
	if f.Lever <= 4 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		f.Encapsulation(format, "ERROR", a...)
	}
	return
}

// Critical
func (f *FilesLog) Critical(format string, a ...interface{}) {
	if f.Lever <= 5 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		f.Encapsulation(format, "CRITICAL", a...)
	}
	return
}

// Close
func (f *FilesLog) HandleClose() {
	f.Handle.Close()
}

// 封装函数
func (f *FilesLog) Encapsulation(format string, Lever string, a ...interface{}) {
	TimeSting := fmt.Sprintf("[%s]", time.Now().Format("2006-01-02 15:04:05.000"))
	// 将引用函数行数加入日志输出
	pc, file, line, ok := runtime.Caller(2)
	pcName := runtime.FuncForPC(pc).Name() // 获取调用函数Name
	filepathname := filepath.Dir(file)
	if !ok {
		panic("找不到对应调度文件")
	}
	// 创建对应引用文件 函数 以及行数字符串
	MisTakeString := fmt.Sprintf("[%s->%s lins:%d]", filepathname, pcName, line)
	// 等级字符串
	LeverString := fmt.Sprintf("[%s]", Lever)
	// 多format融合
	format = TimeSting + " " + LeverString + " " + MisTakeString + "  " + format
	// 写入文件
	fmt.Fprintln(f.Handle, format)
}

// 优化机制
// 判断文件大小 自动进行备份
func (f *FilesLog) JudgeFileSizeBackFile(size int) {
	// 拼接当前备份文件路径
	// 进行检测
	fileInfo, err := f.Handle.Stat()
	if err != nil {
		panic(fmt.Sprintf("检测文件大小失败%s", err))
		return
	}
	if int(fileInfo.Size()) > size {
		NowPostfix := fmt.Sprintf(time.Now().Format("06-01-02 15:04:05"))
		SrcPath := filepath.Join(f.FilePath, f.FileName)
		DstPath := filepath.Join(f.FilePath, f.FileName+NowPostfix+".log")
		f.Handle.Close()                           // 关闭文件
		os.Rename(SrcPath, DstPath)                // 重命名
		f.Handle = NewFile(f.FilePath, f.FileName) // 重新赋值句柄
	}
}

package FLexlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// 多功能日志库(输出至文件版 *并发版)

// FLexLog以声明内容
// Lever 等级类
// 常量 DEBUG INFO WARN ERROR CRITICAL

// 创建AsyncFileLog结构体
type AsyncFileLog struct {
	Lever              Lever
	FilePath, FileName string
	MaxSize            int
	Handle             *os.File
	Buffer             chan *Log
}

// 创建Logging结构体
type Log struct {
	Format string
}

// 创建AsyncFIleLog函数
func NewAsyncFileLog(Lever Lever, FilePath, FileName string, MaxSize, ChanSize int) (NewFileLogObj *AsyncFileLog) {
	// 判断传入等级来指定日志的写入等级
	switch {
	case Lever >= DEBUG && Lever <= CRITICAL:
	default:
		Lever = DEBUG
	}
	// 创建构造AsyncFilesLog 返回指针
	NewFileLogObj = &AsyncFileLog{
		Lever:    Lever,
		FilePath: FilePath,
		FileName: FileName,
		MaxSize:  MaxSize,
	}
	// 调用构造函数创造句柄
	NewFileLogObj.Handle = NewAsyncFile(FilePath, FileName)
	// 使用make创造管道
	NewFileLogObj.Buffer = make(chan *Log, ChanSize)
	return
}

// 创建构造AsyncFileHand函数
func NewAsyncFile(FilePath, FileName string) (NewFileObj *os.File) {
	OpenPath := filepath.Join(FilePath, FileName)
	NewFileObj, err := os.OpenFile(OpenPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(fmt.Sprintf("打开文件/创建文件出现异常:%s", err))
	}
	return
}
// 创建构造Log函数
func NewLogObj(Format string) (LogObj *Log) {
	LogObj = &Log{
		Format: Format,
	}
	return
}

// 创造AsyncFileLog方法 运行至后台进行异步存储日志
// 需要使用 go在主函数启动
func (A *AsyncFileLog) AsyncRun() {
	for {
		for logObj := range A.Buffer {
			// 调用写日志函数
			fmt.Fprintln(A.Handle, logObj.Format)
		}
	}
}

// 创建*AsyncFilesLog结构体的对应方法

// DEBUG
func (A *AsyncFileLog) Debug(format string, a ...interface{}) {
	if A.Lever <= 1 {
		// 调用优化机制
		A.JudgeFileSizeBackFile(A.MaxSize)
		A.Encapsulation(format, "DEBUG", a...)
	}
	return
}

// Info
func (A *AsyncFileLog) Info(format string, a ...interface{}) {
	if A.Lever <= 1 {
		// 调用优化机制
		A.JudgeFileSizeBackFile(A.MaxSize)
		A.Encapsulation(format, "INFO", a...)
	}
	return
}

// Warn
func (A *AsyncFileLog) Warn(format string, a ...interface{}) {
	if A.Lever <= 1 {
		// 调用优化机制
		A.JudgeFileSizeBackFile(A.MaxSize)
		A.Encapsulation(format, "WARN", a...)
	}
	return
}

// Error
func (A *AsyncFileLog) Error(format string, a ...interface{}) {
	if A.Lever <= 1 {
		// 调用优化机制
		A.JudgeFileSizeBackFile(A.MaxSize)
		A.Encapsulation(format, "ERROR", a...)
	}
	return
}

// Critical
func (A *AsyncFileLog) Critical(format string, a ...interface{}) {
	if A.Lever <= 5 {
		// 调用优化机制
		A.JudgeFileSizeBackFile(A.MaxSize)
		A.Encapsulation(format, "Critical", a...)
	}
	return
}

// CLose
// 需要主函数defer注册结束句柄
func (A *AsyncFileLog) HandleClose() {
	A.Handle.Close()
}

// 封装函数
func (A *AsyncFileLog) Encapsulation(format string, Lever string, a ...interface{}) {
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
	// 实例化Log对象
	LogObj := NewLogObj(format)
	// 多路复用 如果可以存放至Buffer中就直接存入 如果不能存入则从Buffer中取出一条后再进行存储
	select {
	case A.Buffer <- LogObj:
	default:
		<-A.Buffer
		A.Buffer <- LogObj
	}
}



// 优化机制
// 判断文件大小 自动进行备份
func (A *AsyncFileLog) JudgeFileSizeBackFile(size int) {
	// 拼接当前备份文件路径
	// 进行检测
	fileInfo, err := A.Handle.Stat()
	if err != nil {
		panic(fmt.Sprintf("检测文件大小失败%s", err))
		return
	}
	if int(fileInfo.Size()) > size {
		NowPostfix := fmt.Sprintf(time.Now().Format("06-01-02 15:04:05"))
		SrcPath := filepath.Join(A.FilePath, A.FileName)
		DstPath := filepath.Join(A.FilePath, A.FileName+NowPostfix+".log")
		A.Handle.Close()                                // 关闭文件
		os.Rename(SrcPath, DstPath)                     // 重命名
		A.Handle = NewAsyncFile(A.FilePath, A.FileName) // 重新赋值句柄
	}
}

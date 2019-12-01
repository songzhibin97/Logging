# 日志库


#### FlexoLog.go 用于输出至文件log内容

1.定义报错等级

```go
// 声明等级类
type Lever uint

const (
	DEBUG Lever = iota + 1
	INFO
	WARN
	ERROR
	CRITICAL
)
```
2.设定结构体

```go
type FilesLog struct {
	// Lever:显示等级 1-5 超过5默认为1
	Lever              Lever
	FilePath, FileName string
	Handle             *os.File
}
```

3.创造构造函数

```go
func NewFileLog(lever Lever, FilePath, FileName string) (NewFileLogObj *FilesLog) {
	// 判断传入等级来指定日志的写入等级
	switch {
	case lever >= DEBUG && lever <= CRITICAL:
		lever = lever
	default:
		lever = DEBUG
	}
	// 创建构造FilesLog 返回指针
	NewFileLogObj = &FilesLog{
		Lever:    lever,
		FilePath: FilePath,
		FileName: FileName}
	NewFileLogObj.Handle = NewFile(FilePath, FileName)
	return
}
```

4.创建构造FileHand函数

```go
func NewFile(FilePath, FileName string) (NewFileObj *os.File) {
	OpenPath := filepath.Join(FilePath, FileName)
	NewFileObj, err := os.OpenFile(OpenPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		panic(fmt.Sprintf("打开文件/创建文件出现异常:%s", err))
	}
	return
}
```

5.创建*FilesLog结构体对应方法

``` go
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
		// 调用封装偏函数
		f.Encapsulation(format, "INFO", a...)
	}
	return
}

// Warn
func (f *FilesLog) Warn(format string, a ...interface{}) {
	if f.Lever <= 3 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		// 调用封装偏函数
		f.Encapsulation(format, "WARN", a...)
	}
	return
}

// Error
func (f *FilesLog) Error(format string, a ...interface{}) {
	if f.Lever <= 4 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		// 调用封装偏函数
		f.Encapsulation(format, "ERROR", a...)
	}
	return
}

// Critical
func (f *FilesLog) Critical(format string, a ...interface{}) {
	if f.Lever <= 5 {
		// 调用优化机制
		f.JudgeFileSizeBackFile(f.MaxSize)
		// 调用封装偏函数
		f.Encapsulation(format, "CRITICAL", a...)
	}
	return
}

// Close
func (f *FilesLog) HandleClose(){
	f.Handle.Close()
}

```
6.封装*FilesLog中对应方法

```go
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
```

7.创建*FilesLog结构体对应方法(优化机制,)

```go
func (f *FilesLog) JudgeFileSizeBackFile(size int) {
	// 拼接当前备份文件路径
	CopyPath := filepath.Join(f.FilePath, f.FileName)
	// 进行检测
	fileInfo, err := os.Stat(CopyPath)
	if err != nil {
		panic(fmt.Sprintf("检测文件大小失败%s", err))
		return
	}
	if int(fileInfo.Size()) > size {
		NowPostfix := fmt.Sprintf(time.Now().Format("06-01-02 15:04:05"))
		SrcPath := filepath.Join(f.FilePath, f.FileName)
		DstPath := filepath.Join(f.FilePath, f.FileName+NowPostfix+".log")
		Srccontent, err := ioutil.ReadFile(SrcPath)
		if err != nil {
			panic(fmt.Sprintf("读取原文件内容失败", err))
			return
		}
		errs := ioutil.WriteFile(DstPath, Srccontent, 0755)
		if errs != nil {
			panic("备份文件失败")
		}
		src, err := os.OpenFile(SrcPath, os.O_TRUNC, 0755)
		defer src.Close()
		if err != nil {
			panic(fmt.Sprintf("清空文件失败%s", err))
		}
	}
}
```

**调用方法**

```go
func main() {
	obj := FLexlog.NewFileLog(FLexlog.WARN, "./", "exe", 1024) // 实例化对象 分别传参 lever(显示等级) filepath(日志存放路径) filename(日志name) maxsize(最大存放大小)
	defer obj.HandleClose()
	obj.Debug(fmt.Sprintf("这是一条%s测试数据", "debug"))
	obj.Info(fmt.Sprintf("这是一条%s测试数据", "info"))
	obj.Warn(fmt.Sprintf("这是一条%s测试数据", "warn"))
	obj.Error(fmt.Sprintf("这是一条%s测试数据", "error"))
	obj.Critical(fmt.Sprintf("这是一条%s测试数据", "critical"))
}
```


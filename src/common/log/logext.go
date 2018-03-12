package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func init() {
	initLogger()
	go RunLogger()
}

// init logger
func initLogger() bool {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("log open path %s error\n", filepath.Dir(os.Args[0]))
		return false
	}
	path = filepath.Join(path, "logs")
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	now := time.Now()

	//file, _ := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	Std = NewLogger(Ldebug, Ldefault, now.Day(), path)

	return true
}

func RunLogger() {
	for {
		Std.WriteFile()
		time.Sleep(1 * time.Second)
	}
}

func OnKill() {
	Std.WriteFile()
}

func (l *Logger) WriteFile() {
	l.mu.Lock()
	defer l.mu.Unlock()
	bufLength := l.buf.Len()
	if bufLength <= 0 {
		return
	}
	name := l.GetFileName()
	file, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	defer file.Close()
	if err != nil {
		fmt.Printf("os openlogfile failed %v \r\n", err)
	} else {
		file.Write(l.buf.Bytes())
	}
	l.buf.Reset()
}

func (l *Logger) GetFileName() string {
	now := time.Now()
	filename := filepath.Join(l.path, GetDateStr(now)) + ".log"
	return filename
}

func SetLvl(lvlStr string) {
	var lvl int
	switch lvlStr {
	case "debug":
		lvl = Ldebug
	case "info":
		lvl = Linfo
	case "warn":
		lvl = Lwarn
	case "error":
		lvl = Lerror
	case "panic":
		lvl = Lpanic
	case "fatal":
		lvl = Lfatal
	}
	Std.lvl = lvl
}

func GetDateStr(date time.Time) string {
	return date.Format("2006-01-02")
}

func NewLogger(lvl, flag int, day int, path string) *Logger {
	return &Logger{lvl: lvl, flag: flag, day: day, path: path}
}

func Debug(msg string) {
	Print(Ldebug, Lskip, msg)
}

func Debugf(format string, v ...interface{}) {
	Print(Ldebug, Lskip, fmt.Sprintf(format, v...))
}

func Info(msg string) {
	Print(Linfo, Lskip, msg)
}

func Infof(format string, v ...interface{}) {
	Print(Linfo, Lskip, fmt.Sprintf(format, v...))
}

func Warn(msg string) {
	Print(Lwarn, Lskip, msg)
}

func Warnf(format string, v ...interface{}) {
	Print(Lwarn, Lskip, fmt.Sprintf(format, v...))
}

func Error(msg string) {
	Print(Lerror, Lskip, msg)
}

func Errorf(format string, v ...interface{}) {
	Print(Lerror, Lskip, fmt.Sprintf(format, v...))
}

func Panic(msg string) {
	Print(Lpanic, Lskip, msg)
	panic(msg)
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Print(Lpanic, Lskip, fmt.Sprintf(format, v...))
	panic(s)
}

func Print(lvl int, calldepth int, msg string) {
	//initLogger()
	Std.Output(lvl, calldepth, msg)
	// out := Std.out
	// file, ok := out.(*os.File)
	// if ok {
	// 	file.Close()
	// } else {
	// 	fmt.Println("not okokokkokookk")
	// }
}

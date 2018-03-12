package log

import (
	"bytes"
	"runtime"
	"sync"
	"time"
)

const (
	Ldate = 1 << iota
	Ltime
	Lmicroseconds
	Lfile
	Lmodule
	Llevel
	LstdFlags = Ldate | Ltime
	Ldefault  = LstdFlags | Lmicroseconds | Lfile | Llevel | Lmodule
)
const (
	Ldebug = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
)

// stack caller skip
const Lskip int = 3

var Std *Logger

var levels = []string{
	"debug",
	"info",
	"warn",
	"error",
	"panic",
	"fatal",
}

type Logger struct {
	buf  bytes.Buffer
	lvl  int
	mu   sync.Mutex
	flag int
	day  int
	path string
}

func itoa(buf *bytes.Buffer, v int, width int) {
	var u uint = uint(v)
	var b [10]byte
	l := len(b)
	for ; width > 0 || u > 0; u /= 10 {
		width--
		l--
		b[l] = byte(u%10) + '0'
	}

	for l < len(b) {
		buf.WriteByte(b[l])
		l++
	}
}

func (l *Logger) WriteHeader(lvl int, t time.Time, file string, line int) {
	// level
	if l.flag&Llevel != 0 {
		l.buf.WriteByte('[')
		l.buf.WriteString(levels[lvl])
		l.buf.WriteByte(']')
	}
	// date
	if l.flag&Ldate != 0 {
		year, month, day := t.Date()
		itoa(&l.buf, year, 4)
		itoa(&l.buf, int(month), 2)
		itoa(&l.buf, day, 2)
	}
	// time
	if l.flag&Ltime != 0 {
		l.buf.WriteByte(' ')
		hour, minute, second := t.Clock()
		itoa(&l.buf, hour, 2)
		l.buf.WriteByte(':')
		itoa(&l.buf, minute, 2)
		l.buf.WriteByte(':')
		itoa(&l.buf, second, 2)
	}
	// file
	if l.flag&Lfile != 0 {
		l.buf.WriteByte(' ')
		l.buf.WriteString(file)
		l.buf.WriteByte(':')
		itoa(&l.buf, line, -1)
		l.buf.WriteString(": ")
	}
}

func (l *Logger) Output(lvl int, calldepth int, msg string) error {
	if lvl < l.lvl {
		return nil
	}
	now := time.Now()

	var file string
	var line int
	if l.flag&Lfile != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	l.WriteHeader(lvl, now, file, line)
	l.buf.WriteString(msg)
	if len(msg) > 0 && msg[len(msg)-1] != '\n' {
		l.buf.WriteByte('\n')
	}

	return nil
}

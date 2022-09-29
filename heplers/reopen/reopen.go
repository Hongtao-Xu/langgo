package reopen

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Reopener interface {
	Reopen() error
}

type Writer interface {
	Reopener
	io.Writer
}

type WriteCloser interface {
	Reopener
	io.WriteCloser
}

type FileWriter struct {
	mu   sync.Mutex
	f    *os.File
	mode os.FileMode
	name string
}

// Fd 返回文件描述符
func (fw *FileWriter) Fd() uintptr {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	err := fw.f.Fd()
	return err
}

// Close 关闭文件
func (fw *FileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	err := fw.f.Close()
	return err
}

// reopen 文件，不加锁版本
func (fw *FileWriter) reopen() error {
	if fw.f != nil {
		fw.f.Close()
		fw.f = nil
	}
	newf, err := os.OpenFile(fw.name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, fw.mode)
	if err != nil {
		fw.f = nil
		return err
	}
	fw.f = newf
	return nil
}

// Reopen the file
func (fw *FileWriter) Reopen() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	err := fw.reopen()
	return err
}

// Write 写入文件
func (fw *FileWriter) Write(p []byte) (int, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	n, err := fw.f.Write(p)
	return n, err
}

// NewFileWriter 构造方法，创建FileWriter，默认权限：rw
func NewFileWriter(name string) (*FileWriter, error) {
	//0666:代表User、Group、及Other的权限分别是6，6，6，即均为rw权限。
	return NewFileWriterMode(name, 0666)
}

// NewFileWriterMode 构造方法，创建FileWriter
func NewFileWriterMode(name string, mode os.FileMode) (*FileWriter, error) {
	writer := FileWriter{
		f:    nil,
		name: name,
		mode: mode,
	}
	err := writer.reopen()
	if err != nil {
		return nil, err
	}
	return &writer, nil
}

type BufferedFileWriter struct {
	mu         sync.Mutex
	quitChan   chan bool
	done       bool
	origWriter *FileWriter
	bufWriter  *bufio.Writer
}

// Reopen bw
func (bw *BufferedFileWriter) Reopen() error {
	bw.mu.Lock()
	bw.bufWriter.Flush() //将缓冲区数据写入io.Writer
	// use non-mutex version since we are using this one
	err := bw.origWriter.reopen()
	//重设io.Writer
	bw.bufWriter.Reset(io.Writer(bw.origWriter))
	bw.mu.Unlock()
	return err
}

// Close bw
func (bw *BufferedFileWriter) Close() error {
	bw.quitChan <- true
	bw.mu.Lock()
	bw.done = true
	bw.bufWriter.Flush()
	bw.origWriter.f.Close()
	bw.mu.Unlock()
	return nil
}

// Write 通过bw,将p写入io.Writer
func (bw *BufferedFileWriter) Write(p []byte) (int, error) {
	bw.mu.Lock()
	n, err := bw.bufWriter.Write(p)
	//当前缓冲区大小<p
	if bw.bufWriter.Buffered() < len(p) {
		bw.bufWriter.Flush()
	}
	bw.mu.Unlock()
	return n, err
}

// Flush bw
func (bw *BufferedFileWriter) Flush() {
	bw.mu.Lock()
	bw.bufWriter.Flush()
	bw.origWriter.f.Sync() //从内存中刷入磁盘
	bw.mu.Unlock()
}

// flushDaemon 定时flush守护进程
func (bw *BufferedFileWriter) flushDaemon(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-bw.quitChan: //终止信号，关闭定时flush任务
			ticker.Stop()
			return
		case <-ticker.C: //每隔interval时间，flush一次
			bw.Flush()
		}
	}
}

const bufferSize = 256 * 1024
const flushInterval = 30 * time.Second

//NewBufferedFileWriter bw默认构造方法
func NewBufferedFileWriter(w *FileWriter) *BufferedFileWriter {
	return NewBufferedFileWriterSize(w, bufferSize, flushInterval)
}

// NewBufferedFileWriterSize bw构造方法
//w:FileWriter；size:缓冲区大小；flush:定时刷盘时间间隔
func NewBufferedFileWriterSize(w *FileWriter, size int, flush time.Duration) *BufferedFileWriter {
	bw := BufferedFileWriter{
		quitChan:   make(chan bool, 1),
		origWriter: w,
		bufWriter:  bufio.NewWriterSize(w, size),
	}
	go bw.flushDaemon(flush)
	return &bw
}

type multiReopenWriter struct {
	writers []Writer
}

func (t *multiReopenWriter) Reopen() error {
	for _, w := range t.writers {
		err := w.Reopen()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *multiReopenWriter) Write(p []byte) (int, error) {
	for _, w := range t.writers {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

func MultiWriter(writers ...Writer) Writer {
	w := make([]Writer, len(writers))
	copy(w, writers) //多个参数，复制到切片
	return &multiReopenWriter{w}
}

type nopReopenWriteCloser struct {
	io.Writer
}

func (nopReopenWriteCloser) Reopen() error {
	return nil
}

func (nopReopenWriteCloser) Close() error {
	return nil
}

func NopWriter(w io.Writer) WriteCloser {
	return nopReopenWriteCloser{w}
}

var (
	Stdout  = NopWriter(os.Stdout)
	Stderr  = NopWriter(os.Stderr)
	Discard = NopWriter(ioutil.Discard)
)

package log

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Writer struct {
	sync.Mutex
	_logger *log.Logger

	file   string
	suffix string
	handle *os.File
}

func (w *Writer) GetLogger() *log.Logger {
	return w._logger
}

func (w *Writer) SetFile(file string) {
	w.file = file
}

func (w *Writer) GetFile() string {
	return w.file
}

func (w *Writer) GetHandle() *os.File {
	return w.handle
}

func (w *Writer) Output() {

}

func (w *Writer) Rotate(rotateType RotateType) error {
	w.Lock()
	defer w.Unlock()

	var suffix string
	if rotateType == RotateHourly {
		suffix = time.Now().Format("2006010215")
	} else if rotateType == RotateDaily {
		suffix = time.Now().Format("20060102")
	}

	fmt.Println("w.suffix:", w.suffix, " w:", suffix)

	if w.suffix != suffix {
		fmt.Println("do rotate")
		if err := w.doRotate(); err != nil {
			return err
		}

		w.suffix = suffix
	}

	return nil
}

func (w *Writer) doRotate() error {
	newFileName := w.file + "." + w.suffix
	if err := os.Rename(w.file, newFileName); err != nil {
		return err
	}

	if err := w.setOutput(); err != nil {
		return err
	}

	return w.handle.Close()
}

func (w *Writer) setOutput() error {
	f, err := os.OpenFile(w.file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	w._logger = log.New(f, w._logger.Prefix(), w._logger.Flags())
	w.handle = f

	return nil
}

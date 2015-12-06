package logger

import (
	"strconv"
	"time"
	"os"
	"io/ioutil"
	"log"
)

// TODO: Make this parameterizable
const capacity = 32768

type Worker struct {
	logFileRoot    string
	buffer         []byte
	bufferPosition int
}

func NewWorker(id int) (w *Worker) {
	return &Worker{
		logFileRoot: strconv.Itoa(id) + "_",
		buffer: make([]byte, capacity),
	}
}

func (w *Worker) Work(channel chan []byte) {
	for {
		event := <- channel
		length := len(event)
		if length > capacity {
			log.Println("message received was too large")
			continue
		}
		if (length + w.bufferPosition) > capacity {
			w.Save()
		}
		copy(w.buffer[w.bufferPosition:], event)
		w.bufferPosition += length
	}
}

func (w *Worker) Save() {
	if w.bufferPosition == 0 { return }
	f, _ := ioutil.TempFile("", "logs_")
	f.Write(w.buffer[0:w.bufferPosition])
	f.Close()
	os.Rename(f.Name(), w.logFileRoot + strconv.FormatInt(time.Now().UnixNano(), 10))
	w.bufferPosition = 0
}

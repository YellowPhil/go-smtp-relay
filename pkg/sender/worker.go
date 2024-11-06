package session

import (
	"time"
)

type SendTask struct {
	To       string
	From     string
	Contents []byte
}

type Worker struct {
	*Client
	retries  int
	sendChan chan *SendTask
}

func NewWorker(_sendChan chan *SendTask) *Worker {
	return &Worker{
		sendChan: _sendChan,
	}
}

func (w *Worker) Submit(rcpt *SendTask) {
	if rcpt != nil {
		w.sendChan <- rcpt
	}
}

func (w *Worker) SendWithRetries() {
	for task := range w.sendChan {
		for range w.retries {
			if err := w.SendMail(task.To, task.From, task.Contents); err != nil {
				time.Sleep(time.Second * 3)
			}
		}
	}
}

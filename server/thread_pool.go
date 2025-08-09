package server

import (
	"log"
)

// element in the queue
type Job struct {
	client *Client
}

// represent the thread in the pool
type Worker struct {
	id      int
	jobChan chan Job
}

// represent the thread pool
type Pool struct {
	jobQueue chan Job
	workers  []*Worker
}

// create a new worker
func NewWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			log.Printf("Worker %d is handling job from %s", w.id, job.client.conn.RemoteAddr())
			job.client.handleConnection()
		}
	}()
}

func NewPool(numOfWorker int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers:  make([]*Worker, numOfWorker),
	}
}

// push job to queue
func (p *Pool) AddJob(client *Client) {
	p.jobQueue <- Job{client: client}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}
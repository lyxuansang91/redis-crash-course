package threadpool

import (
	"log"
)

// Job represents a unit of work to be executed by a worker
type Job struct {
    task func()
}

// Worker represents a worker in the pool
type Worker struct {
    id      int
    jobChan chan Job
}

// Pool represents a pool of workers processing jobs
type Pool struct {
    jobQueue chan Job
    workers  []*Worker
}

// NewWorker creates a new worker bound to the provided job channel
func NewWorker(id int, jobChan chan Job) *Worker {
    return &Worker{
        id:      id,
        jobChan: jobChan,
    }
}

// Start begins the worker loop to process incoming jobs
func (w *Worker) Start() {
    go func() {
        for job := range w.jobChan {
            log.Printf("Worker %d is executing a job ", w.id)
            job.task()
        }
    }()
}

// NewPool creates a pool with the specified number of workers
func NewPool(numOfWorkers int) *Pool {
    return &Pool{
        jobQueue: make(chan Job),
        workers:  make([]*Worker, numOfWorkers),
    }
}

// AddJob enqueues a function to be executed by the pool
func (p *Pool) AddJob(task func()) {
    p.jobQueue <- Job{task: task}
}

// Start initializes and starts all workers
func (p *Pool) Start() {
    for i := 0; i < len(p.workers); i++ {
        worker := NewWorker(i, p.jobQueue)
        p.workers[i] = worker
        worker.Start()
    }
}



package work

import (
	"fmt"
	"sync"
)

type Job interface {
	Do() error
}

// remember to handle panics in the Do method
//func (job *SampleJob) Do() {
//	defer func() {
//		if r := recover(); r != null {
//			fmt.Println("Recovered in Do", r)
//		}
//	}()
//	// existing code
//}

type Worklist struct {
	jobs      chan Job
	closed    bool
	closelock sync.Mutex
}

// New initializes a new Worklist with a specified capacity for the jobs channel
func New(capacity int) *Worklist {
	return &Worklist{
		jobs: make(chan Job, capacity),
	}
}

func NewWorklist() *Worklist {
	return &Worklist{jobs: make(chan Job)}
}

func (w *Worklist) Add(job Job) {
	w.jobs <- job
}

// Next dequeues a job from the worklist, returning false if no jobs are left
func (w *Worklist) Next() (Job, bool) {
	job, ok := <-w.jobs
	return job, ok
}

// Close safely closes the jobs channel, ensuring no more jobs are added
func (w *Worklist) Close() {
	w.closelock.Lock()
	defer w.closelock.Unlock()
	if !w.closed {
		close(w.jobs)
		w.closed = true
	}
}

func (w *Worklist) Jobs() <-chan Job {
	return w.jobs
}

func (w *Worklist) Len() int {
	return len(w.jobs)
}

func (w *Worklist) IsEmpty() bool {
	return w.Len() == 0
}

func (w *Worklist) IsClosed() bool {
	select {
	case _, ok := <-w.jobs:
		return !ok
	default:
		return false
	}
}

func (w *Worklist) Wait() {
	for !w.IsClosed() {
	}
}

func (w *Worklist) DoAll() {
	for job, ok := w.Next(); ok; job, ok = w.Next() {
		job.Do()
	}
}

func (w *Worklist) DoAllConcurrently(wg *sync.WaitGroup) {
	for job, ok := w.Next(); ok; job, ok = w.Next() {
		wg.Add(1)
		go func(job Job) {
			defer wg.Done()
			job.Do()
		}(job)
	}
	wg.Wait() // wait for all goroutines to finish
}

func ExecuteAll(jobs []Job) {
	w := New(len(jobs))
	for _, job := range jobs {
		w.Add(job)
	}
	w.DoAll()
}

func ExecuteAllConcurrently(jobs []Job) {
	w := New(len(jobs))

	var wg sync.WaitGroup
	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(job Job) {
			defer wg.Done() // Ensure that wg.Done() is called after job.Do() finishes
			w.Add(job)
		}(job)
	}

	go func() {
		defer w.Close()
		w.DoAllConcurrently(&wg)
	}()

	wg.Wait()
}

func ExecuteAllConcurrentlyWithProgress(jobs []Job) {
	w := New(len(jobs))
	var wg sync.WaitGroup
	var progressMutex sync.Mutex
	progress := 0

	wg.Add(len(jobs))

	for _, job := range jobs {
		go func(job Job) {
			defer wg.Done() // ensure wg.Done() is called after the job
			err := job.Do()
			if err != nil {
				return
			} // Execute the job
			progressMutex.Lock()
			progress++
			fmt.Printf("\rProgress: %d%%", (progress*100)/len(jobs))
			progressMutex.Unlock()
		}(job)
	}

	go func() {
		defer w.Close()
		wg.Wait()
		fmt.Println("\nAll jobs completed.")
	}()

	wg.Wait()
}

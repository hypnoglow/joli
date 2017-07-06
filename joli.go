// Package joli is a simple job processor library.
//
// It just provides an easy way to define an N-sized queue and
// a K-sized worker pool to concurrently run queued jobs.
//
// As a simple example:
//
//  // Prepare a function to handle errors in jobs.
//  errorHandler := func(err error) {
//      log.Fatalln("Job failed: %s", err.Error())
//  }
//
//  Prepare a job channel as a queue.
//  queue := make(chan joli.Job, QueueMaxSize)
//
//  p := joli.NewProcessor(queue, errorHandler, WorkerPoolSize)
//  p.Start()
//
//  // Create some jobs; then send it directly to the queue.
//  // This code will not block until the queue is not full.
//  for i := 0; i < numberOfJobs; i++ {
//      queue <- job
//  }
//
package joli

// Job describes the job to be run.
type Job interface {
	// Run runs the job.
	Run() error
}

// JobErrorHandler describes a func that is used to handle job error.
type JobErrorHandler func(error)

type empty struct{}

// Processor represents the job processor - it runs jobs.
type Processor struct {
	jobQueue        <-chan Job
	jobErrorHandler JobErrorHandler
	limiter         chan empty
	stop            chan empty
}

// NewProcessor returns a new Processor.
func NewProcessor(queue <-chan Job, errorHandler JobErrorHandler, maxLimiterSize int64) *Processor {
	return &Processor{
		jobQueue:        queue,
		jobErrorHandler: errorHandler,
		limiter:         make(chan empty, maxLimiterSize),
		stop:            make(chan empty),
	}
}

// NumWorkers returns a number of workers currently
// running their jobs.
func (w Processor) NumWorkers() int64 {
	return int64(len(w.limiter))
}

// Start signals the processor to start listening for job requests.
func (w *Processor) Start() {
	go func() {
		for {
			select {
			case job, ok := <-w.jobQueue:
				if !ok {
					return
				}

				w.limiter <- empty{}

				// Spawn a worker goroutine.
				go func() {
					if err := job.Run(); err != nil {
						w.jobErrorHandler(err)
					}
					<-w.limiter
				}()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop signals the processor to stop listening for job requests.
func (w *Processor) Stop() {
	close(w.stop)
}

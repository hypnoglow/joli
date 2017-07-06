package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hypnoglow/joli"
)

const (
	defaultJobChanSize       = 1024
	defaultProcessorPoolSize = 4
	numberOfJobs             = 4096

	processorMonitoringPeriod = time.Millisecond * 500
	senderPeriod              = time.Millisecond * 200
)

func main() {
	rand.Seed(time.Now().Unix())

	// Prepare a function to handle errors in jobs.
	errorHandler := func(err error) {
		log.Fatalln("Job failed: %s", err.Error())
	}

	// Prepare a job channel as a queue.
	jobChan := make(chan joli.Job, defaultJobChanSize)

	p := joli.NewProcessor(jobChan, errorHandler, defaultProcessorPoolSize)

	// Monitoring number of jobs waiting for an available worker.
	go func() {
		for {
			log.Println("[PRCSSR] Number of workers:", p.NumWorkers())
			log.Println("[QUEUE ] Size of job queue:", len(jobChan))
			time.Sleep(processorMonitoringPeriod)
		}
	}()

	sendJobs(jobChan)

	// Just helper code to not to exit early.
	stop := make(chan os.Signal, 8)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func sendJobs(jobChan chan<- joli.Job) {
	// Now you can add any Job to a queue
	// (see printerJob below as an example implementation).
	for i := 0; i < numberOfJobs; i++ {
		log.Printf("[SENDER ] Sending job %d ...\n", i)
		time.Sleep(senderPeriod)

		jobChan <- printerJob{
			number:   i,
			duration: time.Millisecond * 100 * time.Duration(rand.Intn(30)),
		}

		log.Println("[SENDER ] Job sent!")
	}
}

type printerJob struct {
	number   int
	duration time.Duration
}

func (j printerJob) Run() error {
	log.Printf("[JOB #%2d] Doing job for %s...\n", j.number, j.duration.String())
	time.Sleep(j.duration)
	log.Printf("[JOB #%2d] Done job!\n", j.number)
	return nil
}

package executor

import (
	"io"
	"os/exec"
)

// Job is the state of a job
type Job struct {
	Name string
}

// jobRequest is used when requesting job informations, as we
// needs a channel for responses
type jobRequest struct {
	name string
	resp chan []Job
}

type Executor struct {
	cmdStdout io.Writer
	cmdStderr io.Writer
	getJob    chan jobRequest
	addJob    chan Job
	deleteJob chan string
	jobs      map[string]Job
}

// Default create an Executor with default configuration
func Default() Executor {
	exec := Executor{
		cmdStdout: io.Discard,
		cmdStderr: io.Discard,
		getJob:    make(chan jobRequest),
		addJob:    make(chan Job),
		deleteJob: make(chan string),
		jobs:      make(map[string]Job),
	}
	exec.run()
	return exec
}

// WithLogger create an Executor with custom Logger, enabling
// command output logging
func WithLogger(stdout, stderr io.Writer) Executor {
	exec := Executor{
		cmdStdout: stdout,
		cmdStderr: stderr,
		getJob:    make(chan jobRequest),
		addJob:    make(chan Job),
		deleteJob: make(chan string),
		jobs:      make(map[string]Job),
	}
	exec.run()
	return exec
}

// Run maintains and dispatch current running job's state
func (e *Executor) run() {
	go func() {
		for {
			select {
			case req := <-e.getJob:
				resp := filterJobsByName(e.jobs, req.name)
				req.resp <- resp

			case job := <-e.addJob:
				e.jobs[job.Name] = job

			case name := <-e.deleteJob:
				delete(e.jobs, name)
			}
		}
	}()
}

// Execute add a command in the execution queue
func (e *Executor) Execute(service, command, workdir string, environ []string) (err error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}
	cmd.Env = append(cmd.Environ(), environ...)
	cmd.Stdout = e.cmdStdout
	cmd.Stderr = e.cmdStderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	e.addJob <- Job{
		Name: service,
	}

	// wait concurrently for the end of the command
	go func() {
		err = cmd.Wait()
		// remove job from state
		e.deleteJob <- service
	}()

	return nil
}

func (e *Executor) GetJob(service string) []Job {
	resp := make(chan []Job)
	e.getJob <- jobRequest{
		name: service,
		resp: resp,
	}
	jobs := <-resp

	return jobs
}

func filterJobsByName(jobs map[string]Job, name string) []Job {
	var resp []Job

	for _, v := range jobs {
		// if name is empty, acts like a wildcard
		if len(name) == 0 || name == v.Name {
			resp = append(resp, Job{
				Name: v.Name,
			})
		}
	}

	return resp
}

package supervisor

import (
	"time"

	"github.com/docker/containerd/runtime"
)

type StartTask struct {
	baseTask
	platformStartTask
	ID            string
	BundlePath    string
	Stdout        string
	Stderr        string
	Stdin         string
	StartResponse chan StartResponse
	Labels        []string
}

func (s *Supervisor) start(t *StartTask) error {
	start := time.Now()
	container, err := runtime.New(s.stateDir, t.ID, t.BundlePath, s.runtime, t.Labels)
	if err != nil {
		return err
	}
	s.containers[t.ID] = &containerInfo{
		container: container,
	}
	ContainersCounter.Inc(1)
	task := &startTask{
		Err:           t.ErrorCh(),
		Container:     container,
		StartResponse: t.StartResponse,
		Stdin:         t.Stdin,
		Stdout:        t.Stdout,
		Stderr:        t.Stderr,
	}
	task.setTaskCheckpoint(t)

	s.startTasks <- task
	ContainerCreateTimer.UpdateSince(start)
	return errDeferedResponse
}

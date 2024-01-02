package downloader

import (
	"context"
	"errors"
)

const (
	StageCreated Stage = "Created"
)

type Task struct {
	Name          string
	TaskStage     Stage
	Percentage    float64
	DownloadRate  string
	DownloadBytes int64
	TotalBytes    int64
	ActivePeers   int
	CancelFunc    context.CancelFunc

	Config DownloadConfig
}

type TaskPipeline struct {
	tasks map[string]*Task
}

func NewTaskPipeline() *TaskPipeline {
	return &TaskPipeline{
		tasks: make(map[string]*Task),
	}
}

func (p *TaskPipeline) SubmitTask(name string, magnet string, dir string) error {
	if _, ok := p.tasks[name]; ok {
		return errors.New("task name already exists")
	}
	ctx, cancelFunc := context.WithCancel(context.Background())

	t := &Task{
		Name:       name,
		TaskStage:  StageCreated,
		CancelFunc: cancelFunc,
		Config: DownloadConfig{
			Dir:      dir,
			FileName: "",
			Magnet:   magnet,
			Proxy:    "",
			Ctx:      ctx,
		},
	}
	t.Config.Obs = func(stage Stage, readBytes int64, totalBytes int64, activePeers, totalPeers int, downloadRate string, percentage float64) {
		t.TaskStage = stage
		switch stage {
		case StagePrepare:
		case StageDownload:
			// update task
			t.TaskStage = StageDownload
			t.Percentage = percentage
			t.DownloadRate = downloadRate
			t.DownloadBytes = readBytes
			t.TotalBytes = totalBytes
			t.ActivePeers = activePeers

		case StageSuccess:
		case StageCancel:
		case StageError: // never
		}
	}
	p.tasks[name] = t

	// start
	go func() {
		if err := DownloadBlocked(&t.Config); err != nil {
			t.TaskStage = StageError
			cancelFunc()
		}
	}()

	return nil
}

func (p *TaskPipeline) CancelTask(name string) error {
	t := p.GetTask(name)
	if t == nil {
		return errors.New("task not found")
	}
	if t.TaskStage == StageDownload {
		t.CancelFunc()
		return nil
	} else {
		return errors.New("task not in download stage")
	}
}

// Exists return true if task name exists
func (p *TaskPipeline) Exists(name string) bool {
	_, ok := p.tasks[name]
	return ok
}

// GetTask return nil if not found
func (p *TaskPipeline) GetTask(name string) *Task {
	if t, ok := p.tasks[name]; ok {
		return t
	}
	return nil
}

package workerpool

type Tasker interface {
	Process(task any, qid int) []error
	Handle(err []error)
}

type Task struct {
	
}

func (t *Task) Process(task any, qid int) []error {
	panic("implement me")
}

func (t *Task) Handle(err []error) {
	panic("implement me")
}


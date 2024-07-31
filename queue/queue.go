package queue

// empty path indicates no more jobs
type Entry struct {
	Path string
}

type Queue struct {
	jobs chan Entry
}

func (queue *Queue) Add(job Entry) {
	queue.jobs <- job
}

func (queue *Queue) Next() Entry {
	job := <-queue.jobs

	return job
}

func NewQueue(bufferSize int) Queue {
	return Queue{make(chan Entry, bufferSize)}
}

func NewJob(path string) Entry {
	return Entry{path}
}

// add empty message to queue for each worker. Workers terminate after receiving the message, thereafter the programme can continue
func (queue *Queue) Finalise(numberOfWorkers int) {
	for i := 0; i < numberOfWorkers; i++ {
		queue.Add(Entry{""})
	}
}
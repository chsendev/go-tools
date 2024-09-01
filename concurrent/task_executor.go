package concurrent

import (
	"context"
	"fmt"
	"sync"
)

type Result[T any, R any] struct {
	l1      sync.Mutex
	l2      sync.Mutex
	Success []R
	Failed  []T
	Errors  []*Errors[T]
}

type Errors[T any] struct {
	Task  T
	Error error
}

func (r *Result[T, R]) success(v R) {
	r.l1.Lock()
	r.l1.Unlock()
	r.Success = append(r.Success, v)
}

func (r *Result[T, R]) failed(v T, err error) {
	r.l2.Lock()
	r.l2.Unlock()
	r.Failed = append(r.Failed, v)
	r.Errors = append(r.Errors, &Errors[T]{
		Task:  v,
		Error: err,
	})
}

type ConsumerFunc[T any, R any] func(context.Context, T) (R, error)

// TaskExecutor 1生产者 : n消费者
// T 任务 R 结果
type TaskExecutor[T any, R any] struct {
	ctx       context.Context
	cancelCtx context.CancelFunc
	// 任务数
	tasks []T
	// 生产者与消费者通信的channel
	taskChan chan T
	// 消费逻辑
	consumerFunc ConsumerFunc[T, R]
	// 结果
	results     *Result[T, R]
	concurrency int
	wg          *sync.WaitGroup
	errorChan   chan error
}

func NewTaskExecutor[T any, R any](ctx context.Context, tasks []T, consumerFunc ConsumerFunc[T, R]) *TaskExecutor[T, R] {
	concurrency := HalfConcurrency.Concurrency(len(tasks))
	taskChan := make(chan T, concurrency*2)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	res := Result[T, R]{
		Success: make([]R, 0, len(tasks)),
		Failed:  make([]T, 0),
		Errors:  make([]*Errors[T], 0),
	}

	return &TaskExecutor[T, R]{
		ctx:          ctx,
		tasks:        tasks,
		taskChan:     taskChan,
		consumerFunc: consumerFunc,
		concurrency:  concurrency,
		wg:           &wg,
		results:      &res,
		errorChan:    make(chan error),
	}
}
func (p *TaskExecutor[T, R]) WithConcurrency(strategy ConcurrencyStrategy) *TaskExecutor[T, R] {
	concurrency := strategy.Concurrency(len(p.tasks))
	taskChan := make(chan T, concurrency*2)
	var wg sync.WaitGroup
	wg.Add(concurrency)
	p.concurrency = concurrency
	p.wg = &wg
	p.taskChan = taskChan
	return p
}

func (p *TaskExecutor[T, R]) producer() {
	go func() {
		defer p.recover()
		for _, v := range p.tasks {
			p.taskChan <- v
		}
		close(p.taskChan)
	}()
}

func (p *TaskExecutor[T, R]) consumer() {
	for i := 0; i < p.concurrency; i++ {
		go func() {
			defer p.recover()
			defer p.wg.Done()
			for {
				task, ok := <-p.taskChan
				if ok {
					res, err := p.consumerFunc(p.ctx, task)
					if err != nil {
						p.results.failed(task, err)
					} else {
						p.results.success(res)
					}
				} else {
					break
				}
			}
		}()
	}
}

func (p *TaskExecutor[T, R]) Run() *TaskExecutor[T, R] {
	p.producer()
	p.consumer()
	return p
}

func (p *TaskExecutor[T, R]) Wait() *Result[T, R] {
	//go func() {
	//	err := <-p.errorChan
	//	if err != nil {
	//
	//	}
	//}()
	p.wg.Wait()
	return p.results
}

func (p *TaskExecutor[T, R]) recover() {
	if err := recover(); err != nil {
		p.errorChan <- fmt.Errorf("TaskExecutor panic :%v", err)
	}
}

func getConcurrency(taskNum int) int {
	if taskNum < 10 {
		return 5
	} else if taskNum < 20 {
		return 10
	} else if taskNum < 30 {
		return 15
	} else {
		return 20
	}
}

package concurrent

const maxConcurrency = 50 // 最大并发数
var (
	// FixedConcurrency 固定并发策略
	FixedConcurrency = new(fixedConcurrencyStrategy)
	// HalfConcurrency 1/2 并发策略
	HalfConcurrency = new(halfConcurrencyStrategy)
	// QuarterConcurrency 1/4 并发策略
	QuarterConcurrency = new(quarterConcurrencyStrategy)
)

type ConcurrencyStrategy interface {
	Concurrency(taskCount int) int
}

// fixedConcurrencyStrategy 固定并发策略
type fixedConcurrencyStrategy struct {
	concurrency int
}

func NewFixedConcurrency(concurrency int) *fixedConcurrencyStrategy {
	return &fixedConcurrencyStrategy{concurrency}
}

func (s *fixedConcurrencyStrategy) Concurrency(length int) int {
	return s.concurrency
}

// halfConcurrencyStrategy 1/2 并发策略
type halfConcurrencyStrategy struct {
}

func (s *halfConcurrencyStrategy) Concurrency(length int) int {
	return correction(length / 2)
}

// quarterConcurrencyStrategy 1/4 并发策略
type quarterConcurrencyStrategy struct {
}

func (s *quarterConcurrencyStrategy) Concurrency(length int) int {
	return correction(length / 4)
}

func correction(concurrency int) int {
	if concurrency > maxConcurrency {
		return maxConcurrency
	}
	if concurrency <= 0 {
		return 1
	}
	return concurrency
}

package sequencegenerator

import "sync/atomic"

type IsDone bool

const (
	Done    IsDone = true
	NotDone IsDone = false
)

type SequenceGenerator interface {
	Next() (int64, IsDone)
}

type sequenceGenerator struct {
	start int64
	end   int64
	count atomic.Int64
}

func NewSequenceGenerator(start, end int64) *sequenceGenerator {
	rv := &sequenceGenerator{start: start, end: end}
	rv.count.Store(int64(rv.start))
	return rv
}

func (g *sequenceGenerator) Next() (int64, IsDone) {
	i := g.count.Load()
	if i >= g.end {
		return 0, Done
	}
	g.count.Add(1)
	return i, NotDone
}

var _ SequenceGenerator = (*sequenceGenerator)(nil)

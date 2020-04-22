package limit

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"
)

type MemLimit struct {
	ctx       context.Context
	maxMemory uint64
	onLimit   func(string)
}

func WithResourcesLimit(parent context.Context, maxTime time.Duration, maxMem uint64, block func(context.Context)) {
	ctx, cancel := context.WithTimeout(parent, maxTime)
	defer cancel()

	memLimit := &MemLimit{
		ctx:       ctx,
		maxMemory: maxMem,
		onLimit:   func(msg string) { log.Fatal(msg) },
	}

	memLimit.Execute(block)
}

func (m *MemLimit) Execute(block func(context.Context)) error {
	exceeded := make(chan uint64)
	done := make(chan bool)

	go m.sampling(exceeded)
	go m.run(done, block)

	select {
	case <-done:
		return nil
	case <-m.ctx.Done():
		m.onLimit("context done")
	case bytes := <-exceeded:
		m.onLimit(fmt.Sprintf("memory limit exceeded, allocated %d bytes", bytes))
	}

	return errors.New("limits reached")
}

func (m *MemLimit) sampling(exceeded chan<- uint64) {
	var s runtime.MemStats

	for {
		if m.ctx.Err() != nil {
			break
		}

		runtime.ReadMemStats(&s)

		if s.HeapAlloc > m.maxMemory {
			exceeded <- s.HeapAlloc
		}

		time.Sleep(time.Millisecond)
	}
}

func (m *MemLimit) run(done chan<- bool, block func(context.Context)) {
	if m.ctx.Err() != nil {
		return
	}

	block(m.ctx)

	close(done)
}

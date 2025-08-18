package util

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

// TaskGroup represents a group of tasks that can be executed concurrently and waited upon.
type TaskGroup struct {
	name string
	wg   *sync.WaitGroup
}

// NewTaskGroup creates and returns a new TaskGroup with the specified name.
func NewTaskGroup(name string) *TaskGroup {
	return &TaskGroup{name: name, wg: &sync.WaitGroup{}}
}

// GoAndWait starts a new goroutine to execute the given function and waits for its completion or context cancellation.
func (g *TaskGroup) GoAndWait(ctx context.Context, fn func(context.Context) error) error {
	ch := g.Go(ctx, fn)
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return fmt.Errorf("%s wait: %s", g.name, ctx.Err())
	}
}

// Go starts a new goroutine to execute the given function fn with the provided context,
// returning a channel that will receive an error.
func (g *TaskGroup) Go(ctx context.Context, fn func(context.Context) error) <-chan error {
	ech := make(chan error)
	g.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// 修复：只截取 debug.Stack() 的前 4KB，避免 slice 越界
				stack := debug.Stack()
				maxBytes := 4096
				if len(stack) > maxBytes {
					stack = stack[:maxBytes]
				}
				ech <- fmt.Errorf("%v\n%s", err, string(stack))
			}
			g.wg.Done()
		}()
		ech <- fn(ctx)
	}()
	return ech
}

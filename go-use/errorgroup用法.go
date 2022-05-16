package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
//ErrorGroup会返回所有执行任务的goroutine遇到的第一个错误
func main() {
	var eg errgroup.Group
	for i := 0; i < 100; i++ {
		i := i
		eg.Go(func() error {
			time.Sleep(2 * time.Second)
			if i > 90 {
				fmt.Println("Error:", i)
				return fmt.Errorf("Error occurred: %d", i)
			}
			fmt.Println("End:", i)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}


// Go方法单独开启的goroutine在执行参数传递进来的函数时，如果函数返回了错误，会对ErrorGroup持有的err字段进行赋值并及时调用cancel函数，
// 通过上下文通知其他子任务取消执行任务。所以上面更新后的程序会有如下类似的输出。
*/
func main() {
	var eg errgroup.Group
	for i := 0; i < 100; i++ {
		i := i
		eg.Go(func() error {
			time.Sleep(2 * time.Second)
			if i > 90 {
				fmt.Println("Error:", i)
				return fmt.Errorf("Error occurred: %d", i)
			}
			fmt.Println("End:", i)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

// ErrorGroup的实现原理
// ErrorGroup原语的结构体类型errorgroup.Group定义如下：

type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// cancel — 创建 context.Context 时返回的取消函数，用于在多个 goroutine 之间同步取消信号；

// wg — 用于等待一组 goroutine 完成子任务的同步原语；

// errOnce — 用于保证只接收一个子任务返回的错误的同步原语；

// 通过 errgroup.WithContext构造器创建errgroup.Group 结构体：

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// 运行新的并行子任务需要使用errgroup.Group.Go方法，这个方法的执行过程如下：

// 调用 sync.WaitGroup.Add 增加待处理的任务数；
// 创建一个新的 goroutine 并在 goroutine 内部运行子任务；
// 返回错误时及时调用 cancel 并对 err 赋值，只有最早返回的错误才会被上游感知到，后续的错误都会被舍弃：

func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

// 用于等待的errgroup.Group.Wait方法只是调用了 sync.WaitGroup.Wait方法，阻塞地等待所有子任务完成。在子任务全部完成时会通过调用在errorgroup.WithContext创建Group和Context对象时存放在Group.cancel字段里的函数，取消Context对象并返回可能出现的错误。

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// 总结
// Go语言通过errorgroup.Group结构提供的ErrorGroup原语通过封装WaitGroup、Once基本原语结合上下文对象，提供了除同步等待外更加复杂的错误传播和执行任务取消的功能。在使用时，我们也需要注意它的两个特点：

// errgroup.Group在出现错误或者等待结束后都会调用 Context对象 的 cancel 方法同步取消信号。
// 只有第一个出现的错误才会被返回，剩余的错误都会被直接抛弃。

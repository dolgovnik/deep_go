package main

import (
	"testing"
	"container/heap"

	"github.com/stretchr/testify/assert"
)

type Queue []*Task

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i int, j int) bool {
	return q[i].Priority > q[j].Priority
}

func (q Queue) Swap(i int, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(a any) {
	task := a.(*Task)
	*q = append(*q, task)
}

func (q *Queue) Pop() any {
	o := *q
	n := len(o)
	t := o[n-1]
	*q = o[:n-1]
	return t
}

func (q *Queue) Update(taskID int, newPriority int) {
	for i, task := range *q {
		if task.Identifier == taskID {
			task.Priority = newPriority
			heap.Fix(q, i)
			return
		}
	}
}

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap *Queue
}

func NewScheduler() Scheduler {
	q := &Queue{}
	heap.Init(q)

	return Scheduler{heap: q}
}

func (s *Scheduler) AddTask(task Task) {
	t := &task
	heap.Push(s.heap, t)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.heap.Update(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	return *heap.Pop(s.heap).(*Task)
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}

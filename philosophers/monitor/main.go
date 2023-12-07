package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type PrintMessage struct {
	leftFork, philosopher, rightFork int
	acquire                          bool
}

type State int

const (
	Thinking State = iota
	Hungry
	Eating
)

type Philosopher struct {
	id    int
	table *Table
}

type Table struct {
	lock         *sync.Mutex
	conds        []*sync.Cond
	states       []State
	philosophers int
	printChannel chan PrintMessage
}

func NewTable(n int) *Table {
	states := make([]State, n)
	conds := make([]*sync.Cond, n)
	lock := sync.Mutex{}

	for i := range conds {
		conds[i] = sync.NewCond(&lock)
	}

	return &Table{lock: &lock, states: states, conds: conds, philosophers: n, printChannel: make(chan PrintMessage)}
}

func (t *Table) pickup(p *Philosopher) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.states[p.id] = Hungry

	for t.states[(p.id+1)%t.philosophers] == Eating || t.states[(p.id+t.philosophers-1)%t.philosophers] == Eating {
		t.conds[p.id].Wait()
	}

	t.printChannel <- PrintMessage{(p.id + 1) % t.philosophers, p.id, (p.id + t.philosophers - 1) % t.philosophers, true}

	t.states[p.id] = Eating
}

func (t *Table) putdown(p *Philosopher) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.states[p.id] = Thinking
	t.conds[(p.id+1)%t.philosophers].Signal()
	t.conds[(p.id+t.philosophers-1)%t.philosophers].Signal()

	t.printChannel <- PrintMessage{(p.id + 1) % t.philosophers, p.id, (p.id + t.philosophers - 1) % t.philosophers, false}
}

func (t *Table) printMessages() {
	printMessages := make([]PrintMessage, 0)

	for {
		printMessage := <-t.printChannel

		if printMessage.acquire {
			printMessages = append(printMessages, printMessage)
		} else {
			for i, message := range printMessages {
				if message.leftFork == printMessage.leftFork &&
					message.rightFork == printMessage.rightFork &&
					message.philosopher == printMessage.philosopher {
					printMessages = append(printMessages[:i], printMessages[i+1:]...)
					break
				}
			}
		}

		for _, message := range printMessages {
			fmt.Printf("(%d, %d, %d)\n", message.leftFork, message.philosopher, message.rightFork)
		}

		fmt.Println()
	}
}

func (p *Philosopher) eat() {
	for {
		n := rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		p.table.pickup(p)

		n = rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		p.table.putdown(p)
	}
}

func main() {
	n := 5
	table := NewTable(n)
	philosophers := make([]*Philosopher, n)

	go table.printMessages()

	for i := range philosophers {
		philosophers[i] = &Philosopher{id: i, table: table}
		go philosophers[i].eat()
	}

	fmt.Scanln()
}

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Semaphore struct {
	channel chan struct{}
}

func NewSemaphore(count int) *Semaphore {
	return &Semaphore{
		channel: make(chan struct{}, count),
	}
}

func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.channel
}

type PrintMessage struct {
	leftFork, philosopher, rightFork int
	acquire                          bool
}

type Philosopher struct {
	id                      int
	leftFork                *Semaphore
	rightFork               *Semaphore
	leftForkId, rightForkId int
}

func (p *Philosopher) Eat(printChannel chan PrintMessage) {
	p.leftFork.Acquire()
	p.rightFork.Acquire()

	printChannel <- PrintMessage{p.leftForkId, p.id, p.rightForkId, true}

	n := rand.Intn(500)
	time.Sleep(time.Duration(n) * time.Millisecond)

	p.rightFork.Release()
	p.leftFork.Release()

	printChannel <- PrintMessage{p.leftForkId, p.id, p.rightForkId, false}
}

func (p *Philosopher) Think() {
	n := rand.Intn(500)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func PrintMessages(printChannel chan PrintMessage) {
	printMessages := make([]PrintMessage, 0)

	for {
		printMessage := <-printChannel

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

func main() {
	numPhilosophers := 5
	philosophers := make([]*Philosopher, numPhilosophers)
	forks := make([]*Semaphore, numPhilosophers)

	printChannel := make(chan PrintMessage)

	go PrintMessages(printChannel)

	for i := 0; i < numPhilosophers; i++ {
		forks[i] = NewSemaphore(1)
	}

	for i := 0; i < numPhilosophers-1; i++ {
		philosophers[i] = &Philosopher{
			id:          i,
			leftFork:    forks[i],
			rightFork:   forks[(i+1)%numPhilosophers],
			leftForkId:  i,
			rightForkId: (i + 1) % numPhilosophers,
		}

		go func(p *Philosopher) {
			for {
				p.Think()
				p.Eat(printChannel)
			}
		}(philosophers[i])
	}

	philosophers[numPhilosophers-1] = &Philosopher{
		id:          numPhilosophers - 1,
		leftFork:    forks[0],
		rightFork:   forks[numPhilosophers-1],
		leftForkId:  numPhilosophers - 1,
		rightForkId: 0,
	}

	go func(p *Philosopher) {
		for {
			p.Think()
			p.Eat(printChannel)
		}
	}(philosophers[numPhilosophers-1])

	// Wait for philosophers to finish
	fmt.Scanln()
}

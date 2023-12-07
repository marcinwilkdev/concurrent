package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Entity struct {
	id     int
	reader bool
}

type PrintMessage struct {
	entity Entity
	start  bool
}

type ReadersWriters struct {
	readers      int
	lock         *sync.Mutex
	readerCond   *sync.Cond
	writerCond   *sync.Cond
	printChannel chan PrintMessage
}

func (readersWriters *ReadersWriters) StartRead() {
	readersWriters.lock.Lock()
	defer readersWriters.lock.Unlock()

	for readersWriters.readers == -1 {
		readersWriters.readerCond.Wait()
	}

	readersWriters.readers++
}

func (readersWriters *ReadersWriters) EndRead() {
	readersWriters.lock.Lock()
	defer readersWriters.lock.Unlock()

	readersWriters.readers--

	if readersWriters.readers == 0 {
		readersWriters.writerCond.Signal()
	}
}

func (readersWriters *ReadersWriters) StartWrite() {
	readersWriters.lock.Lock()
	defer readersWriters.lock.Unlock()

	for readersWriters.readers != 0 {
		readersWriters.writerCond.Wait()
	}

	readersWriters.readers = -1
}

func (readersWriters *ReadersWriters) EndWrite() {
	readersWriters.lock.Lock()
	defer readersWriters.lock.Unlock()

	readersWriters.readers = 0

	readersWriters.readerCond.Broadcast()
	readersWriters.writerCond.Signal()
}

func NewReadersWriters() *ReadersWriters {
	lock := sync.Mutex{}
	readerCond := sync.NewCond(&lock)
	writerCond := sync.NewCond(&lock)
	printChannel := make(chan PrintMessage)

	return &ReadersWriters{0, &lock, readerCond, writerCond, printChannel}
}

func (readersWriters *ReadersWriters) Reader(id int) {
	for {
		n := 500 + rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		readersWriters.StartRead()

		readersWriters.printChannel <- PrintMessage{Entity{id, true}, true}

		n = rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		readersWriters.EndRead()

		readersWriters.printChannel <- PrintMessage{Entity{id, true}, false}
	}
}

func (readersWriters *ReadersWriters) Writer(id int) {
	for {
		n := 500 + rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		readersWriters.StartWrite()

		readersWriters.printChannel <- PrintMessage{Entity{id, false}, true}

		n = rand.Intn(500)
		time.Sleep(time.Duration(n) * time.Millisecond)

		readersWriters.EndWrite()

		readersWriters.printChannel <- PrintMessage{Entity{id, false}, false}
	}
}

func (readersWriters *ReadersWriters) PrintMessages() {
	entitiesList := make([]Entity, 0)

	for {
		printMessage := <-readersWriters.printChannel

		if printMessage.start {
			entitiesList = append(entitiesList, printMessage.entity)
		} else {
			for i, entity := range entitiesList {
				if entity.id == printMessage.entity.id && entity.reader == printMessage.entity.reader {
					entitiesList = append(entitiesList[:i], entitiesList[i+1:]...)
					break
				}
			}
		}

		for _, entity := range entitiesList {
			if entity.reader {
				fmt.Printf("Reader %d\n", entity.id)
			} else {
				fmt.Printf("Writer %d\n", entity.id)
			}
		}

		fmt.Println()
	}
}

func main() {
	readerWriter := NewReadersWriters()
	readersCount := 5
	writersCount := 2

	go readerWriter.PrintMessages()

	for i := 0; i < readersCount; i++ {
		go readerWriter.Reader(i)
	}

	for i := 0; i < writersCount; i++ {
		go readerWriter.Writer(i)
	}

	fmt.Scanln()
}

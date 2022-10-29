package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// max buffer size
const MAX_BUFFER_SIZE = 10

// consumer producer loop times
const LOOP_TIMES = 100

var (
	// queue buffer
	buffer *list.List
	
	// current buffer index
	index int
	
	wg *sync.WaitGroup
	mutex *sync.Mutex
	cond_var *sync.Cond



	// Console Output Color
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
)


func Producer() {
	defer wg.Done()

	for i := 0; i < LOOP_TIMES; i++ {
		time.Sleep(500 * time.Millisecond)
		
		r := rand.Intn(100)
		
		cond_var.L.Lock()
		
		for buffer.Len() >= MAX_BUFFER_SIZE {
			fmt.Println(Green + "[PRODUCER START WAITING...]" + Reset)
			cond_var.Wait()
		}
		buffer.PushBack(r)
		fmt.Printf(Green + "Producer Create: %d\n" + Reset, r)
		
		cond_var.Signal()
		cond_var.L.Unlock()
	}
}

func Consumer() {
	defer wg.Done()

	for i := 0; i < LOOP_TIMES; i++ {
		time.Sleep(1000 * time.Millisecond)

		cond_var.L.Lock()
		for buffer.Len() == 0 {
			fmt.Println(Red + "[CONSUMER START WAITING...]" + Reset)
			cond_var.Wait()
		}

		n := buffer.Front()
		buffer.Remove(n)
		fmt.Printf(Red + "Consumer Get: %d\n" + Reset, n.Value)
		
		cond_var.Signal()
		cond_var.L.Unlock()
	}
}

func main() {
	// buffer init
	buffer = list.New()
	index = 0
	
	mutex = new(sync.Mutex)
	wg = &sync.WaitGroup{}
	cond_var = sync.NewCond(mutex)

	wg.Add(2)

	go Producer()
	go Consumer()

	wg.Wait()
	fmt.Println("Done")
}

package test_demo

import (
	"fmt"
	"sync"
	"testing"
)

type testMessageStruct struct {
	Number int
	lock   sync.Mutex
}

func TestLock(t *testing.T) {
	return
	fmt.Println("=========")
	var wg sync.WaitGroup
	m := testMessageStruct{Number: 0}
	wg.Add(10)
	for {
		m.lock.Lock()
		if m.Number < 10 {
			go func() {
				defer m.lock.Unlock()
				m.Number = m.Number + 1
				fmt.Println("m.Number === ", m.Number)

				wg.Done()
			}()

		} else {
			break
		}

	}
	wg.Wait()
}
func TestChannel(t *testing.T) {
	 c := make(chan int, 15)
	 go func() {
		 for i := 0 ; i< 6 ; i++ {
			 c <- i
		 }
	 }()

	 go func() {
		 for i := range c {
			 fmt.Println("i === ",i, "\nsize === ", len(c))
		 }
		 fmt.Println("this rang is quit!!")
	 }()
	 var forever chan bool
	 <- forever

}

package main

import (
	"fmt"
	"time"
)

func myProcess(stopChannel chan int) {
	for {
		toots := <-stopChannel
		fmt.Println("The letter", toots)
	}

}

func main() {
	stopChannel := make(chan int, 2)
	go myProcess(stopChannel)
	stopChannel <- 1
	time.Sleep(2 * time.Second)
	stopChannel <- 2
	time.Sleep(time.Second * 2)

	fmt.Println("Main Goroutine exited")
}

// Start Race with key
// 3 seconds to first blue, listen for redlight
//3 seconds to second blue,
// .70 - 1.3 to yellow
// .400 green
// Start timer and wait for Action Key

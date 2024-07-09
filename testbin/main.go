package main

import (
	"fmt"
	"time"
)

func main() {
	active := false
	stage1 := 2
	stage2 := 4
	stage3 := 6
	now := time.Now()

	defer func() {
		fmt.Println(time.Since(now))
	}()

	for stg := 0; stg < 3; stg++ {
		switch {
		case !active:
			ch := make(chan string)
			go func() {
				switch {
				case stg == 0:
					//active = true
					go func() {
						time.Sleep(time.Second * time.Duration(stage1))
						active = false
						ch <- "Stage 1"
					}()

				case stg == 1:
					//active = true
					go func() {
						time.Sleep(time.Second * time.Duration(stage2))
						active = false
						ch <- "Stage 2"
					}()
				case stg == 2:
					//active = true
					go func() {
						time.Sleep(time.Second * time.Duration(stage3))
						active = false
						ch <- "Stage 3"
					}()
				}
			}()
			msg := <-ch
			fmt.Println(msg)
		}
	}
}

package main

import (
	"fmt"
	"time"
)

func main() {
	i := 0
	for i < 999 {
		time.Sleep(1 * time.Second)
		fmt.Println("do some thing")
		fmt.Print("do some thing")
		i++
	}
}

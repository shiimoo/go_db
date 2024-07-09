package main

import "log"

func main() {
	a := make(chan any, 1000)
	a <- 1
	a <- 1
	a <- 1
	a <- 1
	a <- 1
	close(a)
	log.Println(a == nil)
	a <- 1
	for i := range a {
		log.Println(i)
	}
}

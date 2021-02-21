package main

import (
	"../CSP"
	"fmt"
	"sync"
)

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	wg := sync.WaitGroup{}
	m := 3
	length := len(s)
	ch := CSP.MakeChannel(0)
	for i := 0; i < m; i++ {
		CSP.MakeProcessWithChannel(calSum, s[length / m * i : length / m * (i + 1)], &wg, ch)
	}
	sum := 0
	for i := 0; i < m; i++ {
		v := CSP.ChannelInput(ch)
		sum += v
	}
	CSP.EndMainProcess(&wg)
	fmt.Println(sum)
}

func calSum(s []int, wg *sync.WaitGroup, ch chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	CSP.ChannelOutput(ch, sum)
	CSP.EndProcessNotMain(wg)
}
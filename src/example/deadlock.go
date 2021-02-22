package main

import (
	"../CSP"
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	ch0 := CSP.MakeChannel(0)
	ch1 := CSP.MakeChannel(0)
	args := []int{1, 2, 3}
	CSP.MakeProcessWithChannel(f1, args, &wg, ch0)
	CSP.MakeProcessWithChannel(f2, args, &wg, ch1)
	CSP.ChannelOutput(ch1, args[1])
	v := CSP.ChannelInput(ch0)
	CSP.EndMainProcess(&wg)
	fmt.Println(v)
}

func f1(args []int, wg *sync.WaitGroup, ch chan int) {
	v := CSP.ChannelInput(ch)
	CSP.EndProcessNotMain(wg)
	fmt.Println(v)
}

func f2(args []int, wg *sync.WaitGroup, ch chan int) {
	CSP.ChannelOutput(ch, args[0])
	CSP.EndProcessNotMain(wg)
}

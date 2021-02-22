package CSP

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var cspText = ""
var channels = make(map[chan int] string)
var process = make(map[string] []string)
var channelIdx = 0
var processIdx = 0
var now = "MAIN"

func MakeChannel(buffer int) chan int {
	if getMode() {
		fmt.Println("Define a channel named ch" + strconv.Itoa(channelIdx) + "...")
		ch := make(chan int, buffer)
		channels[ch] = "ch" + strconv.Itoa(channelIdx)
		cspText += "channel " + channels[ch] + " " + strconv.Itoa(buffer) + ";\n"
		channelIdx++
		return ch
	} else {
		return make(chan int, buffer)
	}
}

func ChannelOutput(ch chan int, data int) {
	if getMode() {
		fmt.Println("The Process " + now + " Output data to Channel " + channels[ch] + "...")
		process[now] = append(process[now], channels[ch] + "!" + "1")
	} else {
		ch <- data
	}
}

func ChannelInput(ch chan int) int {
	if getMode() {
		fmt.Println("The Process " + now + " Input data from Channel " + channels[ch] + "...")
		process[now] = append(process[now], channels[ch] + "?" + "1")
		return 0
	} else {
		return <- ch
	}
}

func MakeProcessWithOutChannel(f func([]int, *sync.WaitGroup), ar []int, wg *sync.WaitGroup) {
	if !getMode() {
		wg.Add(1)
		go f(ar, wg)
	}
}

func MakeProcessWithChannel(f func([]int, *sync.WaitGroup, chan int), ar []int, wg *sync.WaitGroup, ch chan int) {
	if getMode() {
		fmt.Println("Make a Process named P" + strconv.Itoa(processIdx) + "...")
		now = "P" + strconv.Itoa(processIdx)
		processIdx++
		f(ar, wg, ch)
		now = "MAIN"
	} else {
		wg.Add(1)
		go f(ar, wg, ch)
	}
}

func EndMainProcess(wg *sync.WaitGroup) {
	if getMode() {
		fmt.Println("Generate CSP Script....")
		cspText += "\n"
		var pros [] string
		for n, p := range process {
			s := n + "() = "
			pros = append(pros, n + "()")
			for _, v := range p {
				s += v + " -> "
			}
			s += "Skip;\n"
			cspText += s
		}

		cspText += "\nP() = "
		for idx, p := range pros {
			if idx == 0 {
				cspText += p
			} else {
				cspText += " || " + p
			}
		}
		cspText += ";\n"
		f, err := os.Create("csp/" + getOutFileName())
		if err != nil {
			panic(err)
		}
		f.WriteString(cspText)
		defer f.Close()
		fmt.Println("Done!")
		os.Exit(0)
	}
	wg.Wait()
}

func EndProcessNotMain(wg *sync.WaitGroup)  {
	if !getMode() {
		wg.Done()
	}
}

func getMode() bool {
	for idx, arg := range os.Args {
		if idx == 1 && arg == "-csp" {
			return true
		}
	}
	return false
}

func getOutFileName() string {
	name := ""
	for idx, arg := range os.Args {
		if idx == 2 {
			name = arg
		}
	}
	if name == "" {
		name = "main.csp"
	}
	return name
}

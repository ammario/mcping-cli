package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/ammario/mcping"
	"github.com/fatih/color"
	"os"
	"time"
)

func main() {
	//Parse command line flags
	var (
		_host     = kingpin.Arg("host", "Hostname").Required().String()
		_port     = kingpin.Flag("port", "Port number").Short('p').Default("25565").Int()
		_timeout  = kingpin.Flag("timeout", "Timeout per ping (ms)").Short('t').Default("300").Int()
		_interval = kingpin.Flag("interval", "Interval between pings (ms)").Short('i').Default("1000").Int()
		_count    = kingpin.Flag("count", "Total amount of pings to send").Short('c').Default("4294967295").Int()
		_debug    = kingpin.Flag("debug", "Debug underlying library.").Short('d').Default("false").Bool()
	)

	kingpin.Parse()

	host := *_host
	port := uint16(*_port)
	timeout := time.Duration(*_timeout)
	interval := time.Duration(*_interval)
	count := uint32(*_count)

	if *_debug {
		resp, err := mcping.Ping(host, port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Latency:", resp.Latency, "ms")
		fmt.Println("Online:", resp.Online, " players")
		fmt.Println("Max:", resp.Max, " players")
		fmt.Println("Protocol:", resp.Protocol)
		fmt.Println("Favicon:", len(resp.Favicon), " bytes")
		fmt.Println("Motd:", resp.Motd)
		fmt.Println("Server:", resp.Server)
		fmt.Println("Version:", resp.Version)
		fmt.Println("Player Sample:")
		for i, e := range resp.Sample {
			fmt.Println("#", i)
			fmt.Println(" 	UUID ->", e.UUID)
			fmt.Println(" 	Name ->", e.Name)
		}
		os.Exit(0)
	}

	//Sequentiating ping responses
	var id uint32

	//Output colorer
	failOut := color.New(color.FgRed)

	//Change color based on whether player count changed
	stayOut := color.New(color.FgWhite)
	riseOut := color.New(color.FgGreen)
	dropOut := color.New(color.FgYellow)

	var lastPlayerCount int

	//Output format strings
	failFmt := "(%x) %s; %s\n"
	successFmt := "(%x) %s; latency=%vms players=(%s)\n"

	fullAddr := fmt.Sprint(host, ":", port)
	pendingPings := 0
	for id = 0; id < count; id++ {
		//Have each request asynchronous
		go func(pid uint32) {
			pendingPings = pendingPings + 1
			err := errors.New("")
			timeoutChan := make(chan bool, 1)

			//Default response
			resp := mcping.PingResponse{}

			//Race timeout against ping
			go func() {
				time.Sleep(timeout * time.Millisecond)
				timeoutChan <- true
			}()
			go func() {
				resp, err = mcping.Ping(host, port)
				timeoutChan <- false
			}()
			if <-timeoutChan {
				err = mcping.TimeoutErr
			}
			if err != nil {
				failOut.Printf(failFmt, pid, fullAddr, err)
			} else {
				playerCount := fmt.Sprint(resp.Online, "/", resp.Max)
				latency := resp.Latency
				if resp.Online == lastPlayerCount {
					stayOut.Printf(successFmt, pid, fullAddr, latency, playerCount)
				} else if resp.Online > lastPlayerCount {
					riseOut.Printf(successFmt, pid, fullAddr, latency, playerCount)
				} else if resp.Online < lastPlayerCount {
					dropOut.Printf(successFmt, pid, fullAddr, latency, playerCount)
				}
				lastPlayerCount = resp.Online
			}
			pendingPings = pendingPings - 1
		}(id)
		time.Sleep(interval * time.Millisecond)
	}
}

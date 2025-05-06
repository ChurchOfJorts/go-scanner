package main

import (
	"fmt"
	"net"
	"os"
	"sort"
)

func scan(ports, results chan int, target string) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p) //scanme.nmap.org
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	rhost := os.Args[1]
	ports := make(chan int, 100) //Initialize a channel for 100 ports
	results := make(chan int)    //Make a channel to store the results
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go scan(ports, results, rhost) //Scan the port, and save the output to results
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ { //Loop through the results to find open ports
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results) //Cleanup ports and results
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

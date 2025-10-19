package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	netUdp, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, netUdp)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		delim := '\n'

		line, err := reader.ReadString(byte(delim))
		if err != nil {
			fmt.Println(err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
		}
	}
}

package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "192.168.4.100:4242"

func main() {

	log.Println("dialing quic...")
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("session opened")

	stream, err := session.OpenStreamSync()
	if err != nil {
		panic(err)
	}
	log.Println("synchronized stream opened")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Printf("Client: Sending '%s'\n", scanner.Text())
		start := time.Now()
		_, err = stream.Write([]byte(scanner.Text()))
		if err != nil {
			panic(err)
		}

		buf := make([]byte, len(scanner.Text()))
		_, err = io.ReadFull(stream, buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Client: Got '%s'\n", buf)
		elapsed := time.Since(start)
		log.Printf("Binomial took %s", elapsed)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		panic(err)
	}
}

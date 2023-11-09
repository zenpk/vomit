package main

import (
	"fmt"
	"log"
	"math"
	"net"
)

const (
	port = 9999
	host = "0.0.0.0"
)

func main() {
	go listenTcp()
	go listenUdp()
	log.Printf("ready to vomit on port %v\n", port)
	select {}
}

func listenTcp() {
	const connType = "tcp"
	tcpListener, err := net.Listen(connType, fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		log.Panic(err)
	}
	defer tcpListener.Close()
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Printf("tcpListener error: %v", err)
		}
		defer conn.Close()
		buffer := make([]byte, math.MaxInt32)
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				log.Printf("%s connection error: %v", connType, err)
				break
			}
			printInfo(connType, conn.RemoteAddr(), buffer[:n], err)
		}
	}
}

func listenUdp() {
	const connType = "udp"
	udpAddr, err := net.ResolveUDPAddr(connType, fmt.Sprintf("%s:%v", host, port))
	if err != nil {
		log.Panic(err)
	}
	udpConn, err := net.ListenUDP(connType, udpAddr)
	if err != nil {
		log.Panic(err)
	}
	defer udpConn.Close()
	buffer := make([]byte, math.MaxInt32)
	for {
		n, addr, err := udpConn.ReadFromUDP(buffer)
		printInfo(connType, addr, buffer[:n], err)
	}
}

func printInfo(connType string, addr net.Addr, data []byte, err error) {
	if err != nil {
		log.Printf("%s connection error: %v", connType, err)
		return
	}
	if len(data) <= 0 {
		log.Printf("data length should be larger than 0")
		return
	}
	info := fmt.Sprintf("received a %v packet from %v", connType, addr)
	raw := ""
	for i, b := range data {
		if i > 0 && i%4 == 0 {
			raw += " "
		}
		if i > 0 && i%8 == 0 {
			raw += "\n"
		}
		raw += fmt.Sprintf("%02x ", b)
	}
	const div = "\n==============================\n"
	log.Printf("%v\n%v%s%v%s", info, raw, div, string(data), div)
}

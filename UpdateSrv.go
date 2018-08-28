package main

import (
	"fmt"
	"net"
	"time"
)

func UDP_Updates() {
	/* Lets prepare a address at any address at port 10001*/
	ServerAddr, _ := net.ResolveUDPAddr("udp", ":7778")
	fmt.Println("listening on :7778")

	/* Now listen at selected port */
	ServerConn, _ := net.ListenUDP("udp", ServerAddr)
	defer ServerConn.Close()

	buf := make([]byte, 65535)

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("error: ", err)
		}

		if n == 31380 {
			for i := 0; i < 10460; i++ {
				start := (i * 3)
				end := (i * 3) + 3
				config.LEDs[i][0] = float32(buf[start:end][0])
				config.LEDs[i][1] = float32(buf[start:end][1])
				config.LEDs[i][2] = float32(buf[start:end][2])
			}
		}
	}
}

func SendLed(LEDs [][]float32) {
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3663")
	if err != nil {
		panic(err)
	}

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		panic(err)
	}

	defer Conn.Close()
	i := 0
	msg := LEDs
	i++
	buf := make([]byte, 10460*3)

	for i := 0; i < len(msg); i++ {
		for j := 0; j < len(msg[i]); j++ {
			buf[(i*3)+j] = byte(msg[i][j])
		}
	}

	gg, err := Conn.Write(buf)
	if err != nil {
		fmt.Println(msg, err)
		fmt.Println(gg)
	}
	time.Sleep(time.Second * 1)
}

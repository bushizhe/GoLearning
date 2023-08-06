package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func usage() {
	msg := "Need to run as root"
	fmt.Println(msg)
	os.Exit(0)
}

func sendICMPRequest(icmp ICMP, dstAddr *net.IPAddr) error {
	conn, err := net.DialIP("ip4:icmp", nil, dstAddr)
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)

	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	start := time.Now()
	conn.SetReadDeadline(time.Now().Add(time.Second * 4))
	recv := make([]byte, 1024)
	res, err := conn.Read(recv)
	if err != nil {
		return err
	}
	end := time.Now()
	duration := end.Sub(start).Nanoseconds() / 1e6
	fmt.Printf("%d bytes from %s: seq=%d time=%dms\n", res, dstAddr.String(), icmp.SequenceNum, duration)
	return err
}

func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

// ICMP校验和算法
func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += sum >> 16
	return uint16(^sum)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	host := os.Args[1]
	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		fmt.Printf("Fail to resolve %s, %s\n", host, err)
		return
	}

	fmt.Printf("Ping %s (%s):\n\n", addr.String(), host)

	for i := 1; i < 6; i++ {
		if err = sendICMPRequest(getICMP(uint16(i)), addr); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		time.Sleep(2 * time.Second)
	}
}

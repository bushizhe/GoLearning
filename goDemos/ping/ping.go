package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

// ICMP 定义ICMP报文头结构
type ICMP struct {
	Type        uint8  // ICMP类型，设置为8
	Code        uint8  // 进一步划分ICMP类型，ping使用的是echo类型的ICMP，设置为0
	CheckSum    uint16 // 报文头的校验值，防止数据传输错误，先
	Identifier  uint16 // 表示一个ICMP，可以设置为0
	SequenceNum uint16 // 序列号，发送ICMP报文时依次累加
}

func usage() {
	msg := "Need to run as root"
	fmt.Println(msg)
	os.Exit(0)
}

// 基于序列号生产ICMP报文头
func getICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp) // 网络中传输的数据需要是大端字节序的
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()

	return icmp
}

func sendICMPRequest(icmp ICMP, dstAddr *net.IPAddr) error {
	conn, err := net.DialIP("ip4:icmp", nil, dstAddr) // 创建一个ICMP报文
	if err != nil {
		fmt.Printf("Fail to connect to remote host: %s\n", err)
		return err
	}
	defer conn.Close()

	var buffer bytes.Buffer
	// 填充报文并发送
	binary.Write(&buffer, binary.BigEndian, icmp)
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return err
	}

	start := time.Now()
	conn.SetReadDeadline(time.Now().Add(time.Second * 4))
	// 接收请求
	recv := make([]byte, 1024)
	res, err := conn.Read(recv)
	if err != nil {
		return err
	}
	end := time.Now()
	// 计算发送到接收消耗的时间
	duration := end.Sub(start).Nanoseconds() / 1e6
	fmt.Printf("%d bytes from %s: seq=%d time=%dms\n", res, dstAddr.String(), icmp.SequenceNum, duration)
	return err
}

// CheckSum 计算校验值
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

	fmt.Printf("Ping %s (%s):\n", addr.String(), host)

	for i := 1; i < 6; i++ {
		if err = sendICMPRequest(getICMP(uint16(i)), addr); err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		time.Sleep(2 * time.Second)
	}
}

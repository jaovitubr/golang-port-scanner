package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func incrementIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}

func scanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, time.Second*5)

	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

var analyticsCount int
var commonPorts []int = []int{
	// 80,   // HTTP
	// 443,  // HTTPS
	// 21,   // FTP
	// 22,   // FTPS / SSH
	// 3306, // MySQL
	// 2082, // cPanel
	// 2083, // cPanel SSL
	// 2086, // WHM
	// 2087, // WHM SSL
	25565,
}

func main() {
	simultaneous := 0
	initialIp := net.ParseIP("127.0.0.1")
	f, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)

	go analytics()

	for i := 1; i <= 10000; i++ {
		simultaneous += 1

		ip := incrementIP(initialIp, uint(i))

		go func() {
			for _, port := range commonPorts {
				fmt.Printf("\033[2K\rScanning ip %s at port %d", ip, port)
				open := scanPort("tcp", ip.To4().String(), port)

				if open {
					fmt.Printf("\033[2K\r%s port %d is open!\n", ip, port)
					f.WriteString(ip.To4().String() + ":" + strconv.Itoa(port) + "\n")
				}
			}

			simultaneous -= 1
			analyticsCount += 1
		}()

		time.Sleep(time.Millisecond * 10)

		for simultaneous >= 1024 {
			time.Sleep(time.Millisecond * 10)
		}
	}

	f.Close()
}

func analytics() {
	for true {
		time.Sleep(time.Minute)
		fmt.Printf("\033[2K\rCheck/s: %d\n", analyticsCount)
		analyticsCount = 0
	}
}

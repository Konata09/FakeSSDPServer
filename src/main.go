package main

import (
	"github.com/koron/go-ssdp"
	"os"
	"strconv"
	"time"
)

var mode = "prod"

func initEnv() {
	if os.Getenv("ENSYSLOG") == "true" || os.Getenv("ENSYSLOG") == "TRUE" {
		syslogEnable = true
	}
	if os.Getenv("SYSLOGADDR") != "" {
		syslogAddr = os.Getenv("SYSLOGADDR")
	}
	if os.Getenv("SYSLOGTAG") != "" {
		syslogTag = os.Getenv("SYSLOGTAG")
	}
	if os.Getenv("SYSLOGLEVEL") != "" {
		envSyslogLevel, err := strconv.Atoi(os.Getenv("SYSLOGLEVEL")) // need Integer
		if err != nil {
			syslogLevel = envSyslogLevel
		}
	}
	if os.Getenv("LOGLEVEL") != "" {
		envLogLevel, err := strconv.Atoi(os.Getenv("LOGLEVEL")) // need Integer
		if err != nil {
			logLevel = envLogLevel
		}
	}
	if os.Getenv("LOGTIME") == "false" || os.Getenv("LOGTIME") == "FALSE" {
		logTime = false
	}
	if os.Getenv("MODE") == "dev" {
		mode = "dev"
	}
}

//func test(matrix [][]int) [][]int {
//	if matrix == nil {
//		return nil
//	}
//	m, n := len(matrix), len(matrix[0])
//	res := make([][]int, n)
//	for i := 0; i < n; i++ {
//		res[i] = make([]int, m)
//	}
//	for i := 0; i < m; i++ {
//		for j := 0; j < n; j++ {
//			res[j][i] = matrix[i][j]
//		}
//	}
//	return res
//}
//
//func test2(s string) bool {
//	hash := map[byte]byte{')': '('}
//	stack := make([]byte, 0)
//	if s == "" {
//		return true
//	}
//
//	for i := 0; i < len(s); i++ {
//		if s[i] == '(' {
//			stack = append(stack, s[i])
//		} else if len(stack) > 0 && stack[len(stack)-1] == hash[s[i]] {
//			stack = stack[:len(stack)-1]
//		} else {
//			return false
//		}
//	}
//	return len(stack) == 0
//}

func main() {
	initEnv()
	//test2()
	logInfo("Fake SSDP Server started")
	ssdp.SetMulticastRecvAddrIPv4("239.255.255.250:1900")
	ad, err := ssdp.Advertise(
		"upnp:rootdevice",
		"uuid:935607b0-243b-11e9-8000-e061d673a367::upnp:rootdevice",
		"http://192.168.11.4:10000/",
		"Windows-NT/10.0 UPnP/1.0 Media-Rendering-System/1.0 Fake-SSDP-Server/1.0\r\nX-AV-Physical-Unit-Info: pa=\"IZUMIKONATA: MediaGo\";pl=;",
		1800)
	if err != nil {
		panic(err)
	}
	ad.Alive()
	aliveTick := time.Tick(120 * time.Second)
	for {
		select {
		case <-aliveTick:
			logDebug("Advertise")
			ad.Alive()
		}
	}
}

package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/nameoffnv/time"
)

var (
	errZerotime = fmt.Errorf("time server returns zero time")
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("usage: time-client host [port]")
		os.Exit(1)
	}

	host := os.Args[1]
	port := "37"
	if len(os.Args) == 3 {
		port = os.Args[2]
	}

	ts, err := getTime(host, port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(ts)
}

func getTime(host, port string) (uint32, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return 0, fmt.Errorf("tcp dial: %v", err)
	}
	defer conn.Close()

	b, err := ioutil.ReadAll(conn)
	if err != nil {
		return 0, fmt.Errorf("read from connection: %v", err)
	}

	if len(b) == 0 {
		return 0, errZerotime
	}

	ts := binary.BigEndian.Uint32(b)

	if ts == 0 {
		return 0, errZerotime
	}

	return ts - time.RFCTimeShift, nil
}

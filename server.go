package time

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// RFC868 time is the number of seconds since 00:00 (midnight) 1 January 1900 GMT
	RFCTimeShift = 2208988800
)

type TimeServer struct {
	timeFunc TimeFunc
	done     chan struct{}
	listener net.Listener
}

type TimeFunc func() int64

func RFCTime() int64 {
	t := time.Now().Unix() + RFCTimeShift
	return t
}

func New(timeFunc TimeFunc) *TimeServer {
	return &TimeServer{
		timeFunc: timeFunc,
		done:     make(chan struct{}),
	}
}

func (ts *TimeServer) Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("net listen failed: %v", err)
	}
	ts.listener = listener

	log.Printf("listen %s", addr)

	for {
		select {
		case <-ts.done:
			return nil
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept connection failed: %v", err)
			continue
		}

		go ts.handleConnection(conn)
	}
}

func (ts *TimeServer) Close() error {
	close(ts.done)
	return ts.listener.Close()
}

func (ts *TimeServer) handleConnection(conn net.Conn) {
	log.Printf("incoming connection from %v", conn.RemoteAddr().String())

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(ts.timeFunc()))

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err := conn.Write(b); err != nil {
		log.Printf("write to socket failed: %v", err)
	}

	conn.Close()
}

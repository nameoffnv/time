package time

import (
	"encoding/binary"
	"io/ioutil"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	unixTime := time.Now().Unix()
	rfcTime := unixTime + RFCTimeShift

	cases := []struct {
		test     string
		timeFn   TimeFunc
		excepted int64
	}{
		{
			test: "one-date",
			timeFn: func() int64 {
				return time.Unix(1, 0).Unix()
			},
			excepted: 1,
		},
		{
			test: "current-date-rfc-time",
			timeFn: func() int64 {
				return rfcTime
			},
			excepted: unixTime + RFCTimeShift,
		},
		{
			test: "1 Jan 1970 GMT",
			timeFn: func() int64 {
				return 0 + RFCTimeShift
			},
			excepted: 2208988800,
		},
		{
			test: "1 Jan 1976 GMT",
			timeFn: func() int64 {
				return 189302400 + RFCTimeShift
			},
			excepted: 2398291200,
		},
		{
			test: "17 Nov 1858 GMT",
			timeFn: func() int64 {
				return -3506716800 + RFCTimeShift
			},
			excepted: -1297728000,
		},
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			s := New(c.timeFn)

			writer, reader := net.Pipe()

			go s.handleConnection(writer)

			b, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			utime := binary.BigEndian.Uint32(b)
			if utime != uint32(c.excepted) {
				t.Fatalf("bad time value, excepted %d, actual %d", c.excepted, utime)
			}
		})
	}
}

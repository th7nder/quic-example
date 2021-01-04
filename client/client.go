package client

import (
	"crypto/tls"
	"sync"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/pkg/errors"
	"github.com/th7nder/quic-example/utils"
)

func stream(session quic.Session, wg *sync.WaitGroup) {
	utils.Infof("Waiting for stream")
	s, err := session.OpenStream()
	if err != nil {
		panic(err)
	}
	utils.Infof("[SID: %d] Opened\n", s.StreamID())

	_, err = s.Write([]byte("HELO"))
	if err != nil {
		panic(errors.Wrap(err, "failed to HELO"))
	}

	utils.Infof("[SID: %d] Started data stream", s.StreamID())
	buf := make([]byte, 1024)
	var readBytes int
	start := time.Now()
	for {
		n, err := s.Read(buf)
		if err != nil {
			utils.Infof("[SID: %d] Error: %s", s.StreamID(), err)
			break
		} else {
			readBytes += n
		}
	}
	elapsed := time.Since(start)

	utils.Infof("[SID: %d] Elapsed: %s, bytes read: %d", s.StreamID(), elapsed, readBytes)
	wg.Done()
}

// Client connects to addr, opens n streams and downloads all data it can receive from a stream
func Client(addr string, streams int, multipath bool, game bool) error {
	session, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, &quic.Config{})
	if err != nil {
		return err
	}

	if game {
		streams = 2
	}
	var wg sync.WaitGroup
	for i := 0; i < streams; i++ {
		wg.Add(1)
		go stream(session, &wg)
	}

	wg.Wait()
	return nil
}

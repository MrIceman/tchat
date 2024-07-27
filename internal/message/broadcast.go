package message

import (
	"net"
	"sync"
)

func Broadcast(conns []net.Conn, b []byte) (notReached []net.Conn) {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(len(conns))
	for _, conn := range conns {
		go func(wg *sync.WaitGroup) {
			if err := Transmit(conn, b); err != nil {
				mu.Lock()
				notReached = append(notReached, conn)
				mu.Unlock()
			}
			wg.Done()
		}(&wg)
	}

	wg.Wait()

	return notReached
}

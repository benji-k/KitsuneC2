package listener

import (
	"log"
	"net"
	"strconv"
)

// The listener object represent a listiner of a certain type (e.g. HTTP listener, TCP etc.). The Handler function gets called upon each new
// connection to the listener.
type Listener struct {
	Type      string //for now only TCP-type listeners exist. Might be extended in the future with e.g. HTTP(S), DNS etc.
	Handler   func(net.Conn)
	Network   string //interface on which to listen, leave empty to listen on all available interfaces
	Port      int
	ls        net.Listener
	terminate bool //if set to true, terminates the listener
}

// Start is a non-blocking function. Calling Start will make the listener start listening for connections is a seperate thread.
// Please note that all net.Conn objects that are passed to Handler are expected to be closed by the handler.
func (l *Listener) Start() error {
	log.Printf("[INFO] Attempting to start listener on: %s:%d\n", l.Network, l.Port)
	var err error
	l.ls, err = net.Listen("tcp", l.Network+":"+strconv.Itoa(l.Port))
	if err != nil {
		log.Printf("[ERROR] Failed opening listener on: %s:%d\n		Reason: %s", l.Network, l.Port, err.Error())
		return err
	}
	l.terminate = false

	//Because we want Start() to be non-blocking, we run this in a seperate thread.
	go func() {
		log.Printf("[INFO] Successfully started listener on: %s:%d\n", l.Network, l.Port)
		for !l.terminate {
			conn, err := l.ls.Accept()
			if err == nil {
				go l.Handler(conn)
			}
		}
		log.Printf("[INFO] Successfully stopped listener on: %s:%d\n", l.Network, l.Port)
	}()

	return nil
}

// Stops the listening thread by calling Close() to the listener and making sure the loop in Start() exits.
func (l *Listener) Stop() {
	log.Printf("[INFO] Attempting to stop listener on: %s:%d\n", l.Network, l.Port)
	l.ls.Close()
	l.terminate = true
}

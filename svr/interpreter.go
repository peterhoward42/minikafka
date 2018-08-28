package  svr

import (
    "bufio"
    "encoding/gob"
	"log"
	"net"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
	}
}

func (*Interpreter) Interpret(conn net.Conn) {
	defer conn.Close()
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
    dec := gob.NewDecoder(rw)
    var commandCode CommandCode

	// Switch depending on command code at start of each incoming message.
	for {
        err := dec.Decode(&commandCode) // Blocks waiting for sufficient input.
        if err != nil {
            log.Println("Error decoding GOB data:", err)
            return
        log.Println("Recevied command code: %v", commandCode)
        }
    }
}

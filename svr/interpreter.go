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
    /*
     * Loop forever reading and interpreting incoming Gobs.
     */
	defer conn.Close()
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
    dec := gob.NewDecoder(rw)
    var commandCode CommandCode

	// Switch depending on command code at start of each incoming message.
	for {
        log.Printf("Now call decode()")
        err := dec.Decode(&commandCode) // Blocks waiting for sufficient input.
        if err != nil {
            log.Printf("Error decoding GOB data, with: %v.", err)
            continue
        }
        log.Printf("Decoded command as: %v", commandCode)
    }
}

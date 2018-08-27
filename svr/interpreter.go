package  svr

import (
	"log"
	"net"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
	}
}

type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

func (*Interpreter) Interpret(conn net.Conn) {
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// For each incoming command...
	for {
	    var data complexData
	    dec := gob.NewDecoder(rw)
        err := dec.Decode(&data) // Blocks waiting for sufficient input.
        if err != nil {
            log.Println("Error decoding GOB data:", err)
            return
        }
    }
}

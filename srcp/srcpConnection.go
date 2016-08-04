package srcp

import (
	"bufio"
	"log"
	"net"
)

type SrcpConnection struct {
	connection net.Conn
	reader     *bufio.Reader
}

func (srcpConnection *SrcpConnection) Connect(url string) {
	var error error
	srcpConnection.connection, error = net.Dial("tcp", url)
	if error != nil {
		panic(error)
	}
	srcpConnection.reader = bufio.NewReader(srcpConnection.connection)
}

func (srcpConnection *SrcpConnection) SendAndReceive(request string) string {
	log.Printf(request)
	if _, error := srcpConnection.connection.Write([]byte(request + "\n")); error != nil {
		panic(error)
	}
	return srcpConnection.Receive()
}

func (srcpConnection *SrcpConnection) Receive() string {
	var reply string
	var error error
	if reply, error = srcpConnection.reader.ReadString('\n'); error != nil {
		panic(error)
	}
	log.Printf(reply)
	return reply
}

func (srcpConnection *SrcpConnection) Close() {
	srcpConnection.connection.Close()
}

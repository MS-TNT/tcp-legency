package server

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) StartTcpServer(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetString("port")
	host, _ := cmd.Flags().GetString("ip")
	addr := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	defer listener.Close()
	fmt.Println("Listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		fmt.Println("Accepted a connection")
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		var recvNs int64
		err := binary.Read(conn, binary.BigEndian, &recvNs)
		if err != nil {
			fmt.Println("read not a int type, parse err:", err.Error())
			return
		}
		err = binary.Write(conn, binary.BigEndian, recvNs)
		if err != nil {
			fmt.Println("send error:", err.Error())
			return
		}

	}
}

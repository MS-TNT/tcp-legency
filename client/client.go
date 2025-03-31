package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}
func (c *Client) StartClient(cmd *cobra.Command, args []string) error {
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetString("port")

	conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		return err
	}
	defer conn.Close()
	fmt.Println("Connected to server")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Second)
	count := 0
	sum := 0.0
	printCount := 0
	legacyMsSum := float32(0)
	for {
		select {
		case <-sigChan:
			fmt.Println("\ntotal:", printCount, "avg legacy:", legacyMsSum/float32(printCount), "ms")
			return nil
		case <-ticker.C:
			legacyMs := float32(sum) / float32(count) / (1000 * 1000)
			printCount += 1
			fmt.Println("seq:", printCount, "legacy:", legacyMs, "ms")
			sum = 0.0
			count = 0
			legacyMsSum += legacyMs
		default:
			now := time.Now().UnixNano() //ns
			err = binary.Write(conn, binary.BigEndian, now)
			if err != nil {
				fmt.Println("send error:", err.Error())
				return err
			}

			var recvNs int64
			err = binary.Read(conn, binary.BigEndian, &recvNs)
			if err != nil {
				fmt.Println("read error:", err.Error())
				return err
			}
			periodNs := time.Now().UnixNano() - recvNs
			count += 1
			sum += float64(periodNs) / 2
		}
	}
}

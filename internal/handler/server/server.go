package server

import (
	"bufio"
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	challengeResp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"net"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

func SetupHandlers(ctx context.Context, core application.Core, cfg *config.Config) error {
	ctx = application.WithApp(ctx, core)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))
	if err != nil {
		return err
	}

	// Close the listener when the application closes.
	defer listener.Close()
	fmt.Println("listening", listener.Addr())
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := ProcessRequest(req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMsg(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

func ProcessRequest(msgStr string, clientInfo string) (*challengeResp.Message, error) {
	msg, err := challengeResp.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}
	switch msg.Header {
	case challengeResp.QUIT:
		fmt.Println(challengeResp.QUIT + "s")
	case challengeResp.REQUEST_CHALLENGE:
		fmt.Println(challengeResp.REQUEST_CHALLENGE)
	}
	return nil, nil
}

// sendMsg - send protocol message to connection
func sendMsg(msg challengeResp.Message, conn net.Conn) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	return err
}

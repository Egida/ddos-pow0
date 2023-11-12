package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	challenge_resp "github.com/dwnGnL/ddos-pow/lib/protocol/challenge-resp"
	"io"
	"log"
	"net"
	"net/http"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

type Handler struct {
	conf *config.Config
}

func newHandler(cfg *config.Config) *Handler {
	return &Handler{conf: cfg}
}

/*func SetupHandlers(core application.Core, cfg *config.Config) GracefulStopFuncWithCtx {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	http.HandleFunc("/ping", application.WithApp2(core, Ping))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start srv: %v", err)
		}
	}()

	return srv.Shutdown
}
*/

func SetupHandlers(core application.Core, cfg *config.Config) GracefulStopFuncWithCtx {
	mux := http.NewServeMux()

	handlerRoutes := application.WithApp2(core, mux)

	handler := newHandler(cfg)

	mux.HandleFunc("/ping", handler.Quit)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Client.Host, cfg.Client.Port),
		Handler: handlerRoutes,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start srv: %v", err)
		}
	}()

	return srv.Shutdown
}

func (h Handler) Quit(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("err"))
		return
	}

	address := fmt.Sprintf("%s:%d", h.conf.Server.Host, h.conf.Server.Port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}

	fmt.Println("connected to", address)
	defer conn.Close()

	reader := bufio.NewReader(conn)

	err = sendMsg(challenge_resp.Message{
		Header:  challenge_resp.QUIT,
		Payload: "",
	}, conn)

	// reading and parsing response
	msgStr, err := readConnMsg(reader)
	if err != nil {
		fmt.Errorf("err read msg: %w", err)
		return
	}

	fmt.Println("msgStr = ", msgStr)

	fmt.Println("fmt ping service", app.GetServer().Ping(), address)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong 2"))
}

// readConnMsg - read string message from connection
func readConnMsg(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMsg - send protocol message to connection
func sendMsg(msg challenge_resp.Message, conn io.Writer) error {
	msgStr := fmt.Sprintf("%s\n", msg.Stringify())
	_, err := conn.Write([]byte(msgStr))
	fmt.Println("msg = ", msgStr)
	return err
}

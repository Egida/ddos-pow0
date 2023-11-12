package client

import (
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/application"
	"log"
	"net/http"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

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

	handler := application.WithApp2(core, mux)

	mux.HandleFunc("/ping", Ping)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start srv: %v", err)
		}
	}()

	return srv.Shutdown
}

func Ping(w http.ResponseWriter, r *http.Request) {
	app, err := application.GetAppFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("err"))
		return
	}

	fmt.Println("fmt ping service", app.GetServer().Ping())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong 2"))
}

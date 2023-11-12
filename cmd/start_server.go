package cmd

import (
	"context"
	"fmt"
	"github.com/dwnGnL/ddos-pow/config"
	"github.com/dwnGnL/ddos-pow/internal/handler/server"
	"time"

	"github.com/dwnGnL/ddos-pow/internal/application"
	"github.com/dwnGnL/ddos-pow/internal/service"
)

const (
	gracefulStopServer = 5 * time.Second
)

func StartServer(cfg *config.Config) error {
	ctx := context.Background()
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	s, err := buildServiceClient(cfg)
	if err != nil {
		return fmt.Errorf("build service err:%w", err)
	}

	fmt.Println(s.GetServer().Ping())

	err = server.SetupHandlers(ctx, s, cfg)

	if err != nil {
		return err
	}

	return nil

	//var group errgroup.Group
	//
	//group.Go(func() error {
	//	sigCh := make(chan os.Signal, 1)
	//	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	//	goerrors.Log().Debug("wait for Ctrl-C")
	//	<-sigCh
	//	goerrors.Log().Debug("Ctrl-C signal")
	//	cancelCtx()
	//	shutdownCtx, shutdownCtxFunc := context.WithDeadline(ctx, time.Now().Add(gracefulStop))
	//	defer shutdownCtxFunc()
	//
	//	_ = httpgrpcGracefulStopWithCtx(shutdownCtx)
	//	return nil
	//})
	//
	//if err := group.Wait(); err != nil {
	//	goerrors.Log().WithError(err).Error("Stopping service with error")
	//}
	//return nil
}

func buildServiceClient(conf *config.Config) (application.Core, error) {
	return service.New(conf), nil
}

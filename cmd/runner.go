package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/tahmooress/discount-manager/configs"
	"github.com/tahmooress/discount-manager/logger"
)

type closerWrapper struct {
	closers []io.Closer
	err     error
}

func (c *closerWrapper) add(cl io.Closer) {
	c.closers = append(c.closers, cl)
}

func (c *closerWrapper) Close() error {
	for _, f := range c.closers {
		err := f.Close()
		if err != nil && c.err == nil {
			c.err = err
		}
	}

	return c.err
}

func Runner(ctx context.Context) (closer io.Closer, err error) {
	cfg := configs.Load()

	log, err := logger.New(logger.Config{
		LogFilePath: cfg.LogFilePath,
		LogLevel:    cfg.LogLevel,
	})
	if err != nil {
		panic(err)
	}

	c := new(closerWrapper)

	c.add(log)

	defer func() {
		if err != nil {
			_ = closer.Close()
		}
	}()

	return c, nil
}

func InterruptHook(cancelFunc context.CancelFunc, closer io.Closer) int {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	defer close(c)

	var existStatus int

	select {
	case <-c:
		err := closer.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "terminate signal received -> closer error: %s\n", err)

			existStatus = 1
		} else {
			fmt.Fprintf(os.Stdout, "terminate signal received -> shutdowned cleanly")
		}

		break
	}

	return existStatus
}

package main

import (
	"context"
	"os"

	"github.com/tahmooress/discount-manager/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	closer, err := cmd.Runner(ctx)
	if err != nil {
		panic(err)
	}

	os.Exit(cmd.InterruptHook(cancel, closer))
}

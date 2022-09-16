package main

import (
	"os"

	"github.com/tahmooress/discount-manager/cmd"
)

func main() {
	closer, errChan, err := cmd.FullNodeRunner()
	if err != nil {
		panic(err)
	}

	os.Exit(cmd.Shutdown(errChan, closer))
}

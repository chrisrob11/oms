package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrisrob11/oms/internal/oms"
)

func main() {
	s, err := oms.NewServer()
	if err != nil {
		fmt.Printf("Error occurred starting server: %s", err)
		os.Exit(1)
	}

	ctx := context.Background()

	err = s.Run(ctx)
	if err != nil {
		fmt.Printf("Error occurred running server: %s", err)
		os.Exit(1)
	}

	fmt.Println("Exited")
}

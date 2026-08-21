package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "rbac",
		Usage: "check whether a role is authorized for a scope",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "role", Required: true},
			&cli.StringFlag{Name: "scope", Required: true},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			allowed, err := Authorize(ctx, command.String("role"), command.String("scope"))
			if err != nil {
				return err
			}

			fmt.Println(allowed)
			if !allowed {
				os.Exit(1)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

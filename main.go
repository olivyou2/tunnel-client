package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func createUuid() string {
	uuid := uuid.New()
	return uuid.String()
}

var alias = createUuid()

func main() {
	(&cli.App{
		Name: "molar",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "target",
				Value:       targetHost,
				Usage:       "A target host for tunneling",
				Destination: &targetHost,
			},
			&cli.StringFlag{
				Name:        "tunnel",
				Value:       tunnelHost,
				Usage:       "A tunnel host for serving",
				Destination: &tunnelHost,
			},
			&cli.IntFlag{
				Name:        "logLevel",
				Value:       logLevel,
				Usage:       "Log warning level",
				Destination: &logLevel,
			},
			&cli.BoolFlag{
				Name:        "verbal",
				Value:       verbal,
				Usage:       "Use verbal",
				Destination: &verbal,
			},
			&cli.StringFlag{
				Name:        "hostPort",
				Value:       hostPort,
				Usage:       "Hosting port that you want to host",
				Destination: &hostPort,
			},
		},
		Usage: "a tcp tunnel fits perfectly with you",
		Action: func(context *cli.Context) error {
			tunnelOk := checkHostAvailable(tunnelHost)

			if !tunnelOk {
				fmt.Println("Tunnel is not available.")
				return nil
			}

			targetOk := checkHostAvailable(targetHost)

			if !targetOk {
				fmt.Println("Target is not available.")
				return nil
			}

			tunneling()
			return nil
		},
	}).Run(os.Args)
}

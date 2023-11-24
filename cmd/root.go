package cmd

import (
	"fmt"
	"os"

	"github.com/nyambati/asqe/internal/metrics"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const bannerMsg = `
      _    ____   ___  _____
     / \  / ___| / _ \| ____|
    / _ \ \___ \| | | |  _|
   / ___ \ ___) | |_| | |___
  /_/   \_\____/ \__\_\_____|
`

type Config struct {
	Metrics []struct {
		Service string
		Regions []string
	}
}

func RootCmd() {
	app := cli.NewApp()
	app.Name = "Prometheus service limits exporter"
	app.Usage = "This app exports aws services limits and usage as metric as prometheus metrics"
	app.Description = "This app exports aws services limits and usage as metric as prometheus metrics"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "listen-address",
			Value:  "0.0.0.0",
			Usage:  "address to listen on for metrics",
			EnvVar: "ASQE_LISTEN_ADDRESS",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  9090,
			Usage:  "port to listen on for metrics",
			EnvVar: "ASQE_PORT",
		},
		cli.StringFlag{
			Name:   "config-path",
			Value:  "$HOME/.asqe",
			Usage:  "configuration on how limits should be exported",
			EnvVar: "ASQE_CONFIG_PATH",
		},
		cli.StringFlag{
			Name:   "region",
			Value:  "eu-east-1",
			Usage:  "AWS region",
			EnvVar: "AWS_REGION",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Print(bannerMsg)
		fmt.Println()
		log.Info(fmt.Sprintf("Prometheus exporter starting to listen on %s:%d/metrics", c.String("listen-address"), c.Int("port")))
		metrics.Serve(c)
		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

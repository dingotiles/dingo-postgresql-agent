package main

import (
	"math/rand"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dingotiles/dingo-postgresql-agent/clicmd"
)

func main() {
	rand.Seed(5000)

	app := cli.NewApp()
	app.Name = "dingo-postgresql-agent"
	app.Version = "0.1.0"
	app.Usage = "Agent to configure Dingo PostgreSQL node"
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "run the agent",
			Action: clicmd.RunAgent,
		},
		{
			Name:   "test-api",
			Usage:  "run the sample/test API backend",
			Action: clicmd.RunTestAPI,
		},
	}
	app.Run(os.Args)
}

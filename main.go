package main

import (
	"github.com/codegangsta/cli"
	"os"
	"runtime"
)

func main() {
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu)

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:   "get",
			Usage:  "get objects",
			Action: GetAction,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:   "parallel,p",
					Value:  10,
					EnvVar: "PARALLEL",
				},
			},
		},
	}
	app.Run(os.Args)
}

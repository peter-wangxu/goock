/*
Copyright 2017 The Goock Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/peter-wangxu/goock/client"
	"github.com/urfave/cli"
)

const (
	COMMAND_NAME = "goock"
	VERSION      = `
   ______                 __      ____   ___ ____
  / ____/___  ____  _____/ /__   / __ \ <  // __ \
 / / __/ __ \/ __ \/ ___/ //_/  / / / / / // / / /
/ /_/ / /_/ / /_/ / /__/ ,<    / /_/ / / // /_/ /
\____/\____/\____/\___/_/|_|   \____(_)_(_)____/

   v0.1.0
`
	USAGE = `A easy-to-use block device management tool.`
)

type App struct {
	*cli.App
}

func NewApp() *App {
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = COMMAND_NAME
	app.Usage = USAGE
	// Global switch/flag
	var enableDebug = false
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "enable debug log on the console.",
			Destination: &enableDebug,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "connect",
			Aliases: []string{"c"},
			Usage:   "Connect to a iSCSI or FC device.",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return client.HandleISCSIConnect(c.Args()...)
			},
		},
		{
			Name:    "disconnect",
			Aliases: []string{"d", "clean"},
			Usage:   "Disconnect(cleanup) a device from host.",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return nil
			},
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Query information for the specified target IP or LUNs.",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return nil
			},
		},
	}
	return &App{app}
}

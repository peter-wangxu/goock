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
	"github.com/peter-wangxu/goock/pkg/client"
	"github.com/urfave/cli"
)

const (
	// CommandName sets the cli entry
	CommandName = "goock"
	// Version is the string of cli version
	// Generated by http://patorjk.com/software/taag/#p=display&f=Slant&t=Goock
	// need to be updated along with version update.
	Version = `
    ______                 __      ____   ___ ___ 
    / ____/___  ____  _____/ /__   / __ \ <  /|__ \
   / / __/ __ \/ __ \/ ___/ //_/  / / / / / / __/ /
  / /_/ / /_/ / /_/ / /__/ ,<    / /_/ / / / / __/ 
  \____/\____/\____/\___/_/|_|   \____(_)_(_)____/ 
												   
  
   v0.1.2
`
	// Usage specifies the simple usage
	Usage = `A easy-to-use block device management tool.`
)

// App is the super struct of urfave app
type App struct {
	*cli.App
}

// NewApp initiates the CLI entry.
func NewApp() *App {
	app := cli.NewApp()
	app.Version = Version
	app.Name = CommandName
	app.Usage = Usage
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
				return client.HandleConnect(c.Args()...)
			},
			ArgsUsage: `<target ip>|<wwn> <lun id>`,
			Description: `# Connect a device via iSCSI IP and LUN ID
   goock disconnect 192.168.1.200 25
   # Connect a device via WWn and LUN ID
   goock disconnect 5006016d09200925 25
`,
		},
		{
			Name:    "disconnect",
			Aliases: []string{"d", "clean"},
			Usage:   "Disconnect(cleanup) a device from host.",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return client.HandleISCSIDisconnect(c.Args()...)
			},
			ArgsUsage: `[<device path|device name>|<target ip|wwn> <lun id>]`,
			Description: `# Disconnect a device via local device path
   goock disconnect /dev/sdb
   # Disconnect a device via device alias
   goock disconnect sdb
   Disconnect a device via iSCSI IP and LUN ID
   goock disconnect 192.168.1.200 25
   # Disconnect a device via WWn and LUN ID
   goock disconnect 5006016d09200925 25
`,
		},
		{
			Name:    "extend",
			Aliases: []string{"e", "expand"},
			Usage:   "Extend a device after been extended on storage side.",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return client.HandleExtend(c.Args()...)
			},
			Description: `# Extend a device via local device path
   goock extend /dev/sdb
   # Extend a device via device alias
   goock extend sdb
   # Extend a device via iSCSI IP and LUN ID
   goock extend 192.168.1.200 25
   # Extend a device via WWn and LUN ID
   goock extend 5006016d09200925 25
`,
		},
		{
			Name:    "info",
			Aliases: []string{"i"},
			Usage:   "Query information for host or LUNs",
			Action: func(c *cli.Context) error {
				// Enable debug log from console
				client.InitLog(enableDebug)
				return client.HandleInfo(c.Args()...)
			},
			Description: `# Query host information about iSCSI or FC
   goock info
   # Query LUN information by iSCSI
   goock info lun 192.168.1.200 25
   # Query LUN information by FC
   goock info lun 5006016d09200925 25
`,
		},
	}
	return &App{app}
}

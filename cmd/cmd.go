package cmd

import (
	"go.zoe.im/x/cli"
)

var (
	// root command to contains all sub commands
	cmd = cli.New(
		// set name and description in run function
		cli.Name("proxman"),
		cli.Short("Proxy server."),
		cli.Description(`Proxy man automatically.`),
		cli.Version(Version),
		cli.Run(func(c *cli.Command, args ...string) {
			c.Help()
		}),
	)
)

// Register sub command
func Register(scs ...*cli.Command) {
	cmd.Register(scs...)
}

// Run call the global's command run
func Run(opts ...cli.Option) error {
	return cmd.Run(opts...)
}

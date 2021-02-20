package manual

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/bzimmer/gravl/pkg/commands/encoding"
)

/*
### Create an exchange

**Syntax:**

```
$ buneary create exchange <ADDRESS> <NAME> <TYPE> [flags]
```

**Arguments:**

|Argument|Description|
|-|-|
|`ADDRESS`|The RabbitMQ HTTP API address. If no port is specified, `15672` is used.|
|`NAME`|The desired name of the new exchange.|
|`TYPE`|The exchange type. Has to be one of `direct`, `headers`, `fanout` and `topic`.|

**Flags:**

|Flag|Short|Description|
|-|-|-|
|`--user`|`-u`|The username to connect with. If not specified, you will be asked for it.|
|`--password`|`-p`|The password to authenticate with. If not specified, you will be asked for it.|
|`--auto-delete`||Automatically delete the exchange once there are no bindings left.|
|`--durable`||Make the exchange persistent, surviving server restarts.|
|`--internal`||Make the exchange internal.|

**Example:**

Create a direct exchange called `my-exchange` on a RabbitMQ server running on the local machine.

```
$ buneary create exchange localhost my-exchange direct
```
*/

func manual(cmds []*cli.Command, lineage []*cli.Command) {
	for i := range cmds {
		fmt.Printf("\n### %s\n\n**Syntax:**\n\n", cmds[i].Usage)
		fmt.Printf("```sh\n$ gravl ")
		for j := range lineage {
			fmt.Printf("%s ", lineage[j].Name)
		}
		fmt.Printf("%s\n```", cmds[i].Name)
		fmt.Println()
		if cmds[i].Action != nil && cmds[i].UsageText != "" {
			fmt.Println("\n**Example:**")
			fmt.Println(cmds[i].UsageText)
		}
		manual(cmds[i].Subcommands, append(lineage, cmds[i]))
	}
}

var Manual = &cli.Command{
	Name:    "manual",
	Usage:   "Generate the 'gravl' manual",
	Aliases: []string{"md"},
	Hidden:  true,
	Action: func(c *cli.Context) error {
		fmt.Printf("# %s - %s\n", c.App.Name, c.App.Usage)
		manual(c.App.Commands, nil)
		fmt.Println()
		return nil
	},
}

var Commands = &cli.Command{
	Name:   "commands",
	Usage:  "Return all possible commands",
	Hidden: true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "relative",
			Aliases: []string{"r"},
			Usage:   "Specify the command relative to the current working directory",
		},
	},
	Action: func(c *cli.Context) error {
		var commander func(string, []*cli.Command) []string
		commander = func(prefix string, cmds []*cli.Command) []string {
			var commands []string
			for i := range cmds {
				cmd := fmt.Sprintf("%s %s", prefix, cmds[i].Name)
				if !cmds[i].Hidden && cmds[i].Action != nil {
					commands = append(commands, cmd)
				}
				commands = append(commands, commander(cmd, cmds[i].Subcommands)...)
			}
			return commands
		}
		cmd := c.App.Name
		if c.Bool("relative") {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			cmd, err = os.Executable()
			if err != nil {
				return err
			}
			cmd, err = filepath.Rel(cwd, cmd)
			if err != nil {
				return err
			}
		}
		commands := commander(cmd, c.App.Commands)
		return encoding.Encode(commands)
	},
}

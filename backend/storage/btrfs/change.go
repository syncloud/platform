package btrfs

import (
	"fmt"
	"strings"
)

type Change struct {
	Cmd  string
	Args []string
}

func NewChange(cmd string, args ...string) Change {
	return Change{Cmd: cmd, Args: args}
}

func (c *Change) Append(args ...string) {
	for _, arg := range args {
		c.Args = append(c.Args, arg)
	}
}

func (c *Change) ToString() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}

package cli

import (
	"errors"
)

type Cli struct {
	length int
	args   []string
}

func (c *Cli) StartCli() error {
	if c.length == 0 {
		return errors.New("Illegal call")
	}
	switch c.args[0] {
	case "-v":
	case "erase":
	}
	return nil
}

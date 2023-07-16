package cli

import (
	"fmt"

	"github.com/TwiN/go-color"
)

func (c *Cli) PrintSuccess(msg string) {
	fmt.Println(color.InGreen(msg))
}

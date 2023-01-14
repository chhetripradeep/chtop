package utils

import (
	"fmt"
	"golang.org/x/image/colornames"
)

func ColorNameToHex(color string) string {
	c, ok := colornames.Map[color]
	if !ok {
		return ""
	}
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

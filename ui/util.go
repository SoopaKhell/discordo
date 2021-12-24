package ui

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

func channelToString(c discord.Channel) string {
	var repr string
	if c.Name != "" {
		repr = "#" + c.Name
	} else if len(c.DMRecipients) == 1 {
		rp := c.DMRecipients[0]
		repr = rp.Username + "#" + rp.Discriminator
	} else {
		rps := make([]string, len(c.DMRecipients))
		for i, r := range c.DMRecipients {
			rps[i] = r.Username + "#" + r.Discriminator
		}

		repr = strings.Join(rps, ", ")
	}

	return repr
}

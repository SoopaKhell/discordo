package ui

import (
	"sort"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/rivo/tview"
)

func onGuildsListSelected(app *App, guildIdx int) {
	rootTreeNode := app.ChannelsTreeView.GetRoot()
	rootTreeNode.ClearChildren()
	app.MessagesTextView.
		Clear().
		SetTitle("")
	app.SetFocus(app.ChannelsTreeView)

	if guildIdx == 0 { // Direct Messages
		cs, err := app.State.PrivateChannels()
		if err != nil {
			return
		}
		sort.Slice(cs, func(i, j int) bool {
			return cs[i].LastMessageID > cs[j].LastMessageID
		})

		for _, c := range cs {
			channelTreeNode := tview.NewTreeNode(channelToString(c)).
				SetReference(c.ID)
			rootTreeNode.AddChild(channelTreeNode)
		}
	} else { // Guild
		// The guilds list includes an "extra" item (direct messages); therefore, the guild items start from the index 1 instead of 0.
		guild := app.State.Ready().Guilds[guildIdx-1]
		sort.Slice(guild.Channels, func(i, j int) bool {
			return guild.Channels[i].Position < guild.Channels[j].Position
		})

		for _, c := range guild.Channels {
			if (c.Type == discord.GuildText || c.Type == discord.GuildNews) && (!c.ParentID.IsValid()) {
				channelTreeNode := tview.NewTreeNode(channelToString(c)).
					SetReference(c.ID)
				rootTreeNode.AddChild(channelTreeNode)
			}
		}

	CATEGORY:
		for _, c := range guild.Channels {
			if c.Type == discord.GuildCategory {
				for _, nestedChannel := range guild.Channels {
					if nestedChannel.ParentID == c.ID {
						channelTreeNode := tview.NewTreeNode(c.Name).
							SetReference(c.ID)
						rootTreeNode.AddChild(channelTreeNode)
						continue CATEGORY
					}
				}

				channelTreeNode := tview.NewTreeNode(c.Name).
					SetReference(c.ID)
				rootTreeNode.AddChild(channelTreeNode)
			}
		}

		for _, c := range guild.Channels {
			if (c.Type == discord.GuildText || c.Type == discord.GuildNews) && (c.ParentID.IsValid()) {
				var parentTreeNode *tview.TreeNode
				rootTreeNode.Walk(func(node, _ *tview.TreeNode) bool {
					if node.GetReference() == c.ParentID {
						parentTreeNode = node
						return false
					}

					return true
				})

				if parentTreeNode != nil {
					channelTreeNode := tview.NewTreeNode(channelToString(c)).
						SetReference(c.ID)
					parentTreeNode.AddChild(channelTreeNode)
				}
			}
		}
	}

	app.ChannelsTreeView.SetCurrentNode(rootTreeNode)
}

func onChannelsTreeViewSelected(app *App, node *tview.TreeNode) {
	c, err := app.State.Channel(node.GetReference().(discord.ChannelID))
	if err != nil {
		return
	}

	switch c.Type {
	case discord.GuildCategory:
		node.SetExpanded(!node.IsExpanded())
	case discord.GuildText, discord.GuildNews:
		// TODO
	}
}

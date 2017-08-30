package rolemanager

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/omar-h/snorlax"
	"github.com/omar-h/snorlax/utils"
)

func init() {
	roleHoist := &snorlax.Command{
		Command:    ".rolehoist",
		Alias:      ".rh",
		Desc:       "Role hoist changes whether or not to display a role sperately.",
		ModuleName: moduleName,
		Handler:    roleHoistHandler,
	}

	roleColor := &snorlax.Command{
		Command:    ".rolecolor",
		Alias:      ".rc",
		Desc:       "Role color changes the color of a specified role.",
		ModuleName: moduleName,
		Handler:    roleColorHandler,
	}

	commands[roleHoist.Command] = roleHoist
	commands[roleColor.Command] = roleColor
}

func roleHoistHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		msg := m.Content
		msgParts := strings.Split(msg, " ")

		msgRoleName, parts := utils.GetStringFromParts(msgParts)
		if msgRoleName == "" || len(parts) != 2 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", msgParts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		roles, err := s.Session.GuildRoles(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		var role *discordgo.Role
		for _, checkRole := range roles {
			if !exists && strings.ToLower(checkRole.Name) == strings.ToLower(msgRoleName) {
				exists = true
				role = checkRole
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" does not exist!")
			return
		}

		hoist, err := strconv.ParseBool(parts[1])
		if err != nil {
			s.Session.ChannelMessageSend(m.ChannelID, "Hoist value isn't valid.")
			s.Log.WithField("error", err).Debug("Error parsing hoist value.")
			return
		}

		_, err = s.Session.GuildRoleEdit(channel.GuildID, role.ID, role.Name, role.Color, hoist, role.Permissions, role.Mentionable)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error editing guild role.")
			return
		}

		s.Session.ChannelMessageSend(m.ChannelID, "Hoist value for "+msgRoleName+" set to "+strconv.FormatBool(hoist)+".")
	}
}

func roleColorHandler(s *snorlax.Snorlax, m *discordgo.MessageCreate) {
	permissions, err := s.Session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		s.Log.WithField("error", err).Debug("Error getting user permissions.")
		return
	}

	if permissions&discordgo.PermissionManageRoles != 0 {
		// Get the message content and split it into arguments
		msg := m.Content
		msgParts := strings.Split(msg, " ")

		msgRoleName, parts := utils.GetStringFromParts(msgParts)
		if msgRoleName == "" || len(parts) != 2 {
			s.Log.Debug(fmt.Sprintf("Wrong number of args: %v", msgParts))
			return
		}

		channel, err := s.Session.Channel(m.ChannelID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting channel.")
			return
		}

		roles, err := s.Session.GuildRoles(channel.GuildID)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error getting guild roles.")
			return
		}

		exists := false
		var role *discordgo.Role
		for _, checkRole := range roles {
			if !exists && strings.ToLower(checkRole.Name) == strings.ToLower(msgRoleName) {
				exists = true
				role = checkRole
			}
		}

		if !exists {
			s.Session.ChannelMessageSend(m.ChannelID, "Role "+msgRoleName+" does not exist!")
			return
		}

		colorIsValid, err := regexp.MatchString("^([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", parts[1])
		if err != nil {
			s.Log.WithField("error", err).Debug("Error running regex on colour string.")
			return
		}

		if !colorIsValid {
			s.Session.ChannelMessageSend(m.ChannelID, "Colour isn't valid.")
			return
		}

		color, err := strconv.ParseInt(parts[1], 16, 32)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error parsing colour value.")
			return
		}

		_, err = s.Session.GuildRoleEdit(channel.GuildID, role.ID, role.Name, int(color), role.Hoist, role.Permissions, role.Mentionable)
		if err != nil {
			s.Log.WithField("error", err).Debug("Error editing guild role.")
			return
		}

		s.Session.ChannelMessageSend(m.ChannelID, "msgRoleName"+"'s colour set to "+parts[1]+".")
	}
}

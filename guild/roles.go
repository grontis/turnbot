package guild

// TODO role manager?
// TODO need to propagate error out of this so that we can indicated error in discord
// func createRole(s *discordgo.Session, guildID string, roleName string, permissions int64, color int) {
// 	roleParams := &discordgo.RoleParams{
// 		Name:        roleName,
// 		Permissions: &permissions,
// 		Color:       &color,
// 	}

// 	role, err := s.GuildRoleCreate(guildID, roleParams)
// 	if err != nil {
// 		fmt.Printf("Error creating role: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("Role '%s' created successfully with ID: %s\n", role.Name, role.ID)
// }

// cmdManager.RegisterCommand(&commands.Command{
// 	Name:        "createrole",
// 	Description: "Creates a new role in the server",
// 	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 		createRole(s, i.GuildID, "Foo", discordgo.PermissionManageMessages, 0xFF5733)
// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "Role created successfully!", //TODO conditional data based on success/fail
// 			},
// 		})
// 	},
// })

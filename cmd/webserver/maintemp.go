package main

import (
	"encoding/json"
	"fmt"
	"la_discord_bot/internal/discordacc"
)
import "la_discord_bot/internal/discordgo"

type LazyRequestChannelStep struct {
	Step []int
}

type LazyRequestChannel struct {
	Channel map[string]LazyRequestChannelStep
}

func main() {
	fmt.Println("Hello world!")

	//da, err := discordacc.New("NDYxMDc2ODAxNzU0NzU5MTY5.WzHxOw.HZ8lTozTwyIYfzHAHffRfh_ENXI", "Bot", false)
	da, err := discordacc.New("NDc2NjQ3MTc4NTQzMTY5NTM2.YYvnZA.iSUQ6s78zxSb6-mesy-NfaYvdlA", "Bot", false)
	err = da.Connect()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	da.Session.AddHandler(func(s *discordgo.Session, m *discordgo.Event) {
		fmt.Printf("event: %+v\n", m.Type)
	})

	da.Session.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		fmt.Printf("ready user.id: %+v\n", m.User.ID)
		fmt.Printf("ready guilds.id.count: %v\n", m.Guilds[0].MemberCount)

		for _, c := range m.Guilds[0].Channels {

			fmt.Printf("channel: %v %v %v\n", c.Name, c.ID, c.Type)
		}

		for _, m := range m.Guilds[0].Members {
			fmt.Printf("%v %v %v\n", m.Nick, m.User.ID, m.User.Bot)
		}

		//err = dg.RequestGuildMembers("814273262154416178", "", 100, false)

	})

	//time.Sleep(5 * time.Second)
	err = da.ReloadData()
	//sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	//<-sc

	fmt.Printf("%+v", da)
	b, err := json.Marshal(da.Guilds)
	if err != nil {
		fmt.Println("\n\nerror: ", err)
	}
	fmt.Printf("\n\n%s", string(b))
	da.Session.Close()
	return
	//discordgo.
	//dg, err := discordgo.New("Bot " + "NDYxMDc2ODAxNzU0NzU5MTY5.WzHxOw.HZ8lTozTwyIYfzHAHffRfh_ENXI")
	dg, err := discordgo.New("Bot " + "NDc2NjQ3MTc4NTQzMTY5NTM2.YYvnZA.iSUQ6s78zxSb6-mesy-NfaYvdlA")

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.IntentsAll // discordgo.IntentsDirectMessages | discordgo.IntentsGuildPresences
	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.Event) {
		fmt.Printf("event: %+v\n", m.Type)
	})

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		fmt.Printf("ready user.id: %+v\n", m.User.ID)
		fmt.Printf("ready guilds.id.count: %v\n", m.Guilds[0].MemberCount)

		for _, c := range m.Guilds[0].Channels {

			fmt.Printf("channel: %v %v %v\n", c.Name, c.ID, c.Type)
		}

		for _, m := range m.Guilds[0].Members {
			fmt.Printf("%v %v %v\n", m.Nick, m.User.ID, m.User.Bot)
		}

		//err = dg.RequestGuildMembers("814273262154416178", "", 100, false)

	})

	//dg.Gateway()

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	err = dg.RequestLazyGuildMembers("814273262154416178", "814278005782347786", [][]int{{0, 99}}, true, false, true, []int{})
	if err != nil {
		fmt.Println("RequestLazyGuildMembers error: ", err)
	}

	//data := LazyRequestData{
	//	GuildID:    "814273262154416178",
	//	Typing:     true,
	//	Threads:    false,
	//	Activities: true,
	//	//Members: make(int, 0),
	//	//Channels: []string{814278005782347786},
	//	//Channels: map[string]string{"814278005782347786": []int{0, 99}},
	//	Channels: map[string][][]int{"814278005782347786": [][]int{{0, 99}, {100, 199}}},
	//}
	//fmt.Printf("%+v\n", map[string][][]int{"814278005782347786": [][]int{{0, 99}, {100, 199}}})
	//fmt.Printf("%+v\n", data)

	//body, err := json.Marshal(data)
	//if err != nil {
	//	fmt.Println("errrrr")
	//} else {
	//	fmt.Printf("%s\n", string(body))
	//}
	//data.Channels["814278005782347786"] = [1][1][]

	//dg.RLock()
	//defer dg.RUnlock()
	//if dg.wsConn == nil {
	//	return ErrWSNotFound
	//}
	//
	//s.wsMutex.Lock()
	//err = s.wsConn.WriteJSON(requestGuildMembersOp{8, data})
	//s.wsMutex.Unlock()

	err = dg.RequestGuildMembers("814273262154416178", "", 100, false)

	if err != nil {
		fmt.Println("RequestGuildMembers error: ", err)
	}

	//dg.Gateway()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	fmt.Printf("%+v\n", dg.State.Guilds)

	for _, g := range dg.State.Guilds {
		fmt.Println(g.ID)
		users, err := dg.GuildMembers("814273262154416178", "", 1000)
		if err != nil {
			fmt.Println("error get members,", err)
		} else {
			for _, u := range users {
				fmt.Printf("%+v\n", u)
			}
		}
	}

	//user, err := dg.User("@me")

	//if err != nil {
	//	fmt.Println("error get user @me,", err)
	//} else {
	//	fmt.Printf("%+v", user)
	//}

	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println("error get guilds,", err)
	} else {
		fmt.Printf("%+v\n", guilds)
		for _, id := range guilds {
			fmt.Println(id.ID)
			users, err := dg.GuildMembers(id.ID, "", 1000)
			if err != nil {
				fmt.Println("error get members,", err)
			} else {
				for _, u := range users {
					fmt.Printf("%+v\n", u)
				}
			}
		}
	}

	//guildJson, err := dg.RequestWithBucketID("GET", discordgo.EndpointGuild(strconv.Itoa(guild.ExternalID))+"?with_counts=true", nil, discordgo.EndpointGuild(strconv.Itoa(guild.ExternalID)))
	//guildDiscord := &discordgo.Guild{}
	//
	//if err == nil {
	//	err = json.Unmarshal(guildJson, guildDiscord)
	//}

	//guildDiscord.ApproximateMemberCount

	//sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	//<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	//dg.UserGuilds(dg.User)
	//dg.User()
	//appserver.Run(configsPath)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
//
// It is called whenever a message is created but only when it's sent through a
// server as we did not request IntentsDirectMessages.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("message: " + m.Content)
	fmt.Printf("%+v", m.Message)
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// In this example, we only care about messages that are "ping".
	if m.Content != "ping" {
		return
	}

	// We create the private channel with the user who sent the message.
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}
	// Then we send the message through the channel we created.
	_, err = s.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

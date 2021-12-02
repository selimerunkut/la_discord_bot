package task

import (
	"encoding/json"
	"fmt"
	"la_discord_bot/internal/discordacc"
	"la_discord_bot/internal/discordgo"
	"la_discord_bot/internal/guild"
	"la_discord_bot/internal/helpers"
	"math"
	"time"
)

func (T *Task) ParseMembers() (err error) {
	T.Log.Println("Task Parse Members Starting")
	T.Da, err = discordacc.New(T.Token, "Bot", true)
	defer T.Da.Close()

	if err != nil {
		//T.Log.Printf("ERROR Task Parse Members Starting: " + fmt.Sprint( err ))
		T.SetError(err)
		return err
	}
	T.Log.Println("Logged in. ID: " + T.Da.User.ID)
	T.UserId = T.Da.User.ID
	T.UserName = T.Da.User.Username
	T.GuildMemberCount = T.Da.Guilds[T.GuildId].MemberCount
	T.GuildName = T.Da.Guilds[T.GuildId].Name

	if T.Da.User.Bot {
		T.Log.Println("Bot = true ")
		//T.Log.Printf("%+v", T.Da.Guilds[T.GuildId])
		//T.Log.Printf("%+v", T.Da.Guilds[T.GuildId].MemberCount)

		if T.GuildMemberCount > 0 {
			T.Steps = int(math.Round(float64(T.Da.Guilds[T.GuildId].MemberCount)/ParseBotLimit)) + 1
		}

		T.Log.Printf("Task Steps: %v", T.Steps)
		T.Save()
		err = T.ParseMembersBot()
		if err != nil {
			T.SetError(err)
		}
	} else {
		T.Log.Println("Bot = false ")
		if T.GuildMemberCount > 0 {
			T.Steps = int(math.Round(float64(T.Da.Guilds[T.GuildId].MemberCount)/100)) + 1
		}

		T.Log.Printf("Task Steps: %v", T.Steps)
		T.Save()
		err = T.ParseMembersUser()
		if err != nil {
			T.SetError(err)
		}

	}
	return nil
}

func (T *Task) ParseMembersAll() (err error) {
	T.Log.Println("Parsing All members in the guild")
	T.Da, err = discordacc.New(T.Token, "Bot", true)
	defer T.Da.Close()

	if err != nil {
		T.Log.Printf("ERROR Task Parse All Members Starting: " + fmt.Sprint(err))
		T.SetError(err)
		return err
	}
	T.Log.Println("Logged in. ID: " + T.Da.User.ID)
	T.UserId = T.Da.User.ID
	T.GuildMemberCount = T.Da.Guilds[T.GuildId].MemberCount

	if !T.Da.User.Bot {
		if T.GuildMemberCount > 0 {
			T.Steps = int(math.Round(float64(T.GuildMemberCount)/100)) + 1
		}
		fmt.Println(T.Da.Guilds[T.GuildId].MemberCount)
	}
	return nil
}

func (T *Task) ParseMembersUser() (err error) {

	T.CurrentStep = 0
	guildStore, err := guild.NewStore(T.GuildId, T.GuildName, T.Config)
	err = guildStore.Save()
	if err != nil {
		return err
	}
	err = guildStore.LoadMembers()

	T.Da.Session.AddHandler(func(s *discordgo.Session, m *discordgo.Event) {
		//T.Log.Printf("event: %+v\n", m.Type)
		if m.Type == "GUILD_MEMBER_LIST_UPDATE" {
			var i GuildMemberListUpdateEventAuto
			if err = json.Unmarshal(m.RawData, &i); err != nil {
				T.Log.Println(err)
			}
			memCount := 0
			for _, op := range i.Ops {
				for _, item := range op.Items {
					if item.Member.User.ID != "" {
						//T.ParseMemberIds[item.Member.User.ID] = item.Member.User.ID
						guildStore.MembersIds[item.Member.User.ID] = guild.Member{
							MemberId: item.Member.User.ID,
						}
						T.Log.Printf("%v %v", item.Member.User.ID, item.Member.User.Username)
						memCount++
						T.ParseCount++ // = len(T.ParseMemberIds)
					}
				}
			}

			fmt.Println(memCount)
			if len(i.Ops[0].Range) > 0 {
				T.Log.Printf("Range %+v | memCount %v", i.Ops[0].Range, memCount)
			}

			// save members
			if memCount > 0 {
				//err = T.SaveGuildMembers(T.ParseMemberIds, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
				//err = T.SaveGuildMembers(, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
				err = guildStore.SaveMembers()
				if err != nil {
					T.Log.Println(err)
				}
			}
		}
	})

	err = T.Da.Session.RequestLazyGuildMembers(T.GuildId, T.ChannelId, [][]int{{0, 99}}, true, false, true, []int{})
	if err != nil {
		T.Log.Println("RequestLazyGuildMembers error: ", err)
		return err
	}
	helpers.Sleep(10)
	j := T.CurrentStep
	for {
		if T.Status != StatusWorking {
			return nil
		}
		if j > T.Steps {
			break
		}
		err = T.Da.Session.RequestLazyGuildMembers(T.GuildId, T.ChannelId, [][]int{{0, 99}, {100 * j, 100*j + 99}, {100 * (j + 1), 100*(j+1) + 99}}, false, false, false, []int{})
		T.CurrentStep = j
		T.Log.Printf("Step %v / %v", j, T.Steps)
		err = T.Save()
		if err != nil {
			T.Log.Println(err)
			return err
		}
		helpers.Sleep(1)
		j++
	}
	helpers.Sleep(5)

	err = guildStore.SaveMembers()
	if err != nil {
		T.Log.Println(err)
		return err
	}

	err = guildStore.Save()
	if err != nil {
		return err
	}
	return nil
}

func (T *Task) ParseMembersBot() (err error) {
	guildStore, err := guild.NewStore(T.GuildId, T.GuildName, T.Config)
	err = guildStore.Save()
	if err != nil {
		T.Log.Println(err)
		return err
	}

	after := ""
	for {
		if T.Status != StatusWorking {
			return nil
		}
		err = T.Save()
		if err != nil {
			T.Log.Println(err)
			return err
		}

		T.Log.Printf("Task Step %v / %v", T.CurrentStep, T.Steps)
		members, err := T.Da.Session.GuildMembers(T.GuildId, after, ParseBotLimit)
		if err != nil {
			T.SetError(err)
			return err
		}
		T.Log.Printf("Got %v Members", len(members))

		for _, m := range members {
			if !m.User.Bot {
				//memberIDs[m.User.ID] = m.User.ID
				guildStore.MembersIds[m.User.ID] = guild.Member{
					MemberId: m.User.ID,
				}
				T.ParseCount = len(guildStore.MembersIds)
			}
		}

		// save members
		err = guildStore.SaveMembers()
		if err != nil {
			T.Log.Println(err)
			return err
		}

		after = members[len(members)-1].User.ID

		T.Log.Printf("After = %v", after)
		if len(members) < ParseBotLimit {
			break
		}
		//memberIDs = make(map[string]string)
		T.CurrentStep++
		//discordacc.Sleep(5)
	}

	err = guildStore.Save()
	if err != nil {
		return err
	}

	return nil
}

type GuildMemberListUpdateEventAuto struct {
	Groups []struct {
		Count int    `json:"count"`
		ID    string `json:"id"`
	} `json:"groups"`
	GuildID     string `json:"guild_id"`
	ID          string `json:"id"`
	MemberCount int    `json:"member_count"`
	OnlineCount int    `json:"online_count"`
	Ops         []struct {
		Items []struct {
			Group struct {
				Count int    `json:"count"`
				ID    string `json:"id"`
			} `json:"group,omitempty"`
			Member struct {
				Deaf        bool      `json:"deaf"`
				HoistedRole string    `json:"hoisted_role"`
				JoinedAt    time.Time `json:"joined_at"`
				Mute        bool      `json:"mute"`
				Presence    struct {
					Activities []struct {
						CreatedAt int64 `json:"created_at"`
						Emoji     struct {
							Name string `json:"name"`
						} `json:"emoji,omitempty"`
						ID         string `json:"id"`
						Name       string `json:"name"`
						State      string `json:"state,omitempty"`
						Type       int    `json:"type"`
						Timestamps struct {
							Start int64 `json:"start"`
						} `json:"timestamps,omitempty"`
					} `json:"activities"`
					ClientStatus struct {
						Desktop string `json:"desktop"`
					} `json:"client_status"`
					Status string `json:"status"`
					User   struct {
						ID string `json:"id"`
					} `json:"user"`
				} `json:"presence"`
				Roles []string `json:"roles"`
				User  struct {
					Avatar        string `json:"avatar"`
					Discriminator string `json:"discriminator"`
					ID            string `json:"id"`
					PublicFlags   int    `json:"public_flags"`
					Username      string `json:"username"`
				} `json:"user"`
			} `json:"member,omitempty"`
		} `json:"items"`
		Op    string `json:"op"`
		Range []int  `json:"range"`
	} `json:"ops"`
}

package task

import (
	"la_discord_bot/internal/discordacc"
	"la_discord_bot/internal/guild"
	"la_discord_bot/internal/helpers"
	"strings"
)

func (T *Task) SendMessages() (err error) {
	// TODO join by invite

	T.Log.Println("Task Parse Members Starting")
	T.Da, err = discordacc.New(T.Token, "Bot", true)
	defer T.Da.Close()

	if err != nil {
		//T.Log.Printf("ERROR Task Parse Members Starting: " + fmt.Sprint( err ))
		T.SetError(err)
		return err
	}
	T.Log.Println("Logged in. ID: " + T.Da.User.ID)
	if T.Da.User.Bot {
		T.Log.Println("Bot = true ")
	} else {
		T.Log.Println("Bot = false ")
	}

	T.checkDelayValues()
	guildStore, err := guild.Load(T.GuildId, T.Config)
	if err != nil {
		return err
	}
	T.Log.Printf("Load guild %v data", T.GuildId)
	err = guildStore.LoadMembers()
	if err != nil {
		return err
	}
	T.Log.Printf("Load guild %v members. Count: %v", T.GuildId, len(guildStore.MembersIds))
	skip := false
	if T.SendMemberId != "" {
		skip = true
		T.Log.Printf("Need to skip members to ID %v", T.SendMemberId)
	}

	T.CurrentStep = 0
	T.Steps = len(guildStore.MembersIds)

	T.GuildName = guildStore.GuildName
	if err = T.Save(); err != nil {
		return err
	}

	for _, m := range guildStore.MembersIds {
		if T.Status != StatusWorking {
			return nil
		}

		if skip {
			if m.MemberId == T.SendMemberId {
				skip = false
				T.CurrentStep++
				continue
			} else {
				T.CurrentStep++
				continue
			}
		}

		//T.Log.Printf("%v , %v , %v , %v",
		//	int64(T.SendDelayMin*1000),
		//	int64(T.SendDelayMax*1000),
		//	helpers.MtRand(int64(T.SendDelayMin*1000), int64(T.SendDelayMax*1000)),
		//)
		delay := helpers.MtRandFloat(T.SendDelayMin, T.SendDelayMax)
		delay1 := helpers.MtRandFloat(0.1, delay/2)
		delay2 := delay - delay1
		//T.Log.Printf("Delay: %v | Delay1: %v | Delay2: %v ", delay, delay1, delay2)
		//T.Log.Printf(" helpers.SleepFloat(delay1) %v", int64(delay1*1000))
		helpers.SleepFloat(delay1)

		if m.MemberId != "" {
			channel, err := T.Da.Session.UserChannelCreate(m.MemberId)
			if err != nil {
				T.Log.Printf("MemberID %v error creating channel: %v", m.MemberId, err)
				T.SendFail++
				T.SendMemberId = m.MemberId
				T.CurrentStep++
				if err = T.Save(); err != nil {
					return err
				}
				continue
			}

			helpers.SleepFloat(delay2)

			message := T.SendMessage
			if strings.Contains(message, "<user>") {
				ping := "<@" + m.MemberId + ">"
				message = strings.ReplaceAll(message, "<user>", ping)
			}

			_, err = T.Da.Session.ChannelMessageSend(channel.ID, message)
			if err != nil {
				T.Log.Printf("MemberID %v error sending DM message: %v", m.MemberId, err)
				T.SendFail++
				T.SendMemberId = m.MemberId
				T.CurrentStep++
				if err = T.Save(); err != nil {
					return err
				}
				continue
			}

		} else {
			_, err = helpers.UserChannelCreateTest(m.MemberId)

			if err != nil {
				T.Log.Printf("MemberID %v error creating channel: %v", m.MemberId, err)
				T.SendFail++
				T.SendMemberId = m.MemberId
				T.CurrentStep++
				if err = T.Save(); err != nil {
					return err
				}
				continue
			}
			helpers.SleepFloat(delay2)
			_, err = helpers.ChannelMessageSendTest("", T.SendMessage)
			if err != nil {
				T.Log.Printf("MemberID %v error sending DM message: %v", m.MemberId, err)
				T.SendFail++
				T.SendMemberId = m.MemberId
				T.CurrentStep++
				if err = T.Save(); err != nil {
					return err
				}
				continue
			}

		}

		T.Log.Printf("MemberID %v success sending DM message", m.MemberId)
		T.SendSuccess++
		T.CurrentStep++
		T.SendMemberId = m.MemberId
		if err = T.Save(); err != nil {
			return err
		}
	}

	return nil
}

func (T *Task) checkDelayValues() {
	//if f, err := strconv.ParseFloat(T.SendDelayMin, 32); err == nil {
	//fmt.Printf("Delay %v min  < %v . Set delay min to %v",
	//	T.SendDelayMin, T.Config.DelayMin, T.Config.DelayMin)

	if T.SendDelayMin <= 0 {
		T.SendDelayMin = T.Config.DelayMin
	}

	if T.SendDelayMax <= 0 {
		T.SendDelayMax = T.Config.DelayMax
	}

	if T.SendDelayMin < T.Config.DelayMin {
		T.Log.Printf("Delay %v min  < %v . Set delay min to %v",
			T.SendDelayMin, T.Config.DelayMin, T.Config.DelayMin)
		T.SendDelayMin = T.Config.DelayMin
	}

	if T.SendDelayMax > T.Config.DelayMax {
		T.Log.Printf("Delay %v max  < %v . Set delay max to %v",
			T.SendDelayMax, T.Config.DelayMax, T.Config.DelayMax)
		T.SendDelayMax = T.Config.DelayMax
	}
}

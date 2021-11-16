package task

import "la_discord_bot/internal/guild"

func (T *Task) SendMessages() (err error) {
	// load memberId

	return nil

	guildStore, err := guild.Load(T.GuildId, T.Config)
	if err != nil {
		return err
	}
	err = guildStore.LoadMembers()
	if err != nil {
		return err
	}
	skip := false
	if T.SendMemberId != "" {
		skip = true
	}

	for _, m := range guildStore.MembersIds {
		if T.Status != StatusWorking {
			return nil
		}

		if skip {
			if m.MemberId == T.SendMemberId {
				skip = false
				continue
			} else {
				continue
			}
		}

		channel, err := T.Da.Session.UserChannelCreate(m.MemberId)

		if err != nil {
			T.Log.Println("error creating channel:", err)
			T.SendFail++
			T.SendMemberId = m.MemberId
			if err = T.Save(); err != nil {
				return err
			}
			continue
		}

		_, err = T.Da.Session.ChannelMessageSend(channel.ID, T.SendMessage)

		if err != nil {
			T.Log.Println("error sending DM message:", err)
			T.SendFail++
			T.SendMemberId = m.MemberId
			if err = T.Save(); err != nil {
				return err
			}
			continue
		}

		T.SendSuccess++
		T.SendMemberId = m.MemberId
		if err = T.Save(); err != nil {
			return err
		}
		// TODO sleep
		//SleepInterval()
	}

	return nil
}

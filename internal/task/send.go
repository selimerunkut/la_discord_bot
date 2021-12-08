package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"la_discord_bot/internal/discordacc"
	"la_discord_bot/internal/guild"
	"la_discord_bot/internal/helpers"
	"net/http"
	"strings"
)

type jsonDiscordResponse struct {
	StatusCode int    `json:"responseCode"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (T *Task) SendMessages() (err error) {
	cookie := helpers.Cookie{}
	fingerprint := helpers.Fingerprintx{}
	T.Log.Println("Task Send Members Starting")
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
		if cookie, err = helpers.GetCookie(); err != nil {
			//T.Log.Printf("%v", err)
			return err
		}
		if fingerprint, err = helpers.GetFingerprint(); err != nil {
			return err
		}

		if T.SendInvite != "" {
			if err = helpers.JoinGuild(T.SendInvite, T.Token, fingerprint, cookie); err != nil {
				return err
			}
		}
	}

	T.checkDelayValues()
	membersIDs := []guild.Member{}
	if T.GuildId != "" {
		guildStore, err := guild.Load(T.GuildId, T.Config)
		if err != nil {
			return err
		}
		T.Log.Printf("Load guild %v data", T.GuildId)
		//err = guildStore.LoadMembers()
		membersIDs, err = guildStore.LoadMembersSlice()
		T.GuildName = guildStore.GuildName
	}
	if len(T.SendMembersIDS) > 0 {
		membersIDs = append(membersIDs, T.SendMembersIDS...)
		T.GuildName = T.GuildName + " " + membersIDs[0].MemberId
	}
	if err != nil {
		return err
	}
	T.Log.Printf("Load guild %v members. Count: %v", T.GuildId, len(membersIDs))
	skip := false
	if T.SendMemberId != "" {
		skip = true
		T.Log.Printf("Need to skip members to ID %v", T.SendMemberId)
	}

	T.CurrentStep = 0
	T.Steps = len(membersIDs)

	if err = T.Save(); err != nil {
		return err
	}

	for _, m := range membersIDs {
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
		delay := helpers.MtRandFloat(T.SendDelayMin, T.SendDelayMax)

		if T.Da.User.Bot {
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
			//T.Log.Println("%v", m.MemberId)
			//T.Log.Println("%v", fingerprint)
			//T.Log.Println("%v", cookie)
			channelID, err := OpenChannelUser(T.Token, m.MemberId, fingerprint, cookie)
			T.Log.Printf("%v", channelID)

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

			jsonResp, err := SendMessageUser(T.Token, channelID, T.SendMessage, m.MemberId, fingerprint, cookie)
			if err != nil {
				T.Log.Printf("MemberID %v error sending DM message: %v", m.MemberId, err)
				//T.Log.Printf("%v", DiscordResponseToError(jsonResp))
				// TODO need long delay
				if jsonResp.StatusCode == 403 && jsonResp.Code == 40003 {
					longDelay := helpers.MtRand(10*60, 15*60)
					T.Log.Printf("Long delay sleep %v seconds", longDelay)
					helpers.Sleep(longDelay)
				}
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
		helpers.SleepFloat(delay)
	}

	return nil
}

func OpenChannelUser(token string, recepientUID string, fingerprint helpers.Fingerprintx, cookie helpers.Cookie) (channelID string, err error) {
	url := "https://discord.com/api/v9/users/@me/channels"

	jsonData := []byte("{\"recipients\":[\"" + recepientUID + "\"]}")

	if fingerprint.ToString() == "" {
		return channelID, errors.New("empty fingerprint")
	}
	if cookie.ToString() == "" {
		return channelID, errors.New("empty cookie")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error while making request")
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Cookie", cookie.ToString())
	req.Header.Set("x-fingerprint", fingerprint.ToString())
	httpClient := &http.Client{}
	resp, err := httpClient.Do(CommonHeaders(req))

	if err != nil {
		return "", fmt.Errorf("Error while sending Open channel request  %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("invalid Status Code while sending request %v \n", resp.StatusCode)
		return "", err
	}
	type responseBody struct {
		ID string `json:"id,omitempty"`
	}
	//fmt.Printf("%+v", body)
	//fmt.Printf("%+v", string(body))
	var channelSnowflake responseBody
	if err = json.Unmarshal(body, &channelSnowflake); err != nil {
		return
	}

	return channelSnowflake.ID, nil
}

func SendMessageUser(token string, channelSnowflake string, message string, memberId string, fingerprint helpers.Fingerprintx, cookie helpers.Cookie) (response jsonDiscordResponse, err error) {

	if fingerprint.ToString() == "" {
		return response, errors.New("empty fingerprint")
	}
	if cookie.ToString() == "" {
		return response, errors.New("empty cookie")
	}

	if strings.Contains(message, "<user>") {
		ping := "<@" + memberId + ">"
		message = strings.ReplaceAll(message, "<user>", ping)
	}

	body, err := json.Marshal(&map[string]interface{}{
		"content": message,
		"embeds":  false,
		"tts":     false,
		"nonce":   helpers.Snowflake(),
	})

	if err != nil {
		return
	}

	url := "https://discord.com/api/v9/channels/" + channelSnowflake + "/messages"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))

	if err != nil {
		//log.Panicf("Error while making HTTP request")
		return
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("referer", "https://discord.com/channels/@me/"+channelSnowflake)
	req.Header.Set("Cookie", cookie.ToString())
	req.Header.Set("x-fingerprint", fingerprint.ToString())
	httpClient := &http.Client{}
	res, err := httpClient.Do(CommonHeaders(req))

	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &response)
	response.StatusCode = res.StatusCode
	return
}

func DiscordResponseToError(response jsonDiscordResponse) error {
	switch response.StatusCode {
	case 200:
		return nil
	case 403:
		switch response.Code {
		case 40003:
			// need long delay
			return fmt.Errorf("need long delay")
		case 50007:
			return fmt.Errorf("user has either closed DMs or is not in a mutual server or has blocked the token")
		case 40002:
			return fmt.Errorf("token is locked")
		case 50009:
			return fmt.Errorf("can't DM maybe user accept DM's only from friends")
		}
	case 401:
		return fmt.Errorf("token is wrong or disabled")
	case 405:
		return fmt.Errorf("token might be phone locked or disabled or may not have a mutual server")
	default:
		return fmt.Errorf("failed to send DM")
	}

	return nil
}

func CommonHeaders(req *http.Request) *http.Request {

	req.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	req.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-GB")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	return req
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

	if T.SendDelayMax < T.SendDelayMin {
		tt := T.SendDelayMin
		T.SendDelayMin = T.SendDelayMax
		T.SendDelayMax = tt
	}

}

package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"la_discord_bot/internal/discordgo"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Fingerprintx struct {
	Fingerprint string `json:"fingerprint"`
}

func (f *Fingerprintx) ToString() string {
	return f.Fingerprint
}

func GetFingerprint() (fingerprint Fingerprintx, err error) {

	resp, err := http.Get("https://discordapp.com/api/v9/experiments")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var fingerprinty Fingerprintx
	err = json.Unmarshal(body, &fingerprinty)
	if err != nil {
		return
	}
	return fingerprinty, nil

}

type Cookie struct {
	Dcfduid  string `json:"dcfduid"`
	Sdcfduid string `json:"sdcfduid"`
}

func (c *Cookie) ToString() string {
	if c.Dcfduid == "" && c.Sdcfduid == "" {
		return ""
	}
	return "__dcfduid=" + c.Dcfduid + "; " + "__sdcfduid=" + c.Sdcfduid + "; " + " locale=us" + "; __cfruid=d2f75b0a2c63c38e6b3ab5226909e5184b1acb3e-1634536904"
}

func GetCookie() (c Cookie, err error) {
	resp, err := http.Get("https://discord.com")
	if err != nil {
		return Cookie{}, err
	}
	defer resp.Body.Close()

	Cookie := Cookie{}
	if resp.Cookies() != nil {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "__dcfduid" {
				Cookie.Dcfduid = cookie.Value
			}
			if cookie.Name == "__sdcfduid" {
				Cookie.Sdcfduid = cookie.Value
			}
		}
	}
	return Cookie, nil
}

func JoinGuild(inviteCode string, token string, fingerprintx Fingerprintx, cookie Cookie) (err error) {
	url := "https://discord.com/api/v9/invites/" + inviteCode
	if fingerprintx.ToString() == "" {
		fingerprintx, err = GetFingerprint()
		if err != nil {
			return err
		}
	}

	if cookie.ToString() == "" {
		cookie, err = GetCookie()
		if err != nil {
			return err
		}
	}

	if cookie.ToString() == "" {
		return errors.New("empty cookie")
	}

	var headers struct{}
	requestBytes, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		return errors.New("error while creating request")
	}

	req.Header.Set("cookie", cookie.ToString())
	req.Header.Set("authorization", token)

	httpClient := http.Client{}
	resp, err := httpClient.Do(CommonHeaders(req, fingerprintx.ToString()))
	if err != nil {
		return errors.New("error while sending request")
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("unexpected Status code %v while joining token %v", resp.StatusCode, token)

}

func CommonHeaders(req *http.Request, fingerprint string) *http.Request {
	req.Header.Set("accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("sec-ch-ua", "'Chromium';v='92', ' Not A;Brand';v='99', 'Google Chrome';v='92'")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-context-properties", "eyJsb2NhdGlvbiI6IkpvaW4gR3VpbGQiLCJsb2NhdGlvbl9ndWlsZF9pZCI6Ijg4NTkwNzE3MjMwNTgwOTUxOSIsImxvY2F0aW9uX2NoYW5uZWxfaWQiOiI4ODU5MDcxNzIzMDU4MDk1MjUiLCJsb2NhdGlvbl9jaGFubmVsX3R5cGUiOjB9")
	req.Header.Set("x-fingerprint", fingerprint)
	req.Header.Set("x-super-properties", "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiRmlyZWZveCIsImRldmljZSI6IiIsInN5c3RlbV9sb2NhbGUiOiJlbi1VUyIsImJyb3dzZXJfdXNlcl9hZ2VudCI6Ik1vemlsbGEvNS4wIChXaW5kb3dzIE5UIDEwLjA7IFdpbjY0OyB4NjQ7IHJ2OjkzLjApIEdlY2tvLzIwMTAwMTAxIEZpcmVmb3gvOTMuMCIsImJyb3dzZXJfdmVyc2lvbiI6IjkzLjAiLCJvc192ZXJzaW9uIjoiMTAiLCJyZWZlcnJlciI6IiIsInJlZmVycmluZ19kb21haW4iOiIiLCJyZWZlcnJlcl9jdXJyZW50IjoiIiwicmVmZXJyaW5nX2RvbWFpbl9jdXJyZW50IjoiIiwicmVsZWFzZV9jaGFubmVsIjoic3RhYmxlIiwiY2xpZW50X2J1aWxkX251bWJlciI6MTAwODA0LCJjbGllbnRfZXZlbnRfc291cmNlIjpudWxsfQ==")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("origin", "https://discord.com")
	req.Header.Set("referer", "https://discord.com/channels/@me")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) discord/0.0.16 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36")
	req.Header.Set("te", "trailers")
	return req
}

func Snowflake() int64 {
	snowflake := strconv.FormatInt((time.Now().UTC().UnixNano()/1000000)-1420070400000, 2) + "0000000000000000000000"
	nonce, _ := strconv.ParseInt(snowflake, 2, 64)
	return nonce
}

func FilePutContents(filename string, data string, mode os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), mode)
}

func FileGetContents(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	return string(data), err
}

func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func Sleep(seconds int64) {
	time.Sleep(time.Duration(seconds) * time.Second)
	return
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func CountLinesInFile(filename string) (count int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return count, err
	}

	count, err = LineCounter(file)

	if err != nil {
		return count, err
	}

	return count, nil
}

func SleepInterval(min, max float64) {
	minMilliseconds := int64(min * 1000)
	maxMilliseconds := int64(max * 1000)
	duration := MtRand(minMilliseconds, maxMilliseconds)
	sleepTime := time.Duration(duration) * time.Millisecond
	time.Sleep(sleepTime)
	return
}

func SleepFloat(seconds float64) {
	//time.Sleep(time.Duration(int64(seconds*1000)) * time.Millisecond)
	time.Sleep(time.Duration(int64(seconds*1000)) * time.Millisecond)
	return
}

func MtRandFloat(min, max float64) float64 {
	minMilliseconds := int64(min * 1000)
	maxMilliseconds := int64(max * 1000)
	r := MtRand(minMilliseconds, maxMilliseconds)
	return float64(r) / 1000
}

func MtRand(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min+1) + min
}

func UserChannelCreateTest(recipientID string) (st *discordgo.Channel, err error) {
	if MtRand(0, 1000) > 500 {
		return nil, nil
	}
	return nil, fmt.Errorf("test error UserChannelCreateTest")
}

func ChannelMessageSendTest(channelID string, content string) (*discordgo.Message, error) {
	if MtRand(0, 1000) > 500 {
		return nil, nil
	}
	return nil, fmt.Errorf("test error ChannelMessageSendTest")
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

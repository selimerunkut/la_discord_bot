package discordacc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getFingerprint() (fingerprint string, err error) {
	//log.SetOutput(ioutil.Discard)
	resp, err := http.Get("https://discordapp.com/api/v9/experiments")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type Fingerprintx struct {
		Fingerprint string `json:"fingerprint"`
	}
	var fingerprinty Fingerprintx
	err = json.Unmarshal(body, &fingerprinty)
	if err != nil {
		return "", err
	}
	//color.Green("INFO: Obtained Fingerprint: " + fingerprinty.Fingerprint)
	return fingerprinty.Fingerprint, nil

}

type cookie struct {
	Dcfduid  string
	Sdcfduid string
}

func getCookie() (c cookie, err error) {
	//log.SetOutput(ioutil.Discard)
	resp, err := http.Get("https://discord.com")
	if err != nil {
		return cookie{}, err
	}
	defer resp.Body.Close()

	Cookie := cookie{}
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
	//color.Yellow("INFO: Obtained Cookies: " + "__dcfduid= " + Cookie.Dcfduid + " " + "__sdcfduid= " + Cookie.Sdcfduid)
	return Cookie, nil
}

func joinGuild(inviteCode string, token string) (err error) {
	url := "https://discord.com/api/v9/invites/" + inviteCode
	fingerprint, err := getFingerprint()
	if err != nil {
		return err
	}

	Cookie, err := getCookie()
	if err != nil {
		return err
	}

	if Cookie.Dcfduid == "" && Cookie.Sdcfduid == "" {
		return errors.New("empty cookie")
	}

	Cookies := "__dcfduid=" + Cookie.Dcfduid + "; " + "__sdcfduid=" + Cookie.Sdcfduid + "; " + "locale=us"

	var headers struct{}
	requestBytes, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBytes))
	if err != nil {
		return errors.New("error while creating request")
	}

	req.Header.Set("cookie", Cookies)
	req.Header.Set("authorization", token)

	httpClient := http.Client{}
	resp, err := httpClient.Do(commonHeaders(req, fingerprint))
	if err != nil {
		return errors.New("error while sending request")
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("unexpected Status code %v while joining token %v", resp.StatusCode, token)

}

func commonHeaders(req *http.Request, fingerprint string) *http.Request {
	req.Header.Set("accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB")
	req.Header.Set("content-type", "application/json")
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

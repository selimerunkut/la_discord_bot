package guild

import (
	"bufio"
	"encoding/json"
	"fmt"
	"la_discord_bot/internal/config"
	"la_discord_bot/internal/helpers"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Member struct {
	MemberId string `json:"member_id"`
}

type Store struct {
	GuildId      string            `json:"guild_id"`
	GuildName    string            `json:"guild_name"`
	MembersCount int               `json:"members_count"`
	MembersIds   map[string]Member `json:"-"`
	Config       *config.Config    `json:"-"`
}

type Stores struct {
	Stores map[string]Store
}

func LoadAll(c *config.Config) (ss Stores, err error) {
	ss = Stores{
		Stores: map[string]Store{},
	}

	files, err := filepath.Glob(c.PathToStorage + "guilds/*/*")

	if err != nil {
		return ss, err
	}

	for _, f := range files {
		//regexp.Match(`/(\d+)\.members`, []byte(f))
		if strings.Index(f, ".members") == -1 {
			continue
		}
		id := filepath.Base(f)
		id = strings.Replace(id, ".members", "", 1)
		store, err := Load(id, c)
		//log.Println(store)
		if err != nil {
			return ss, err
		}
		ss.Stores[store.GuildId] = store
	}

	return ss, nil
}

func Load(guildId string, config *config.Config) (store Store, err error) {

	filename := config.PathToStorage + "guilds/" + guildId + "/" + guildId + ".guild"

	if helpers.FileExists(filename) {
		data := ""
		if data, err = helpers.FileGetContents(filename); err != nil {
			return store, err
		}
		store = Store{}
		if err = json.Unmarshal([]byte(data), &store); err != nil {
			return store, err
		}
		store.MembersIds = map[string]Member{}
		store.Config = config
		return store, nil
	}

	if helpers.FileExists(config.PathToStorage + "guilds/" + guildId + "/" + guildId + ".members") {
		store = Store{
			GuildId:      guildId,
			GuildName:    "",
			MembersCount: 0,
			MembersIds:   map[string]Member{},
			Config:       config,
		}
		return store, nil
	}

	return store, fmt.Errorf("guild store " + guildId + " not found")
}

func NewStore(guildId string, guildName string, config *config.Config) (store Store, err error) {

	if guildId == "" {
		return store, fmt.Errorf("guild ID can't be empty")
	}

	store = Store{
		GuildId:      guildId,
		GuildName:    guildName,
		MembersCount: 0,
		MembersIds:   map[string]Member{},
		Config:       config,
	}

	return store, nil
}

func (s *Store) Save() (err error) {
	dir := s.Config.PathToStorage + "guilds/" + s.GuildId
	filename := dir + "/" + s.GuildId + ".guild"

	if !helpers.FileExists(dir) {
		err := helpers.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	s.MembersCount = len(s.MembersIds)
	data, err := json.Marshal(s)
	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func (s *Store) LoadMembers() (err error) {
	dir := s.Config.PathToStorage + "guilds/" + s.GuildId
	filename := dir + "/" + s.GuildId + ".members"

	if !helpers.FileExists(dir) {
		err = helpers.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}
	file, _ := os.Open(filename)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		id := fileScanner.Text()
		if _, err = regexp.Match(`\d+`, []byte(id)); err != nil {
			s.MembersIds[id] = Member{
				MemberId: id,
			}
		}
	}
	return nil
}

func (s *Store) SaveMembers() (err error) {
	dir := s.Config.PathToStorage + "guilds/" + s.GuildId
	filename := dir + "/" + s.GuildId + ".members"

	if !helpers.FileExists(dir) {
		err = helpers.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	//f, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	for _, i := range s.MembersIds {
		if _, err = f.WriteString(i.MemberId + "\n"); err != nil {
			panic(err)
		}
	}

	return nil
}

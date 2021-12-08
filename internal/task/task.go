package task

import (
	"encoding/json"
	"fmt"
	"la_discord_bot/internal/config"
	"la_discord_bot/internal/discordacc"
	"la_discord_bot/internal/guild"
	"la_discord_bot/internal/helpers"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	StatusCreated = 1
	StatusWorking = 2
	StatusStopped = 3
	StatusDone    = 4
	StatusError   = 5
)

const (
	TypeTaskParse    = 1
	TypeTaskSend     = 2
	TypeTaskParseAll = 3
)

const ParseBotLimit = 1000

type Task struct {
	Id               string `json:"id"`
	Token            string `json:"token"`
	TypeTask         int    `json:"type_task"`
	GuildId          string `json:"guild_id"`
	GuildName        string `json:"guild_name"`
	GuildMemberCount int    `json:"guild_member_count"`
	UserId           string `json:"user_id"`
	UserName         string `json:"user_name"`
	ChannelId        string `json:"channel_id"`
	ChannelName      string `json:"channel_name"`

	ParseCount     int               `json:"parse_count"`
	ParseMemberIds map[string]string `json:"-"`

	//SendMemberIds []string `json:"-"`
	SendMemberId   string         `json:"current_member_id"`
	SendSuccess    int            `json:"send_success"`
	SendFail       int            `json:"send_fail"`
	SendMessage    string         `json:"send_message"`
	SendDelayMin   float64        `json:"send_delay_min"`
	SendDelayMax   float64        `json:"send_delay_max"`
	SendInvite     string         `json:"send_invite"`
	SendMembersIDS []guild.Member `json:"send_members_ids"`

	Created     int64                 `json:"created"`
	Status      int                   `json:"status"`
	Steps       int                   `json:"steps"`
	CurrentStep int                   `json:"current_step"`
	Da          discordacc.DiscordAcc `json:"-"`
	Err         string                `json:"error"`
	Log         *log.Logger           `json:"-"`
	Config      *config.Config        `json:"-"`
}

func NewTask(token string, type_task int, guild_id string, channel_id string, conf *config.Config) (task *Task, err error) {
	//id := guild_id + "-" + channel_id + "-" + fmt.Sprint(type_task)
	id := uuid.New()

	task = &Task{
		Id:             id.String(),
		Token:          token,
		TypeTask:       type_task,
		GuildId:        guild_id,
		ChannelId:      channel_id,
		CurrentStep:    0,
		ParseMemberIds: map[string]string{},
		Config:         conf,
	}

	f, err := os.OpenFile(task.Filename()+".log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//defer f.Close()
	task.Log = log.New(f, "", log.Ldate|log.Ltime)
	task.Created = time.Now().Unix()
	task.Status = StatusCreated
	if err := task.Save(); err != nil {
		task.Log.Println(err)
	}

	return
}

func (T *Task) Resume() (err error) {

	return nil
}

func (T *Task) Delete() (err error) {
	if T.Status == StatusWorking {
		return fmt.Errorf("can't delete. task working")
	}

	os.Remove(T.FilenameTask())
	os.Remove(T.FilenameLog())

	return nil
}

func (T *Task) Filename() (filename string) {
	return T.Config.PathToStorage + "tasks/" + T.Id
}

func (T *Task) FilenameLog() (filename string) {
	return T.Config.PathToStorage + "tasks/" + T.Id + ".log"
}

func (T *Task) FilenameTask() (filename string) {
	return T.Config.PathToStorage + "tasks/" + T.Id + ".task"
}

func (T *Task) FilenameMembersIds() (filename string) {
	return T.Config.PathToStorage + "tasks/" + T.Id + ".members"
}

func (T *Task) Start() (err error) {

	if T.Status == StatusWorking {
		T.Log.Println("Task already started. Break.")
		return fmt.Errorf("task already started. break")
	}

	if T.Status == StatusDone {
		T.Log.Println("Task already Done. Break.")
		return fmt.Errorf("task already done. break")
	}

	T.Log.Println("Task Starting")
	if T.TypeTask == TypeTaskParse {
		T.Log.Println("Task Parse")
		// Clear members file
		if T.Status == StatusCreated {
			dir := T.Config.PathToStorage + "guilds/" + T.GuildId
			filename := dir + "/" + T.GuildId + ".members"
			if helpers.FileExists(filename) {
				if err = os.Remove(filename); err != nil {
					T.SetError(err)
					T.Stop()
					return err
				}
			}
		}

		T.Status = StatusWorking
		if err = T.ParseMembers(); err != nil {
			T.SetError(err)
			T.Stop()
			//T.Save()
			return err
		}
	}

	if T.TypeTask == TypeTaskParseAll {
		T.Log.Println("Task Parse all")
		if T.Status == StatusCreated {
			//
		}

		T.Status = StatusWorking
		if err = T.ParseMembersAll(); err != nil {
			T.SetError(err)
			T.Stop()
			return err
		}
	}

	if T.TypeTask == TypeTaskSend {
		T.Log.Println("Task Send")
		if T.Status == StatusCreated {
			//
		}
		T.Status = StatusWorking
		T.Save()
		if err = T.SendMessages(); err != nil {
			T.SetError(err)
			T.Stop()
			return err
		}
	}

	if T.Status == StatusWorking {
		T.Done()
	}
	return nil
}

func (T *Task) Stop() (err error) {
	if T.Status == StatusWorking {
		T.Status = StatusStopped
	}
	T.Log.Println("Task stopped")
	T.Save()
	return nil
}

func (T *Task) SetError(err error) {
	T.Status = StatusError
	T.Log.Println("Task Error: " + fmt.Sprint(err))
	T.Err = fmt.Sprintf("%v", err)
	T.Stop()
}
func (T *Task) Done() {
	T.Status = StatusDone
	T.Log.Println("Task Done")
	T.Save()
}

func (T *Task) Save() error {
	filename := T.Filename() + ".task"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		T.SetError(err)
		return err
	}
	defer f.Close()
	bytes, err := json.Marshal(T)
	if err != nil {
		T.SetError(err)
		return err
	}
	if _, err = f.Write(bytes); err != nil {
		T.SetError(err)
		return err
	}
	//T.Log.Println("Task Saved")
	return nil
}

func (T *Task) SaveGuildMembers(members map[string]string, flag int) error {
	dir := T.Config.PathToStorage + "guilds/" + T.GuildId
	filename := dir + "/" + T.GuildId + ".members"

	if !helpers.FileExists(dir) {
		err := helpers.Mkdir(dir, 0777)
		if err != nil {
			return err
		}
	}

	//f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	f, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, s := range members {
		if _, err = f.WriteString(s + "\n"); err != nil {
			panic(err)
		}
	}

	return nil
}

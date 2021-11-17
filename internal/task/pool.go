package task

import (
	"encoding/json"
	"fmt"
	"la_discord_bot/internal/config"
	"la_discord_bot/internal/helpers"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Pool struct {
	Tasks   []*Task
	AppConf *config.Config
}

func (P *Pool) LoadTasks(c *config.Config) (err error) {
	files, err := filepath.Glob(c.PathToStorage + "tasks/*")
	//fmt.Printf("%+v\n", files)
	//fmt.Printf("%+v\n", err)
	for _, f := range files {
		if strings.Index(f, ".task") == -1 {
			continue
		}
		data, err := helpers.FileGetContents(f)
		if err != nil {
			continue
		}
		task := &Task{
			Config: c,
		}
		err = json.Unmarshal([]byte(data), task)
		if err != nil {
			continue
		}

		f, err := os.OpenFile(task.Filename()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		task.ParseMemberIds = map[string]string{}
		task.Log = log.New(f, "", log.Ldate|log.Ltime)

		if P.AppConf.AutoStartTasks == "1" && task.Status == StatusWorking {
			task.Status = StatusStopped
			P.Tasks = append(P.Tasks, task)
			P.Start(task.Id)
		} else {
			if task.Status == StatusWorking {
				task.Status = StatusStopped
			}
			P.Tasks = append(P.Tasks, task)
		}

	}
	return nil
}

func (P *Pool) NewTask(token string, type_task int, guild_id string, channel_id string, c *config.Config) (task *Task, err error) {

	if token == "" {
		return nil, fmt.Errorf("token can't be empty")
	}

	if guild_id == "" {
		return nil, fmt.Errorf("guild id can't be empty")
	}

	for _, t := range P.Tasks {
		//if t.TypeTask == type_task && t.GuildId == guild_id && t.ChannelId == channel_id {
		if t.TypeTask == type_task && t.GuildId == guild_id {
			return nil, fmt.Errorf("task already exists")
		}
	}

	task, err = NewTask(token, type_task, guild_id, channel_id, c)
	if err != nil {
		return nil, err
	}

	P.Tasks = append(P.Tasks, task)
	//return task, nil
	return
}

func (P *Pool) Start(taskId string) (err error) {
	for _, t := range P.Tasks {
		if t.Id == taskId {
			go func() {
				t.Start()
			}()
		}
	}
	return
}

func (P *Pool) Stop(taskId string) (err error) {
	for _, t := range P.Tasks {
		if t.Id == taskId {
			t.Stop()
			helpers.Sleep(3)
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func (P *Pool) Delete(taskId string) (err error) {
	for i, t := range P.Tasks {
		if t.Id == taskId {
			if err = t.Delete(); err != nil {
				return err
			}
			P.Tasks = append(P.Tasks[:i], P.Tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func (P *Pool) Resume(taskId string) (err error) {
	for _, t := range P.Tasks {
		if t.Id == taskId {
			t.Start()
			//discordacc.Sleep(3)
			return nil
		}
	}
	return fmt.Errorf("task not found")
}

func (P *Pool) Get(taskId string) (task *Task, err error) {
	for _, t := range P.Tasks {
		if t.Id == taskId {
			return t, nil
		}
	}
	return nil, fmt.Errorf("task not found")
}

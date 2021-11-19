package webserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"la_discord_bot/internal/config"
	"la_discord_bot/internal/discordacc"
	"la_discord_bot/internal/guild"
	"la_discord_bot/internal/helpers"
	"la_discord_bot/internal/task"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type DiscordError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type jsonSendTask struct {
	Token      string `json:"token"`
	GuildId    string `json:"guild_id"`
	MembersIDS string `json:"members_ids"`
	Invite     string `json:"invite"`
	DelayMin   string `json:"delay_min"`
	DelayMax   string `json:"delay_max"`
	Message    string `json:"message"`
}

var Tasks = task.Pool{}

var AppConf = config.Config{}

func Init(c config.Config) (err error) {

	AppConf = c
	Tasks.AppConf = &AppConf
	if err = Tasks.LoadTasks(&AppConf); err != nil {
		log.Println("Error Load Tasks")
	}
	gin.SetMode(c.GinMode)
	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		AppConf.Login: AppConf.Password,
	}))

	authorized.StaticFile("/", AppConf.PathToWWW+"index.html")
	authorized.Static("/css", AppConf.PathToWWW+"css")
	authorized.Static("/js", AppConf.PathToWWW+"js")

	authorized.GET("/api/discord/UserByToken", func(c *gin.Context) {
		auth := c.Query("auth")
		token := c.Query("token")
		da, err := discordacc.New(token, auth, true)
		defer da.Close()
		if err != nil {
			fmt.Printf("%v\n", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": da.User, "guilds": da.Guilds})
	})

	authorized.GET("/api/discord/joinGuild", func(c *gin.Context) {
		inviteCode := c.Query("invite_code")
		inviteCode = strings.Replace(inviteCode, "http://", "", 1)
		inviteCode = strings.Replace(inviteCode, "https://", "", 1)
		inviteCode = strings.Replace(inviteCode, "discord.gg/", "", 1)
		inviteCode = strings.Replace(inviteCode, "discord.com/invite/", "", 1)
		// discord.gg/dbd
		// https://discord.com/invite/dbd
		token := c.Query("token")
		err := helpers.JoinGuild(inviteCode, token, helpers.Fingerprintx{}, helpers.Cookie{})
		if err != nil {
			//log.Printf("%v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			da, err := discordacc.New(token, "Bot", true)
			defer da.Close()
			if err != nil {
				//fmt.Printf("%v\n", err)
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, gin.H{"user": da.User, "guilds": da.Guilds})
		}

	})

	authorized.GET("/api/discord/task/parse", func(c *gin.Context) {
		guildId := c.Query("guild_id")
		channelId := c.Query("channel_id")
		token := c.Query("token")

		task, err := Tasks.NewTask(token, task.TypeTaskParse, guildId, channelId, &AppConf)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		err = Tasks.Start(task.Id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"task_id": task.Id})

	})

	authorized.GET("/api/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"tasks": Tasks.TasksSlice()})
	})

	authorized.GET("/api/tasks/stop", func(c *gin.Context) {
		taskId := c.Query("task_id")
		err = Tasks.Stop(taskId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": Tasks.TasksSlice()})
	})

	authorized.GET("/api/tasks/delete", func(c *gin.Context) {
		taskId := c.Query("task_id")
		err = Tasks.Delete(taskId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": Tasks.TasksSlice()})
	})

	authorized.GET("/api/tasks/resume", func(c *gin.Context) {
		taskId := c.Query("task_id")
		err = Tasks.Resume(taskId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tasks": Tasks.TasksSlice()})
	})

	authorized.GET("/api/guilds/file", func(c *gin.Context) {
		guildId := c.Query("guild_id")
		download := c.Query("download")
		filename := AppConf.PathToStorage + "guilds/" + guildId + "/" + guildId + ".members"

		if !helpers.FileExists(filename) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("file not found")})
			return
		}

		header := c.Writer.Header()
		header["Content-type"] = []string{"application/octet-stream"}
		if download != "" {
			header["Content-Disposition"] = []string{"attachment; filename=" + guildId + "members.txt"}
		}

		file, err := os.Open(filename)
		if err != nil {
			c.String(http.StatusOK, "%v", err)
			return
		}
		defer file.Close()

		io.Copy(c.Writer, file)
	})

	authorized.GET("/api/tasks/file", func(c *gin.Context) {
		taskId := c.Query("task_id")
		fileType := c.Query("type")
		download := c.Query("download")

		T, err := Tasks.Get(taskId)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		filename := T.FilenameLog()
		ext := ".log.txt"
		switch fileType {
		case "log":
			filename = T.FilenameLog()
			ext = ".log.txt"
		case "members":
			filename = T.FilenameMembersIds()
			ext = ".members.txt"
		case "task":
			filename = T.FilenameTask()
			ext = ".task.txt"
		default:
			filename = T.FilenameLog()
			ext = ".log.txt"
		}

		header := c.Writer.Header()
		header["Content-type"] = []string{"application/octet-stream"}

		if download != "" {
			header["Content-Disposition"] = []string{"attachment; filename=" + T.Id + ext}
		}

		file, err := os.Open(filename)
		if err != nil {
			c.String(http.StatusOK, "%v", err)
			return
		}
		defer file.Close()

		io.Copy(c.Writer, file)

	})

	authorized.GET("/api/guilds/store", func(c *gin.Context) {
		stores, err := guild.LoadAll(&AppConf)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"stores": stores.Stores})
	})

	authorized.POST("/api/tasks/send", func(c *gin.Context) {
		sendTask := jsonSendTask{}

		if err := c.ShouldBindJSON(&sendTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if strings.TrimSpace(sendTask.Token) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token can't be empty"})
			return
		}

		if strings.TrimSpace(sendTask.Message) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message can't be empty"})
			return
		}

		if strings.TrimSpace(sendTask.GuildId) == "" && strings.TrimSpace(sendTask.MembersIDS) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Guild and members can't be empty"})
			return
		}

		t, err := Tasks.NewTask(sendTask.Token, task.TypeTaskSend, sendTask.GuildId, "", &AppConf)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if strings.TrimSpace(sendTask.MembersIDS) != "" {
			for _, s := range helpers.Explode("\n", strings.TrimSpace(sendTask.MembersIDS)) {
				if strings.TrimSpace(s) == "" {
					continue
				}
				t.SendMembersIDS = append(t.SendMembersIDS, guild.Member{
					MemberId: s,
				})

			}
		}

		t.SendInvite = strings.TrimSpace(sendTask.Invite)
		t.SendMessage = strings.TrimSpace(sendTask.Message)

		if f, err := strconv.ParseFloat(strings.TrimSpace(sendTask.DelayMin), 32); err == nil {
			if f < AppConf.DelayMin {
				t.Log.Printf("Delay %v min  < %v . Set delay min to %v", f, AppConf.DelayMin, AppConf.DelayMin)
				f = AppConf.DelayMin
			}
			t.SendDelayMin = f
		} else {
			t.Log.Printf("Delay min %v not float. Set %v", sendTask.DelayMin, AppConf.DelayMin)
			t.SendDelayMin = AppConf.DelayMin
		}

		if f, err := strconv.ParseFloat(strings.TrimSpace(sendTask.DelayMax), 32); err == nil {
			if f > AppConf.DelayMax {
				t.Log.Printf("Delay %v max  > %v . Set delay max to %v", f, AppConf.DelayMax, AppConf.DelayMax)
				f = AppConf.DelayMax
			}
			t.SendDelayMax = f
		} else {
			t.Log.Printf("Delay max %v not float. Set %v", sendTask.DelayMax, AppConf.DelayMax)
			t.SendDelayMax = AppConf.DelayMax
		}

		if err = t.Save(); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		err = Tasks.Start(t.Id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"task_id": t.Id})
	})

	//authorized.POST("/api/tasks/sendonce", func(c *gin.Context) {
	//	sendTask := jsonSendTask{}
	//	if err := c.ShouldBindJSON(&sendTask); err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//		return
	//	}
	//
	//	t, err := Tasks.NewTask(sendTask.Token, task.TypeTaskSend, sendTask.GuildId, "", &AppConf)
	//	if err != nil {
	//		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	//		return
	//	}
	//
	//	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	//})

	// TODO Upload memberIds file
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	err = r.Run(":" + port)
	if err != nil {
		return err
	}

	return nil
}

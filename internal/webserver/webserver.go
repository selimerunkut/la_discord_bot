package webserver

import (
	"github.com/gin-gonic/gin"
	"la_discord_bot/internal/discordacc"
	"net/http"
	"time"
)

func Init() (err error) {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/discord/UserByToken", func(c *gin.Context) {
		auth := c.Query("auth")
		token := c.Query("token")
		da, err := discordacc.New(token, auth, true)
		defer da.Close()
		if err != nil {
			c.JSON(http.StatusGone, gin.H{"error": err})
		}
		if err = da.GetMembers("814273262154416178", "814278005782347786", 0); err != nil {
			c.JSON(http.StatusGone, gin.H{"error": err})
		}
		time.Sleep(10 * time.Second)
		c.JSON(http.StatusOK, gin.H{"user": da.User, "guilds": da.Guilds})
	})

	err = r.Run()
	if err != nil {
		return err
	}

	return nil
}

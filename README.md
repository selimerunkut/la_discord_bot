# la_discord_bot

### Compile

```shell
go build -o bin/webserver -v cmd/webserver/main.go
```

### Run local
```shell
./bin/webserver
```

app should be visible here:

```
http://localhost:8080/
```

## Instruction

### Token
1. Paste Discord token to "Token".
2. Click "Get Info"
3. If the token is valid you'll see user info.

### Join \ Invite
If it's a User Token you can join to the server by invite.
Paste to "Invite:" needed invite code or invite URL and click Join.

Examples of invites:
```
https://discord.gg/dbd
https://discord.com/invite/dbd
dbd
```


### Parse Members
1. Choose guild to parse.
2. Choose a channel to parse. (Better choose a channel that available for all users, like "Rules", "Info", etc)
3. Click "Parse members"

Parse task will be created.

Note. To prevent double parsing from one server (guild) you can have only one Parse Task per Server.

If you need to reparse members, you need to stop the task (if it running), then delete it.


### Send Message
1. Choose Members to send message
   1. "Parsed Guilds" get members from guilds that were parsed.
   2. "Manual Member IDs" - paste IDs list (one ID per line)
2. "Join By Invite Before Send" paste invite code, if need joins to the server before sending messages.
3. "Delay Between Send Messages" sets the random delay between sending messages.
4. "Message" - message to send. You can use <user> to mention the user in the message.
5. Click "Send Message"

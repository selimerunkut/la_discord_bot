function form() {
    return {
        taskStatusToText: {
            1: "Created",
            2: "Working",
            3: "Stopped",
            4: "Done",
            5: "Error"
        },
        taskStatusToClass: {
            1: "",
            2: "alert-primary",
            3: "alert-secondary",
            4: "alert-success",
            5: "alert-danger"
        },
        taskTypeToText: {
            1: "Parse",
            2: "Send"
        },
        rs: false,
        send: {
            token: "",
            guild_id: "",
            members_ids: "",
            // members_ids: [],
            invite: "",
            delay_min: "",
            delay_max: "",
            message: "",
        },
        token: "",
        guildStores: [],

        successMessage: false,
        errorGetInfo: false,
        errorJoin: false,
        errorParse: false,
        errorSend: false,
        tasks: false,
        acc: false,
        guilds: false,
        currentGuildId: false,
        channels: false,
        currentChannelId: false,
        clientId: "",
        invite: "",
        sendTaskButtonClick: async function (event) {
            this.errorSend = false
            this.send.token = this.token
            let response = await fetch('/api/tasks/send', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json;charset=utf-8'
                },
                body: JSON.stringify(this.send)
            });
            if (response.ok) {
                let json = await response.json();
                this.updateTasks()
                this.successMessage = "Send Task Created"
            } else {
                let json = { error: "" }
                try {
                    json = await response.json();
                } catch (e) {
                    json = { error: "505 Internal Server Error"}
                }
                console.log(json)
                this.errorSend = json.error
            }

        },
        getInfoButtonClick: async function (event) {
            // console.log(event)
            this.rs = true
            this.acc = false
            this.guilds = false
            this.currentGuildId = false
            this.currentChannelId = false
            this.errorGetInfo = false
            this.errorJoin = false
            this.errorParse = false
            this.errorSend = false
            this.successMessage = false
            this.clientId = ""

            let auth = "Bot" //document.getElementById("auth").value
            let token = document.getElementById("token").value
            let response = await fetch("/api/discord/UserByToken?auth=" + auth + "&token=" + token);

            if (response.ok) {
                let json = await response.json();
                this.acc = json.user
                // console.log(this.acc)
                this.guilds = json.guilds
                // console.log(this.guilds)
                // this.selectGuildsChange()
            } else {
                let json = { error: "" }
                try {
                    json = await response.json();
                } catch (e) {
                    json = { error: "505 Internal Server Error"}
                }
                console.log(json)
                this.errorGetInfo = json.error
                // alert("HTTP-Error: " + response.status);
            }
            this.rs = false
        },
        selectGuildsChange: function (event) {
            this.currentGuildId = document.getElementById("guilds").value
        },
        parseMembersButtonClick: async function (event) {
            this.rs = true
            this.errorParse = false
            this.successMessage = false

            let channel = ""
            // console.log(channel)
            try {
                // console.log(document.getElementById("channels").value)
                channel = document.getElementById("channels").value
            }
            catch (e) {
                channel = ""
                // console.log(e)
            }
            // console.log(channel)
            let token = document.getElementById("token").value
            // let response = await fetch("/api/discord/ParseMembers?token=" + token + "&guild_id=" + this.currentGuildId);
            let response = await fetch("/api/discord/task/parse?token=" + token + "&guild_id=" + this.currentGuildId + "&channel_id=" + channel);
            if (response.ok) {
                let json = await response.json();
                await this.updateTasks()
                this.successMessage = "Parse Task Created"
                // console.log(json)
            } else {
                let json = { error: "" }
                try {
                    json = await response.json();
                } catch (e) {
                    json = { error: "505 Internal Server Error"}
                }
                console.log(json)
                this.errorParse = json.error
                // alert("HTTP-Error: " + response.status );
            }
            this.rs = false
        },
        joinButtonClick: async function (event) {
            this.rs = true
            this.errorJoin = false

            let token = document.getElementById("token").value
            if (this.invite !== "") {
                let response = await fetch("/api/discord/joinGuild?token=" + token + "&invite_code=" + this.invite);
                if (response.ok) {
                    let json = await response.json();
                    this.acc = json.user
                    // console.log(this.acc)
                    this.guilds = json.guilds
                    console.log(json)
                    this.invite = ""
                    this.successMessage = "Join Success"
                } else {
                    let json = { error: "" }
                    try {
                        json = await response.json();
                    } catch (e) {
                        json = { error: "505 Internal Server Error"}
                    }
                    console.log(json)
                    this.errorJoin = json.error
                    // alert("HTTP-Error: " + response.status );
                }
                this.rs = false
            }
        },
        init() {
            this.updateTasks()
            setInterval(function () {document.getElementById('refreshTasksButton').click()}, 10000)
        },
        refreshTasksButtonClick: function (event) {
            this.updateTasks()
        },
        deleteTaskButtonClick: function (task_id) {
            // console.log(task_id)
            let result = confirm("Are you sure to delete task?")
            if (result) {
                self = this
                fetch('/api/tasks/delete?task_id='+task_id)
                    .then(response => response.json())
                    .then(tasks => {
                        self.tasks = tasks.tasks
                        // console.log(tasks)
                    })
            }
        },
        stopTaskButtonClick: function (task_id) {
            // console.log(task_id)
            let result = confirm("Are you sure to stop task?")
            if (result) {
                self = this
                fetch('/api/tasks/stop?task_id='+task_id)
                    .then(response => response.json())
                    .then(tasks => {
                        self.tasks = tasks.tasks
                        // console.log(tasks)
                    })
            }
        },
        resumeTaskButtonClick: function (task_id) {
            // console.log(task_id)
            let result = confirm("Are you sure to resume task?")
            if (result) {
                self = this
                fetch('/api/tasks/resume?task_id='+task_id)
                    .then(response => response.json())
                    .then(tasks => {
                        self.tasks = tasks.tasks
                        // console.log(tasks)
                    })
            }
        },
        updateTasks: function () {
            self = this
            fetch('/api/tasks')
                .then(response => response.json())
                .then(tasks => {
                    self.tasks = tasks.tasks
                    console.log(tasks)
                })
        },
        updateGuildStores: function () {
            self = this
            fetch('/api/guilds/store')
                .then(response => response.json())
                .then(stores => {
                    self.guildStores = stores.stores
                    console.log(stores)
                })
        },
    };
}
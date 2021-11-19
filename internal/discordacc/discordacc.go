package discordacc

import (
	"la_discord_bot/internal/discordgo"
	"la_discord_bot/internal/helpers"
)

type Guild struct {
	// The ID of the guild.
	ID string `json:"id"`

	// The name of the guild. (2–100 characters)
	Name string `json:"name"`

	// The hash of the guild's icon. Use Session.GuildIcon
	// to retrieve the icon itself.
	Icon string `json:"icon"`

	// The voice region of the guild.
	Region string `json:"region"`

	// The ID of the AFK voice channel.
	AfkChannelID string `json:"afk_channel_id"`

	// The user ID of the owner of the guild.
	OwnerID string `json:"owner_id"`

	// If we are the owner of the guild
	Owner bool `json:"owner"`

	// The time at which the current user joined the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	JoinedAt discordgo.Timestamp `json:"joined_at"`

	// The hash of the guild's discovery splash.
	DiscoverySplash string `json:"discovery_splash"`

	// The hash of the guild's splash.
	Splash string `json:"splash"`

	// The timeout, in seconds, before a user is considered AFK in voice.
	AfkTimeout int `json:"afk_timeout"`

	// The number of members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	MemberCount int `json:"member_count"`

	// The verification level required for the guild.
	VerificationLevel discordgo.VerificationLevel `json:"verification_level"`

	// Whether the guild is considered large. This is
	// determined by a member threshold in the identify packet,
	// and is currently hard-coded at 250 members in the library.
	Large bool `json:"large"`

	// The default message notification setting for the guild.
	DefaultMessageNotifications discordgo.MessageNotifications `json:"default_message_notifications"`

	// A list of roles in the guild.
	Roles []discordgo.Role `json:"roles"`

	// A list of the custom emojis present in the guild.
	Emojis []discordgo.Emoji `json:"emojis"`

	// A list of the members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Members []discordgo.Member `json:"members"`

	// A list of partial presence objects for members in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Presences []discordgo.Presence `json:"presences"`

	// The maximum number of presences for the guild (the default value, currently 25000, is in effect when null is returned)
	MaxPresences int `json:"max_presences"`

	// The maximum number of members for the guild
	MaxMembers int `json:"max_members"`

	// A list of channels in the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Channels map[string]discordgo.Channel `json:"channels"`

	// A list of voice states for the guild.
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	VoiceStates []discordgo.VoiceState `json:"voice_states"`

	// Whether this guild is currently unavailable (most likely due to outage).
	// This field is only present in GUILD_CREATE events and websocket
	// update events, and thus is only present in state-cached guilds.
	Unavailable bool `json:"unavailable"`

	// The explicit content filter level
	ExplicitContentFilter discordgo.ExplicitContentFilterLevel `json:"explicit_content_filter"`

	// The list of enabled guild features
	Features []string `json:"features"`

	// Required MFA level for the guild
	MfaLevel discordgo.MfaLevel `json:"mfa_level"`

	// The application id of the guild if bot created.
	ApplicationID string `json:"application_id"`

	// Whether or not the Server Widget is enabled
	WidgetEnabled bool `json:"widget_enabled"`

	// The Channel ID for the Server Widget
	WidgetChannelID string `json:"widget_channel_id"`

	// The Channel ID to which system messages are sent (eg join and leave messages)
	SystemChannelID string `json:"system_channel_id"`

	// The System channel flags
	SystemChannelFlags discordgo.SystemChannelFlag `json:"system_channel_flags"`

	// The ID of the rules channel ID, used for rules.
	RulesChannelID string `json:"rules_channel_id"`

	// the vanity url code for the guild
	VanityURLCode string `json:"vanity_url_code"`

	// the description for the guild
	Description string `json:"description"`

	// The hash of the guild's banner
	Banner string `json:"banner"`

	// The premium tier of the guild
	PremiumTier discordgo.PremiumTier `json:"premium_tier"`

	// The total number of users currently boosting this server
	PremiumSubscriptionCount int `json:"premium_subscription_count"`

	// The preferred locale of a guild with the "PUBLIC" feature; used in server discovery and notices from Discord; defaults to "en-US"
	PreferredLocale string `json:"preferred_locale"`

	// The id of the channel where admins and moderators of guilds with the "PUBLIC" feature receive notices from Discord
	PublicUpdatesChannelID string `json:"public_updates_channel_id"`

	// The maximum amount of users in a video channel
	MaxVideoChannelUsers int `json:"max_video_channel_users"`

	// Approximate number of members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximateMemberCount int `json:"approximate_member_count"`

	// Approximate number of non-offline members in this guild, returned from the GET /guild/<id> endpoint when with_counts is true
	ApproximatePresenceCount int `json:"approximate_presence_count"`

	// Permissions of our user
	Permissions int64 `json:"permissions,string"`
}

// A Channel holds all data related to an individual Discord channel.
type Channel struct {
	// The ID of the channel.
	ID string `json:"id"`

	// The ID of the guild to which the channel belongs, if it is in a guild.
	// Else, this ID is empty (e.g. DM channels).
	GuildID string `json:"guild_id"`

	// The name of the channel.
	Name string `json:"name"`

	// The topic of the channel.
	Topic string `json:"topic"`

	// The type of the channel.
	Type discordgo.ChannelType `json:"type"`

	// The ID of the last message sent in the channel. This is not
	// guaranteed to be an ID of a valid message.
	LastMessageID string `json:"last_message_id"`

	// The timestamp of the last pinned message in the channel.
	// Empty if the channel has no pinned messages.
	LastPinTimestamp discordgo.Timestamp `json:"last_pin_timestamp"`

	// Whether the channel is marked as NSFW.
	NSFW bool `json:"nsfw"`

	// Icon of the group DM channel.
	Icon string `json:"icon"`

	// The position of the channel, used for sorting in client.
	Position int `json:"position"`

	// The bitrate of the channel, if it is a voice channel.
	Bitrate int `json:"bitrate"`

	// The recipients of the channel. This is only populated in DM channels.
	Recipients []discordgo.User `json:"recipients"`

	// The messages in the channel. This is only present in state-cached channels,
	// and State.MaxMessageCount must be non-zero.
	Messages []discordgo.Message `json:"-"`

	// A list of permission overwrites present for the channel.
	PermissionOverwrites []discordgo.PermissionOverwrite `json:"permission_overwrites"`

	// The user limit of the voice channel.
	UserLimit int `json:"user_limit"`

	// The ID of the parent channel, if the channel is under a category
	ParentID string `json:"parent_id"`

	// Amount of seconds a user has to wait before sending another message (0-21600)
	// bots, as well as users with the permission manage_messages or manage_channel, are unaffected
	RateLimitPerUser int `json:"rate_limit_per_user"`

	// ID of the DM creator Zeroed if guild channel
	OwnerID string `json:"owner_id"`

	// ApplicationID of the DM creator Zeroed if guild channel or not a bot user
	ApplicationID string `json:"application_id"`
}

type DiscordAcc struct {
	Auth    string
	Token   string
	Session *discordgo.Session
	User    discordgo.User
	Guilds  map[string]Guild
}

func New(token string, auth string, auto bool) (dac DiscordAcc, err error) {
	newAcc := DiscordAcc{
		Token:  token,
		Auth:   auth,
		Guilds: make(map[string]Guild),
	}

	newAcc.Session, err = discordgo.New(newAcc.Auth + " " + newAcc.Token)
	if err != nil {
		return newAcc, err
	}
	//newAcc.Session.Identify.Intents = discordgo.IntentsAll
	if auto {
		err = newAcc.Connect()

		if err != nil {
			return newAcc, err
		}

		err = newAcc.ReloadData()
		if err != nil {
			return newAcc, err
		}
	}

	return newAcc, nil
}

func (a *DiscordAcc) Connect() (err error) {
	err = a.Session.Open()
	if err != nil {
		return err
	}
	return
}

func (a *DiscordAcc) ReloadData() (err error) {
	if a.Session.State.User != nil {
		a.User = *a.Session.State.User
	}

	// load guilds
	for _, g := range a.Session.State.Guilds {
		if a.User.Bot {
			g, err = a.Session.Guild(g.ID)
			if err != nil {
				return err
			}
			g.Channels, err = a.Session.GuildChannels(g.ID)
			if err != nil {
				return err
			}
			//a.Session.GuildMembers()
		}
		tg := Guild{
			ID:                          g.ID,
			Name:                        g.Name,
			Icon:                        g.Icon,
			Region:                      g.Region,
			AfkChannelID:                g.AfkChannelID,
			OwnerID:                     g.OwnerID,
			Owner:                       g.Owner,
			JoinedAt:                    g.JoinedAt,
			DiscoverySplash:             g.DiscoverySplash,
			Splash:                      g.Splash,
			AfkTimeout:                  g.AfkTimeout,
			MemberCount:                 g.MemberCount,
			VerificationLevel:           g.VerificationLevel,
			Large:                       g.Large,
			DefaultMessageNotifications: g.DefaultMessageNotifications,
			//Roles:                       *g.Roles,
			//Emojis:                      nil,
			//Members:                     nil,
			//Presences:                   nil,
			MaxPresences: g.MaxPresences,
			MaxMembers:   g.MaxMembers,
			Channels:     make(map[string]discordgo.Channel),
			//VoiceStates:                 nil,
			Unavailable:           g.Unavailable,
			ExplicitContentFilter: g.ExplicitContentFilter,
			//Features:                    g.Features,
			MfaLevel:                 g.MfaLevel,
			ApplicationID:            g.ApplicationID,
			WidgetEnabled:            g.WidgetEnabled,
			WidgetChannelID:          g.WidgetChannelID,
			SystemChannelID:          g.SystemChannelID,
			SystemChannelFlags:       g.SystemChannelFlags,
			RulesChannelID:           g.RulesChannelID,
			VanityURLCode:            g.VanityURLCode,
			Description:              g.Description,
			Banner:                   g.Banner,
			PremiumTier:              g.PremiumTier,
			PremiumSubscriptionCount: g.PremiumSubscriptionCount,
			PreferredLocale:          g.PreferredLocale,
			PublicUpdatesChannelID:   g.PublicUpdatesChannelID,
			MaxVideoChannelUsers:     g.MaxVideoChannelUsers,
			ApproximateMemberCount:   g.ApproximateMemberCount,
			ApproximatePresenceCount: g.ApproximatePresenceCount,
			Permissions:              g.Permissions,
		}

		for _, c := range g.Channels {
			tg.Channels[c.ID] = *c
		}

		for _, m := range g.Members {
			tg.Members = append(tg.Members, *m)
		}

		a.Guilds[g.ID] = tg
	}
	return
}

func (a *DiscordAcc) Close() (err error) {
	return a.Session.Close()
}

func (a *DiscordAcc) Join(invite string) (err error) {
	err = helpers.JoinGuild(invite, a.Token, helpers.Fingerprintx{}, helpers.Cookie{})
	if err != nil {
		return err
	}
	return nil
}

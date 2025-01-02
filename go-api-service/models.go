package main

// Models
// {"response":{"steamid":"76561198305662842","success":1}}
type SteamId struct {
	Response struct {
		Steamid  string `json:"steamid"`
		Success  int    `json:"success"`
		Username string `json:"username"`
	} `json:"response"`
}

type PlayerSummaryResponse struct {
	Response struct {
		Player []PlayerSummary `json:"players"`
	} `json:"response"`
}

type PlayerSummary struct {
	SteamId     string `json:"steamid"`
	PersonaName string `json:"personaname"`
	ProfileUrl  string `json:"profileurl"`
	Avatar      string `json:"avatar"`
	AvatarFull  string `json:"avatarfull"`
}

type OwnedGame struct {
	AppId     int    `json:"appid"`
	GameName  string `json:"name"`
	Playtime  int    `json:"playtime_forever"` // Total # of minutes played
	ImageIcon string `json:"img_icon_url"`
	ImageLogo string `json:"img_logo_url"`
}

type GameData struct {
	Response struct {
		GameCount  int         `json:"game_count"`
		OwnedGames []OwnedGame `json:"games"`
	} `json:"response"`
}

type Data struct {
	PersonaName string      `json:"persona_name"`
	ProfileUrl  string      `json:"profile_url"`
	Avatar      string      `json:"avatar_url"`
	AvatarFull  string      `json:"avatar_full_url"`
	SteamId     string      `json:"steamid"`
	GameCount   int         `json:"game_count"`
	Games       []OwnedGame `json:"games"`
}

type Playtime struct {
	TotalPlaytime  int       `json:"total_playtime"` // in minutes
	MostPlayedGame OwnedGame `json:"most_played_game"`
}

type LifetimePlaytime struct {
	SteamId            string  `json:"steamid"`
	PersonaName        string  `json:"persona_name"`
	AppId              int     `json:"appid"`
	Name               string  `json:"name"`
	Playtime           int     `json:"playtime"`
	LifetimePlaytime   int     `json:"lifetime_playtime"`
	PlaytimePercentage float64 `json:"playtime_percentage"`
}

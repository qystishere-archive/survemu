package web

type Pops struct {
	Eu string `json:"eu"`
	Na string `json:"na"`
	Sa string `json:"sa"`
	As string `json:"as"`
}

type YouTube struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

type Twitch struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Img     string `json:"img"`
	Viewers int    `json:"viewers"`
}

type SiteInfoResult struct {
	Pops          *Pops      `json:"pops"`
	Youtube       []*YouTube `json:"youtube"`
	Twitch        []*Twitch  `json:"twitch"`
	PromptConsent bool       `json:"promptConsent"`
}

type Game struct {
	Zone   string   `json:"zone"`
	GameID string   `json:"gameId"`
	Hosts  []string `json:"hosts"`
	Addrs  []string `json:"addrs"`
	Data   string   `json:"data"`
}

type FindGameResult struct {
	Res []*Game `json:"res"`
}

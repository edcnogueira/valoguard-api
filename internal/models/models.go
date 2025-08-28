package models

import "encoding/json"

type APIResponse struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
}

type Account struct {
	Puuid        string `json:"puuid"`
	Name         string `json:"name"`
	Tag          string `json:"tag"`
	AccountLevel int    `json:"account_level"`
}

type MMR struct {
	CurrentTierPatched string `json:"currenttierpatched"`
}

type Player struct {
	Puuid string `json:"puuid"`
	Team  string `json:"team_id"`
	Stats struct {
		Kills   int `json:"kills"`
		Deaths  int `json:"deaths"`
		Assists int `json:"assists"`
		Head    int `json:"headshots"`
		Body    int `json:"bodyshots"`
		Leg     int `json:"legshots"`
		Score   int `json:"score"`
		Damage  struct {
			Dealt int `json:"dealt"`
		} `json:"damage"`
	} `json:"stats"`
}

type Match struct {
	Metadata struct {
		MatchID string `json:"match_id"`
	} `json:"metadata"`
	Players []Player `json:"players"`
	Teams   []struct {
		TeamID     string `json:"team_id"`
		Won        bool   `json:"won"`
		RoundsWon  int    `json:"rounds_won"`
		RoundsLost int    `json:"rounds_lost"`
	} `json:"teams"`
}

type rawMatch struct {
	Metadata struct {
		MatchID string `json:"match_id"`
	} `json:"metadata"`
	Players map[string]Player `json:"players"`
	Teams   []struct {
		TeamID     string `json:"team_id"`
		Won        bool   `json:"won"`
		RoundsWon  int    `json:"rounds_won"`
		RoundsLost int    `json:"rounds_lost"`
	} `json:"teams"`
}

func (m *Match) UnmarshalJSON(data []byte) error {
	var raw rawMatch
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.Metadata = raw.Metadata
	m.Teams = raw.Teams

	m.Players = make([]Player, 0, len(raw.Players))
	for _, player := range raw.Players {
		m.Players = append(m.Players, player)
	}

	return nil
}

type PlayerStats struct {
	KD           float64 `json:"kd"`
	HSPercent    float64 `json:"hs_percent"`
	WinRate      float64 `json:"win_rate"`
	AvgDamage    float64 `json:"avg_damage"`
	AccountLevel int     `json:"account_level"`
	Rank         string  `json:"rank"`
}

type MatchSummary struct {
	MatchID string `json:"match_id"`
	Score   int    `json:"score"`
}

type AnalyzeRequest struct {
	Name   string
	Tag    string
	Region string
}

type AnalyzeResponse struct {
	Probability  int            `json:"probability"`
	Stats        PlayerStats    `json:"stats"`
	MatchSummary []MatchSummary `json:"matches_summary"`
}

package models

import (
	"encoding/json"
	"time"
)

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

type MMRCurrentData struct {
	CurrentTier        int    `json:"currenttier"`
	CurrentTierPatched string `json:"currenttierpatched"`
	Elo                int    `json:"elo"`
}

type MMRHighestRank struct {
	Tier        int    `json:"tier"`
	PatchedTier string `json:"patched_tier"`
	Season      string `json:"season"`
}

type SeasonData struct {
	Wins             int    `json:"wins"`
	NumberOfGames    int    `json:"number_of_games"`
	FinalRank        int    `json:"final_rank"`
	FinalRankPatched string `json:"final_rank_patched"`
}

type MMRBySeason map[string]SeasonData

type MMR struct {
	CurrentData MMRCurrentData `json:"current_data"`
	HighestRank MMRHighestRank `json:"highest_rank"`
	BySeason    MMRBySeason    `json:"by_season"`
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
	Metadata MatchMetadata `json:"metadata"`
	Players  []Player      `json:"players"`
	Rounds   []Round       `json:"rounds"`
	Teams    []Team        `json:"teams"`
}

type MatchMetadata struct {
	MatchID         string           `json:"match_id"`
	Map             Map              `json:"map"`
	GameVersion     string           `json:"game_version"`
	GameLengthInMs  int              `json:"game_length_in_ms"`
	StartedAt       time.Time        `json:"started_at"`
	IsCompleted     bool             `json:"is_completed"`
	Queue           Queue            `json:"queue"`
	Season          Season           `json:"season"`
	Platform        string           `json:"platform"`
	Premier         interface{}      `json:"premier"`
	PartyRrPenaltys []PartyRrPenalty `json:"party_rr_penaltys"`
	Region          string           `json:"region"`
	Cluster         string           `json:"cluster"`
}

type Map struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Queue struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ModeType string `json:"mode_type"`
}

type Season struct {
	ID    string `json:"id"`
	Short string `json:"short"`
}

type PartyRrPenalty struct {
	PartyID string  `json:"party_id"`
	Penalty float64 `json:"penalty"`
}

type Agent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlayerStats struct {
	KD           float64 `json:"kd"`
	HSPercent    float64 `json:"hs_percent"`
	WinRate      float64 `json:"win_rate"`
	AvgDamage    float64 `json:"avg_damage"`
	AccountLevel int     `json:"account_level"`
	Rank         string  `json:"rank"`
}

type Damage struct {
	Dealt    int `json:"dealt"`
	Received int `json:"received"`
}

type AbilityCasts struct {
	Grenade  int `json:"grenade"`
	Ability1 int `json:"ability1"`
	Ability2 int `json:"ability2"`
	Ultimate int `json:"ultimate"`
}

type Tier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Customization struct {
	Card                 string `json:"card"`
	Title                string `json:"title"`
	PreferredLevelBorder string `json:"preferred_level_border"`
}

type Behavior struct {
	AfkRounds     int          `json:"afk_rounds"`
	FriendlyFire  FriendlyFire `json:"friendly_fire"`
	RoundsInSpawn int          `json:"rounds_in_spawn"`
}

type FriendlyFire struct {
	Incoming int `json:"incoming"`
	Outgoing int `json:"outgoing"`
}

type Economy struct {
	Spent        Spent        `json:"spent"`
	LoadoutValue LoadoutValue `json:"loadout_value"`
}

type Spent struct {
	Overall float64 `json:"overall"`
	Average float64 `json:"average"`
}

type LoadoutValue struct {
	Overall float64 `json:"overall"`
	Average float64 `json:"average"`
}

type Round struct {
	Number       int                `json:"number"`
	WinningTeam  string             `json:"winning_team"`
	EndReason    string             `json:"end_reason"`
	BombPlanted  bool               `json:"bomb_planted"`
	BombDefused  bool               `json:"bomb_defused"`
	PlantEvents  []PlantEvent       `json:"plant_events"`
	DefuseEvents []DefuseEvent      `json:"defuse_events"`
	PlayerStats  []RoundPlayerStats `json:"player_stats"`
	Economy      RoundEconomy       `json:"economy"`
	DamageEvents []DamageEvent      `json:"damage_events"`
}

type PlantEvent struct {
	PlantLocation Location   `json:"plant_location"`
	PlantedBy     PlayerInfo `json:"planted_by"`
	PlantSite     string     `json:"plant_site"`
	PlantTimeInMs int        `json:"plant_time_in_ms"`
}

type DefuseEvent struct {
	DefuseLocation Location   `json:"defuse_location"`
	DefusedBy      PlayerInfo `json:"defused_by"`
	DefuseTimeInMs int        `json:"defuse_time_in_ms"`
}

type Location struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerInfo struct {
	Puuid string `json:"puuid"`
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Team  string `json:"team"`
}

type RoundPlayerStats struct {
	Player          PlayerInfo         `json:"player"`
	AbilityCasts    AbilityCasts       `json:"ability_casts"`
	DamageEvents    []DamageEvent      `json:"damage_events"`
	Stats           PlayerStatsSummary `json:"stats"`
	Economy         RoundPlayerEconomy `json:"economy"`
	WasAfk          bool               `json:"was_afk"`
	ReceivedPenalty bool               `json:"received_penalty"`
	StayedInSpawn   bool               `json:"stayed_in_spawn"`
}

type DamageEvent struct {
	Player    PlayerInfo `json:"player"`
	Bodyshots int        `json:"bodyshots"`
	Headshots int        `json:"headshots"`
	Legshots  int        `json:"legshots"`
	Damage    int        `json:"damage"`
}

type PlayerStatsSummary struct {
	Score     int `json:"score"`
	Kills     int `json:"kills"`
	Headshots int `json:"headshots"`
	Bodyshots int `json:"bodyshots"`
	Legshots  int `json:"legshots"`
}

type RoundPlayerEconomy struct {
	LoadoutValue int    `json:"loadout_value"`
	Remaining    int    `json:"remaining"`
	Weapon       Weapon `json:"weapon"`
	Armor        Armor  `json:"armor"`
}

type Weapon struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Armor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoundEconomy struct {
	LoadoutValue int    `json:"loadout_value"`
	Weapon       Weapon `json:"weapon"`
	Armor        Armor  `json:"armor"`
}

type Team struct {
	TeamID     string `json:"team_id"`
	Won        bool   `json:"won"`
	RoundsWon  int    `json:"rounds_won"`
	RoundsLost int    `json:"rounds_lost"`
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

type PlayerFeatures struct {
	KD                      float64 `json:"kd"`
	HSPercent               float64 `json:"hs_percent"`
	WinRate                 float64 `json:"win_rate"`
	AvgDamageDealt          float64 `json:"avg_damage_dealt"`
	AvgScore                float64 `json:"avg_score"`
	AccountLevel            int     `json:"account_level"`
	RankNum                 int     `json:"rank_num"`
	SD_KD                   float64 `json:"sd_kd"`
	SD_HSPercent            float64 `json:"sd_hs_percent"`
	SD_WinRate              float64 `json:"sd_win_rate"`
	AvgAbilityCastsPerRound float64 `json:"avg_ability_casts_per_round"`
	RankProgression         int     `json:"rank_progression"`
}

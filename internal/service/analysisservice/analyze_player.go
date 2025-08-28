package analysisservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/edcnogueira/valoguard-api/internal/models"
)

var rankMap = map[string]int{
	"Iron 1":      3,
	"Iron 2":      4,
	"Iron 3":      5,
	"Bronze 1":    6,
	"Bronze 2":    7,
	"Bronze 3":    8,
	"Silver 1":    9,
	"Silver 2":    10,
	"Silver 3":    11,
	"Gold 1":      12,
	"Gold 2":      13,
	"Gold 3":      14,
	"Platinum 1":  15,
	"Platinum 2":  16,
	"Platinum 3":  17,
	"Diamond 1":   18,
	"Diamond 2":   19,
	"Diamond 3":   20,
	"Ascendant 1": 21,
	"Ascendant 2": 22,
	"Ascendant 3": 23,
	"Immortal 1":  24,
	"Immortal 2":  25,
	"Immortal 3":  26,
	"Radiant":     27,
}

func (s *Service) AnalyzePlayer(ctx context.Context, req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	region := req.Region
	if region == "" {
		region = "br"
	}

	account, err := s.client.GetAccount(ctx, req.Name, req.Tag)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}

	mmr, err := s.client.GetMMR(ctx, region, req.Name, req.Tag)
	if err != nil {
		return nil, fmt.Errorf("error getting MMR: %w", err)
	}

	matches, err := s.client.GetMatches(ctx, region, req.Name, req.Tag)
	if err != nil {
		return nil, fmt.Errorf("error getting matches: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("nenhuma partida encontrada")
	}

	stats, cheatScore := s.calculateStatistics(matches, account, mmr)

	matchSummaries := s.createMatchSummaries(matches, account.Puuid)

	response := &models.AnalyzeResponse{
		Probability:  cheatScore,
		Stats:        stats,
		MatchSummary: matchSummaries,
	}

	return response, nil
}

func (s *Service) calculateStatistics(matches []models.Match, account *models.Account, mmr *models.MMR) (models.PlayerStats, int) {
	var totalKills, totalDeaths, totalHead, totalBody, totalLeg, totalScore, totalDamage, wins int

	for _, match := range matches {
		for _, player := range match.Players {
			if player.Puuid == account.Puuid {
				totalKills += player.Stats.Kills
				totalDeaths += player.Stats.Deaths
				totalHead += player.Stats.Head
				totalBody += player.Stats.Body
				totalLeg += player.Stats.Leg
				totalScore += player.Stats.Score
				totalDamage += player.Stats.Damage.Dealt

				for _, team := range match.Teams {
					if team.TeamID == player.Team && team.Won {
						wins++
						break
					}
				}
				break
			}
		}
	}

	matchCount := len(matches)

	kd := float64(totalKills) / float64(totalDeaths+1) // +1 to avoid division by zero
	hsPercent := float64(0)
	if totalHead+totalBody+totalLeg > 0 {
		hsPercent = float64(totalHead) / float64(totalHead+totalBody+totalLeg) * 100
	}
	winRate := float64(wins) / float64(matchCount) * 100
	avgDamage := float64(totalDamage) / float64(matchCount)

	stats := models.PlayerStats{
		KD:           kd,
		HSPercent:    hsPercent,
		WinRate:      winRate,
		AvgDamage:    avgDamage,
		AccountLevel: account.AccountLevel,
		Rank:         mmr.CurrentTierPatched,
	}

	cheatScore := 0
	if hsPercent > 35 {
		cheatScore += 30
	}
	if kd > 2.5 {
		cheatScore += 20
	}
	if winRate > 70 {
		cheatScore += 20
	}

	rankNum, exists := rankMap[strings.TrimSpace(mmr.CurrentTierPatched)]
	if exists && account.AccountLevel < 50 && rankNum > 12 { // > Gold 1
		cheatScore += 30
	}

	return stats, cheatScore
}

func (s *Service) createMatchSummaries(matches []models.Match, puuid string) []models.MatchSummary {
	matchSummaries := []models.MatchSummary{}

	for _, match := range matches {
		score := s.findPlayerScore(match, puuid)

		matchSummaries = append(matchSummaries, models.MatchSummary{
			MatchID: match.Metadata.MatchID,
			Score:   score,
		})
	}

	return matchSummaries
}

func (s *Service) findPlayerScore(match models.Match, puuid string) int {
	for _, player := range match.Players {
		if player.Puuid == puuid {
			return player.Stats.Score
		}
	}
	return 0
}

package analysisservice

import (
	"context"
	"fmt"
	"math"
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

	features, err := s.calculateFeatures(matches, account, mmr)
	if err != nil {
		return nil, fmt.Errorf("error calculating features: %w", err)
	}

	probability, err := s.predictProbability(features)
	if err != nil {
		return nil, fmt.Errorf("error predicting probability: %w", err)
	}

	matchSummaries := s.createMatchSummaries(matches, account.Puuid)

	response := &models.AnalyzeResponse{
		Probability: int(probability * 100),
		Stats: models.PlayerStats{
			KD:           features.KD,
			HSPercent:    features.HSPercent,
			WinRate:      features.WinRate,
			AvgDamage:    features.AvgDamageDealt,
			AccountLevel: features.AccountLevel,
			Rank:         mmr.CurrentData.CurrentTierPatched,
		},
		MatchSummary: matchSummaries,
	}

	return response, nil
}

func (s *Service) calculateFeatures(matches []models.Match, account *models.Account, mmr *models.MMR) (models.PlayerFeatures, error) {
	var totalKills, totalDeaths, totalAssists, totalHeadshots, totalBodyshots, totalLegshots, totalScore, totalDamageDealt, _, totalRoundsPlayed, totalAbilityCasts int
	var wins int

	var kds, hsPercents, winRates []float64

	for _, match := range matches {
		for _, player := range match.Players {
			if player.Puuid == account.Puuid {
				totalKills += player.Stats.Kills
				totalDeaths += player.Stats.Deaths
				totalAssists += player.Stats.Assists
				totalHeadshots += player.Stats.Head
				totalBodyshots += player.Stats.Body
				totalLegshots += player.Stats.Leg
				totalScore += player.Stats.Score
				totalDamageDealt += player.Stats.Damage.Dealt

				matchKD := float64(player.Stats.Kills) / float64(player.Stats.Deaths+1)
				matchHSPercent := float64(0)
				if player.Stats.Head+player.Stats.Body+player.Stats.Leg > 0 {
					matchHSPercent = float64(player.Stats.Head) / float64(player.Stats.Head+player.Stats.Body+player.Stats.Leg) * 100
				}
				kds = append(kds, matchKD)
				hsPercents = append(hsPercents, matchHSPercent)

				matchWon := false
				for _, team := range match.Teams {
					if team.TeamID == player.Team && team.Won {
						matchWon = true
						wins++
						break
					}
				}
				if matchWon {
					winRates = append(winRates, 1.0)
				} else {
					winRates = append(winRates, 0.0)
				}

				for _, round := range match.Rounds {
					for _, roundPlayerStats := range round.PlayerStats {
						if roundPlayerStats.Player.Puuid == account.Puuid {
							totalAbilityCasts += roundPlayerStats.AbilityCasts.Grenade + roundPlayerStats.AbilityCasts.Ability1 + roundPlayerStats.AbilityCasts.Ability2 + roundPlayerStats.AbilityCasts.Ultimate
							totalRoundsPlayed++
							break
						}
					}
				}
				break
			}
		}
	}

	matchCount := len(matches)
	if matchCount == 0 {
		return models.PlayerFeatures{}, fmt.Errorf("nenhuma partida encontrada para calcular features")
	}

	kd := float64(totalKills) / float64(totalDeaths+1)
	hsPercent := float64(0)
	if totalHeadshots+totalBodyshots+totalLegshots > 0 {
		hsPercent = float64(totalHeadshots) / float64(totalHeadshots+totalBodyshots+totalLegshots) * 100
	}
	winRate := float64(wins) / float64(matchCount) * 100
	avgDamageDealt := float64(totalDamageDealt) / float64(matchCount)
	avgScore := float64(totalScore) / float64(matchCount)

	sdKD := calculateStandardDeviation(kds)
	sdHSPercent := calculateStandardDeviation(hsPercents)
	sdWinRate := calculateStandardDeviation(winRates)

	avgAbilityCastsPerRound := float64(0)
	if totalRoundsPlayed > 0 {
		avgAbilityCastsPerRound = float64(totalAbilityCasts) / float64(totalRoundsPlayed)
	}

	rankNum, exists := rankMap[strings.TrimSpace(mmr.CurrentData.CurrentTierPatched)]
	rankProgression := 0
	if exists && mmr.BySeason != nil {
	}

	features := models.PlayerFeatures{
		KD:                      kd,
		HSPercent:               hsPercent,
		WinRate:                 winRate,
		AvgDamageDealt:          avgDamageDealt,
		AvgScore:                avgScore,
		AccountLevel:            account.AccountLevel,
		RankNum:                 rankNum,
		SD_KD:                   sdKD,
		SD_HSPercent:            sdHSPercent,
		SD_WinRate:              sdWinRate,
		AvgAbilityCastsPerRound: avgAbilityCastsPerRound,
		RankProgression:         rankProgression,
	}

	return features, nil
}

func calculateStandardDeviation(data []float64) float64 {
	if len(data) < 2 {
		return 0.0
	}

	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(len(data))

	variance := 0.0
	for _, v := range data {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(data) - 1)

	return math.Sqrt(variance)
}

func (s *Service) predictProbability(features models.PlayerFeatures) (float64, error) {
	probability := 0.0

	if features.HSPercent > 40 && features.SD_HSPercent < 5 {
		probability += 0.3
	}
	if features.KD > 3.0 && features.SD_KD < 0.5 {
		probability += 0.2
	}
	if features.WinRate > 80 && features.SD_WinRate < 10 {
		probability += 0.2
	}
	if features.AccountLevel < 50 && features.RankNum > 15 {
		probability += 0.3
	}

	if probability > 1.0 {
		probability = 1.0
	}

	return probability, nil
}

func (s *Service) createMatchSummaries(matches []models.Match, puuid string) []models.MatchSummary {
	var matchSummaries []models.MatchSummary

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

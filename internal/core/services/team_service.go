package services

import (
	"errors"
	"soccer-service/internal/core/models"
	"soccer-service/internal/repository"
	"sort"
	"strconv"
	"strings"
)

type TeamService struct {
	playerRepo *repository.PlayerRepository
	matchRepo  *repository.MatchRepository
	eventRepo  *repository.EventRepository
}

func NewTeamService(playerRepo *repository.PlayerRepository, matchRepo *repository.MatchRepository, eventRepo *repository.EventRepository) *TeamService {
	return &TeamService{
		playerRepo: playerRepo,
		matchRepo:  matchRepo,
		eventRepo:  eventRepo,
	}
}

type playerScore struct {
	player models.Player
	score  float64
}

func (s *TeamService) GenerateBalancedTeams(players []models.Player) (models.TeamComposition, models.TeamComposition, error) {
	if len(players) < 2 {
		return models.TeamComposition{}, models.TeamComposition{}, errors.New("not enough players")
	}

	// Count players by gender
	maleCount, femaleCount := countPlayersByGender(players)
	if maleCount == 0 || femaleCount == 0 {
		// Single gender match, use regular balancing
		return s.generateRegularTeams(players)
	}

	// Mixed match, ensure gender balance
	return s.generateGenderBalancedTeams(players, maleCount, femaleCount)
}

func (s *TeamService) generateGenderBalancedTeams(players []models.Player, maleCount, femaleCount int) (models.TeamComposition, models.TeamComposition, error) {
	// Separate players by gender
	malePlayersByPos := make(map[models.Position][]models.Player)
	femalePlayersByPos := make(map[models.Position][]models.Player)

	for _, p := range players {
		if p.Gender == models.Male {
			malePlayersByPos[p.Position] = append(malePlayersByPos[p.Position], p)
		} else {
			femalePlayersByPos[p.Position] = append(femalePlayersByPos[p.Position], p)
		}
	}

	// Calculate optimal formation
	formation := determineOptimalFormation(malePlayersByPos, femalePlayersByPos)

	// Sort players by comprehensive score within their groups
	for _, group := range malePlayersByPos {
		sortPlayersByScore(group)
	}
	for _, group := range femalePlayersByPos {
		sortPlayersByScore(group)
	}

	// Initialize teams
	var team1, team2 models.TeamComposition
	team1.Formation = formation
	team2.Formation = formation

	// Distribute players ensuring gender and position balance
	team1.Players, team2.Players = distributePlayersWithGenderBalance(
		malePlayersByPos,
		femalePlayersByPos,
		formation,
		maleCount/2,
		femaleCount/2,
	)

	// Calculate final ratings and chemistry
	team1.AvgRating = calculateTeamRating(team1.Players, players)
	team2.AvgRating = calculateTeamRating(team2.Players, players)
	team1.Chemistry = s.calculateTeamChemistry(getPlayersFromRoles(team1.Players, players))
	team2.Chemistry = s.calculateTeamChemistry(getPlayersFromRoles(team2.Players, players))

	return team1, team2, nil
}

func countPlayersByGender(players []models.Player) (males, females int) {
	for _, p := range players {
		if p.Gender == models.Male {
			males++
		} else {
			females++
		}
	}
	return
}

func calculatePlayerScore(p models.Player, forPosition models.Position) float64 {
	// Base score from stats
	baseScore := p.Stats.AverageRating * 0.4 // Reduced weight of general stats

	// Calculate position-specific ability score
	abilityScore := calculatePositionAbilityScore(p, forPosition)

	// Calculate physical condition score
	physicalScore := calculatePhysicalScore(p)

	// Calculate mental/tactical score
	tacticalScore := calculateTacticalScore(p)

	// Apply versatility bonus if playing in secondary position
	versatilityBonus := calculateVersatilityBonus(p, forPosition)

	// Weighted combination of all scores
	return (baseScore * 0.35) + (abilityScore * 0.3) + (physicalScore * 0.15) +
		(tacticalScore * 0.15) + (versatilityBonus * 0.05)
}

func calculateVersatilityBonus(p models.Player, forPosition models.Position) float64 {
	if p.Position == forPosition {
		return 1.0 // Playing in primary position
	}

	// Check if this is a secondary position
	for _, pos := range p.Versatility {
		if pos == forPosition {
			return 0.8 // Playing in secondary position
		}
	}

	return 0.5 // Playing out of position
}

func calculatePositionAbilityScore(p models.Player, pos models.Position) float64 {
	totalScore := 0.0
	totalWeight := 0.0

	for _, ability := range p.Abilities {
		if info, exists := models.CommonAbilities[ability.Name]; exists {
			weight := info.Weights[pos]
			score := models.AbilityLevelWeight(ability.Level)

			totalScore += weight * score
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.5 // Default score if no relevant abilities found
	}

	return totalScore / totalWeight
}

func calculatePhysicalScore(p models.Player) float64 {
	physicalAbilities := []string{"Speed", "Stamina", "Strength", "Agility", "Jump"}

	totalScore := 0.0
	count := 0

	for _, ability := range p.Abilities {
		for _, physicalAbility := range physicalAbilities {
			if ability.Name == physicalAbility {
				totalScore += models.AbilityLevelWeight(ability.Level)
				count++
				break
			}
		}
	}

	if count == 0 {
		return 0.5
	}

	return totalScore / float64(count)
}

func calculateTacticalScore(p models.Player) float64 {
	tacticalAbilities := []string{
		"Game Reading", "Decision Making", "Teamwork",
		"Communication", "Positioning",
	}

	totalScore := 0.0
	count := 0

	for _, ability := range p.Abilities {
		for _, tacticalAbility := range tacticalAbilities {
			if ability.Name == tacticalAbility {
				totalScore += models.AbilityLevelWeight(ability.Level)
				count++
				break
			}
		}
	}

	if count == 0 {
		return 0.5
	}

	return totalScore / float64(count)
}

func sortPlayersByScore(players []models.Player) {
	sort.Slice(players, func(i, j int) bool {
		scoreI := calculatePlayerScore(players[i], players[i].Position)
		scoreJ := calculatePlayerScore(players[j], players[j].Position)
		return scoreI > scoreJ
	})
}

func distributePlayersWithGenderBalance(
	malePlayersByPos map[models.Position][]models.Player,
	femalePlayersByPos map[models.Position][]models.Player,
	formation string,
	targetMalesPerTeam, targetFemalesPerTeam int,
) ([]models.PlayerRole, []models.PlayerRole) {
	var team1, team2 []models.PlayerRole
	team1Males, team2Males := 0, 0
	team1Females, team2Females := 0, 0

	// Helper function to add player to team
	addToTeam := func(p models.Player, preferTeam1 bool) {
		role := models.PlayerRole{
			PlayerID: p.ID,
			Position: p.Position,
			Role:     determineRole(p.Position),
		}

		// Determine which team needs this player more
		addToTeam1 := preferTeam1
		if p.Gender == models.Male {
			if team1Males >= targetMalesPerTeam {
				addToTeam1 = false
			} else if team2Males >= targetMalesPerTeam {
				addToTeam1 = true
			}
		} else {
			if team1Females >= targetFemalesPerTeam {
				addToTeam1 = false
			} else if team2Females >= targetFemalesPerTeam {
				addToTeam1 = true
			}
		}

		if addToTeam1 {
			team1 = append(team1, role)
			if p.Gender == models.Male {
				team1Males++
			} else {
				team1Females++
			}
		} else {
			team2 = append(team2, role)
			if p.Gender == models.Male {
				team2Males++
			} else {
				team2Females++
			}
		}
	}

	// Distribute players by position while maintaining gender balance
	positions := []models.Position{models.Goalkeeper, models.Defender, models.Midfielder, models.Forward}
	for _, pos := range positions {
		// Distribute male players
		for _, p := range malePlayersByPos[pos] {
			addToTeam(p, len(team1) <= len(team2))
		}
		// Distribute female players
		for _, p := range femalePlayersByPos[pos] {
			addToTeam(p, len(team1) <= len(team2))
		}
	}

	return team1, team2
}

func (s *TeamService) CalculatePlayerAffinity(player1ID, player2ID string) (float64, error) {
	// Get historical match data
	matches, err := s.matchRepo.GetMatchesWithPlayers(player1ID, player2ID)
	if err != nil {
		return 0, err
	}

	if len(matches) == 0 {
		return 0.5, nil // Neutral affinity for no history
	}

	var totalScore float64
	successfulInteractions := 0
	totalInteractions := 0

	for _, match := range matches {
		// Get events involving both players
		events, err := s.eventRepo.GetByMatch(match.ID)
		if err != nil {
			continue
		}

		for _, event := range events {
			if isPlayerInteraction(event, player1ID, player2ID) {
				totalInteractions++
				if event.Success {
					successfulInteractions++
					totalScore += event.Impact
				}
			}
		}
	}

	if totalInteractions == 0 {
		return 0.5, nil
	}

	// Calculate affinity based on success rate and impact
	successRate := float64(successfulInteractions) / float64(totalInteractions)
	avgImpact := totalScore / float64(successfulInteractions)

	return (successRate*0.6 + avgImpact*0.4), nil
}

func (s *TeamService) UpdatePlayerAffinities(matchID string) error {
	events, err := s.eventRepo.GetByMatch(matchID)
	if err != nil {
		return err
	}

	// Track player interactions
	interactions := make(map[string]map[string]*InteractionStats)

	// Analyze events for player interactions
	for _, event := range events {
		if event.TargetPlayerID == "" {
			continue
		}

		// Initialize interaction tracking if needed
		if _, exists := interactions[event.PlayerID]; !exists {
			interactions[event.PlayerID] = make(map[string]*InteractionStats)
		}
		if _, exists := interactions[event.PlayerID][event.TargetPlayerID]; !exists {
			interactions[event.PlayerID][event.TargetPlayerID] = &InteractionStats{}
		}

		// Update interaction stats
		stats := interactions[event.PlayerID][event.TargetPlayerID]
		stats.TotalInteractions++
		if event.Success {
			stats.SuccessfulInteractions++
			stats.TotalImpact += event.Impact
		}
	}

	// Update affinities in the database
	for player1ID, playerInteractions := range interactions {
		for player2ID, stats := range playerInteractions {
			affinity := calculateAffinityScore(stats)
			if err := s.playerRepo.UpdateAffinity(player1ID, player2ID, affinity); err != nil {
				return err
			}
		}
	}

	return nil
}

type InteractionStats struct {
	TotalInteractions      int
	SuccessfulInteractions int
	TotalImpact            float64
}

func calculateAffinityScore(stats *InteractionStats) float64 {
	if stats.TotalInteractions == 0 {
		return 0.5
	}

	successRate := float64(stats.SuccessfulInteractions) / float64(stats.TotalInteractions)
	avgImpact := stats.TotalImpact / float64(stats.SuccessfulInteractions)

	return (successRate*0.6 + avgImpact*0.4)
}

func groupPlayersByPosition(players []models.Player) map[models.Position][]models.Player {
	groups := make(map[models.Position][]models.Player)
	for _, player := range players {
		groups[player.Position] = append(groups[player.Position], player)
	}
	return groups
}

func determineOptimalFormation(malePlayersByPos, femalePlayersByPos map[models.Position][]models.Player) string {
	// Combine player counts from both genders
	totalByPosition := make(map[models.Position]int)
	for pos := range malePlayersByPos {
		totalByPosition[pos] = len(malePlayersByPos[pos]) + len(femalePlayersByPos[pos])
	}

	// Common formations
	formations := []string{"4-4-2", "4-3-3", "3-5-2", "4-5-1", "3-4-3"}

	bestFormation := formations[0]
	bestScore := -1.0

	for _, formation := range formations {
		score := evaluateFormationFit(formation, totalByPosition)
		if score > bestScore {
			bestScore = score
			bestFormation = formation
		}
	}

	return bestFormation
}

func evaluateFormationFit(formation string, playersByPos map[models.Position]int) float64 {
	required := parseFormation(formation)

	// Calculate how well the available players match the formation requirements
	var score float64
	for pos, count := range required {
		available := playersByPos[pos]
		if available < count {
			// Penalize not having enough players
			score -= float64(count - available)
		} else if available > count {
			// Small penalty for having too many players
			score -= float64(available-count) * 0.5
		} else {
			// Bonus for perfect match
			score += 1.0
		}
	}

	return score
}

func parseFormation(formation string) map[models.Position]int {
	parts := strings.Split(formation, "-")
	required := make(map[models.Position]int)

	// Always need one goalkeeper
	required[models.Goalkeeper] = 1

	if len(parts) >= 3 {
		defenders, _ := strconv.Atoi(parts[0])
		midfielders, _ := strconv.Atoi(parts[1])
		forwards, _ := strconv.Atoi(parts[2])

		required[models.Defender] = defenders
		required[models.Midfielder] = midfielders
		required[models.Forward] = forwards
	}

	return required
}

func distributePlayersOptimally(groups map[models.Position][]models.Player, formation string) ([]models.PlayerRole, []models.PlayerRole) {
	var team1, team2 []models.PlayerRole
	formationReqs := parseFormation(formation)

	// First, distribute key positions (goalkeeper and central positions)
	distributeKeyPositions(&team1, &team2, groups, formationReqs)

	// Then distribute remaining players considering versatility and complementary abilities
	distributeRemainingPlayers(&team1, &team2, groups, formationReqs)

	return team1, team2
}

func distributeKeyPositions(team1, team2 *[]models.PlayerRole, groups map[models.Position][]models.Player, formationReqs map[models.Position]int) {
	// Distribute goalkeepers first
	distributePositionPlayers(team1, team2, groups[models.Goalkeeper], models.Goalkeeper, formationReqs[models.Goalkeeper])

	// Distribute central defenders and midfielders
	distributePositionPlayers(team1, team2, groups[models.Defender], models.Defender, formationReqs[models.Defender])
	distributePositionPlayers(team1, team2, groups[models.Midfielder], models.Midfielder, formationReqs[models.Midfielder])
}

func distributeRemainingPlayers(team1, team2 *[]models.PlayerRole, groups map[models.Position][]models.Player, formationReqs map[models.Position]int) {
	// Create pools of remaining players by primary and secondary positions
	remainingPlayers := getAllRemainingPlayers(groups)

	// Fill remaining positions considering versatility
	for pos, required := range formationReqs {
		team1Required := required/2 - countPositionInTeam(*team1, pos)
		team2Required := required/2 - countPositionInTeam(*team2, pos)

		if team1Required > 0 || team2Required > 0 {
			candidates := findCandidatesForPosition(remainingPlayers, pos)
			sortPlayersByPositionScore(candidates, pos)

			// Distribute to teams based on needs and scores
			for _, player := range candidates {
				if team1Required > 0 && (team1Required >= team2Required || len(*team1) < len(*team2)) {
					addPlayerToTeam(team1, player, pos)
					team1Required--
					removePlayerFromPool(&remainingPlayers, player)
				} else if team2Required > 0 {
					addPlayerToTeam(team2, player, pos)
					team2Required--
					removePlayerFromPool(&remainingPlayers, player)
				}
			}
		}
	}

	// Distribute any remaining players to balance team sizes
	distributeRemainingToBalance(team1, team2, remainingPlayers)
}

func getAllRemainingPlayers(groups map[models.Position][]models.Player) []models.Player {
	var players []models.Player
	for _, posPlayers := range groups {
		players = append(players, posPlayers...)
	}
	return players
}

func findCandidatesForPosition(players []models.Player, pos models.Position) []models.Player {
	var candidates []models.Player
	for _, p := range players {
		if p.Position == pos {
			candidates = append(candidates, p)
			continue
		}
		// Check secondary positions
		for _, secPos := range p.Versatility {
			if secPos == pos {
				candidates = append(candidates, p)
				break
			}
		}
	}
	return candidates
}

func sortPlayersByPositionScore(players []models.Player, pos models.Position) {
	sort.Slice(players, func(i, j int) bool {
		scoreI := calculatePlayerScore(players[i], pos)
		scoreJ := calculatePlayerScore(players[j], pos)
		return scoreI > scoreJ
	})
}

func countPositionInTeam(team []models.PlayerRole, pos models.Position) int {
	count := 0
	for _, role := range team {
		if role.Position == pos {
			count++
		}
	}
	return count
}

func removePlayerFromPool(pool *[]models.Player, player models.Player) {
	for i, p := range *pool {
		if p.ID == player.ID {
			*pool = append((*pool)[:i], (*pool)[i+1:]...)
			return
		}
	}
}

func distributeRemainingToBalance(team1, team2 *[]models.PlayerRole, players []models.Player) {
	for _, player := range players {
		if len(*team1) <= len(*team2) {
			addPlayerToTeam(team1, player, player.Position)
		} else {
			addPlayerToTeam(team2, player, player.Position)
		}
	}
}

func distributePositionPlayers(team1, team2 *[]models.PlayerRole, players []models.Player, pos models.Position, required int) {
	if len(players) < required {
		return
	}

	// Find players with best relevant abilities for central positions
	sort.Slice(players, func(i, j int) bool {
		scoreI := calculatePositionAbilityScore(players[i], pos)
		scoreJ := calculatePositionAbilityScore(players[j], pos)
		return scoreI > scoreJ
	})

	// Distribute top players evenly
	for i := 0; i < required; i++ {
		addPlayerToTeam(team1, players[i], pos)
		addPlayerToTeam(team2, players[i], pos)
	}
}

func addPlayerToTeam(team *[]models.PlayerRole, player models.Player, pos models.Position) {
	*team = append(*team, models.PlayerRole{
		PlayerID: player.ID,
		Position: pos,
		Role:     determineRole(pos),
	})
}

func isPlayerDistributed(player models.Player, team []models.PlayerRole) bool {
	for _, role := range team {
		if role.PlayerID == player.ID {
			return true
		}
	}
	return false
}

func getTeamAbilities(team []models.PlayerRole) []models.Ability {
	var abilities []models.Ability
	// This function should be implemented to get abilities of all players in the team
	// You'll need to have access to the full player data
	return abilities
}

func getPlayersFromRoles(roles []models.PlayerRole, allPlayers []models.Player) []models.Player {
	var result []models.Player
	playerMap := make(map[string]models.Player)
	for _, p := range allPlayers {
		playerMap[p.ID] = p
	}

	for _, role := range roles {
		if player, exists := playerMap[role.PlayerID]; exists {
			result = append(result, player)
		}
	}

	return result
}

func (s *TeamService) calculateTeamChemistry(players []models.Player) float64 {
	if len(players) < 2 {
		return 1.0
	}

	var totalAffinity float64
	connections := 0

	// Calculate affinity between all player pairs
	for i := 0; i < len(players)-1; i++ {
		for j := i + 1; j < len(players); j++ {
			affinity, err := s.CalculatePlayerAffinity(players[i].ID, players[j].ID)
			if err == nil {
				totalAffinity += affinity
				connections++
			}
		}
	}

	if connections == 0 {
		return 0.5
	}

	return totalAffinity / float64(connections)
}

func isPlayerInteraction(event models.Event, player1ID, player2ID string) bool {
	return (event.PlayerID == player1ID && event.TargetPlayerID == player2ID) ||
		(event.PlayerID == player2ID && event.TargetPlayerID == player1ID)
}

func determineRole(pos models.Position) string {
	switch pos {
	case models.Forward:
		return "striker"
	case models.Midfielder:
		return "central_mid"
	case models.Defender:
		return "center_back"
	case models.Goalkeeper:
		return "goalkeeper"
	default:
		return "undefined"
	}
}

func (s *TeamService) generateRegularTeams(players []models.Player) (models.TeamComposition, models.TeamComposition, error) {
	// Group players by position
	positionGroups := groupPlayersByPosition(players)

	// For single-gender teams, we'll use an empty map for the other gender
	emptyMap := make(map[models.Position][]models.Player)

	// Calculate optimal formation based on available players
	formation := determineOptimalFormation(positionGroups, emptyMap)

	// Sort players by rating within their positions
	for _, group := range positionGroups {
		sort.Slice(group, func(i, j int) bool {
			return calculatePlayerScore(group[i], group[i].Position) > calculatePlayerScore(group[j], group[j].Position)
		})
	}

	// Initialize teams
	var team1, team2 models.TeamComposition
	team1.Formation = formation
	team2.Formation = formation

	// Distribute players ensuring position balance and chemistry
	team1.Players, team2.Players = distributePlayersOptimally(positionGroups, formation)

	// Calculate team ratings and chemistry
	team1.AvgRating = calculateTeamRating(team1.Players, players)
	team2.AvgRating = calculateTeamRating(team2.Players, players)
	team1.Chemistry = s.calculateTeamChemistry(getPlayersFromRoles(team1.Players, players))
	team2.Chemistry = s.calculateTeamChemistry(getPlayersFromRoles(team2.Players, players))

	return team1, team2, nil
}

func calculateTeamRating(roles []models.PlayerRole, allPlayers []models.Player) float64 {
	if len(roles) == 0 {
		return 0
	}

	var totalRating float64
	playerMap := make(map[string]models.Player)
	for _, p := range allPlayers {
		playerMap[p.ID] = p
	}

	for _, role := range roles {
		if player, exists := playerMap[role.PlayerID]; exists {
			// Use the enhanced player scoring for team rating
			totalRating += calculatePlayerScore(player, role.Position)
		}
	}

	return totalRating / float64(len(roles))
}

func areComplementaryAbilities(ability1, ability2 string) bool {
	complementaryPairs := map[string][]string{
		"Speed":         {"Game Reading", "Ball Control", "Positioning"},
		"Short Passing": {"Game Reading", "First Touch", "Positioning"},
		"Strength":      {"Agility", "Ball Control", "Speed"},
		"Tackling":      {"Positioning", "Game Reading", "Speed"},
		"Dribbling":     {"Speed", "Agility", "Decision Making"},
		"Ball Control":  {"Speed", "Agility", "First Touch"},
		"First Touch":   {"Ball Control", "Short Passing", "Dribbling"},
		"Game Reading":  {"Decision Making", "Positioning", "Interception"},
		"Communication": {"Teamwork", "Game Reading", "Positioning"},
		"Heading":       {"Jump", "Positioning", "Strength"},
	}

	if complements, exists := complementaryPairs[ability1]; exists {
		for _, complement := range complements {
			if complement == ability2 {
				return true
			}
		}
	}

	return false
}

func calculateTeamFit(player models.Player, team []models.PlayerRole) float64 {
	// Base fit score
	fit := 0.0

	// Consider complementary abilities
	teamAbilities := getTeamAbilities(team)
	playerAbilities := player.Abilities

	// Reward diverse ability combinations
	for _, playerAbility := range playerAbilities {
		hasComplement := false
		for _, teamAbility := range teamAbilities {
			if areComplementaryAbilities(playerAbility.Name, teamAbility.Name) {
				hasComplement = true
				fit += models.AbilityLevelWeight(playerAbility.Level)
			}
		}
		if !hasComplement {
			fit += 0.5 // Bonus for bringing new abilities
		}
	}

	return fit
}

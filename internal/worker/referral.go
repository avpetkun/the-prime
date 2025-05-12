package worker

import "context"

type ReferralBonus struct {
	UserID int64
	Points int64
	Level  int
}

func (s *Service) calculateRefRewards(ctx context.Context, userID, taskPoints int64) (results []ReferralBonus) {
	// Early return if no referral bonus configured or task has no points
	if s.cfg.Referral.TotalBonus == 0 || taskPoints <= 0 || userID <= 0 {
		return
	}

	totalBonusPoints := (taskPoints * s.cfg.Referral.TotalBonus) / 100

	for i, level := range s.cfg.Referral.Levels {
		// Get parent for current user
		parentID, err := s.db.GetUserRefID(ctx, userID)
		if err != nil {
			s.log.Warn().Err(err).Int64("user_id", userID).Msg("[service] failed to get parent for user")
			break
		}
		if parentID <= 0 || parentID == userID {
			break
		}

		bonusPoints := (totalBonusPoints * level) / 100
		if bonusPoints == 0 {
			break
		}

		results = append(results, ReferralBonus{
			UserID: parentID,
			Points: bonusPoints,
			Level:  i + 1,
		})
		userID = parentID
	}
	return
}

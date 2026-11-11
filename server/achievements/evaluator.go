package achievements

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/abhinash-kml/nova/server/messagebroker"
	"go.uber.org/zap"
)

type Evaluator struct {
	service      Service
	statsChannel string
	achievements []Achievement
	mapping      map[int][]Achievement
	ctx          context.Context
	cancel       context.CancelFunc
	logger       *zap.Logger
}

func NewEvaluator(ctx context.Context, s Service, channel string, l *zap.Logger) *Evaluator {
	internalCtx, cancelFunc := context.WithCancel(ctx)
	return &Evaluator{
		service:      s,
		mapping:      make(map[int][]Achievement),
		statsChannel: channel,
		ctx:          internalCtx,
		cancel:       cancelFunc,
		logger:       l,
	}
}

func (e *Evaluator) Initialize() {
	achievements, err := e.service.GetAchievementsForEvaluator(context.Background())
	if err != nil {
		e.logger.Error("Failed to get achievements for evaluator", zap.Error(err))
		return
	}
	e.achievements = append(e.achievements, achievements...)
}

func (e *Evaluator) BuildMapping() {
	e.buildMapping()
}

func (e *Evaluator) Run() {
	messageBroker := messagebroker.GetGlobal()
	eventsCh := messageBroker.Consume(e.statsChannel)
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				e.logger.Info("Achievement evaluator stopped due to context cancellation")
				return
			case event := <-eventsCh:
				userId := event["user_id"].(string)
				statId, _ := strconv.Atoi(event["stat_id"].(string))
				statValue, _ := strconv.Atoi(event["current_value"].(string))

				e.Evaluate(userId, statId, statValue)
			}
		}
	}()
}

func (e *Evaluator) buildMapping() {
	for _, achievement := range e.achievements {
		e.mapping[achievement.StatId] = append(e.mapping[achievement.StatId], achievement)
	}

	fmt.Println(e.mapping)
}

func (e *Evaluator) Evaluate(userid string, statid int, statCurrentvalue int) {
	e.logger.Debug(
		"Incoming values",
		zap.String("user_id", userid),
		zap.Int("stat_id", statid),
		zap.Int("current_value", statCurrentvalue),
	)
	achievements := e.mapping[statid]

	// Fetch user's completed achievements
	completedUserAchievements, err := e.service.GetCompletedAchievementsOfUser(context.Background(), CompletedAchievementRequest{UserId: userid})
	if err != nil {
		return
	}

	// Loop through all the achievents based on incoming stat value
	for _, a := range achievements {
		// If the current achievement is completed skip it
		if slices.Contains(completedUserAchievements, a.AchievementId) {
			e.logger.Debug("Achievement already completed", zap.Int("achievement_id", a.AchievementId))
			return
		}

		// Check if current stat value of achievement is greated than or equal to required stat value, if so complete achievement
		if statCurrentvalue >= a.TargetValue {
			err = e.service.CompleteAchievement(context.Background(), CompleteAchievementRequest{UserId: userid, AchievementId: a.AchievementId})
			if err != nil {
				e.logger.Error("Failed to get completed achievements of user", zap.Error(err))
				return
			}
			e.logger.Debug("User completed achievement", zap.String("user_id", userid), zap.Int("achievement_id", a.AchievementId))

			// Send notification to realtime hub, it will forward further
			messageBroker := messagebroker.GetGlobal()
			messageBroker.Produce("user:notifications", map[string]any{
				"type":           "achievement_complete",
				"user_id":        userid,
				"achievement_id": a.AchievementId,
			})
			e.logger.Debug("Sent achievement completed notification to realtime hub")
		}
	}

}

package achievements

import (
	"context"
	"fmt"
	"strconv"

	"github.com/abhinash-kml/nova/server/events"
)

type Evaluator struct {
	service      Service
	statsChannel string
	achievements []Achievement
	mapping      map[int][]Achievement
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewEvaluator(ctx context.Context, s Service, channel string) *Evaluator {
	internalCtx, cancelFunc := context.WithCancel(ctx)
	return &Evaluator{
		service:      s,
		mapping:      make(map[int][]Achievement),
		statsChannel: channel,
		ctx:          internalCtx,
		cancel:       cancelFunc,
	}
}

func (e *Evaluator) Initialize() {
	achievements, err := e.service.GetAchievementsForEvaluator(context.Background())
	if err != nil {
		fmt.Println("Failed to get achievements for evaluator. Error: %w", err)
		return
	}
	e.achievements = append(e.achievements, achievements...)
}

func (e *Evaluator) BuildMapping() {
	e.buildMapping()
}

func (e *Evaluator) Run() {
	eventsManager := events.GetGlobalEventsManager()
	eventsCh := eventsManager.Consume(e.statsChannel)
	go func() {
		for {
			select {
			case <-e.ctx.Done():
				fmt.Println("Evaluator stopped due to context cancellation")
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
	fmt.Println("STAT_ID:", statid, "USER_ID:", userid, "CURRENT_VALUE:", statCurrentvalue)
	achievements := e.mapping[statid]

	for _, a := range achievements {
		if statCurrentvalue >= a.TargetValue {
			fmt.Printf("UserID: %s completed achievement %s\n", userid, a.Key)
		}
	}

}

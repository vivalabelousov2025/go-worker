package calc

import (
	"context"
	"fmt"
	"time"

	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
)

func CalcTeam(ctx context.Context, teams []dto.Team, resp *dto.Response) (*dto.Team, error) {
	if len(teams) == 0 {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Пустой массив")
	}

	const dateFormat = "2006-01-02"

	var earliestTeam dto.Team
	earliestCompletionTime := time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC)

	for i := range teams {
		team := teams[i]
		fmt.Println(team.NextFreeDate, "team next free date")
		parsedTime, err := time.Parse(dateFormat, team.NextFreeDate)
		if err != nil {
			return nil, err
		}

		if parsedTime.Before(earliestCompletionTime) {
			earliestCompletionTime = parsedTime
			earliestTeam = team
		}
	}

	resp.DateStart = earliestTeam.NextFreeDate
	resp.TeamID = earliestTeam.TeamID

	return &earliestTeam, nil
}

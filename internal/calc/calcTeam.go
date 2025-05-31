package calc

import (
	"context"
	"time"

	"github.com/vivalabelousov2025/go-worker/internal/dto"
	"github.com/vivalabelousov2025/go-worker/pkg/logger"
)

func CalcTeam(ctx context.Context, teams []dto.Team) (*dto.Team, error) {
	if len(teams) == 0 {
		logger.GetLoggerFromCtx(ctx).Fatal(ctx, "Пустой массив")
	}

	var earliestTeam *dto.Team
	earliestCompletionTime := time.Time{}

	for i, v := range teams {
		if i == 0 || v.NextFreeDate.Before(earliestCompletionTime) {
			earliestCompletionTime = v.NextFreeDate
			earliestTeam = &teams[i]
		}
	}

	return earliestTeam, nil
}

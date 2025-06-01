package calc

import (
	"github.com/vivalabelousov2025/go-worker/internal/dto"
)

const basePrice float64 = 10000

func CalcPrice(order *dto.Order, team *dto.Team, hard float64) (float64, error) {
	var price float64

	ratioExpirience := 1 + (float64(team.Experience) / 10)
	ratioMembers := 1 - (float64(team.MembersCount) / 10)

	diff := order.EstimatedEndDate.Sub(order.EstimatedStartDate)
	days := int(diff.Hours() / 24)

	var ratioUrgency float64
	if days > 7 {
		ratioUrgency = 1
	} else if days < 3 {
		ratioUrgency = 2
	} else {
		ratioUrgency = 1.5
	}

	price = basePrice * ratioExpirience * ratioMembers * ratioUrgency * hard

	return price, nil
}

package calc

import (
	"github.com/vivalabelousov2025/go-worker/internal/dto"
)

const basePrice float64 = 10000

func CalcPrice(o *dto.Order, t *dto.Team, hard float64) (float64, error) {
	var price float64

	ratioExpirience := 1 + (float64(t.Experience) / 10)
	ratioMembers := 1 - (float64(t.MembersCount) / 10)

	diff := o.EstimatedEndDate.Sub(o.EstimatedStartDate)
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

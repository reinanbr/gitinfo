package query

import (
	"fmt"
)


func BuildContributionGraphQuery(user string, year int) string {
	start := fmt.Sprintf("%d-01-01T00:00:00Z", year)
	end := fmt.Sprintf("%d-12-31T23:59:59Z", year)
	return fmt.Sprintf(`query { user(login: "%s") { createdAt contributionsCollection(from: "%s", to: "%s") { contributionCalendar { weeks { contributionDays { contributionCount date } } } } } }`, user, start, end)
}

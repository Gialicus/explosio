package tree

import "explosio/lib/domain"

// CountActivities conta il numero totale di attività nell'albero radicato in root.
func CountActivities(root *domain.Activity) int {
	if root == nil {
		return 0
	}
	count := 1
	for _, sub := range root.SubActivities {
		count += CountActivities(sub)
	}
	return count
}

// Walk attraversa l'albero in pre-order chiamando f su ogni attività.
func Walk(root *domain.Activity, f func(*domain.Activity)) {
	if root == nil {
		return
	}
	f(root)
	for _, sub := range root.SubActivities {
		Walk(sub, f)
	}
}

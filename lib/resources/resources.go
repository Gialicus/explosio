package resources

import (
	"explosio/lib/domain"
	"explosio/lib/tree"
)

// ForEachResource chiama f per ogni risorsa (Human, Material, Asset) dell'attività a.
func ForEachResource(a *domain.Activity, f func(domain.Resource)) {
	if a == nil {
		return
	}
	for i := range a.Humans {
		f(a.Humans[i])
	}
	for i := range a.Materials {
		f(a.Materials[i])
	}
	for i := range a.Assets {
		f(a.Assets[i])
	}
}

// WalkResources attraversa l'albero e per ogni attività chiama f(activity, resource) per ogni risorsa.
func WalkResources(root *domain.Activity, f func(*domain.Activity, domain.Resource)) {
	tree.Walk(root, func(a *domain.Activity) {
		ForEachResource(a, func(r domain.Resource) {
			f(a, r)
		})
	})
}

// ResourceDisplayName restituisce un nome di visualizzazione per la risorsa (per messaggi di errore o report).
func ResourceDisplayName(r domain.Resource) string {
	switch x := r.(type) {
	case domain.HumanResource:
		return x.Role
	case domain.MaterialResource:
		return x.Name
	case domain.Asset:
		return x.Name
	default:
		return ""
	}
}

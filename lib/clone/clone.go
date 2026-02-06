package clone

import "explosio/lib/domain"

// CloneActivity restituisce una copia profonda dell'albero di attività radicato in root.
func CloneActivity(root *domain.Activity) *domain.Activity {
	if root == nil {
		return nil
	}
	cl := &domain.Activity{
		ID:            root.ID,
		Name:          root.Name,
		Description:   root.Description,
		Duration:      root.Duration,
		MinDuration:   root.MinDuration,
		CrashCostStep: root.CrashCostStep,
		ES:            root.ES,
		EF:            root.EF,
		LS:            root.LS,
		LF:            root.LF,
		Slack:         root.Slack,
	}
	if len(root.Humans) > 0 {
		cl.Humans = make([]domain.HumanResource, len(root.Humans))
		copy(cl.Humans, root.Humans)
	}
	if len(root.Materials) > 0 {
		cl.Materials = make([]domain.MaterialResource, len(root.Materials))
		copy(cl.Materials, root.Materials)
	}
	if len(root.Assets) > 0 {
		cl.Assets = make([]domain.Asset, len(root.Assets))
		copy(cl.Assets, root.Assets)
	}
	if len(root.Next) > 0 {
		cl.Next = make([]string, len(root.Next))
		copy(cl.Next, root.Next)
	}
	if len(root.SubActivities) > 0 {
		cl.SubActivities = make([]*domain.Activity, len(root.SubActivities))
		for i, sub := range root.SubActivities {
			cl.SubActivities[i] = CloneActivity(sub)
		}
	}
	return cl
}

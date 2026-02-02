package lib

// CloneActivity restituisce una copia profonda dell'albero di attività radicato in root.
// Il progetto originale non viene modificato.
func CloneActivity(root *Activity) *Activity {
	if root == nil {
		return nil
	}
	clone := &Activity{
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
		clone.Humans = make([]HumanResource, len(root.Humans))
		copy(clone.Humans, root.Humans)
	}
	if len(root.Materials) > 0 {
		clone.Materials = make([]MaterialResource, len(root.Materials))
		copy(clone.Materials, root.Materials)
	}
	if len(root.Assets) > 0 {
		clone.Assets = make([]Asset, len(root.Assets))
		copy(clone.Assets, root.Assets)
	}
	if len(root.Next) > 0 {
		clone.Next = make([]string, len(root.Next))
		copy(clone.Next, root.Next)
	}
	if len(root.SubActivities) > 0 {
		clone.SubActivities = make([]*Activity, len(root.SubActivities))
		for i, sub := range root.SubActivities {
			clone.SubActivities[i] = CloneActivity(sub)
		}
	}
	return clone
}

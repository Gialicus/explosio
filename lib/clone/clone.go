package clone

import "explosio/lib/domain"

func cloneSupplier(s *domain.Supplier) *domain.Supplier {
	if s == nil {
		return nil
	}
	return &domain.Supplier{
		Name:              s.Name,
		Description:       s.Description,
		UnitCost:          s.UnitCost,
		AvailableQuantity: s.AvailableQuantity,
		Period:            s.Period,
	}
}

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
		for i, h := range root.Humans {
			cl.Humans[i] = h
			if h.Supplier != nil {
				cl.Humans[i].Supplier = cloneSupplier(h.Supplier)
			}
		}
	}
	if len(root.Materials) > 0 {
		cl.Materials = make([]domain.MaterialResource, len(root.Materials))
		for i, m := range root.Materials {
			cl.Materials[i] = m
			if m.Supplier != nil {
				cl.Materials[i].Supplier = cloneSupplier(m.Supplier)
			}
		}
	}
	if len(root.Assets) > 0 {
		cl.Assets = make([]domain.Asset, len(root.Assets))
		for i, as := range root.Assets {
			cl.Assets[i] = as
			if as.Supplier != nil {
				cl.Assets[i].Supplier = cloneSupplier(as.Supplier)
			}
		}
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

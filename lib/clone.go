package lib

// cloneSupplier crea una copia profonda di un fornitore
func cloneSupplier(s *Supplier) *Supplier {
	if s == nil {
		return nil
	}
	return &Supplier{
		Name:              s.Name,
		Description:       s.Description,
		UnitCost:          s.UnitCost,
		AvailableQuantity: s.AvailableQuantity,
		Period:            s.Period,
	}
}

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
		for i, h := range root.Humans {
			clone.Humans[i] = h
			if h.Supplier != nil {
				// Clona il fornitore per evitare mutazioni condivise
				clone.Humans[i].Supplier = cloneSupplier(h.Supplier)
			}
		}
	}
	if len(root.Materials) > 0 {
		clone.Materials = make([]MaterialResource, len(root.Materials))
		for i, m := range root.Materials {
			clone.Materials[i] = m
			if m.Supplier != nil {
				// Clona il fornitore per evitare mutazioni condivise
				clone.Materials[i].Supplier = cloneSupplier(m.Supplier)
			}
		}
	}
	if len(root.Assets) > 0 {
		clone.Assets = make([]Asset, len(root.Assets))
		for i, as := range root.Assets {
			clone.Assets[i] = as
			if as.Supplier != nil {
				// Clona il fornitore per evitare mutazioni condivise
				clone.Assets[i].Supplier = cloneSupplier(as.Supplier)
			}
		}
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

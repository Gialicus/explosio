package lib

import (
	"encoding/json"
	"fmt"
)

// SerializableProject rappresenta un progetto serializzabile in JSON
type SerializableProject struct {
	Root      *SerializableActivity `json:"root"`
	Suppliers []*Supplier            `json:"suppliers,omitempty"`
}

// SerializableActivity rappresenta un'attività serializzabile
type SerializableActivity struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Duration    int                      `json:"duration"`
	MinDuration int                      `json:"minDuration"`
	CrashCost   float64                 `json:"crashCost,omitempty"`
	Humans      []SerializableHuman      `json:"humans,omitempty"`
	Materials   []SerializableMaterial   `json:"materials,omitempty"`
	Assets      []SerializableAsset      `json:"assets,omitempty"`
	SubActivities []*SerializableActivity `json:"subActivities,omitempty"`
}

// SerializableHuman rappresenta una risorsa umana serializzabile
type SerializableHuman struct {
	Role        string   `json:"role"`
	Description string   `json:"description"`
	CostPerH    float64  `json:"costPerH"`
	Quantity    float64  `json:"quantity"`
	SupplierRef string   `json:"supplierRef,omitempty"`
}

// SerializableMaterial rappresenta un materiale serializzabile
type SerializableMaterial struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	UnitCost    float64  `json:"unitCost"`
	Quantity    float64  `json:"quantity"`
	SupplierRef string   `json:"supplierRef,omitempty"`
}

// SerializableAsset rappresenta un asset serializzabile
type SerializableAsset struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CostPerUse  float64  `json:"costPerUse"`
	Quantity    float64  `json:"quantity"`
	SupplierRef string   `json:"supplierRef,omitempty"`
}

// SerializeProject converte un progetto in JSON
func SerializeProject(project *Project) ([]byte, error) {
	if project == nil || project.Root == nil {
		return nil, fmt.Errorf("project or root is nil")
	}

	// Raccogli tutti i fornitori unici dal progetto
	supplierMap := make(map[string]*Supplier)
	collectSuppliers(project.Root, supplierMap)
	
	suppliers := make([]*Supplier, 0, len(supplierMap))
	for _, s := range supplierMap {
		suppliers = append(suppliers, s)
	}

	serializable := &SerializableProject{
		Root:      serializeActivity(project.Root, supplierMap),
		Suppliers: suppliers,
	}

	return json.MarshalIndent(serializable, "", "  ")
}

// collectSuppliers raccoglie tutti i fornitori unici dall'albero di attività
func collectSuppliers(a *Activity, supplierMap map[string]*Supplier) {
	if a == nil {
		return
	}

	for _, h := range a.Humans {
		if h.Supplier != nil {
			supplierMap[h.Supplier.Name] = h.Supplier
		}
	}
	for _, m := range a.Materials {
		if m.Supplier != nil {
			supplierMap[m.Supplier.Name] = m.Supplier
		}
	}
	for _, as := range a.Assets {
		if as.Supplier != nil {
			supplierMap[as.Supplier.Name] = as.Supplier
		}
	}

	for _, sub := range a.SubActivities {
		collectSuppliers(sub, supplierMap)
	}
}

// serializeActivity converte un'Activity in SerializableActivity
func serializeActivity(a *Activity, supplierMap map[string]*Supplier) *SerializableActivity {
	if a == nil {
		return nil
	}

	serializable := &SerializableActivity{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Duration:    a.Duration,
		MinDuration:  a.MinDuration,
		CrashCost:   a.CrashCostStep,
	}

	// Serializza risorse umane
	if len(a.Humans) > 0 {
		serializable.Humans = make([]SerializableHuman, len(a.Humans))
		for i, h := range a.Humans {
			serializable.Humans[i] = SerializableHuman{
				Role:        h.Role,
				Description: h.Description,
				CostPerH:    h.CostPerH,
				Quantity:    h.Quantity,
			}
			if h.Supplier != nil {
				serializable.Humans[i].SupplierRef = h.Supplier.Name
			}
		}
	}

	// Serializza materiali
	if len(a.Materials) > 0 {
		serializable.Materials = make([]SerializableMaterial, len(a.Materials))
		for i, m := range a.Materials {
			serializable.Materials[i] = SerializableMaterial{
				Name:        m.Name,
				Description: m.Description,
				UnitCost:    m.UnitCost,
				Quantity:    m.Quantity,
			}
			if m.Supplier != nil {
				serializable.Materials[i].SupplierRef = m.Supplier.Name
			}
		}
	}

	// Serializza asset
	if len(a.Assets) > 0 {
		serializable.Assets = make([]SerializableAsset, len(a.Assets))
		for i, as := range a.Assets {
			serializable.Assets[i] = SerializableAsset{
				Name:        as.Name,
				Description: as.Description,
				CostPerUse:  as.CostPerUse,
				Quantity:    as.Quantity,
			}
			if as.Supplier != nil {
				serializable.Assets[i].SupplierRef = as.Supplier.Name
			}
		}
	}

	// Serializza sotto-attività
	if len(a.SubActivities) > 0 {
		serializable.SubActivities = make([]*SerializableActivity, len(a.SubActivities))
		for i, sub := range a.SubActivities {
			serializable.SubActivities[i] = serializeActivity(sub, supplierMap)
		}
	}

	return serializable
}

// DeserializeProject carica un progetto da JSON
func DeserializeProject(data []byte) (*Project, error) {
	var serializable SerializableProject
	if err := json.Unmarshal(data, &serializable); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project: %w", err)
	}

	if serializable.Root == nil {
		return nil, fmt.Errorf("project root is nil")
	}

	// Crea mappa fornitori per riferimento rapido
	supplierMap := make(map[string]*Supplier)
	for _, s := range serializable.Suppliers {
		supplierMap[s.Name] = s
	}

	// Ricostruisci il progetto
	project := &Project{}
	project.Root = deserializeActivity(serializable.Root, supplierMap)

	return project, nil
}

// deserializeActivity converte SerializableActivity in Activity
func deserializeActivity(sa *SerializableActivity, supplierMap map[string]*Supplier) *Activity {
	if sa == nil {
		return nil
	}

	activity := &Activity{
		ID:            sa.ID,
		Name:          sa.Name,
		Description:   sa.Description,
		Duration:      sa.Duration,
		MinDuration:   sa.MinDuration,
		CrashCostStep: sa.CrashCost,
	}

	// Deserializza risorse umane
	if len(sa.Humans) > 0 {
		activity.Humans = make([]HumanResource, len(sa.Humans))
		for i, sh := range sa.Humans {
			activity.Humans[i] = HumanResource{
				Role:        sh.Role,
				Description: sh.Description,
				CostPerH:    sh.CostPerH,
				Quantity:    sh.Quantity,
			}
			if sh.SupplierRef != "" {
				if supplier, ok := supplierMap[sh.SupplierRef]; ok {
					activity.Humans[i].Supplier = supplier
				}
			}
		}
	}

	// Deserializza materiali
	if len(sa.Materials) > 0 {
		activity.Materials = make([]MaterialResource, len(sa.Materials))
		for i, sm := range sa.Materials {
			activity.Materials[i] = MaterialResource{
				Name:        sm.Name,
				Description: sm.Description,
				UnitCost:    sm.UnitCost,
				Quantity:    sm.Quantity,
			}
			if sm.SupplierRef != "" {
				if supplier, ok := supplierMap[sm.SupplierRef]; ok {
					activity.Materials[i].Supplier = supplier
				}
			}
		}
	}

	// Deserializza asset
	if len(sa.Assets) > 0 {
		activity.Assets = make([]Asset, len(sa.Assets))
		for i, sa := range sa.Assets {
			activity.Assets[i] = Asset{
				Name:        sa.Name,
				Description: sa.Description,
				CostPerUse:  sa.CostPerUse,
				Quantity:    sa.Quantity,
			}
			if sa.SupplierRef != "" {
				if supplier, ok := supplierMap[sa.SupplierRef]; ok {
					activity.Assets[i].Supplier = supplier
				}
			}
		}
	}

	// Deserializza sotto-attività e ricostruisci dipendenze
	if len(sa.SubActivities) > 0 {
		activity.SubActivities = make([]*Activity, len(sa.SubActivities))
		for i, ssub := range sa.SubActivities {
			subActivity := deserializeActivity(ssub, supplierMap)
			activity.SubActivities[i] = subActivity
			// Ricostruisci la relazione Next
			subActivity.Next = append(subActivity.Next, activity.ID)
		}
	}

	return activity
}

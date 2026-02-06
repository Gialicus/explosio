package serialize

import (
	"encoding/json"
	"explosio/lib/domain"
	"explosio/lib/resources"
	"explosio/lib/tree"
	"fmt"
)

// SerializableProject rappresenta un progetto serializzabile in JSON
type SerializableProject struct {
	Root      *SerializableActivity `json:"root"`
	Suppliers []*domain.Supplier    `json:"suppliers,omitempty"`
}

// SerializableActivity rappresenta un'attività serializzabile
type SerializableActivity struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Duration       int                      `json:"duration"`
	MinDuration    int                      `json:"minDuration"`
	CrashCost      float64                  `json:"crashCost,omitempty"`
	Humans         []SerializableHuman      `json:"humans,omitempty"`
	Materials      []SerializableMaterial   `json:"materials,omitempty"`
	Assets         []SerializableAsset      `json:"assets,omitempty"`
	SubActivities  []*SerializableActivity `json:"subActivities,omitempty"`
}

// SerializableHuman rappresenta una risorsa umana serializzabile
type SerializableHuman struct {
	Role        string  `json:"role"`
	Description string  `json:"description"`
	CostPerH    float64 `json:"costPerH"`
	Quantity    float64 `json:"quantity"`
	SupplierRef string  `json:"supplierRef,omitempty"`
}

// SerializableMaterial rappresenta un materiale serializzabile
type SerializableMaterial struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	UnitCost    float64 `json:"unitCost"`
	Quantity    float64 `json:"quantity"`
	SupplierRef string  `json:"supplierRef,omitempty"`
}

// SerializableAsset rappresenta un asset serializzabile
type SerializableAsset struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CostPerUse  float64 `json:"costPerUse"`
	Quantity    float64 `json:"quantity"`
	SupplierRef string  `json:"supplierRef,omitempty"`
}

// SerializeProject converte un progetto in JSON
func SerializeProject(project *domain.Project) ([]byte, error) {
	if project == nil || project.Root == nil {
		return nil, fmt.Errorf("project or root is nil")
	}
	supplierMap := make(map[string]*domain.Supplier)
	tree.Walk(project.Root, func(a *domain.Activity) {
		resources.ForEachResourceWithSupplier(a, func(r domain.Resource, s *domain.Supplier) {
			supplierMap[s.Name] = s
		})
	})
	suppliers := make([]*domain.Supplier, 0, len(supplierMap))
	for _, s := range supplierMap {
		suppliers = append(suppliers, s)
	}
	serializable := &SerializableProject{
		Root:      serializeActivity(project.Root, supplierMap),
		Suppliers: suppliers,
	}
	return json.MarshalIndent(serializable, "", "  ")
}

func serializeActivity(a *domain.Activity, supplierMap map[string]*domain.Supplier) *SerializableActivity {
	if a == nil {
		return nil
	}
	serializable := &SerializableActivity{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Duration:    a.Duration,
		MinDuration: a.MinDuration,
		CrashCost:   a.CrashCostStep,
	}
	var sh []SerializableHuman
	var sm []SerializableMaterial
	var sas []SerializableAsset
	resources.ForEachResource(a, func(r domain.Resource) {
		switch x := r.(type) {
		case domain.HumanResource:
			s := SerializableHuman{Role: x.Role, Description: x.Description, CostPerH: x.CostPerH, Quantity: x.Quantity}
			if x.Supplier != nil {
				s.SupplierRef = x.Supplier.Name
			}
			sh = append(sh, s)
		case domain.MaterialResource:
			s := SerializableMaterial{Name: x.Name, Description: x.Description, UnitCost: x.UnitCost, Quantity: x.Quantity}
			if x.Supplier != nil {
				s.SupplierRef = x.Supplier.Name
			}
			sm = append(sm, s)
		case domain.Asset:
			s := SerializableAsset{Name: x.Name, Description: x.Description, CostPerUse: x.CostPerUse, Quantity: x.Quantity}
			if x.Supplier != nil {
				s.SupplierRef = x.Supplier.Name
			}
			sas = append(sas, s)
		}
	})
	serializable.Humans = sh
	serializable.Materials = sm
	serializable.Assets = sas
	if len(a.SubActivities) > 0 {
		serializable.SubActivities = make([]*SerializableActivity, len(a.SubActivities))
		for i, sub := range a.SubActivities {
			serializable.SubActivities[i] = serializeActivity(sub, supplierMap)
		}
	}
	return serializable
}

// DeserializeProject carica un progetto da JSON
func DeserializeProject(data []byte) (*domain.Project, error) {
	var serializable SerializableProject
	if err := json.Unmarshal(data, &serializable); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project: %w", err)
	}
	if serializable.Root == nil {
		return nil, fmt.Errorf("project root is nil")
	}
	supplierMap := make(map[string]*domain.Supplier)
	for _, s := range serializable.Suppliers {
		supplierMap[s.Name] = s
	}
	project := &domain.Project{}
	project.Root = deserializeActivity(serializable.Root, supplierMap)
	return project, nil
}

func deserializeActivity(sa *SerializableActivity, supplierMap map[string]*domain.Supplier) *domain.Activity {
	if sa == nil {
		return nil
	}
	activity := &domain.Activity{
		ID:            sa.ID,
		Name:          sa.Name,
		Description:   sa.Description,
		Duration:      sa.Duration,
		MinDuration:   sa.MinDuration,
		CrashCostStep: sa.CrashCost,
	}
	if len(sa.Humans) > 0 {
		activity.Humans = make([]domain.HumanResource, len(sa.Humans))
		for i, sh := range sa.Humans {
			activity.Humans[i] = domain.HumanResource{
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
	if len(sa.Materials) > 0 {
		activity.Materials = make([]domain.MaterialResource, len(sa.Materials))
		for i, sm := range sa.Materials {
			activity.Materials[i] = domain.MaterialResource{
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
	if len(sa.Assets) > 0 {
		activity.Assets = make([]domain.Asset, len(sa.Assets))
		for i, sas := range sa.Assets {
			activity.Assets[i] = domain.Asset{
				Name:        sas.Name,
				Description: sas.Description,
				CostPerUse:  sas.CostPerUse,
				Quantity:    sas.Quantity,
			}
			if sas.SupplierRef != "" {
				if supplier, ok := supplierMap[sas.SupplierRef]; ok {
					activity.Assets[i].Supplier = supplier
				}
			}
		}
	}
	if len(sa.SubActivities) > 0 {
		activity.SubActivities = make([]*domain.Activity, len(sa.SubActivities))
		for i, ssub := range sa.SubActivities {
			subActivity := deserializeActivity(ssub, supplierMap)
			activity.SubActivities[i] = subActivity
			subActivity.Next = append(subActivity.Next, activity.ID)
		}
	}
	return activity
}

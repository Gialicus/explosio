package main

import (
	"explosio/lib"
	"fmt"
)

func main() {
	builder := lib.NewProject()
	engine := &lib.AnalysisEngine{}

	// Creazione Fornitori (riutilizzabili)
	// Fornitore caffè: 50 sacchi/mese (ogni sacco = 1kg = 1000g)
	// Capacità: 50kg/mese = 50000g/mese
	fornitoreCaffe := lib.NewSupplier("Fornitore Caffè", "Fornitore principale", 22.0, 50000, lib.PeriodMonth)

	// Fornitore latte: 200 litri/settimana
	fornitoreLatte := lib.NewSupplier("Fornitore Latte", "Fornitore latte fresco", 1.0, 200000, lib.PeriodWeek) // 200L = 200000ml

	// Creazione Progetto
	// Dati realistici:
	// - Miscela: 18g per cappuccino (7-8g espresso + sprechi) a 0.40€ = ~22€/kg
	// - Latte: 150ml per cappuccino a 0.15€ = 1€/L
	// - Tempi: 1 min espresso, 3 min montatura latte, 1 min assemblaggio, 2 min servizio
	cappuccino := builder.Start("Vendita Cappuccino", "Consegna al tavolo", 2).
		WithHuman("Cameriere", "Servizio al tavolo", 15, 1).
		DependsOn(
			builder.Node("Assemblaggio", "Mix ingredienti", 1).
				CanCrash(0, 5.0).
				DependsOn(
					builder.Node("Preparazione Caffè", "Espresso", 1).
						WithMaterialFromSupplier("Miscela", "Caffè Arabica", 0.40, 18, fornitoreCaffe).
						CanCrash(1, 10.0),
					builder.Node("Montatura Latte", "Vapore e schiuma", 3).
						WithMaterialFromSupplier("Latte", "Intero fresco", 0.15, 150, fornitoreLatte).
						WithHuman("Barista", "Specialista caffè", 20, 1),
				),
		)

	// Esecuzione Calcoli
	engine.ComputeCPM(cappuccino)

	fmt.Println("=== GERARCHIA FILIERA ===")
	lib.PrintReport(cappuccino, 0, true, "")

	fmt.Println("\n=== DIAGRAMMA MERMAID ===")
	lib.PrintMermaid(cappuccino)
	if err := lib.WriteMermaidToFile(cappuccino, "cappuccino.mmd"); err != nil {
		fmt.Println("Errore salvataggio Mermaid:", err)
	} else {
		fmt.Println("Diagramma salvato in cappuccino.mmd")
	}

	// Analisi Finanziaria
	targetPrice := 4.50
	fin := engine.GetFinancials(cappuccino, targetPrice)

	fmt.Println("\n=== REPORT ECONOMICO ===")
	fmt.Printf("Costo di produzione: €%.2f\n", fin.TotalCost)
	fmt.Printf("Prezzo consigliato:  €%.2f\n", targetPrice)
	fmt.Printf("Margine unitario:    €%.2f\n", fin.Margin)
	fmt.Printf("Markup:              %.1f%%\n", fin.Markup)

	// Analisi Ottimizzazione
	savedTime, extraCost := engine.GetMaxCrashPotential(cappuccino)
	fmt.Println("\n=== OTTIMIZZAZIONE (CRASHING) ===")
	fmt.Printf("Tempo risparmiabile: %d min\n", savedTime)
	fmt.Printf("Costo extra stimato: €%.2f\n", extraCost)
	fmt.Printf("Costo nuovo (accelerato): €%.2f\n", fin.TotalCost+extraCost)

	// What-if: confronto scenari
	whatIf := lib.NewWhatIfEngine()
	materialFactor := 1.5
	scenarios := []lib.Scenario{
		{Name: "Base", SellPrice: 4.50, Overrides: nil},
		{Name: "Prezzo 5€", SellPrice: 5.0, Overrides: nil},
		{Name: "Caffè più caro (+50%)", SellPrice: 4.50, Overrides: map[string]lib.ActivityOverride{
			"ACT-003": {MaterialCostFactor: materialFactor},
		}},
	}
	results := whatIf.RunScenarios(cappuccino, scenarios)
	fmt.Println("\n=== CONFRONTO SCENARI WHAT-IF ===")
	fmt.Printf("%-25s %8s %10s %10s %8s\n", "Scenario", "Durata", "Costo", "Margine", "Markup%")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, r := range results {
		fmt.Printf("%-25s %6d min %9.2f€ %9.2f€ %7.1f%%\n",
			r.ScenarioName, r.TotalDuration, r.TotalCost, r.Margin, r.Markup)
	}

	// Analisi scenario complesso: quanti fornitori servono per 1000 cappuccini/giorno?
	fmt.Println("\n=== ANALISI SCENARIO COMPLESSO ===")
	fmt.Println("Scenario: Produzione di 1000 cappuccini al giorno")
	productionTarget := 1000.0 // 1000 cappuccini al giorno
	targetPeriod := lib.PeriodDay

	requirements := engine.CalculateSupplierRequirements(cappuccino, productionTarget, targetPeriod)
	lib.PrintSupplierRequirements(requirements)

	// Calcolo dettagliato
	if len(requirements) > 0 {
		fmt.Println("\n=== DETTAGLIO CALCOLO ===")
		for _, req := range requirements {
			fmt.Printf("\nFornitore: %s\n", req.SupplierName)
			fmt.Printf("  Produzione target: %.0f cappuccini/%s\n", productionTarget, targetPeriod.String())

			// Mostra la conversione in unità più leggibili
			switch req.SupplierName {
			case "Fornitore Caffè":
				kgPerMonth := req.RequiredQuantity / 1000.0 // converti grammi in kg
				fmt.Printf("  Quantità richiesta: %.2f kg/%s (%.0f g/%s)\n", kgPerMonth, req.SupplierPeriod.String(), req.RequiredQuantity, req.SupplierPeriod.String())
				// Capacità fornitore: 50000g/mese = 50kg/mese = 50 sacchi da 1kg
				fmt.Printf("  Capacità fornitore: 50 kg/%s (50 sacchi da 1kg)\n", req.SupplierPeriod.String())
			case "Fornitore Latte":
				litersPerWeek := req.RequiredQuantity / 1000.0 // converti ml in litri
				fmt.Printf("  Quantità richiesta: %.2f L/%s (%.0f ml/%s)\n", litersPerWeek, req.SupplierPeriod.String(), req.RequiredQuantity, req.SupplierPeriod.String())
				// Capacità fornitore: 200000ml/settimana = 200L/settimana
				fmt.Printf("  Capacità fornitore: 200 L/%s\n", req.SupplierPeriod.String())
			default:
				fmt.Printf("  Quantità richiesta: %.2f unità/%s\n", req.RequiredQuantity, req.SupplierPeriod.String())
			}

			fmt.Printf("  Fornitori necessari: %.2f\n", req.SuppliersNeeded)
			if req.SuppliersNeeded > 0 {
				rounded := int(req.SuppliersNeeded + 0.5) // arrotondamento per eccesso
				fmt.Printf("  (arrotondato: %d fornitori)\n", rounded)
			}
		}
	}
}

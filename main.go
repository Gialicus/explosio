// Package main is the entry point of the explosio application.
package main

import (
	"explosio/core"
	"explosio/core/unit"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		runDemo()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "run":
		runDemo()
	case "load":
		runLoad(os.Args[2:])
	case "export":
		runExport(os.Args[2:])
	case "query":
		runQuery(os.Args[2:])
	case "gantt":
		runGantt(os.Args[2:])
	case "validate":
		runValidate(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		// Default: run demo (backward compatibility)
		runDemo()
	}
}

func printUsage() {
	fmt.Print(`explosio - Activity tree modelling with CPM

Usage:
  explosio              Run demo (default)
  explosio run          Run demo project
  explosio load <file>  Load project from JSON or YAML file and print
  explosio export       Export project to JSON or YAML
    -input <file>       Input file (JSON or YAML)
    -output <file>      Output file (default: stdout)
    -format json|yaml   Output format (default: json)
  explosio query       Filter activities by criteria
    -input <file>       Input file (required)
    -price-range min-max  Filter by price range (e.g. 100-1000)
    -name <pattern>     Filter by name (substring match)
  explosio gantt       Print ASCII Gantt chart
    [-input <file>]     Input file (default: demo)
    [-start YYYY-MM-DD] Project start date (default: today)
  explosio validate    Validate project (circular deps, references, warnings)
    [-input <file>]     Input file (default: demo)
`)
}

func runDemo() {
	homeRenovation := BuildDemoTree()
	core.PrettyPrint([]*core.Activity{homeRenovation}, homeRenovation.CalculateCriticalPath())
	fmt.Printf("\nTotal price: %.2f EUR\n", homeRenovation.CalculatePrice())
	fmt.Printf("Total duration: %.0f days\n", homeRenovation.CalculateDuration())
	meas := homeRenovation.GetMeasurableMaterials()
	countable := homeRenovation.GetCountableMaterials()
	complexMat := homeRenovation.GetComplexMaterials()
	hr := homeRenovation.GetHumanResources()
	assets := homeRenovation.GetAssets()
	fmt.Printf("Materials: %d measurable, %d countable, %d complex | Human resources: %d | Assets: %d\n", len(meas), len(countable), len(complexMat), len(hr), len(assets))
	cb := homeRenovation.CostBreakdown()
	fmt.Printf("Cost breakdown: Activities %.2f | Materials %.2f | Human %.2f | Assets %.2f\n", cb.Activities, cb.Materials, cb.Human, cb.Assets)

	path := homeRenovation.CalculateCriticalPath()
	var pathNames []string
	for _, a := range path {
		pathNames = append(pathNames, a.Name)
	}
	fmt.Printf("Critical path: %s\n", strings.Join(pathNames, " -> "))

	r := homeRenovation.Validate()
	if r.Valid() {
		fmt.Println("Validation: OK")
	} else {
		for _, e := range r.Errors {
			fmt.Printf("Validation error: %s\n", e.Error())
		}
	}
	fmt.Println("Run 'explosio gantt' for timeline chart")
}

func runLoad(args []string) {
	fs := flag.NewFlagSet("load", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: explosio load <file>")
	}
	_ = fs.Parse(args)
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: load requires a file path")
		fs.Usage()
		os.Exit(1)
	}
	path := fs.Arg(0)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	proj, err := readProject(path, f)
	if err != nil {
		log.Fatalf("load %s: %v", path, err)
	}

	root := proj.Root
	core.PrettyPrint([]*core.Activity{root}, root.CalculateCriticalPath())
	fmt.Printf("\nTotal price: %.2f %s\n", root.CalculatePrice(), root.Price.Currency)
	fmt.Printf("Total duration: %.0f %s\n", root.CalculateDuration(), root.Duration.Unit)
	meas := root.GetMeasurableMaterials()
	countable := root.GetCountableMaterials()
	complexMat := root.GetComplexMaterials()
	hr := root.GetHumanResources()
	assets := root.GetAssets()
	fmt.Printf("Materials: %d measurable, %d countable, %d complex | Human resources: %d | Assets: %d\n", len(meas), len(countable), len(complexMat), len(hr), len(assets))
	cb := root.CostBreakdown()
	fmt.Printf("Cost breakdown: Activities %.2f | Materials %.2f | Human %.2f | Assets %.2f\n", cb.Activities, cb.Materials, cb.Human, cb.Assets)
}

func runExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	input := fs.String("input", "", "Input file (JSON or YAML); if empty, exports demo")
	output := fs.String("output", "", "Output file (default: stdout)")
	format := fs.String("format", "json", "Output format: json or yaml")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: explosio export [-input <file>] [-output <file>] [-format json|yaml]")
	}
	_ = fs.Parse(args)

	var proj *core.Project
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			log.Fatalf("open %s: %v", *input, err)
		}
		defer f.Close()
		proj, err = readProject(*input, f)
		if err != nil {
			log.Fatalf("load %s: %v", *input, err)
		}
	} else {
		proj = core.NewProject(BuildDemoTree())
	}

	out := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("create %s: %v", *output, err)
		}
		defer f.Close()
		out = f
	}

	switch strings.ToLower(*format) {
	case "json":
		if err := proj.WriteJSON(out); err != nil {
			log.Fatalf("write JSON: %v", err)
		}
	case "yaml":
		if err := proj.WriteYAML(out); err != nil {
			log.Fatalf("write YAML: %v", err)
		}
	default:
		log.Fatalf("unsupported format: %s (use json or yaml)", *format)
	}
}

func runQuery(args []string) {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	input := fs.String("input", "", "Input file (JSON or YAML)")
	priceRange := fs.String("price-range", "", "Filter by price range (e.g. 100-1000)")
	name := fs.String("name", "", "Filter by name (substring match)")
	nameRegex := fs.String("name-regex", "", "Filter by name (regex)")
	material := fs.String("material", "", "Filter activities using material (substring)")
	resource := fs.String("resource", "", "Filter activities using human resource (substring)")
	sortBy := fs.String("sort", "", "Sort by: name, price, duration")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: explosio query -input <file> [-price-range min-max] [-name <pattern>] [-material <name>] [-resource <name>] [-sort name|price|duration]")
	}
	_ = fs.Parse(args)

	if *input == "" {
		fmt.Fprintln(os.Stderr, "Error: -input is required")
		fs.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*input)
	if err != nil {
		log.Fatalf("open %s: %v", *input, err)
	}
	defer f.Close()

	proj, err := readProject(*input, f)
	if err != nil {
		log.Fatalf("load %s: %v", *input, err)
	}

	activities := proj.Root.GetActivities()
	filtered := core.FilterActivities(activities, core.FilterOptions{
		PriceMin:     parsePriceRangeMin(*priceRange),
		PriceMax:     parsePriceRangeMax(*priceRange),
		Name:         *name,
		NameRegex:    *nameRegex,
		MaterialName: *material,
		ResourceName: *resource,
	})

	switch strings.ToLower(*sortBy) {
	case "price":
		core.SortActivities(filtered, core.SortByPrice)
	case "duration":
		core.SortActivities(filtered, core.SortByDuration)
	case "name":
		core.SortActivities(filtered, core.SortByName)
	}

	for _, a := range filtered {
		fmt.Printf("%s: %.2f %s\n", a.Name, a.CalculatePrice(), a.Price.Currency)
	}
}

func runGantt(args []string) {
	fs := flag.NewFlagSet("gantt", flag.ExitOnError)
	input := fs.String("input", "", "Input file (default: demo)")
	startStr := fs.String("start", "", "Project start date YYYY-MM-DD (default: today)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: explosio gantt [-input <file>] [-start YYYY-MM-DD]")
	}
	_ = fs.Parse(args)

	var root *core.Activity
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			log.Fatalf("open %s: %v", *input, err)
		}
		defer f.Close()
		proj, err := readProject(*input, f)
		if err != nil {
			log.Fatalf("load %s: %v", *input, err)
		}
		root = proj.Root
	} else {
		root = BuildDemoTree()
	}

	projectStart := unit.NewDate(time.Now().Year(), time.Now().Month(), time.Now().Day())
	if *startStr != "" {
		var y, m, d int
		if _, err := fmt.Sscanf(*startStr, "%d-%d-%d", &y, &m, &d); err == nil {
			projectStart = unit.NewDate(y, time.Month(m), d)
		}
	}

	root.PrintGantt(core.GanttConfig{
		ProjectStart: projectStart,
		Width:        50,
		ShowDates:    true,
	})
}

func runValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	input := fs.String("input", "", "Input file (default: demo)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: explosio validate [-input <file>]")
	}
	_ = fs.Parse(args)

	var root *core.Activity
	if *input != "" {
		f, err := os.Open(*input)
		if err != nil {
			log.Fatalf("open %s: %v", *input, err)
		}
		defer f.Close()
		proj, err := readProject(*input, f)
		if err != nil {
			log.Fatalf("load %s: %v", *input, err)
		}
		root = proj.Root
	} else {
		root = BuildDemoTree()
	}

	r := root.Validate()
	for _, e := range r.Errors {
		fmt.Printf("Error: %s\n", e.Error())
	}
	for _, w := range r.Warnings {
		fmt.Printf("Warning: %s\n", w.Error())
	}
	if r.Valid() {
		fmt.Println("Validation passed.")
	} else {
		os.Exit(1)
	}
}

func readProject(path string, r *os.File) (*core.Project, error) {
	if strings.HasSuffix(strings.ToLower(path), ".yaml") || strings.HasSuffix(strings.ToLower(path), ".yml") {
		return core.ReadYAML(r)
	}
	return core.ReadJSON(r)
}

func parsePriceRangeMin(s string) float64 {
	if s == "" {
		return 0
	}
	parts := strings.Split(s, "-")
	if len(parts) < 2 {
		return 0
	}
	var min float64
	_, _ = fmt.Sscanf(parts[0], "%f", &min)
	return min
}

func parsePriceRangeMax(s string) float64 {
	if s == "" {
		return 0
	}
	parts := strings.Split(s, "-")
	if len(parts) < 2 {
		return 0
	}
	var max float64
	_, _ = fmt.Sscanf(parts[1], "%f", &max)
	return max
}

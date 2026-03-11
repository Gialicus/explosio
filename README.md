# explosio

Application for modelling activity trees (e.g. renovations) with durations, prices, materials, and CPM-style critical path calculation.

## Requirements

- Go 1.21 or later (compatible with the version in `go.mod`).

## Build and run

```bash
go build -o explosio .
./explosio
```

## CLI commands

- `explosio` or `explosio run` — Run demo project
- `explosio load <file>` — Load project from JSON or YAML and print
- `explosio export [-input <file>] [-output <file>] [-format json|yaml]` — Export project
- `explosio query -input <file> [-price-range min-max] [-name <pattern>] [-material <name>] [-resource <name>] [-sort name|price|duration]` — Filter activities
- `explosio gantt [-input <file>] [-start YYYY-MM-DD]` — Print ASCII Gantt chart
- `explosio validate [-input <file>]` — Validate project (circular deps, references, warnings)
- `explosio help` — Show usage

## Project structure

- **main.go**, **demo.go**: Entry point and demo tree
- **core/**: Activity model, CPM, calculations, serialization, Gantt, validation
- **core/material/**: Material types (complex, countable, measurable)
- **core/unit/**: Types for durations, prices, dates, measurable quantities
- **core/resource/**: Human resources and assets

## Features

- Activity tree with materials, human resources, assets
- CPM critical path and slack (float) calculation
- Explicit dependencies (`DependsOn`) for cross-branch CPM
- Cost breakdown by category (activities, materials, human, assets)
- Milestones (zero-duration activities)
- JSON/YAML persistence
- ASCII Gantt chart with dates
- Filter and sort activities
- Clone for scenario comparison
- Validation (circular dependencies, references, warnings)

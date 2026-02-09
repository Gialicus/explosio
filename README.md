# explosio

Application for modelling activity trees (e.g. renovations) with durations, prices, materials, and CPM-style critical path calculation.

## Requirements

- Go 1.21 or later (compatible with the version in `go.mod`).

## Build and run

```bash
go build -o explosio .
```

Or, to build all packages:

```bash
go build ./...
```

To run:

```bash
go run .
```

Or, after building:

```bash
./explosio
```

## Project structure

- **main.go**: Entry point; builds a sample 5-level activity tree (home renovation), prints it, and shows totals and critical path.
- **core/**: `Activity` model, calculation methods (price, duration, quantity, critical path), and formatted tree printing.
- **core/material/**: Material types (complex, countable, measurable).
- **core/unit/**: Types for durations, prices, and measurable quantities.

## Output

The program prints the activity tree with price and duration (own and total) for each node; activities on the critical path are highlighted. At the end it shows total price, total duration, a summary of materials and resources (counts), and the critical path.

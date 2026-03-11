// Package core provides Gantt chart output for activity scheduling.
package core

import (
	"fmt"
	"strings"

	"explosio/core/unit"
)

// Schedule holds computed start/end dates for an activity.
type Schedule struct {
	Activity  *Activity
	StartDate unit.Date
	EndDate   unit.Date
	ES        float64 // Early start (hours from project start)
	EF        float64 // Early finish (hours from project start)
}

// ComputeSchedule computes start/end dates for all activities given a project start date.
func (a *Activity) ComputeSchedule(projectStart unit.Date) map[*Activity]Schedule {
	slackMap := a.CalculateSlack()
	result := make(map[*Activity]Schedule)
	for act, info := range slackMap {
		start := projectStart.AddHours(info.ES)
		end := projectStart.AddHours(info.EF)
		result[act] = Schedule{
			Activity:  act,
			StartDate: start,
			EndDate:   end,
			ES:        info.ES,
			EF:        info.EF,
		}
	}
	return result
}

// GanttConfig holds options for Gantt output.
type GanttConfig struct {
	ProjectStart unit.Date
	Width        int  // Character width for the bar (default 40)
	ShowDates    bool // Show date labels (default true)
}

// PrintGantt prints an ASCII Gantt chart for the activity tree.
func (a *Activity) PrintGantt(cfg GanttConfig) {
	if cfg.Width <= 0 {
		cfg.Width = 40
	}
	schedule := a.ComputeSchedule(cfg.ProjectStart)
	var projectEnd float64
	if a.hasExplicitDependencies() {
		_, projectEnd = a.cpmForwardBackward()
	} else {
		_, projectEnd = a.criticalPathAndDuration()
	}
	totalHours := projectEnd
	if totalHours <= 0 {
		totalHours = 1
	}

	activities := a.GetActivities()
	fmt.Println("--------------------------------")
	fmt.Println("   Gantt Chart")
	fmt.Println("--------------------------------")
	if cfg.ShowDates {
		fmt.Printf("Project start: %s | Duration: %.0f hours\n", cfg.ProjectStart.String(), totalHours)
		fmt.Println("--------------------------------")
	}

	for _, act := range activities {
		sched, ok := schedule[act]
		if !ok {
			continue
		}
		startPct := sched.ES / totalHours
		endPct := sched.EF / totalHours
		barLen := cfg.Width
		startIdx := int(startPct * float64(barLen))
		endIdx := int(endPct * float64(barLen))
		if endIdx <= startIdx && sched.EF > sched.ES {
			endIdx = startIdx + 1
		}
		if endIdx > barLen {
			endIdx = barLen
		}

		var bar strings.Builder
		for i := 0; i < barLen; i++ {
			if i >= startIdx && i < endIdx {
				bar.WriteString("█")
			} else {
				bar.WriteString("░")
			}
		}
		name := act.Name
		if len(name) > 20 {
			name = name[:17] + "..."
		}
		fmt.Printf("%-20s |%s|\n", name, bar.String())
		if cfg.ShowDates {
			fmt.Printf("%20s   %s - %s\n", "", sched.StartDate.String(), sched.EndDate.String())
		}
	}
}

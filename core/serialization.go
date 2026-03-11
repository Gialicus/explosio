// Package core provides serialization (JSON/YAML) for Activity trees.
package core

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const ProjectVersion = "1.0"

// Project wraps a root activity for file persistence. Allows adding metadata (version, name) later.
type Project struct {
	Version string    `json:"version" yaml:"version"`
	Root    *Activity `json:"root" yaml:"root"`
}

// NewProject creates a project with the given root activity.
func NewProject(root *Activity) *Project {
	return &Project{
		Version: ProjectVersion,
		Root:    root,
	}
}

// WriteJSON writes the project to w as formatted JSON.
func (p *Project) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}

// ReadJSON reads a project from r (JSON format).
func ReadJSON(r io.Reader) (*Project, error) {
	var p Project
	if err := json.NewDecoder(r).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	if p.Root == nil {
		return nil, fmt.Errorf("project root is nil")
	}
	return &p, nil
}

// WriteYAML writes the project to w as YAML.
func (p *Project) WriteYAML(w io.Writer) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// ReadYAML reads a project from r (YAML format).
func ReadYAML(r io.Reader) (*Project, error) {
	var p Project
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&p); err != nil {
		return nil, fmt.Errorf("decode YAML: %w", err)
	}
	if p.Root == nil {
		return nil, fmt.Errorf("project root is nil")
	}
	return &p, nil
}

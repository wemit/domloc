package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Route struct {
	Domain   string `json:"domain"`
	Port     int    `json:"port"`
	HTTPS    bool   `json:"https"`
	Wildcard bool   `json:"wildcard"`
}

type Registry struct {
	Routes       []Route `json:"routes"`
	HTTPSDefault *bool   `json:"https_default,omitempty"`
	path         string
}

func Load() (*Registry, error) {
	path, err := routesPath()
	if err != nil {
		return nil, err
	}
	r := &Registry{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return r, nil
	}
	if err != nil {
		return nil, err
	}
	return r, json.Unmarshal(data, r)
}

func (r *Registry) Save() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0644)
}

func (r *Registry) Add(route Route) {
	for i, existing := range r.Routes {
		if existing.Domain == route.Domain {
			r.Routes[i] = route
			return
		}
	}
	r.Routes = append(r.Routes, route)
}

func (r *Registry) Remove(domain string) bool {
	for i, route := range r.Routes {
		if route.Domain == domain {
			r.Routes = append(r.Routes[:i], r.Routes[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Registry) Find(domain string) (Route, bool) {
	for _, route := range r.Routes {
		if route.Domain == domain {
			return route, true
		}
	}
	return Route{}, false
}

func (r *Registry) IsHTTPSDefault() bool {
	if r.HTTPSDefault == nil {
		return true
	}
	return *r.HTTPSDefault
}

func BoolPtr(b bool) *bool { return &b }

func routesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "domloc", "routes.json"), nil
}

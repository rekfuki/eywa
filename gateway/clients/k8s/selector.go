package k8s

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// Selector represents requirement used to filter k8s labels
type Selector struct {
	labels.Selector
}

// LabelSelector returns empty requirements to be used to filter labels
func LabelSelector() Selector {
	return Selector{labels.NewSelector()}
}

// In ...
func (s Selector) In(field string, values []string) Selector {
	r, _ := labels.NewRequirement(field, selection.In, values)
	s.Selector = s.Add(*r)
	return s
}

// Equals ...
func (s Selector) Equals(field string, value string) Selector {
	r, _ := labels.NewRequirement(field, selection.Equals, []string{value})
	s.Selector = s.Add(*r)
	return s
}

// Exists ...
func (s Selector) Exists(field string) Selector {
	r, _ := labels.NewRequirement(field, selection.Exists, []string{})
	s.Selector = s.Add(*r)
	return s
}

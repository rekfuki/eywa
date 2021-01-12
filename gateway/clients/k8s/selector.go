package k8s

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

// Selector represents requirement used to filter k8s labels
type Selector []labels.Requirement

// LabelSelector returns empty requirements to be used to filter labels
func LabelSelector() Selector {
	return Selector{}
}

// In ...
func (rs Selector) In(field string, values []string) Selector {
	r, _ := labels.NewRequirement(field, selection.In, values)
	rs = append(rs, *r)
	return rs
}

// Equals ...
func (rs Selector) Equals(field string, value string) Selector {
	r, _ := labels.NewRequirement(field, selection.Equals, []string{value})
	rs = append(rs, *r)
	return rs
}

// Exists ...
func (rs Selector) Exists(field string) Selector {
	r, _ := labels.NewRequirement(field, selection.Exists, []string{})
	rs = append(rs, *r)
	return rs
}

// Package filter provides service-level filtering for drift detection runs.
//
// It allows callers to narrow the set of services that driftcheck inspects,
// either by an explicit allow-list of service names or by requiring a set of
// Docker / Kubernetes labels to be present on the running container.
//
// Typical usage:
//
//	services := filter.Filter(allServices, filter.Options{
//		Services: cfg.Services,
//		Labels:   cfg.RequiredLabels,
//	})
//
// The filter is applied before any manifest comparison so that irrelevant
// services never reach the detector, keeping output focused and fast.
package filter

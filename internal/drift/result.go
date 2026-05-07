// Package drift provides types and helpers shared across drift detection.
package drift

import "fmt"

// DriftResult captures the outcome of comparing a single service against
// its running container counterpart.
type DriftResult struct {
	// ServiceName is the name declared in the manifest.
	ServiceName string

	// ContainerID is the short ID of the matched running container, if any.
	ContainerID string

	// Drifted is true when the running state differs from the manifest.
	Drifted bool

	// Reason is a human-readable explanation of the drift, if any.
	Reason string

	// ExpectedImage is the image declared in the manifest.
	ExpectedImage string

	// ActualImage is the image currently running in the container.
	ActualImage string
}

// String returns a concise single-line summary of the result.
func (r DriftResult) String() string {
	if !r.Drifted {
		return fmt.Sprintf("[OK]      %s", r.ServiceName)
	}
	return fmt.Sprintf("[DRIFTED] %s — %s", r.ServiceName, r.Reason)
}

// ContainerInfo holds the subset of runtime container data used by detectors.
type ContainerInfo struct {
	// ID is the full container ID.
	ID string

	// Name is the primary container name (may include a leading '/').
	Name string

	// Image is the image tag or digest the container is running.
	Image string
}

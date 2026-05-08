// Package plugin defines the Checker interface and Registry that allow
// third-party or built-in extensions to participate in the driftcheck
// pipeline.
//
// # Extension point
//
// Implement the Checker interface and register your implementation with
// a Registry before running a scan:
//
//	reg := plugin.NewRegistry()
//	_ = reg.Register(myChecker)
//	results, err = reg.Run(results)
//
// Checkers are executed in registration order. Each checker receives the
// full slice of drift.Result values and may modify, annotate or extend it.
//
// # Built-in checkers
//
// Ready-to-use checkers live in the builtin sub-package and can be
// registered without writing any additional code.
package plugin

// Package ignore implements drift-ignore rule loading and evaluation.
//
// A drift-ignore file is a YAML file with the following structure:
//
//	ignore:
//	  - service: web
//	    fields:
//	      - image
//	      - env
//	  - service: worker
//	    fields:
//	      - ports
//
// Rules are matched by exact service name. The fields correspond to drift
// categories produced by the detector (e.g. "image", "env", "ports").
// Any matching field for a matching service is considered suppressed and
// should be excluded from drift reports.
//
// If no ignore file path is provided, or the file does not exist, Load returns
// an empty Set so callers do not need to handle the missing-file case
// explicitly.
package ignore

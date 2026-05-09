// Package trend provides utilities for analysing drift patterns across
// multiple historical scan results.
//
// Given a chronologically ordered slice of [Entry] values — each pairing a
// timestamp with the drift results from one scan — [Analyse] returns a
// per-service [ServiceTrend] that captures:
//
//   - TotalScans: how many scans included the service
//   - DriftedScans: how many of those scans detected drift
//   - DriftRate: the fraction of drifted scans (0.0–1.0)
//   - Worsening: true when the most recent scan is drifted but the
//     previous scan was clean, indicating a newly introduced drift
//
// Results are sorted by DriftRate descending so the most problematic
// services appear first.
package trend

// Package scanner provides functionality to detect and analyze WiFi hotspots
// and network access points in the surrounding area.
package scanner

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

// AccessPoint represents a detected WiFi access point or network hotspot.
type AccessPoint struct {
	SSID       string
	BSSID      string
	Signal     int // Signal strength in dBm
	Channel    int
	Encryption string
	LastSeen   time.Time
}

// String returns a human-readable representation of the access point.
func (ap AccessPoint) String() string {
	return fmt.Sprintf("%-32s %-20s %4d dBm  ch%-3d  %s",
		ap.SSID, ap.BSSID, ap.Signal, ap.Channel, ap.Encryption)
}

// Scanner manages the discovery and tracking of nearby access points.
type Scanner struct {
	mu      sync.RWMutex
	points  map[string]*AccessPoint
	iface   string
	timeout time.Duration
}

// New creates a new Scanner instance bound to the given network interface.
func New(iface string, timeout time.Duration) *Scanner {
	return &Scanner{
		points:  make(map[string]*AccessPoint),
		iface:   iface,
		timeout: timeout,
	}
}

// Scan performs a single scan pass and updates the internal access point list.
// Returns an error if the scan could not be initiated.
func (s *Scanner) Scan() error {
	// Validate that the interface exists before attempting a scan.
	if _, err := net.InterfaceByName(s.iface); err != nil {
		return fmt.Errorf("interface %q not found: %w", s.iface, err)
	}

	// TODO: integrate platform-specific scanning (e.g. iwlist / CoreWLAN).
	// For now, return a not-implemented sentinel so callers can handle gracefully.
	return fmt.Errorf("platform scan not yet implemented for interface %q", s.iface)
}

// Add inserts or updates an access point in the scanner's registry.
func (s *Scanner) Add(ap AccessPoint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ap.LastSeen = time.Now()
	s.points[ap.BSSID] = &ap
}

// Prune removes access points that have not been seen within the scanner timeout.
func (s *Scanner) Prune() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-s.timeout)
	removed := 0

	for bssid, ap := range s.points {
		if ap.LastSeen.Before(cutoff) {
			delete(s.points, bssid)
			removed++
		}
	}

	return removed
}

// List returns a snapshot of all currently tracked access points,
// sorted by signal strength (strongest first).
func (s *Scanner) List() []AccessPoint {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]AccessPoint, 0, len(s.points))
	for _, ap := range s.points {
		result = append(result, *ap)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Signal > result[j].Signal
	})

	return result
}

// Count returns the number of access points currently tracked.
func (s *Scanner) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.points)
}

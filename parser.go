package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// parsePhase tracks which section of the input we are currently in.
type parsePhase int

const (
	phaseAnts  parsePhase = iota // expecting ant count (skipping comments)
	phaseRooms                   // reading rooms (and ##start/##end directives)
	phaseLinks                   // reading links
)

// ParseInput reads a lem-in input file and returns a fully validated Farm.
// It never prints; all problems are reported via the returned error.
func ParseInput(fileName string) (*Farm, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	return ParseReader(file)
}

// ParseReader reads lem-in input from any io.Reader and returns a fully
// validated Farm. This is the core parsing logic.
func ParseReader(r io.Reader) (*Farm, error) {
	scanner := bufio.NewScanner(r)
	farm := &Farm{
		Rooms: make(map[string]*Room),
	}

	phase := phaseAnts
	lineNum := 0
	pendingStart := false // ##start waiting for  rooms
	pendingEnd := false   // ##end waiting for  rooms
	hasStart := false
	hasEnd := false
	linkSet := make(map[string]bool) // forbid dublicate links

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// ── Handle directives (##start / ##end) ──────────────────────
		if strings.HasPrefix(line, "##") {
			if line == "##start" || line == "##end" {
				if phase == phaseAnts {
					return nil, fmt.Errorf("line %d: %s directive cannot appear before ant count", lineNum, line)
				}
				if phase == phaseLinks {
					return nil, fmt.Errorf("line %d: %s directive found in links section", lineNum, line)
				}
				if line == "##start" {
					if hasStart || pendingStart { // if start already exists or pending
						return nil, fmt.Errorf("line %d: Duplicate ##start directive", lineNum)
					}
					pendingStart = true
				} else {
					if hasEnd || pendingEnd { // if end already exists or pending
						return nil, fmt.Errorf("line %d: duplicate ##end directive", lineNum)
					}
					pendingEnd = true
				}
				continue
			}
			// Unknown ## directive — treat as a comment.
			continue
		}

		// ── Regular comments ─────────────────────────────────────────
		if strings.HasPrefix(line, "#") {
			continue
		}

		// ── Empty lines ──────────────────────────────────────────────
		if strings.TrimSpace(line) == "" {
			return nil, fmt.Errorf("line %d: empty line is not allowed", lineNum)
		}

		// ── Phase: ant count ─────────────────────────────────────────
		if phase == phaseAnts {
			ants, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid number of ants %q", lineNum, line)
			}
			if ants <= 0 {
				return nil, fmt.Errorf("line %d: number of ants must be > 0, got %d", lineNum, ants)
			}
			farm.AntsCount = ants
			phase = phaseRooms // change status to rooms
			continue
		}

		// ── Determine if this line is a room or a link ───────────────
		tokens := strings.Fields(line)

		if len(tokens) == 3 {
			// ── Room line ────────────────────────────────────────────
			if phase == phaseLinks {
				return nil, fmt.Errorf("line %d: room definition %q found after links section started", lineNum, line)
			}

			name := tokens[0]

			// Validate room name.
			if strings.HasPrefix(name, "L") {
				return nil, fmt.Errorf("line %d: room name %q cannot start with 'L'", lineNum, name)
			}
			if strings.Contains(name, "-") {
				return nil, fmt.Errorf("line %d: room name %q cannot contain '-'", lineNum, name)
			}

			// Parse coordinates.
			x, err := strconv.Atoi(tokens[1])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid x coordinate %q for room %q", lineNum, tokens[1], name)
			}
			y, err := strconv.Atoi(tokens[2])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid y coordinate %q for room %q", lineNum, tokens[2], name)
			}

			// Check duplicate.
			if _, exists := farm.Rooms[name]; exists {
				return nil, fmt.Errorf("line %d: duplicate room name %q", lineNum, name)
			}

			room := &Room{Name: name, x: x, y: y}
			farm.Rooms[name] = room

			// Apply pending ##start / ##end directives.
			if pendingStart {
				farm.Start = name
				hasStart = true
				pendingStart = false
			}
			if pendingEnd {
				farm.End = name
				hasEnd = true
				pendingEnd = false
			}
			continue
		}

		// If a directive is pending but the next non-comment line is NOT a room,
		// that's an error.
		if pendingStart || pendingEnd {
			pending := "##start"
			if pendingEnd {
				pending = "##end"
			}
			return nil, fmt.Errorf("line %d: expected a room after %s, got %q", lineNum, pending, line)
		}

		// ── Link line ────────────────────────────────────────────────
		// A link has exactly one dash separating two non-empty names.
		// We use strings.Cut for a single split.
		left, right, ok := strings.Cut(line, "-")
		if !ok || left == "" || right == "" {
			return nil, fmt.Errorf("line %d: invalid link format %q", lineNum, line)
		}
		// Ensure there's no extra content (no spaces allowed in a link line).
		if strings.ContainsAny(line, " \t") {
			return nil, fmt.Errorf("line %d: link line %q must not contain spaces", lineNum, line)
		}

		phase = phaseLinks

		// Validate rooms exist.
		if _, ok := farm.Rooms[left]; !ok {
			return nil, fmt.Errorf("line %d: link references unknown room %q", lineNum, left)
		}
		if _, ok := farm.Rooms[right]; !ok {
			return nil, fmt.Errorf("line %d: link references unknown room %q", lineNum, right)
		}

		// Self-link check.
		if left == right {
			return nil, fmt.Errorf("line %d: self-link %q is not allowed", lineNum, line)
		}

		// Canonical key for dedup (alphabetical order).
		if left > right {
			left, right = right, left
		}
		key := left + "|" + right

		if linkSet[key] {
			// Duplicate link — silently skip (no duplicate adjacency).
			continue
		}
		linkSet[key] = true

		farm.Rooms[left].Links = append(farm.Rooms[left].Links, right)
		farm.Rooms[right].Links = append(farm.Rooms[right].Links, left)
	}

	// ── Post-scan checks ─────────────────────────────────────────────
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	if farm.AntsCount == 0 {
		return nil, fmt.Errorf("no ant count found in input")
	}
	if pendingStart {
		return nil, fmt.Errorf("unexpected end of input: ##start without a following room")
	}
	if pendingEnd {
		return nil, fmt.Errorf("unexpected end of input: ##end without a following room")
	}
	if !hasStart {
		return nil, fmt.Errorf("no ##start room defined")
	}
	if !hasEnd {
		return nil, fmt.Errorf("no ##end room defined")
	}

	if farm.Start == farm.End {
		return nil, fmt.Errorf("invalid map: start and end rooms cannot be the same")
	}

	return farm, nil
}

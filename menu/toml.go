package menu

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tamnd/vigo/event"
)

// LoadTOML reads a menu tree from r. The format is a minimal TOML
// subset: each top-level menu starts with [[menu]], each row inside a
// menu starts with [[menu.items]]. Recognized keys are title, hotkey,
// cmd, shortcut, sep. Unknown keys are ignored.
//
//	[[menu]]
//	title = "File"
//	hotkey = 0
//
//	  [[menu.items]]
//	  title = "New"
//	  cmd = 101
//	  shortcut = "F3"
//
//	  [[menu.items]]
//	  sep = true
//
// Returns the slice of top-level Items (each with Children populated)
// or an error pointing at the first malformed line.
func LoadTOML(r io.Reader) ([]Item, error) {
	scanner := bufio.NewScanner(r)
	var menus []Item
	var current *Item
	var item *Item
	var inItems bool
	lineno := 0

	flush := func() {
		if item != nil && current != nil {
			current.Children = append(current.Children, *item)
			item = nil
		}
	}
	flushMenu := func() {
		flush()
		if current != nil {
			menus = append(menus, *current)
			current = nil
		}
	}

	for scanner.Scan() {
		lineno++
		raw := scanner.Text()
		line := stripComment(raw)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		switch line {
		case "[[menu]]":
			flushMenu()
			current = &Item{}
			inItems = false
		case "[[menu.items]]":
			if current == nil {
				return nil, fmt.Errorf("line %d: [[menu.items]] before [[menu]]", lineno)
			}
			flush()
			item = &Item{}
			inItems = true
		default:
			key, val, err := parseKV(line)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineno, err)
			}
			target := current
			if inItems {
				target = item
			}
			if target == nil {
				return nil, fmt.Errorf("line %d: %s outside [[menu]] or [[menu.items]]", lineno, key)
			}
			if err := assignField(target, key, val); err != nil {
				return nil, fmt.Errorf("line %d: %w", lineno, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	flushMenu()
	return menus, nil
}

// SaveTOML writes menus to w in the format LoadTOML accepts.
func SaveTOML(w io.Writer, menus []Item) error {
	bw := bufio.NewWriter(w)
	for i, m := range menus {
		if i > 0 {
			if _, err := bw.WriteString("\n"); err != nil {
				return err
			}
		}
		if _, err := bw.WriteString("[[menu]]\n"); err != nil {
			return err
		}
		if err := writeItemFields(bw, m, false); err != nil {
			return err
		}
		for _, child := range m.Children {
			if _, err := bw.WriteString("\n  [[menu.items]]\n"); err != nil {
				return err
			}
			if err := writeItemFields(bw, child, true); err != nil {
				return err
			}
		}
	}
	return bw.Flush()
}

func writeItemFields(bw *bufio.Writer, it Item, indent bool) error {
	prefix := ""
	if indent {
		prefix = "  "
	}
	if it.Sep {
		_, err := fmt.Fprintf(bw, "%ssep = true\n", prefix)
		return err
	}
	if _, err := fmt.Fprintf(bw, "%stitle = %q\n", prefix, it.Title); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(bw, "%shotkey = %d\n", prefix, it.Hotkey); err != nil {
		return err
	}
	if it.Cmd != event.CmdNone {
		if _, err := fmt.Fprintf(bw, "%scmd = %d\n", prefix, uint16(it.Cmd)); err != nil {
			return err
		}
	}
	if it.Shortcut != "" {
		if _, err := fmt.Fprintf(bw, "%sshortcut = %q\n", prefix, it.Shortcut); err != nil {
			return err
		}
	}
	return nil
}

func stripComment(s string) string {
	inString := false
	for i, r := range s {
		switch r {
		case '"':
			inString = !inString
		case '#':
			if !inString {
				return s[:i]
			}
		}
	}
	return s
}

func parseKV(line string) (string, string, error) {
	k, v, ok := strings.Cut(line, "=")
	if !ok {
		return "", "", fmt.Errorf("expected key = value, got %q", line)
	}
	key := strings.TrimSpace(k)
	val := strings.TrimSpace(v)
	if key == "" {
		return "", "", fmt.Errorf("empty key in %q", line)
	}
	return key, val, nil
}

func assignField(it *Item, key, val string) error {
	switch key {
	case "title":
		s, err := parseString(val)
		if err != nil {
			return err
		}
		it.Title = s
	case "hotkey":
		n, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("hotkey: %w", err)
		}
		it.Hotkey = n
	case "cmd":
		n, err := strconv.ParseUint(val, 10, 16)
		if err != nil {
			return fmt.Errorf("cmd: %w", err)
		}
		it.Cmd = event.CommandID(n)
	case "shortcut":
		s, err := parseString(val)
		if err != nil {
			return err
		}
		it.Shortcut = s
	case "sep":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("sep: %w", err)
		}
		it.Sep = b
	default:
		// Unknown keys are ignored to keep the format forward-compatible.
	}
	return nil
}

func parseString(val string) (string, error) {
	if len(val) < 2 || val[0] != '"' || val[len(val)-1] != '"' {
		return "", fmt.Errorf("expected quoted string, got %q", val)
	}
	return val[1 : len(val)-1], nil
}

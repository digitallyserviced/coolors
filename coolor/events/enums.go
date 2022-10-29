package events

import "strconv"

// EnumName associates an enum value with its name for printing.
type EnumName struct {
	V uint32
	S string
}

// EnumString converts a flag-based enum value into its string representation.
func EnumString(v uint32, names []EnumName, goSyntax bool) string {
	s := ""
	for _, n := range names {
		if v&n.V == n.V && (n.V != 0 || v == 0) {
			if len(s) > 0 {
				s += "+"
			}
			if goSyntax {
				s += "imap."
			}
			s += n.S
			if v &= ^n.V; v == 0 {
				return s
			}
		}
	}
	if len(s) > 0 {
		s += "+"
	}
	return s + "0x" + strconv.FormatUint(uint64(v), 16)
}

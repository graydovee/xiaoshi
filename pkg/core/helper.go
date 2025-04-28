package core

import (
	"github.com/CuteReimu/onebot"
	"strings"
)

func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n':
		return true
	}
	return false
}

type argType int

const (
	argNo argType = iota
	argSingle
	argQuoted
)

func ExtractPlainText(m onebot.MessageChain) string {
	sb := strings.Builder{}
	for _, val := range m {
		if val.GetMessageType() == "text" {
			text := val.(*onebot.Text)
			sb.WriteString(text.String())
		}
	}
	return sb.String()
}

// ParseShell 将指令转换为指令参数.
// modified from https://github.com/mattn/go-shellwords
func ParseShell(s string) []string {
	var args []string
	buf := strings.Builder{}
	var escaped, doubleQuoted, singleQuoted, backQuote bool
	backtick := ""

	got := argNo

	for _, r := range s {
		if escaped {
			buf.WriteRune(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf.WriteRune(r)
			} else {
				escaped = true
			}
			continue
		}

		if isSpace(r) {
			if singleQuoted || doubleQuoted || backQuote {
				buf.WriteRune(r)
				backtick += string(r)
			} else if got != argNo {
				args = append(args, buf.String())
				buf.Reset()
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted {
				backtick = ""
				backQuote = !backQuote
			}
		case '"':
			if !singleQuoted {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
			}
		case '\'':
			if !doubleQuoted {
				if singleQuoted {
					got = argSingle
				}
				singleQuoted = !singleQuoted
			}
		default:
			got = argSingle
			buf.WriteRune(r)
			if backQuote {
				backtick += string(r)
			}
		}
	}

	if got != argNo {
		args = append(args, buf.String())
	}

	return args
}

package text

import (
	"bufio"
	"strings"
)

type Alignment struct {
	sep          string
	isWhiteSpace bool
	linePrefix   string
}

func NewAlignment(sep string, prefix string) *Alignment {
	return &Alignment{
		sep:          sep,
		isWhiteSpace: strings.TrimSpace(sep) == "",
		linePrefix:   prefix,
	}
}

func (a *Alignment) Format(text string) string {
	if text == "" {
		return ""
	}

	var builder strings.Builder
	lines := []string{}

	err := a.ForeachLine(text, func(line string) {
		if line != "" {
			lines = append(lines, line)
			return
		}
		// 按照空行分隔为多块处理
		if len(lines) > 0 {
			builder.WriteString(a.FormatLines(lines))
			lines = []string{}
			builder.WriteString("\n\n") // 前面不是空行，输出2个'\n'
		} else {
			builder.WriteByte('\n') // 前面是空行，当前也是空行，输出1个'\n'
		}
	})
	if err != nil {
		return text
	}
	if len(lines) > 0 {
		builder.WriteString(a.FormatLines(lines))
		if text[len(text)-1] == '\n' {
			// if a.charIs(text, len(text)-1, '\n') {
			builder.WriteByte('\n') // 保留最后的'\n'
		}
	}
	return builder.String()
}

func (a *Alignment) charIs(s string, index int, ch byte) bool {
	return s[index] == ch
}

func (a *Alignment) ForeachLine(text string, fn func(string)) error {
	r := strings.NewReader(text)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fn(strings.TrimSpace(scanner.Text()))
	}
	return scanner.Err()
}

func (a *Alignment) FormatLines(lines []string) string {
	maxP1Len := 0
	content := make([]*Line, 0, len(lines))
	for _, line := range lines {
		p1, p2, hasSep := a.split(line, a.sep)
		if len(p1) > maxP1Len {
			maxP1Len = len(p1)
		}
		content = append(content, NewLine(p1, p2, hasSep, a))
	}

	var builder strings.Builder
	end := len(content) - 1
	for index, line := range content {
		builder.WriteString(line.format(maxP1Len))
		if index < end {
			builder.WriteByte('\n')
		}
	}
	return builder.String()
}

func (a *Alignment) split(s, sep string) (string, string, bool) {
	index := strings.Index(s, sep)
	if index == -1 {
		return strings.TrimSpace(s), "", false
	} else {
		return strings.TrimSpace(s[:index]), strings.TrimSpace(s[index+len(sep):]), true
	}
}

type Line struct {
	p1     string
	p2     string
	hasSep bool
	align  *Alignment
}

func NewLine(p1, p2 string, b bool, a *Alignment) *Line {
	return &Line{p1, p2, b, a}
}

func (l *Line) format(maxsize int) string {
	if !l.hasSep {
		return l.p1
	}

	var builder strings.Builder
	// prefix
	builder.WriteString(l.align.linePrefix)
	// p1
	builder.WriteString(l.p1)
	for i := maxsize - len(l.p1); i > 0; i-- {
		builder.WriteByte(' ')
	}
	// sep
	if !l.align.isWhiteSpace {
		builder.WriteByte(' ')
	}
	builder.WriteString(l.align.sep)
	if !l.align.isWhiteSpace {
		builder.WriteByte(' ')
	}
	// p2
	builder.WriteString(l.p2)
	return builder.String()
}

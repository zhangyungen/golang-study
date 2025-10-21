package tslog

import (
	"fmt"
	"runtime"
	"strings"
)

func GetSimplifiedStack(err error) string {
	stack := fmt.Sprintf("%+v", err.Error())
	lines := strings.Split(stack, "\n")
	if len(lines) < 3 {
		const depth = 32
		var pcs [depth]uintptr
		n := runtime.Callers(2, pcs[:]) // 跳过当前函数和调用者
		if n == 0 {
			return ""
		}
		frames := runtime.CallersFrames(pcs[:n])
		stack = formatFrames(frames)
		lines = make([]string, 0)
		lines = append(lines, err.Error())
		lines1 := strings.Split(stack, "\n")
		lines = append(lines, lines1...)
	}
	var simplified []string
	for _, line := range lines {
		if !strings.Contains(line, "/middleware") && !strings.Contains(line, "goexit") && !strings.Contains(line, "/gin-gonic") && !strings.Contains(line, runtime.GOROOT()) {
			simplified = append(simplified, line)
		}
	}

	if len(simplified) > 32 {
		simplified = simplified[:32]
	}
	return strings.Join(simplified, "\n")
}

func formatFrames(frames *runtime.Frames) string {
	var sb strings.Builder
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		sb.WriteString(fmt.Sprintf("\n%s:%d %s", frame.File, frame.Line, frame.Function))
	}
	return sb.String()
}

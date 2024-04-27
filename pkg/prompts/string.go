package prompts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func StringPrompt(label string, stdin io.Reader) string {
	var s string
	r := bufio.NewReader(stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}

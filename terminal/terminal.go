package terminal

import (
	"fmt"
	"strings"
)

type Terminal struct {
	RootStage *TerminalStage
}

func NewTerminal(guideText string, initFunc ...InitFunc) *Terminal {
	if len(initFunc) == 0 {
		initFunc = append(initFunc, nil)
	}
	return &Terminal{
		RootStage: NewTerminalStage(
			"",
			"",
			"",
			WithInitFunc(initFunc[0]),
			WithGuideText(guideText),
		),
	}
}

// AddStage will add a stage to terminal.
// @param usagePath: the usage path of this stage, likes: "/a/b/c"
func (t *Terminal) AddStage(stage *TerminalStage, usagePath ...string) error {
	path := "/"
	if len(usagePath) > 0 {
		path = usagePath[0]
	}
	path = reviseUsagePath(path)[1:] // remove the first "/"
	usages := strings.Split(path, "/")

	parent := t.RootStage
	for _, usage := range usages[1:] {
		if child, ok := parent.Children[usage]; ok {
			parent = child
		} else {
			return fmt.Errorf("cannot find stage: %s", path)
		}
	}
	return parent.AddChild(stage)
}

func (t *Terminal) Run() error {
	_, err := t.RootStage.Run(nil)
	if err != nil {
		return err
	}
	return nil
}

func reviseUsagePath(s string) string {
	s = strings.TrimSpace(s)
	for strings.Contains(s, "//") {
		s = strings.Replace(s, "//", "/", -1)
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	if strings.HasSuffix(s, "/") && s != "/" {
		s = s[:len(s)-1]
	}
	return s
}

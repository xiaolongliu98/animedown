package terminal

import (
	"fmt"
	"maps"
	"strings"
)

type Terminal struct {
	rootStage   *TerminalStage
	globalStage map[string]func() *TerminalStage
}

func NewTerminal(guideText string, initFunc ...InitFunc) *Terminal {
	if len(initFunc) == 0 {
		initFunc = append(initFunc, nil)
	}
	return &Terminal{
		rootStage: NewTerminalStage(
			"",
			"",
			"",
			WithInitFunc(initFunc[0]),
			WithGuideText(guideText),
		),
		globalStage: maps.Clone(DefaultStageFactory),
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

	parent := t.rootStage
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
	if err := t.applyGlobalStage(); err != nil {
		return err
	}
	_, err := t.rootStage.Run(nil)
	if err != nil {
		return err
	}
	return nil
}

// AddGlobalStage will add a global stage to terminal.
func (t *Terminal) AddGlobalStage(usage string, stageGenerator func() *TerminalStage) {
	t.globalStage[usage] = stageGenerator
}

func (t *Terminal) applyGlobalStage() error {
	q := []*TerminalStage{t.rootStage}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		for _, child := range cur.Children {
			if !child.IsLeaf {
				q = append(q, child)
			}
		}

		for _, f := range t.globalStage {
			stage := f()
			stage.isGlobal = true
			if err := cur.AddChild(stage); err != nil {
				return err
			}
		}
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

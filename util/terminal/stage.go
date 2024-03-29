package terminal

import (
	"animedown/util"
	"fmt"
	"strings"
)

// Hook Func Turn: Run -> GuideFunc -> InitFunc -> Workloop(if not leaf) -> ExitFunc
// 													/  |  \
//                                              Child.Run1,2,3...

type (
	ExitCode int

	InitFunc  func(this *TerminalStage, args []string) (ExitCode, error)
	GuideFunc func(this *TerminalStage)
	ExitFunc  func(this *TerminalStage, exitCode ExitCode, err error) (ExitCode, error)
)

const (
	ExitCodeFatal        ExitCode = -2 // fatal error, will exit program.
	ExitCodeError        ExitCode = -1 // error, will return to parent stage, and print error.
	ExitCodeOK           ExitCode = 0  // ok, will return to parent stage.
	ExitCodeQuitParent   ExitCode = 1  // quit parent, will return to parent's parent stage.
	ExitCodeReturnToRoot ExitCode = 2  // return to root, will return to root stage.
	ExitCodeRecallSelf   ExitCode = 3  // return to parent's parent stage, and then recall parent stage with new args.
	ExitCodeRecallChild  ExitCode = 4  // return to parent stage, and then recall current stage with new args.
	ExitCodeExit         ExitCode = 5  // exit program, will exit program.
	ExitCodeCurrentQuit  ExitCode = 6  // stage quit signal.

	RecallSelfArgsKey = "recallSelfArgs"
)

type TerminalStage struct {
	Usage        string // like: "s"
	UsageExplain string // like: "Search anime by name, and then download it."
	Format       string // like: "s <anime name>"

	Children map[string]*TerminalStage
	Parent   *TerminalStage

	// if IsLeaf, then this stage is a leaf stage,
	// and it will not have any children.
	// When user input this stage's usage, it will call InitFunc at once,
	// and then return to its parent stage.
	IsLeaf   bool
	isGlobal bool
	// InitFunc will be called when user input this stage's usage.
	// You can use this func to do something, likes: print guide text, do your task, etc.
	InitFunc
	GuideFunc
	ExitFunc
	GuideText string

	noEntryGuide                     bool
	noClearScreen                    bool
	noPrintDefaultGuideFuncGuideText bool
	noPrintDefaultGuideFuncUsageText bool
	noPrintDefaultCMDUsage           bool
	data                             map[string]interface{}
}

// NewTerminalStage will create a new TerminalStage.
// @guideFunc: empty means default guide func, if you want to modify it, you can pass a func here.
func NewTerminalStage(usage string, usageExplain string, format string, opts ...StageOption) *TerminalStage {
	t := &TerminalStage{
		Usage:        usage,
		UsageExplain: usageExplain,
		Format:       format,
		Children:     make(map[string]*TerminalStage),
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.GuideFunc == nil {
		t.GuideFunc = defaultGuideFunc
	}
	return t
}

func (c *TerminalStage) AddChild(child ...*TerminalStage) error {
	for _, ch := range child {
		if err := c.addChild(ch); err != nil {
			return err
		}
	}
	return nil
}

func (c *TerminalStage) addChild(child *TerminalStage) error {
	if c.IsLeaf {
		return fmt.Errorf("cannot add child to a leaf stage")
	}

	c.Children[child.Usage] = child
	child.Parent = c
	return nil
}

// AddRecallSelfChild will add a self child stage to this stage.
// This means that in the current stage, you can recall itself by inputting its usage.
func (c *TerminalStage) AddRecallSelfChild() error {
	t := NewTerminalStage(c.Usage, c.UsageExplain, c.Format,
		WithInitFunc(func(this *TerminalStage, args []string) (ExitCode, error) {
			// add cmd usage to args for recall
			args = append([]string{this.Usage}, args...)
			// save args to data
			this.Set(RecallSelfArgsKey, args)
			return ExitCodeRecallSelf, nil
		}),
		WithLeafStage(),
	)
	return c.addChild(t)
}

// Run will workloop this stage and its children. If this stage is a leaf stage,
// it will only call its InitFunc, and then return to its parent stage.
// @return: exit code, error
func (c *TerminalStage) Run(initArgs []string) (ExitCode, error) {
	// 1. init data
	if c.IsRoot() {
		// root stage, init data
		if c.data == nil {
			c.data = make(map[string]interface{})
		}
	} else {
		// not root stage, inherit data
		c.data = c.Parent.data
	}

	// 2. run leaf stage
	if c.IsLeaf {
		return runInitFuncSafely(c, initArgs)
	}

	// 3. clean screen
	if !c.noClearScreen {
		util.ClearScreen()
	}
	// 4. print guide text
	runGuideFuncSafely(c)
	// 5. print default guide text
	exitCode, err := runInitFuncSafely(c, initArgs)
	if exitCode != ExitCodeOK {
		// 6. exit func
		return runExitFuncSafely(c, exitCode, err)
	}
	// 6. workloop
	exitCode, err = c.workloop()

	// 7. exit func
	return runExitFuncSafely(c, exitCode, err)
}

func (c *TerminalStage) workloop() (ExitCode, error) {
	for {
		var (
			line string
			err  error
		)
		line, err = util.ReadLine()
		if err != nil {
			return 0, err
		}
		line = strings.TrimSpace(line)
		line = util.StandardizeCMDString(line)

		if line == "" {
			continue
		}

		args := strings.Split(line, " ")

	TagRecall: // here args should be set
		usage := args[0]
		args = args[1:]

		if child, ok := c.Children[usage]; ok {
			var exitCode ExitCode
			exitCode, err = child.Run(args)
			switch exitCode {
			case ExitCodeFatal:
				return ExitCodeFatal, err

			case ExitCodeError:
				fmt.Println(err.Error())
				continue

			case ExitCodeQuitParent:
				return ExitCodeCurrentQuit, err
			case ExitCodeCurrentQuit:
				c.printGuideTextForChildQuit(child)

			case ExitCodeOK:
				c.printGuideTextForChildQuit(child)

			case ExitCodeReturnToRoot:
				if !c.IsRoot() {
					return ExitCodeReturnToRoot, err
				}
				// is root
				c.printGuideTextForChildQuit(child)

			case ExitCodeRecallSelf:
				return ExitCodeRecallChild, err
			case ExitCodeRecallChild:
				// reset args
				args = child.GetDelete(RecallSelfArgsKey).([]string)
				goto TagRecall

			case ExitCodeExit:
				return ExitCodeExit, err

			default:
				panic("unknown exit code")
			}
		} else {
			err = fmt.Errorf("unknown command: %s", usage)
			fmt.Println(err.Error())
		}
	}
}

func (c *TerminalStage) printGuideTextForChildQuit(child *TerminalStage) {
	if child.IsLeaf {
		return
	}
	if !c.noClearScreen {
		util.ClearScreen()
	}
	// print guide text
	runGuideFuncSafely(c)
}

func (c *TerminalStage) PrintChildrenUsage(noLimit ...bool) {
	limit := true
	if len(noLimit) > 0 && noLimit[0] {
		limit = false
	}

	var (
		subs   []*TerminalStage
		leaves []*TerminalStage
	)
	for _, child := range c.Children {
		if child.IsLeaf {
			leaves = append(leaves, child)
		} else {
			subs = append(subs, child)
		}
	}

	// Firstly, print sub children's usage.
	if len(subs) > 0 {
		fmt.Println(" - Sub-Terminals:")
		for _, child := range subs {
			fmt.Print("   ")
			fmt.Printf("[-] %s \t %s\n", child.Format, child.UsageExplain)
		}
	}
	// Secondly, print leaf children's usage.
	if len(leaves) > 0 {
		fmt.Println()
		fmt.Println(" - Commands:")
		countNonDefaultStage := 0
		for _, child := range leaves {
			if !child.IsGlobalStage() {
				fmt.Print("   ")
				fmt.Printf("[*] %s \t %s\n", child.Format, child.UsageExplain)
				countNonDefaultStage++
			}
		}

		if len(leaves)-countNonDefaultStage == 0 {
			return
		}

		if countNonDefaultStage > 0 {
			fmt.Println()
		}

		for _, child := range leaves {
			if child.IsGlobalStage() {
				if limit && c.noPrintDefaultCMDUsage && child.IsDefaultStage() {
					continue
				}
				fmt.Print("   ")
				fmt.Printf("[*] %s \t %s\n", child.Format, child.UsageExplain)
			}
		}
	}
}

func (c *TerminalStage) RunDefaultGuideFunc() {
	defaultGuideFunc(c)
}

func (c *TerminalStage) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *TerminalStage) Get(key string) interface{} {
	return c.data[key]
}

// GetDelete will get value and delete it from data.
func (c *TerminalStage) GetDelete(key string) interface{} {
	value := c.data[key]
	delete(c.data, key)
	return value
}

func (c *TerminalStage) Exist(key string) bool {
	_, ok := c.data[key]
	return ok
}

func (c *TerminalStage) Delete(key string) {
	delete(c.data, key)
}

func (c *TerminalStage) ClearScreen() {
	util.ClearScreen()
}

func (c *TerminalStage) IsRoot() bool {
	return c.Parent == nil
}

// IsGlobalStage will return true if this stage is a global stage.
func (c *TerminalStage) IsGlobalStage() bool {
	return c.isGlobal
}

// IsDefaultStage will return true if this stage is a default stage.
func (c *TerminalStage) IsDefaultStage() bool {
	_, ok := DefaultStages[c.Usage]
	return ok
}

func runInitFuncSafely(stage *TerminalStage, args []string) (ExitCode, error) {
	if stage.InitFunc != nil {
		return stage.InitFunc(stage, args)
	}
	return ExitCodeOK, nil
}

func runGuideFuncSafely(stage *TerminalStage) {
	if stage.GuideFunc != nil && !stage.noEntryGuide {
		stage.GuideFunc(stage)
	}
}

func defaultGuideFunc(this *TerminalStage) {
	if this.GuideText != "" && !this.noPrintDefaultGuideFuncGuideText {
		fmt.Println(this.GuideText)
	}
	if !this.noPrintDefaultGuideFuncUsageText {
		this.PrintChildrenUsage()
	}
}

func runExitFuncSafely(stage *TerminalStage, exitCode ExitCode, err error) (ExitCode, error) {
	if stage.ExitFunc != nil {
		return stage.ExitFunc(stage, exitCode, err)
	}
	return exitCode, err
}

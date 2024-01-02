package terminal

import "fmt"

const (
	DefaultCMDQuit = "q"    // quit current stage.
	DefaultCMDHelp = "help" // print help text.
	DefaultCMDRoot = "rr"   // return to root stage.
	DefaultCMDExit = "exit" // exit program.
)

var (
	DefaultStageFactory = map[string]func() *TerminalStage{
		DefaultCMDQuit: NewDefaultCMDQuitStage,
		DefaultCMDHelp: NewDefaultCMDHelpStage,
		DefaultCMDRoot: NewDefaultCMDRootStage,
		DefaultCMDExit: NewDefaultCMDExitStage,
	}

	DefaultStages = map[string]struct{}{
		DefaultCMDQuit: {},
		DefaultCMDHelp: {},
		DefaultCMDRoot: {},
		DefaultCMDExit: {},
	}
)

func NewDefaultCMDQuitStage() *TerminalStage {
	t := &TerminalStage{
		Usage:        DefaultCMDQuit,
		UsageExplain: "Quit current stage.",
		Format:       DefaultCMDQuit,
		IsLeaf:       true,
		InitFunc: func(this *TerminalStage, args []string) (ExitCode, error) {
			return ExitCodeQuitParent, nil
		},
	}
	return t
}

func NewDefaultCMDHelpStage() *TerminalStage {
	t := &TerminalStage{
		Usage:        DefaultCMDHelp,
		UsageExplain: "Print help text.",
		Format:       DefaultCMDHelp,
		IsLeaf:       true,
		InitFunc: func(this *TerminalStage, args []string) (ExitCode, error) {
			if this.Parent.GuideText != "" {
				fmt.Println(this.Parent.GuideText)
			}
			this.Parent.PrintChildrenUsage(true)
			return ExitCodeOK, nil
		},
	}
	return t
}

func NewDefaultCMDRootStage() *TerminalStage {
	t := &TerminalStage{
		Usage:        DefaultCMDRoot,
		UsageExplain: "Return to root stage.",
		Format:       DefaultCMDRoot,
		IsLeaf:       true,
		InitFunc: func(this *TerminalStage, args []string) (ExitCode, error) {
			return ExitCodeReturnToRoot, nil
		},
	}
	return t
}

func NewDefaultCMDExitStage() *TerminalStage {
	t := &TerminalStage{
		Usage:        DefaultCMDExit,
		UsageExplain: "Exit program.",
		Format:       DefaultCMDExit,
		IsLeaf:       true,
		InitFunc: func(this *TerminalStage, args []string) (ExitCode, error) {
			return ExitCodeExit, nil
		},
	}
	return t
}

package terminal

type StageOption func(*TerminalStage)

func WithCustomGuideFunc(f GuideFunc) StageOption {
	return func(stage *TerminalStage) {
		stage.GuideFunc = f
	}
}

func WithInitFunc(f InitFunc) StageOption {
	return func(stage *TerminalStage) {
		stage.InitFunc = f
	}
}

func WithNoClearScreenForDefaultGuideFunc() StageOption {
	return func(stage *TerminalStage) {
		stage.noClearScreen = true
	}
}

func WithNoPrintDefaultUsageText() StageOption {
	return func(stage *TerminalStage) {
		stage.noPrintDefaultUsageText = true
	}
}

func WithLeafStage() StageOption {
	return func(stage *TerminalStage) {
		stage.IsLeaf = true
	}
}

func WithGuideText(guideText string) StageOption {
	return func(stage *TerminalStage) {
		stage.GuideText = guideText
	}
}

func WithExitFunc(f ExitFunc) StageOption {
	return func(stage *TerminalStage) {
		stage.ExitFunc = f
	}
}

func WithNoEntryGuide() StageOption {
	return func(stage *TerminalStage) {
		stage.noEntryGuide = true
	}
}

func WithNoPrintDefaultGuideText() StageOption {
	return func(stage *TerminalStage) {
		stage.noPrintDefaultGuideText = true
	}
}

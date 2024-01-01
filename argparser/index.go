package argparser

import "fmt"

type ArgContext struct {
	Flags map[string]string // -x : <value>
	Args  []string          // <value>
}

func (a *ArgContext) Check(minNumArgs int, mustFlags ...string) error {
	if len(a.Args) < minNumArgs {
		return fmt.Errorf("invalid args, at least %d args", minNumArgs)
	}

	for _, flag := range mustFlags {
		if _, ok := a.Flags[flag]; !ok {
			return fmt.Errorf("invalid args, flag %s not found", flag)
		}
	}
	return nil
}

func Parse(args []string, hasCommand bool) (*ArgContext, error) {
	if hasCommand {
		args = args[1:]
	}

	ctx := &ArgContext{
		Flags: make(map[string]string),
		Args:  make([]string, 0),
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg[0] == '-' {
			// flag
			// check valid
			if len(args) <= i+1 || args[i+1][0] == '-' {
				return nil, fmt.Errorf("invalid flag %s, value unset", arg)
			}
			ctx.Flags[arg] = args[i+1]
			i++

		} else {
			// arg
			ctx.Args = append(ctx.Args, arg)
		}
	}

	return ctx, nil
}

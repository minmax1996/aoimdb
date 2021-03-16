package commands

import "errors"

//BaseCommand base command fimpemeting commander interface
type BaseCommand struct {
	name         string
	desc         string
	example      string
	minArgsCount int
	maxArgsCount int
	callback     func(string, ...string) error
}

//NewBaseCommand constructor for private fields
func NewBaseCommand(name string, minArgsCount, maxArgsCount int, desc, example string, callback func(string, ...string) error) *BaseCommand {
	return &BaseCommand{
		name:         name,
		desc:         desc,
		example:      example,
		minArgsCount: minArgsCount,
		maxArgsCount: maxArgsCount,
		callback:     callback,
	}
}

//Name returns Name associated with command
func (bc BaseCommand) Name() string {
	return bc.name
}

//ShowDescription returns description associated with command
func (bc BaseCommand) ShowDescription() string {
	return bc.desc
}

//ShowExample returns example of use this command
func (bc BaseCommand) ShowExample() string {
	return bc.example
}

//ValidateUserInput validates args for this command, in base command only checks for count (args including name)
func (bc BaseCommand) ValidateUserInput(args []string) error {
	if bc.minArgsCount <= len(args) && len(args) <= bc.maxArgsCount {
		return nil
	}
	return errors.New("wrong args count")
}

//CallWithArgs invokes callback associated with command
func (bc BaseCommand) CallWithArgs(args ...string) error {
	return bc.callback(bc.name, args...)
}

package tools

func (t Tool) LaunchCommand(model string) string {
	return t.Command + " --model " + model
}
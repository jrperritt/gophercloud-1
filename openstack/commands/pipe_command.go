package commands

type Pipeable struct{}

func (c *Pipeable) HandlePipe(item string) (interface{}, error) {
	return item, nil
}

func (c *Pipeable) PipeFieldOptions() []string {
	return []string{"id"}
}

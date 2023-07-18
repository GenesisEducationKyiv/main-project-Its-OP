package commands

type LogCommand struct {
	LogData  string
	LogLevel string
}

func NewLogCommand(data string, level string) *LogCommand {
	return &LogCommand{LogData: data, LogLevel: level}
}

func (c *LogCommand) Reset() {
	c.LogData = ""
}

func (c *LogCommand) String() string {
	return c.LogData
}

func (c *LogCommand) ProtoMessage() {}

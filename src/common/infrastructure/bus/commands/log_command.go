package commands

type LogCommand struct {
	Log string `protobuf:"bytes,1,opt,name=Log"`
}

func (c *LogCommand) Reset() {
	c.Log = ""
}

func (c *LogCommand) String() string {
	return c.Log
}

func (c *LogCommand) ProtoMessage() {}

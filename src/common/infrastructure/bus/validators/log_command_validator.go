package validators

import (
	"btcRate/common/infrastructure/bus/commands"
	"fmt"
)

type NewLogCommandValidator struct {
}

func Validate(command commands.LogCommand) error {
	if len(command.LogAttributes)%2 != 0 {
		return fmt.Errorf("attributes array must have an even length: %s", command.LogAttributes)
	}

	for index, attr := range command.LogAttributes {
		if index%2 == 1 {
			if _, ok := attr.(string); ok {
				return fmt.Errorf("attribute keys must be strings: %s", command.LogAttributes)
			}
		}
	}
}

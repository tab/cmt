package helpers

import (
	"fmt"

	"cmt/internal/app/cli/model"
	"cmt/internal/config"
)

// IsVersionCmd checks if the command is a version command
func IsVersionCmd(cmd string) bool {
	return model.Contains(model.CmdVersion, cmd)
}

func RenderVersion() {
	fmt.Println(config.Version + " (" + config.AppName + ")")
}

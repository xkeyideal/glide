package action

import (
	"github.com/xkeyideal/glide/msg"
)

// Name prints the name of the package, according to the glide.yaml file.
func Name() {
	conf := EnsureConfig()
	msg.Puts(conf.Name)
}

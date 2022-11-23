package action

import (
	"github.com/moby/buildkit/client/llb"
)

type Action interface {
	Execute(llb.State) llb.State
}

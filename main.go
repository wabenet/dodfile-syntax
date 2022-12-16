package main

import (
	"fmt"
	"os"

	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/wabenet/dodfile-syntax/pkg/build"
)

func main() {
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), build.Build); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

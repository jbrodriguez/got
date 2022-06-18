package cmd

import "got/cli/core/recorder"

type On struct {
}

func (o *On) Run(ctx *Context) error {
	r := recorder.CreateRecorder(ctx.DataDir)
	return r.AddOn()
}

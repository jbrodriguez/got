package cmd

import (
	"got/cli/core/recorder"
)

type Add struct {
	Activity string `arg:"" help:"activity you want to add"`
}

func (a *Add) Run(ctx *Context) error {
	r := recorder.CreateRecorder(ctx.DataDir)
	return r.AddActivity(a.Activity)
}

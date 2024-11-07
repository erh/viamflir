package main

import (
	"context"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/module"

	"github.com/erh/viamflir"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {

	ctx := context.Background()

	myMod, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	err = myMod.AddModelFromRegistry(ctx, camera.API, viamflir.FlirModel)
	if err != nil {
		return err
	}

	err = myMod.Start(ctx)
	defer myMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

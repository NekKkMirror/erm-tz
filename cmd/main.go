package main

import (
	"github.com/NekKkMirror/erm-tz/internal/app"
)

func main() {
	application := app.Init()
	application.Run()
}

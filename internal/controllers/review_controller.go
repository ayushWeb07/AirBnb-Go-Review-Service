package controllers

import (
	renderPkg "github.com/unrolled/render"
)

var render *renderPkg.Render

func init() {
	render = renderPkg.New()
}

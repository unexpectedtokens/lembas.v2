package view

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func RenderView(viewToRender templ.Component, ctx context.Context, w io.Writer) error {
	return viewToRender.Render(ctx, w)
}

// func RenderComponent(viewToRender templ.Component, ctx context.Context, w io.Writer) error {}

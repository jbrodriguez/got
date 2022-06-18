package model

type Renderer interface {
	Add(a *Activity)
	Render()
}

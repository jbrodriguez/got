package model

type Section interface {
	Add(a *Activity)
	Render()
}

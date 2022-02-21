package entity

import (
	"cranky/video"
	"cranky/video/texture"
	"cranky/video/window"
)

type EntityType interface {
	// you should never have to call this manually.
	Init(win *window.Window)
	// draw it!!
	Draw(win *window.Window)
}

type Entity struct {
	Mesh    video.Mesh
	Texture texture.Texture
}

func (e Entity) Init(win *window.Window) {
	println(e.Mesh.Vao.IndexCount)
}

func (e Entity) Draw(win *window.Window) {
}

func NewEntity(msh video.Mesh, tex texture.Texture) Entity {
	return Entity{Mesh: msh, Texture: tex}
}

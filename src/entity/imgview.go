package entity

import (
	"hornyaf/src/data/fs"
	"hornyaf/src/event"
	"hornyaf/src/util"
	"hornyaf/src/video"
	"hornyaf/src/video/texture"
	"hornyaf/src/video/window"
	"math/rand"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type ImageViewer struct {
	EntityType
	Mesh    video.Mesh
	Texture texture.Texture
	Dir     fs.Directory
}

func NewImageViewer(file string) ImageViewer {
	prg, err := video.LoadShaderProgram("resources/shaders/default.vert", "resources/shaders/default.frag")
	util.Assert(err)
	tex, err := texture.TextureFromFile(file)
	util.Assert(err)
	imgView := ImageViewer{
		Mesh: video.MeshFromZmsh(
			"resources/mesh/square.zmsh",
			*prg),
		Texture: tex,
		Dir:     fs.OpenDirectory(fs.GetParentDir(file)),
	}
	return imgView
}

func (e ImageViewer) Init(win *window.Window) {
	println(e.Mesh.Vao.IndexCount)
}

func (e *ImageViewer) SetTexture(newTex texture.Texture) {
	e.Texture.Destroy()
	e.Texture = newTex
}
func (e ImageViewer) resizeTexture(win window.Window) {
	tex := e.Texture
	texW, texH := tex.GetSize()
	imgRect := util.NewRect(0, 0, float32(texW), float32(texH))
	winW, winH := win.Win.GetFramebufferSize()
	winRect := util.NewRect(0, 0, float32(winW), float32(winH))
	newRect := imgRect.CenterAndFit(&winRect)
	xSize := (newRect.W) / winRect.W
	ySize := (newRect.H) / winRect.H
	gl.Uniform2f(e.Mesh.Program.GetUniformLocation("textureSize"), xSize, ySize)
}

func (e *ImageViewer) Draw(win window.Window) {
	e.resizeTexture(win)
	switch t := e.Texture.(type) {
	case texture.AnimatedTexture:
		delta := win.GetDelta()
		t.Animate(delta)
		if int(t.FrameNum) > len(t.Frames) {
			t.Frames = append(t.Frames, *t.DecodeFrame(int(t.FrameNum)))
		}
		t.Bind(gl.TEXTURE0)
		t.SetUniform(e.Mesh.Program.GetUniformLocation("spriteTexture"))
		e.Mesh.Draw()
		t.UnBind()
		//win.Win.SetTitle(fmt.Sprintf("Time: %f %f Frame: %d / %d", delta, t.AnimTime, t.FrameNum, len(t.Frames)))
		e.Texture = t
	case texture.StaticTexture:
		t.Bind(gl.TEXTURE0)
		t.SetUniform(e.Mesh.Program.GetUniformLocation("spriteTexture"))
		e.Mesh.Draw()
		//tW, tH := t.GetSize()
		//win.Win.SetTitle(fmt.Sprintf("SIZE: %d,%d, ID: %d", tW, tH, t.Id))
		t.UnBind()
	default:
		win.Win.SetTitle("NO IMAGE")
		gl.ClearColor(1.0, 0.0, 1.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	}
}

func (e *ImageViewer) EventHandle(evt event.Event) {
	switch t := evt.(type) {
	case event.DropEvent:
		if len(t.Files) > 0 {

			if fs.IsDir(t.Files[0]) {
				e.Dir = fs.OpenDirectory(t.Files[0])
				newTex, err := texture.TextureFromFile(e.Dir.RootPath + "/" + e.Dir.GetFile())
				if err == nil {
					e.SetTexture(newTex)
				} else {
					panic(err)
				}
			} else {
				dirPaths := strings.Split(t.Files[0], "/")
				dirPath := strings.Join(dirPaths[:len(dirPaths)-1], "/")

				newTex, err := texture.TextureFromFile(t.Files[0])
				if err != nil {
					panic(err)
				}
				println(newTex)
				e.SetTexture(newTex)
				e.Dir = fs.OpenDirectory(dirPath)
				e.Dir.SetIndex(e.Dir.IndexOf(dirPaths[len(dirPaths)-1]))
			}
		}
	case event.KeyEvent:
		if t.Action == 0 {
			break
		}
		switch t.Key {
		case glfw.KeyRight:
			e.Dir.SetIndex(e.Dir.Index + 1)
			println(e.Dir.Index)
			newTex, err := texture.TextureFromFile(e.Dir.RootPath + "/" + e.Dir.GetFile())
			if err == nil {
				e.SetTexture(newTex)
			} else {
				println(err)
			}
		case glfw.KeyLeft:
			e.Dir.SetIndex(e.Dir.Index - 1)
			println(e.Dir.Index)
			newTex, err := texture.TextureFromFile(e.Dir.RootPath + "/" + e.Dir.GetFile())
			if err == nil {
				e.SetTexture(newTex)
			} else {
				println(err)
			}
		case glfw.KeyR:
			e.Dir.SetIndex(rand.Intn(len(e.Dir.Files) - 1))
			println(e.Dir.Index)
			newTex, err := texture.TextureFromFile(e.Dir.RootPath + "/" + e.Dir.GetFile())
			if err == nil {
				e.SetTexture(newTex)
			} else {
				println(err)
			}
		}
	}
}

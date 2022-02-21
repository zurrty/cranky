package window

import (
	"github.com/zurrty/cranky/data"
	"github.com/zurrty/cranky/event"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	Win           *glfw.Window
	Running       bool
	lastFrameTime float64
	deltaTime     float64
	Config        data.ConfigFile
}

/*
func (win Window) SetViewport(rect util.Rect) {
	gl.Viewport(int32(rect.X), int32(rect.Y), int32(rect.W), int32(rect.H))
}*/

func InitOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	println("OpenGL version: " + version)
}

func (win *Window) CreateEventLoop() *event.EventLoop {
	evl := event.EventLoop{}

	win.Win.SetDropCallback(evl.OnDrop)
	win.Win.SetRefreshCallback(evl.OnRefresh)
	win.Win.SetKeyCallback(evl.OnKey)
	return &evl
}

func (win *Window) GetTime() float64 {
	return glfw.GetTime()
}

func (win *Window) GetDelta() float64 {
	return win.deltaTime //win.GetTime() - win.lastFrameTime
}

// Swaps GL buffers and also resets delta time
func (win *Window) SwapBuffers() {
	win.Win.SwapBuffers()
	t := win.GetTime()
	win.deltaTime = t - win.lastFrameTime
	win.lastFrameTime = t
}

func (win *Window) OnResize() {
	winW, winH := win.Win.GetFramebufferSize()
	gl.Viewport(0, 0, int32(winW), int32(winH))
}

func CreateWindow(title string, width int, height int) (Window, error) {
	glfw.Init()
	glfw.WindowHint(glfw.Samples, 8)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	win.MakeContextCurrent()
	InitOpenGL()
	return Window{Win: win, lastFrameTime: 0}, nil
}

func (win Window) Destroy() {
	win.Win.Destroy()
}

func (win *Window) Quit() {
	win.Running = false
	println("Quitting...")
}

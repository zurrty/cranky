package event

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Event interface{}

type RefreshEvent struct {
	Win *glfw.Window
}

type DropEvent struct {
	Win   *glfw.Window
	Files []string
}

type KeyEvent struct {
	Win      *glfw.Window
	Key      glfw.Key
	Scancode int
	Action   glfw.Action
	Mods     glfw.ModifierKey
}

type EventLoop struct {
	Events []Event
}

func (evl *EventLoop) Poll() []Event {
	glfw.PollEvents()
	defer evl.clearQueue()
	return evl.Events
}

func (evl *EventLoop) clearQueue() {
	evl.Events = []Event{}
}

func (evl *EventLoop) queueEvent(evt Event) {
	evl.Events = append(evl.Events, evt)
}

func (evl *EventLoop) OnRefresh(w *glfw.Window) {
	evl.queueEvent(RefreshEvent{w})
}
func (evl *EventLoop) OnDrop(w *glfw.Window, files []string) {
	evl.queueEvent(DropEvent{w, files})
	println(files)
}

func (evl *EventLoop) OnKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	evl.queueEvent(KeyEvent{w, key, scancode, action, mods})
}

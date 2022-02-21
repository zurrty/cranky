package texture

import (
	"image"
	"image/draw"
	"image/gif"
	"io"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type AnimatedTexture struct {
	Texture
	Width     uint32
	Height    uint32
	FrameNum  uint32
	AnimTime  float64
	Duration  uint32
	bgFrame   int
	gifData   *gif.GIF
	prevFrame *image.RGBA
	Frames    []StaticTexture
	Delays    []int
}

func createAnimatedTexture(r io.Reader) AnimatedTexture {
	gifData, err := gif.DecodeAll(r)
	if err != nil {
		panic(err)
	}
	texture := AnimatedTexture{
		gifData:  gifData,
		Frames:   []StaticTexture{},
		Delays:   gifData.Delay,
		FrameNum: 0,
		Duration: 0,
		bgFrame:  int(gifData.BackgroundIndex),
		Width:    uint32(gifData.Config.Width),
		Height:   uint32(gifData.Config.Height),
		AnimTime: 0.0,
	}
	texture.Frames = []StaticTexture{*texture.DecodeFrame(0)}
	return texture
}

func (tex *AnimatedTexture) GetFrame() *StaticTexture {
	if int(tex.FrameNum) >= len(tex.Frames) {
		newFrame := tex.DecodeFrame(int(tex.FrameNum))
		tex.appendFrame(*newFrame)
		return newFrame
	}
	return &tex.Frames[tex.FrameNum]
}

func (tex AnimatedTexture) GetSize() (uint32, uint32) {
	return tex.Width, tex.Height
}
func (tex *AnimatedTexture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	frm := tex.GetFrame()
	gl.BindTexture(frm.Target, frm.Id)
	frm.TexUnit = texUnit
}
func (tex *AnimatedTexture) appendFrame(frame StaticTexture) {
	tex.Frames = append(tex.Frames, frame)
}
func (tex *AnimatedTexture) DecodeFrame(idx int) *StaticTexture {
	if idx > len(tex.gifData.Image) {
		return nil
	}
	var frame *image.RGBA
	thisImg := tex.gifData.Image[idx]
	println(tex.bgFrame, tex.gifData.Disposal[idx])

	switch tex.gifData.Disposal[idx] {
	case 0:
		if tex.bgFrame < len(tex.Delays) {
			img := tex.gifData.Image[tex.bgFrame]
			frame = image.NewRGBA(img.Bounds())
			draw.Src.Draw(frame, img.Bounds(), img, img.Bounds().Min)
			draw.FloydSteinberg.Draw(frame, thisImg.Bounds(), thisImg, thisImg.Bounds().Min)
		}
	case 1, 3:
		if tex.prevFrame != nil {
			frame = image.NewRGBA(tex.prevFrame.Bounds())
			draw.Src.Draw(frame, tex.prevFrame.Bounds(), tex.prevFrame, tex.prevFrame.Bounds().Min)
			draw.Over.Draw(frame, thisImg.Bounds(), thisImg, thisImg.Bounds().Min)
		}
	}
	if frame == nil {
		frame = image.NewRGBA(thisImg.Bounds())
		draw.Src.Draw(frame, thisImg.Bounds(), thisImg, image.Point{0, 0})
	}
	tex.prevFrame = frame
	tex.Duration += uint32(tex.gifData.Delay[idx])
	newTex := staticFromRGBA(frame)
	return &newTex
}

func (tex *AnimatedTexture) Animate(delta float64) {
	tex.AnimTime += delta
	if int(tex.AnimTime*100) > tex.Delays[tex.FrameNum] {
		newFrame := tex.FrameNum + 1
		tex.AnimTime = 0
		if int(newFrame) >= len(tex.Frames) {
			if int(newFrame) >= len(tex.gifData.Image) {
				newFrame = 0
				tex.AnimTime = 0.0 //math.Max(0.0, tex.AnimTime-tex.Duration)
			} else {
				tex.DecodeFrame(int(newFrame))
			}
		}

		tex.FrameNum = newFrame
	}
}

func (tex *AnimatedTexture) UnBind() {
	frm := tex.GetFrame()
	frm.TexUnit = 0
	gl.BindTexture(frm.Target, 0)
}

func (tex AnimatedTexture) SetUniform(uniformLoc int32) error {
	frm := tex.GetFrame()
	return frm.SetUniform(uniformLoc)
}

func (tex AnimatedTexture) Destroy() {
	for _, img := range tex.Frames {
		img.Destroy()
	}
}

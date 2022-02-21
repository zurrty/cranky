package texture

import (
	"bytes"
	"hornyaf/src/data/fs"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture interface {
	SetUniform(uniformLoc int32) error
	GetSize() (uint32, uint32)
	Destroy()
}
type StaticTexture struct {
	Id      uint32
	Target  uint32
	TexUnit uint32
	Width   uint32
	Height  uint32
}

func decodeStaticImage(r io.Reader) image.Image {
	img, _, _ := image.Decode(r)
	return img
}
func staticFromRGBA(rgba *image.RGBA) StaticTexture {
	var texId uint32
	gl.GenTextures(1, &texId)
	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Bounds().Size().X)
	height := int32(rgba.Bounds().Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := StaticTexture{
		Id:     texId,
		Target: target,
		Width:  uint32(width),
		Height: uint32(height),
	}

	texture.Bind(gl.TEXTURE0)

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	gl.TexParameteri(texture.Target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.Target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.Target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(texture.Target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	// gl.GenerateMipmap(texture.id)
	texture.UnBind()

	return texture
}
func createStaticTexture(img image.Image) StaticTexture {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	return staticFromRGBA(rgba)
}

func TextureFromFile(filePath string) (Texture, error) {
	var tex Texture
	println(filePath)
	buf, err := fs.BReadFile(filePath)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(buf)
	filePath = strings.ToLower(filePath)
	if strings.HasSuffix(filePath, ".gif") {
		tex = createAnimatedTexture(reader)
		println(tex.GetSize())
	} else {
		tex = createStaticTexture(decodeStaticImage(reader))
		println(tex.GetSize())
	}
	return tex, nil
}

func (tex StaticTexture) GetSize() (uint32, uint32) {
	return tex.Width, tex.Height
}
func (tex *StaticTexture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.Target, tex.Id)
	tex.TexUnit = texUnit
}

func (tex *StaticTexture) UnBind() {
	tex.TexUnit = 0
	gl.BindTexture(tex.Target, 0)
}

func (tex StaticTexture) SetUniform(uniformLoc int32) error {
	if tex.TexUnit == 0 {
		println("failed to bind texture")
		return nil
	}
	gl.Uniform1i(uniformLoc, int32(tex.TexUnit-gl.TEXTURE0))
	return nil
}

func (tex StaticTexture) Destroy() {
	gl.DeleteTextures(1, &tex.Id)
}

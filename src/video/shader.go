package video

import (
	"fmt"
	"hornyaf/src/data/fs"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
)

// https://github.com/cstegel/opengl-samples-golang/blob/master/basic-shaders/gfx/shader.go
// ^^^ thank u :D :D :D :D :D :D

type Shader struct {
	id uint32
}

type ShaderProgram struct {
	id      uint32
	shaders []*Shader
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.id)
}

func (prog ShaderProgram) Delete() {
	for _, shader := range prog.shaders {
		shader.Delete()
	}
	gl.DeleteProgram(prog.id)
}

func (prog *ShaderProgram) Attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(prog.id, shader.id)
		prog.shaders = append(prog.shaders, shader)
	}
}

func (prog ShaderProgram) Use() {
	gl.UseProgram(prog.id)
}

func (prog ShaderProgram) Link() error {
	gl.LinkProgram(prog.id)
	return getGlError(prog.id, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog, "PROGRAM::LINKING_FAILURE")
}

func (prog ShaderProgram) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(prog.id, gl.Str(name+"\x00"))
}

func CreateShader(src string, shaderType uint32) (*Shader, error) {
	sId := gl.CreateShader(shaderType)
	glSrc, freeFn := gl.Strs(src + "\x00")
	defer freeFn()
	gl.ShaderSource(sId, 1, glSrc, nil)
	gl.CompileShader(sId)
	err := getGlError(sId, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog, "SHADER::COMPILE_FAILURE::")
	if err != nil {
		return nil, err
	}
	return &Shader{id: sId}, nil
}
func ShaderFromFile(filePath string, shaderType uint32) (*Shader, error) {
	src, err := fs.ReadFile(filePath)
	if err != nil {
		panic(err)
		//return nil, err
	}
	return CreateShader(src, shaderType)
}

func CreateShaderProgram(shaders ...*Shader) (*ShaderProgram, error) {
	prog := ShaderProgram{id: gl.CreateProgram()}
	prog.Attach(shaders...)
	if err := prog.Link(); err != nil {
		return nil, err
	}
	return &prog, nil
}

func ReadShaderProgram(vertSrc string, fragSrc string) (*ShaderProgram, error) {
	prog := &ShaderProgram{id: gl.CreateProgram()}
	vert, err := CreateShader(vertSrc, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	frag, err := CreateShader(fragSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	prog.Attach(vert, frag)
	if err := prog.Link(); err != nil {
		return nil, err
	}
	return prog, nil
}

func LoadShaderProgram(vertSrcPath string, fragSrcPath string) (*ShaderProgram, error) {
	vert, err := ShaderFromFile(vertSrcPath, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	frag, err := ShaderFromFile(fragSrcPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	return CreateShaderProgram(vert, frag)
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glId uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glId, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glId, gl.INFO_LOG_LENGTH, &logLength)
		if int(logLength) > 0 {
			log := gl.Str(strings.Repeat("\x00", int(logLength)))
			getObjInfoLogFn(glId, logLength, nil, log)

			return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
		}
		return fmt.Errorf("%s: %s", failMsg, "")
	}

	return nil
}

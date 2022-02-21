package video

import (
	"cranky/data/fs"
	"strconv"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type VertexArray struct {
	IndexCount uint32
	Id         uint32
}
type Mesh struct {
	Vao     VertexArray
	Program ShaderProgram
}

func (msh Mesh) Draw() {
	msh.Program.Use()
	gl.BindVertexArray(msh.Vao.Id)
	gl.DrawElements(gl.TRIANGLES, int32(msh.Vao.IndexCount), gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.BindVertexArray(0)
}

func MakeMesh(vao VertexArray, prgm ShaderProgram) Mesh {
	//println("vao id: ", vao.Id)

	return Mesh{Vao: vao, Program: prgm}
}
func MakeVAO(points []float32, indices []uint32) VertexArray {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 2*4
	var offset int = 0

	// position
	gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, stride, uintptr(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4
	gl.BindVertexArray(0)
	return VertexArray{Id: VAO, IndexCount: uint32(len(indices))}
}

func ReadZmsh(filePath string) ([]float32, []uint32) {
	src, _ := fs.ReadFile(filePath)
	arr := strings.Split(src, "\n")
	mode := 0 // 0 = verts 1 = indices
	vrts := []float32{}
	idxs := []uint32{}

	for l := 1; l < len(arr); l++ {
		line := strings.Trim(arr[l], " \r")
		if strings.Contains(line, "[VRTS]") {
			mode = 0
			continue
		}
		if strings.Contains(line, "[IDXS]") {
			mode = 1
			continue
		}
		//println(line)
		if mode == 1 {
			idx_str := strings.Split(line, ",")
			for i := 0; i < len(idx_str); i++ {
				val, err := strconv.ParseUint(idx_str[i], 10, 32)
				if err != nil {
					println(err)
				}
				idxs = append(idxs, uint32(val))
			}
		} else {
			vert_str := strings.Split(line, ",")
			for i := 0; i < len(vert_str); i++ {
				val, err := strconv.ParseFloat(vert_str[i], 64)
				if err != nil {
					panic(err)
				}
				//print(val, " = ", vert_str[i], " | ")
				vrts = append(vrts, float32(val))
			}
			//print("\n")
		}
	}
	return vrts, idxs
}

func MeshFromZmsh(filePath string, prgm ShaderProgram) Mesh {
	verts, indices := ReadZmsh(filePath)
	return MakeMesh(MakeVAO(verts, indices), prgm)
}

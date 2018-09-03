package util

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strconv"
	"unsafe"
)

type Vertex struct {
	Position  mgl32.Vec3
	Normal    mgl32.Vec3
	TexCoords mgl32.Vec2
	Tangent   mgl32.Vec3
	Bitangent mgl32.Vec3
}
type Texture struct {
	Id   uint32
	Type string
	Path string
}

type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
	Textures []Texture
	VAO      uint32
	vbo      uint32
	ebo      uint32
}

func NewMesh(vertices []Vertex, indices []uint32, textures []Texture) *Mesh {
	mesh := Mesh{}
	mesh.Vertices = vertices
	mesh.Indices = indices
	mesh.Textures = textures
	mesh.setupMesh()
	return &mesh
}
func (mesh *Mesh) Draw(shader Shader) {
	diffuseNr := uint32(1)
	specularNr := uint32(1)
	normalNr := uint32(1)
	heightNr := uint32(1)
	for i := 0; i < len(mesh.Textures); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		var number string
		name := mesh.Textures[i].Type
		if name == "texture_diffuse" {
			diffuseNr++
			number = strconv.Itoa(int(diffuseNr))
		} else if name == "texture_specular" {
			specularNr++
			number = strconv.Itoa(int(specularNr))
		} else if name == "texture_normal" {
			normalNr++
			number = strconv.Itoa(int(normalNr))
		} else if name == "texture_height" {
			heightNr++
			number = strconv.Itoa(int(heightNr))
		}
		shader.SetInt(name+number, int32(i))
		gl.BindTexture(gl.TEXTURE_2D, mesh.Textures[i].Id)
	}
}

func (mesh *Mesh) setupMesh() {

	gl.GenVertexArrays(1, &mesh.VAO)
	gl.GenBuffers(1, &mesh.vbo)
	gl.GenBuffers(1, &mesh.ebo)

	gl.BindVertexArray(mesh.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Vertices) * 56, &mesh.Vertices[0], gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, &mesh.Indices[0], gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 56, nil)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 56, unsafe.Pointer(uintptr(12)))

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 56, unsafe.Pointer(uintptr(24)))

	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 3, gl.FLOAT, false, 56, unsafe.Pointer(uintptr(32)))

	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 3, gl.FLOAT, false, 56, unsafe.Pointer(uintptr(44)))

	gl.BindVertexArray(0)
}

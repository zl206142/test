package util

import (
	"io/ioutil"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexPath string, fragmentPath string, geometryPath string) *Shader {
	s := Shader{}
	var geometry uint32
	vertex := compileShaderFromFile(vertexPath, gl.VERTEX_SHADER)
	fragment := compileShaderFromFile(fragmentPath, gl.FRAGMENT_SHADER)
	if geometryPath != "" {
		geometry = compileShaderFromFile(geometryPath, gl.GEOMETRY_SHADER)
	}
	s.ID = gl.CreateProgram()
	gl.AttachShader(s.ID, vertex)
	gl.AttachShader(s.ID, fragment)
	if geometryPath != "" {
		gl.AttachShader(s.ID, geometry)
	}
	gl.LinkProgram(s.ID)
	checkCompileErrors(s.ID, "PROGRAM")
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)
	if geometryPath != "" {
		gl.DeleteShader(geometry)
	}
	return &s
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) getLocation(name string) int32 {
	return gl.GetUniformLocation(s.ID, gl.Str(name+"\x00"))
}

func (s *Shader) SetBool(name string, value bool) {
	i := int32(0)
	if value {
		i = 1
	}
	gl.Uniform1i(s.getLocation(name), i)
}

func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(s.getLocation(name), value)
}
func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(s.getLocation(name), value)
}
func (s *Shader) SetMat4(name string, mat *mgl32.Mat4)  {
	gl.UniformMatrix4fv(s.getLocation(name),1,false,&mat[0])
}
func (s *Shader) SetVec3v(name string, vec3s mgl32.Vec3) {
	gl.Uniform3fv(s.getLocation(name),1,&vec3s[0])
}
func (s *Shader) SetVec3(name string, x,y,z float32) {
	gl.Uniform3f(s.getLocation(name),x,y,z)
}

func readString(path string) string {
	if "" == path {
		return ""
	}
	buffer, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(buffer)
}
func compileShaderFromFile(path string, shaderType uint32) uint32 {
	source:=readString(path)
	return compileShader(source+"\x00",shaderType)
}
func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	sources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, sources, nil)
	free()
	gl.CompileShader(shader)
	checkCompileErrors(shader, "VERTEX&FRAGMENT")
	return shader
}
func checkCompileErrors(shader uint32, tp string) {
	var success, logLength int32
	if tp != "PROGRAM" {
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
		if success == gl.FALSE {
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
			fmt.Println("ERROR::SHADER_COMPILATION_ERROR of type: ", tp, log)
		}
	} else {
		gl.GetProgramiv(shader, gl.LINK_STATUS, &success)
		if success == gl.FALSE {
			gl.GetProgramiv(shader, gl.INFO_LOG_LENGTH, &logLength)
			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetProgramInfoLog(shader, logLength, nil, gl.Str(log))
			fmt.Println("ERROR::PROGRAM_LINKING_ERROR of type: ", tp, log)
		}
	}

}

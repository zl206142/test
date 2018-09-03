package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
	"./util"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
	"math"
)

var (
	width    = 800
	height   = 600
	vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
		-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
		-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	}

	deltaTime = float64(0)
	lastFrame = float64(0)

	camera     = util.NewCamera1(mgl32.Vec3{0, 0, 3})
	lastX      = float32(width / 2)
	lastY      = float32(height / 2)
	firstMouse = true

	lightPos = mgl32.Vec3{1.2, 1, 2}
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
	fmt.Println("init")
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window := initWindow()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	//sr := util.NewShader("4.3.shader.vs", "4.3.shader.fs", "")
	lightingShader := util.NewShader("materials.vs", "materials.fs", "")
	lampShader := util.NewShader("lamp.vs", "lamp.fs", "")
	VBO, cubeVAO, lightVAO := vaboo(vertices)

	for !window.ShouldClose() {

		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		processInput(window)

		gl.ClearColor(.1, .1, .1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		lightingShader.Use()
		lightingShader.SetVec3v("light.position", lightPos)
		lightingShader.SetVec3v("viewPos", camera.Position)

		lightColor := mgl32.Vec3{
			float32(math.Sin(glfw.GetTime() * 2)),
			float32(math.Sin(glfw.GetTime() * .7)),
			float32(math.Sin(glfw.GetTime() * 1.3)),
		}
		diffuseColor := lightColor.Mul(.5)
		ambientColor := lightColor.Mul(.2)
		lightingShader.SetVec3v("light.ambient", ambientColor)
		lightingShader.SetVec3v("light.diffuse", diffuseColor)
		lightingShader.SetVec3("light.specular", 1.0, 1.0, 1.0)

		// material proSerties
		lightingShader.SetVec3("material.ambient", 1.0, 0.5, 0.31)
		lightingShader.SetVec3("material.diffuse", 1.0, 0.5, 0.31)
		lightingShader.SetVec3("material.specular", 0.5, 0.5, 0.5) // specular lighting doesn't have full effect on this object's material
		lightingShader.SetFloat("material.shininess", 32.0)

		projection := mgl32.Perspective(float32(mgl32.DegToRad(camera.Zoom)), float32(width)/float32(height), .1, 100)
		view := camera.GetViewMatrix()

		lightingShader.SetMat4("projection",&projection)
		lightingShader.SetMat4("view",&view)


		model:=mgl32.Ident4()
		lightingShader.SetMat4("model",&model)

		gl.BindVertexArray(cubeVAO)
		gl.DrawArrays(gl.TRIANGLES,0,36)

		lampShader.Use()
		lampShader.SetMat4("projection",&projection)
		lampShader.SetMat4("view",&view)
		translate:=mgl32.Translate3D(lightPos[0],lightPos[1],lightPos[2])
		scale:=mgl32.Scale3D(.2,.2,.2)
		transform:= translate.Mul4(scale)
		lampShader.SetMat4("model",&transform)

		gl.BindVertexArray(lightVAO)
		gl.DrawArrays(gl.TRIANGLES,0,36)

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &cubeVAO)
	gl.DeleteVertexArrays(1, &lightVAO)
	gl.DeleteBuffers(1, &VBO)
}

func initWindow() *glfw.Window {
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(width, height, "learnOpenGL", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		x := float32(xpos)
		y := float32(ypos)
		if firstMouse {
			lastX = x
			lastY = y
			firstMouse = false
		}
		xoffset := x - lastX
		yoffset := lastY - y
		lastX = x
		lastY = y

		camera.ProcessMouseMovement2(xoffset, yoffset)
	})
	window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		camera.ProcessMouseScroll(float32(yoff))
	})
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	return window
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.ProcessKeyboard(util.FORWARD, float32(deltaTime))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.ProcessKeyboard(util.BACKWARD, float32(deltaTime))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.ProcessKeyboard(util.LEFT, float32(deltaTime))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.ProcessKeyboard(util.RIGHT, float32(deltaTime))
	}
}
func vaboo(vertices []float32) (uint32, uint32, uint32) {
	var VBO, cubeVAO, lightVAO uint32
	gl.GenVertexArrays(1, &cubeVAO)
	gl.GenBuffers(1, &VBO)

	gl.BindVertexArray(cubeVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 24, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 24, unsafe.Pointer(uintptr(12)))
	gl.EnableVertexAttribArray(1)

	gl.GenVertexArrays(1, &lightVAO)
	gl.BindVertexArray(lightVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 24, nil)
	gl.EnableVertexAttribArray(0)

	return VBO, cubeVAO, lightVAO
}

func getTexture() (uint32, uint32) {
	return util.CreateTexture("container.jpg", gl.TEXTURE0), util.CreateTexture("awesomeface.png", gl.TEXTURE1)
}

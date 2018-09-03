package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"unsafe"
	"./util"
	"os"
	"image"
	"image/draw"
	"strings"
	"image/jpeg"
	"image/png"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)

var (
	width    = 800
	height   = 600
	vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
	cubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}
	deltaTime  = float64(0)
	lastFrame  = float64(0)
	camera     = util.NewCamera1(mgl32.Vec3{0, 0, 3})
	lastX      = float32(width / 2)
	lastY      = float32(height / 2)
	firstMouse = true
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

	sr := util.NewShader("4.3.shader.vs", "4.3.shader.fs", "")
	vbo, vao, _ := vaboo(vertices)
	texture1, texture2 := getTexture()
	sr.Use()
	sr.SetInt("texture1\x00", 0)
	sr.SetInt("texture2\x00", 1)
	modeLoc := gl.GetUniformLocation(sr.ID, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(sr.ID, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(sr.ID, gl.Str("projection\x00"))

	for !window.ShouldClose() {

		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture1)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)

		sr.Use()

		p := mgl32.Perspective(float32(mgl32.DegToRad(camera.Zoom)), float32(width)/float32(height), .1, 100)
		gl.UniformMatrix4fv(projectionLoc, 1, false, &p[0])

		v := camera.GetViewMatrix()

		gl.UniformMatrix4fv(viewLoc, 1, false, &v[0])

		gl.BindVertexArray(vao)

		for i := 0; i < 10; i++ {
			angle := 20 * i
			m := mgl32.HomogRotate3D(float32(mgl32.DegToRad(float32(angle))), mgl32.Vec3{1, .3, .5})
			m = mgl32.Translate3D(cubePositions[i][0], cubePositions[i][1], cubePositions[i][2]).Mul4(m)
			gl.UniformMatrix4fv(modeLoc, 1, false, &m[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
	gl.DeleteVertexArrays(1, &vao)
	gl.DeleteBuffers(1, &vbo)
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
	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, unsafe.Pointer(uintptr(12)))
	gl.EnableVertexAttribArray(1)

	return vbo, vao, ebo
}

func getTexture() (uint32, uint32) {
	return createTexture("container.jpg", gl.TEXTURE0), createTexture("awesomeface.png", gl.TEXTURE1)
}
func createTexture(path string, n uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(n)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	m := readTexture(path)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(m.Rect.Size().X),
		int32(m.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(m.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture
}

func readTexture(path string) *image.RGBA {
	ss := strings.Split(path, ".")
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	var img image.Image
	if len(ss) > 1 {
		switch ss[len(ss)-1] {
		case "jpg":
			fallthrough
		case "jpeg":
			img, err = jpeg.Decode(f)
		case "png":
			img, err = png.Decode(f)
		default:
			img, _, err = image.Decode(f)
		}
	} else {
		img, _, err = image.Decode(f)
	}
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}

package util

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera_Movement int

const (
	FORWARD  = iota
	BACKWARD
	LEFT
	RIGHT
)

const YAW float32 = -90.0
const PITCH float32 = 0.0
const SPEED float32 = 2.5
const SENSITIVITY float32 = 0.1
const ZOOM float32 = 45.0

type Camera struct {
	Position         mgl32.Vec3
	Front            mgl32.Vec3
	Up               mgl32.Vec3
	Right            mgl32.Vec3
	WorldUp          mgl32.Vec3
	Yaw              float32
	Pitch            float32
	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func NewCamera4(position mgl32.Vec3, up mgl32.Vec3, yaw float32, pitch float32) *Camera {
	c := &Camera{Front: mgl32.Vec3{0, 0, -1}, MovementSpeed: SPEED, MouseSensitivity: SENSITIVITY, Zoom: ZOOM}
	c.Position = position
	c.WorldUp = up
	c.Yaw = yaw
	c.Pitch = pitch
	c.updateCameraVectors()
	return c
}
func NewCamera3(position mgl32.Vec3, up mgl32.Vec3, yaw float32) *Camera {
	return NewCamera4(position, up, yaw, PITCH)
}
func NewCamera2(position mgl32.Vec3, up mgl32.Vec3) *Camera {
	return NewCamera3(position, up, YAW)
}
func NewCamera1(position mgl32.Vec3) *Camera {
	return NewCamera2(position, mgl32.Vec3{0, 1, 0})
}
func NewCamera0() *Camera {
	return NewCamera1(mgl32.Vec3{0, 0, 0})
}
func NewCamera(posX, posY, posZ, upX, upY, upZ, yaw, pitch float32, ) *Camera {
	return NewCamera4(
		mgl32.Vec3{posX, posY, posZ},
		mgl32.Vec3{upX, upY, upZ},
		yaw,
		pitch)
}
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}
func (c *Camera) ProcessKeyboard(direction Camera_Movement, deltaTime float32) {
	velocity := c.MovementSpeed * deltaTime
	switch direction {
	case FORWARD:
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	case BACKWARD:
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	case LEFT:
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	case RIGHT:
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity
	c.Yaw += xoffset
	c.Pitch += yoffset
	if constrainPitch {
		if c.Pitch > 89 {
			c.Pitch = 89
		}
		if c.Pitch < -89 {
			c.Pitch = -89
		}
	}
	c.updateCameraVectors()
}
func (c *Camera) ProcessMouseMovement2(xoffset, yoffset float32) {
	c.ProcessMouseMovement(xoffset, yoffset, true)
}

func (c *Camera) ProcessMouseScroll(yoffset float32) {
	if c.Zoom >= 1 && c.Zoom <= 45 {
		c.Zoom -= yoffset
	}
	if c.Zoom <= 1 {
		c.Zoom = 1
	}
	if c.Zoom >= 45 {
		c.Zoom = 45
	}

}

func (c *Camera) updateCameraVectors() {
	var front mgl32.Vec3
	front[0] = float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch))))
	front[1] = float32(math.Sin(float64(mgl32.DegToRad(c.Pitch))))
	front[2] = float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch))))
	c.Front = front.Normalize()
	c.Right = c.Front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

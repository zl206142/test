package util

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
	"os"
	"github.com/go-gl/gl/v3.3-core/gl"
)

func CreateTexture(path string, n uint32) uint32 {
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

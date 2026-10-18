package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type userStatePlayer struct {
	Name     string
	Head     image.Image
	JoinedAt time.Time
}

// Render separately from skin downloads so the complete layout works offline.
func renderUserStateImage(players []userStatePlayer, now time.Time) (image.Image, error) {
	regular, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	bold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return nil, err
	}
	nameFace := truetype.NewFace(bold, &truetype.Options{Size: 22})
	small := truetype.NewFace(regular, &truetype.Options{Size: 14})
	defer nameFace.Close()
	defer small.Close()
	const width, margin, gap, cardHeight = 960, 32, 16, 96
	const cardWidth = (width - margin*2 - gap) / 2
	rows := (len(players) + 1) / 2
	height := 80 + rows*(cardHeight+gap)
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// A subtle indigo glow gives the background depth without distracting from names.
			glow := (width - x) * (height - y) * 12 / width / height
			canvas.SetRGBA(x, y, color.RGBA{uint8(15 + glow), uint8(20 + glow/2), uint8(32 + glow*2), 255})
		}
	}
	white := color.RGBA{237, 243, 251, 255}
	muted := color.RGBA{147, 163, 184, 255}
	green := color.RGBA{113, 235, 176, 255}
	roundedRect(canvas, image.Rect(margin, 20, margin+182, 60), 20, color.RGBA{29, 57, 51, 255})
	roundedRect(canvas, image.Rect(margin+19, 35, margin+29, 45), 5, green)
	stateText(canvas, small, fmt.Sprintf("%02d ONLINE", len(players)), margin+42, 20, 40, green)
	for i, player := range players {
		x := margin + (i%2)*(cardWidth+gap)
		y := 80 + (i/2)*(cardHeight+gap)
		roundedRect(canvas, image.Rect(x, y+4, x+cardWidth, y+cardHeight+4), 16, color.RGBA{9, 14, 25, 255})
		roundedRect(canvas, image.Rect(x, y, x+cardWidth, y+cardHeight), 16, color.RGBA{53, 65, 84, 255})
		roundedRect(canvas, image.Rect(x+1, y+1, x+cardWidth-1, y+cardHeight-1), 15, color.RGBA{29, 39, 56, 255})
		roundedRect(canvas, image.Rect(x+16, y+16, x+80, y+80), 12, color.RGBA{50, 65, 84, 255})
		if player.Head != nil {
			xdraw.NearestNeighbor.Scale(canvas, image.Rect(x+20, y+20, x+76, y+76), player.Head, player.Head.Bounds(), draw.Over, nil)
		} else {
			initials := []rune(strings.ToUpper(player.Name))
			if len(initials) > 2 {
				initials = initials[:2]
			}
			stateText(canvas, nameFace, string(initials), x+29, y+16, 64, white)
		}
		roundedRect(canvas, image.Rect(x+68, y+68, x+84, y+84), 8, color.RGBA{29, 39, 56, 255})
		roundedRect(canvas, image.Rect(x+71, y+71, x+81, y+81), 5, green)
		name := fitStateText(nameFace, player.Name, cardWidth-120)
		stateText(canvas, nameFace, name, x+100, y+20, 32, white)
		stateText(canvas, small, formatLoginDuration(player.JoinedAt, now), x+100, y+54, 22, muted)
	}
	return canvas, nil
}

func userStateTextBaseline(face font.Face, height int) fixed.Int26_6 {
	m := face.Metrics()
	return (fixed.I(height)-m.Ascent-m.Descent)/2 + m.Ascent
}

func stateText(dst draw.Image, face font.Face, text string, x, y, height int, c color.Color) {
	d := font.Drawer{Dst: dst, Src: image.NewUniform(c), Face: face,
		Dot: fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y) + userStateTextBaseline(face, height)}}
	d.DrawString(text)
}

func fitStateText(face font.Face, text string, maxWidth int) string {
	if font.MeasureString(face, text).Ceil() <= maxWidth {
		return text
	}
	runes := []rune(text)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		if font.MeasureString(face, string(runes)+"…").Ceil() <= maxWidth {
			return string(runes) + "…"
		}
	}
	return ""
}

func roundedRect(dst *image.RGBA, rect image.Rectangle, radius int, c color.RGBA) {
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			dx, dy := 0, 0
			if x < rect.Min.X+radius {
				dx = rect.Min.X + radius - x - 1
			}
			if x >= rect.Max.X-radius {
				dx = x - (rect.Max.X - radius)
			}
			if y < rect.Min.Y+radius {
				dy = rect.Min.Y + radius - y - 1
			}
			if y >= rect.Max.Y-radius {
				dy = y - (rect.Max.Y - radius)
			}
			if dx*dx+dy*dy <= radius*radius {
				dst.SetRGBA(x, y, c)
			}
		}
	}
}

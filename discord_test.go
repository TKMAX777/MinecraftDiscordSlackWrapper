package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

func TestUserStateTextFitsRow(t *testing.T) {
	ft, err := truetype.Parse(gomono.TTF)
	if err != nil {
		t.Fatal(err)
	}
	face := truetype.NewFace(ft, &truetype.Options{Size: 30})
	defer face.Close()
	const rowHeight = 40
	for _, text := range []string{"<Logon Users>", "tkmax777", "gjpqy_Player"} {
		t.Run(text, func(t *testing.T) {
			// Render onto a taller canvas to expose any pixels that a row would clip.
			canvas := image.NewAlpha(image.Rect(0, 0, 800, rowHeight*3))
			dr := font.Drawer{Dst: canvas, Src: image.White, Face: face,
				Dot: fixed.Point26_6{X: fixed.I(5), Y: fixed.I(rowHeight) + userStateTextBaseline(face, rowHeight)}}
			dr.DrawString(text)
			pixels := 0
			for y := 0; y < canvas.Bounds().Dy(); y++ {
				for x := 0; x < canvas.Bounds().Dx(); x++ {
					if canvas.AlphaAt(x, y).A == 0 {
						continue
					}
					pixels++
					if y < rowHeight || y >= rowHeight*2 {
						t.Fatalf("text clipped at y=%d", y-rowHeight)
					}
				}
			}
			if pixels == 0 {
				t.Fatal("text was not rendered")
			}
		})
	}
}

func TestUserStateCards(t *testing.T) {
	players := []userStatePlayer{{Name: "tkmax777"}, {Name: "gjpqy_Player"}, {Name: "A_very_long_player_name_that_must_fit"}}
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	for i := range players {
		players[i].JoinedAt = now.Add(-time.Duration(i*73+5) * time.Minute)
		// Synthetic pixel faces for an offline layout preview.
		head := image.NewRGBA(image.Rect(0, 0, 8, 8))
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				c := color.RGBA{215, uint8(160 + i*10), 125, 255}
				if y < 2 {
					c = color.RGBA{uint8(60 + i*30), 52, 65, 255}
				}
				if y == 4 && (x == 2 || x == 5) {
					c = color.RGBA{48, 55, 85, 255}
				}
				head.SetRGBA(x, y, c)
			}
		}
		players[i].Head = head
	}
	for _, count := range []int{0, 1, 2, len(players)} {
		img, err := renderUserStateImage(players[:count], now)
		if err != nil {
			t.Fatal(err)
		}
		rows := (count + 1) / 2
		if img.Bounds() != image.Rect(0, 0, 960, 80+rows*112) {
			t.Fatalf("unexpected bounds: %v", img.Bounds())
		}
		if count == len(players) && os.Getenv("ROSTER_PREVIEW_PATH") != "" {
			out, err := os.Create(os.Getenv("ROSTER_PREVIEW_PATH"))
			if err != nil {
				t.Fatal(err)
			}
			if err := png.Encode(out, img); err != nil {
				out.Close()
				t.Fatal(err)
			}
			if err := out.Close(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

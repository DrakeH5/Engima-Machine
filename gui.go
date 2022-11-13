package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	//"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var keys = [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

var rotorImg *ebiten.Image
var topbgImg *ebiten.Image

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	rotorImg, _, err = ebitenutil.NewImageFromFile("rotors.png")
	op.GeoM.Scale(-0.40, 0.40)
	topbgImg, _, err = ebitenutil.NewImageFromFile("topbg.png")
}

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if len(g.keys) == 0 {
		keyReleased = true
	}
	return nil
}

var keyReleased bool

var plugBoardLetters []string

var movingRotor bool

var op = &ebiten.DrawImageOptions{}

var oldMouseX int
var oldMouseY int

func (g *Game) Draw(screen *ebiten.Image) {
	{
		for i := 0; i < 26; i++ {
			for _, j := range g.keys {
				if keyReleased == true {
					//SEND LETTER TO ENIGMA
					keyReleased = false
				}
				if ebiten.Key.String(j) == strings.ToUpper(keys[i]) {
					//vector.DrawFilledCircle(screen, 400, 400, 100, color.RGBA{0x80, 0x00, 0x80, 0x80})
					text.Draw(screen, keys[i], mplusBigFont, i*25, 310, color.RGBA{255, 255, 0, 0xff})
				} else {
					text.Draw(screen, keys[i], mplusNormalFont, i*25, 310, color.Gray16{0xffff})
				}
			}
			if len(g.keys) == 0 {
				text.Draw(screen, keys[i], mplusNormalFont, i*25, 310, color.Gray16{0xffff})
			}
		}
		for i := 0; i < 26; i++ {
			var yPos int = i / 9
			var x int = i - int(math.Floor(float64(i/9)))*9
			var xPos int = x * 75
			text.Draw(screen, keys[i], mplusNormalFont, xPos, int(400+40*math.Floor(float64(yPos))), color.Gray16{0xffff})
			for k := 0; k < len(plugBoardLetters); k++ {
				if plugBoardLetters[k] == keys[i] {
					q := int(math.Floor(float64(k / 2)))
					r := uint8((q * 90) + 100)
					g := uint8((q * 70) + 10)
					b := uint8((q * 55))
					text.Draw(screen, keys[i], mplusNormalFont, xPos, int(400+40*math.Floor(float64(yPos))), color.RGBA{r, g, b, 0xff})
				}
			}
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true {
			xPos, yPos := ebiten.CursorPosition()
			keyXCord := (int(xPos*9) / screenWidth)
			keyYCord := (yPos / 40) - 9
			keyPosInArray := (keyYCord * 9) + keyXCord
			if keyPosInArray > 0 && keyPosInArray < len(keys) {
				plugBoardLetters = append(plugBoardLetters, keys[keyPosInArray])
			}
			fmt.Println(plugBoardLetters)
			movingRotor = true
		} else {
			movingRotor = false
		}

		for i := 0; i < 3; i++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64(i*160), 1)
			screen.DrawImage(topbgImg, options)
		}

		if movingRotor == true {
			mouseX, mouseY := ebiten.CursorPosition()
			if mouseY < 200 && mouseX < (3*160) {
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Scale(-0.40, 0.40)
				op.GeoM.Translate(float64((math.Floor(float64(mouseX/160)*160) + 105)), float64(10))
				mouseX, mouseY = int((math.Floor(float64(mouseX/160)*160) + 105)), 10
			} else {
				op.GeoM.Translate(float64(mouseX-oldMouseX), float64(mouseY-oldMouseY))
			}
			oldMouseX = mouseX
			oldMouseY = mouseY
		}
		screen.DrawImage(rotorImg, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Enigma Machine")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

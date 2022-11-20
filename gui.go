package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strconv"
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
var reflectortopImg *ebiten.Image
var emptyRotorSlotImg *ebiten.Image

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
	topbgImg, _, err = ebitenutil.NewImageFromFile("topbg.png")
	reflectortopImg, _, err = ebitenutil.NewImageFromFile("reflectortop.png")
	emptyRotorSlotImg, _, err = ebitenutil.NewImageFromFile("emptyrotorSlot.png")
	for i := 0; i < 5; i++ {
		rotorOptions[i].GeoM.Scale(-0.40-(math.Floor(float64(i/3))*-0.15), 0.40-(math.Floor(float64(i/3))*0.15))
		rotorOptions[i].GeoM.Translate(float64((i*160)+150)-(math.Floor(float64(i/4))*160), float64(math.Floor(float64(i/4))*100))
	}
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
var keyPressed string

var plugBoardLetters []string

var movingRotor bool

var rotorOptions = [5]*ebiten.DrawImageOptions{&ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}, &ebiten.DrawImageOptions{}}

var oldMouseX int
var oldMouseY int

var rotorInMotion int

var rotorNbms = [5]string{"1", "2", "3", "4", "5"}

var selectedRotor string
var err error

var rotorsRotationAmounts = [5]int{0, 0, 0, 0, 0}

var plugBoard = map[interface{}]interface{}{
	"a": " ",
	"b": " ",
	"c": " ",
	"d": " ",
	"e": " ",
	"f": " ",
	"g": " ",
	"h": " ",
	"i": " ",
	"j": " ",
	"k": " ",
	"l": " ",
	"m": " ",
	"n": " ",
	"o": " ",
	"p": " ",
	"q": " ",
	"r": " ",
	"s": " ",
	"t": " ",
	"u": " ",
	"v": " ",
	"w": " ",
	"x": " ",
	"y": " ",
	"z": " ",
}

var rightMouseClicked bool

var movingRotorStartingPosSelected int
var movingRotorStartRotationInitialValue int

func (g *Game) Draw(screen *ebiten.Image) {
	{
		for i := 0; i < 26; i++ {
			for _, j := range g.keys {
				if keyReleased == true {
					keyPressed = encrypt(strings.ToLower(ebiten.Key.String(j)))
					fmt.Println(keyPressed)
					keyReleased = false
				}
				encryptedKey := keyPressed
				if encryptedKey == keys[i] {
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
				if len(plugBoardLetters)%2 == 0 {
					plugBoard[plugBoardLetters[len(plugBoardLetters)-1]] = plugBoardLetters[len(plugBoardLetters)-2]
					plugBoard[plugBoardLetters[len(plugBoardLetters)-2]] = plugBoardLetters[len(plugBoardLetters)-1]
				}
			}
		}

		screen.DrawImage(reflectortopImg, nil)
		for i := 0; i < 3; i++ {
			options := &ebiten.DrawImageOptions{}
			options.GeoM.Translate(float64((i*160)+50), 1)
			screen.DrawImage(topbgImg, options)
			text.Draw(screen, strconv.Itoa(rotorsRotationAmounts[i]), mplusNormalFont, (i*160)+130, 230, color.White)
		}
		option := &ebiten.DrawImageOptions{}
		option.GeoM.Scale(1, 0.5)
		option.GeoM.Translate(530, 1)
		screen.DrawImage(emptyRotorSlotImg, option)
		option.GeoM.Translate(0, 110)
		screen.DrawImage(emptyRotorSlotImg, option)

		if movingRotor == true {
			mouseX, mouseY := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true {
				if mouseY < 200 {
					var clickedSlot int
					if mouseX < (3*160)+25 {
						clickedSlot = int(math.Floor(float64(mouseX-50) / 160))
					} else {
						clickedSlot = int(math.Floor(float64(mouseY/100))) + 3
					}
					if mouseX < (3*160)+25 {
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.40, 0.40)
						rotorOptions[rotorInMotion].GeoM.Translate(float64((math.Floor(float64((mouseX-50)/160)*160) + 155)), float64(10))
						oldMouseX, oldMouseY = int((math.Floor(float64((mouseX-50)/160)*160) + 155)), 10
					} else {
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.25, 0.25)
						rotorOptions[rotorInMotion].GeoM.Translate(float64((math.Floor(float64((mouseX-50)/160)*160) + 155)), math.Floor(float64(mouseY/100))*100)
						oldMouseX, oldMouseY = int((math.Floor(float64((mouseX-50)/160)*160) + 155)), int(math.Floor(float64(mouseY/100))*100)
					}
					if clickedSlot == rotorInMotion {
						movingRotor = false
						rotorNbms[clickedSlot] = selectedRotor
					} else {
						rotorOptions[clickedSlot] = rotorOptions[rotorInMotion]
						shortTermSelctedRotor := rotorNbms[clickedSlot]
						rotorNbms[clickedSlot] = selectedRotor
						selectedRotor = shortTermSelctedRotor
						rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
						rotorOptions[rotorInMotion].GeoM.Scale(-0.50, 0.50)
						rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX), float64(mouseY))
						oldRotorRotationShortTermStore := rotorsRotationAmounts[clickedSlot]
						rotorsRotationAmounts[clickedSlot] = rotorsRotationAmounts[rotorInMotion]
						rotorsRotationAmounts[rotorInMotion] = oldRotorRotationShortTermStore
					}
				}
			} else {
				rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX-oldMouseX), float64(mouseY-oldMouseY))
				oldMouseX = mouseX
				oldMouseY = mouseY
			}
		} else {
			mouseX, mouseY := ebiten.CursorPosition()
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonLeft)) == true && mouseY < 200 && mouseX < screenWidth {
				movingRotor = true
				if mouseX < (3*160)+25 {
					rotorInMotion = int(math.Floor(float64(mouseX-50) / 160))
				} else {
					rotorInMotion = int(math.Floor(float64(mouseY/100))) + 3
				}
				rotorOptions[rotorInMotion] = &ebiten.DrawImageOptions{}
				rotorOptions[rotorInMotion].GeoM.Scale(-0.50, 0.50)
				rotorOptions[rotorInMotion].GeoM.Translate(float64(mouseX), float64(mouseY))
				oldMouseX = mouseX
				oldMouseY = mouseY
				selectedRotor = rotorNbms[rotorInMotion]
				rotorNbms[rotorInMotion] = ""
			}
		}
		screen.DrawImage(rotorImg, rotorOptions[0])
		screen.DrawImage(rotorImg, rotorOptions[1])
		screen.DrawImage(rotorImg, rotorOptions[2])
		screen.DrawImage(rotorImg, rotorOptions[3])
		screen.DrawImage(rotorImg, rotorOptions[4])

		for i := 0; i < 5; i++ {
			if i < 3 {
				text.Draw(screen, rotorNbms[i], mplusNormalFont, (i*160)+150, 200, color.RGBA{10, 100, 10, 0xff})
			} else {
				text.Draw(screen, rotorNbms[i], mplusNormalFont, 550, (int(math.Floor(float64(i/4)))*100)+70, color.RGBA{100, 100, 10, 0xff})
			}
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton(ebiten.MouseButtonRight)) == true {
			xPos, yPos := ebiten.CursorPosition()
			if yPos > 205 && yPos < 235 {
				if xPos > 110 && xPos < 150 {
					rightMouseClicked = true
					movingRotorStartingPosSelected = 0
					movingRotorStartRotationInitialValue = 0
				} else if xPos > 270 && xPos < 310 {
					rightMouseClicked = true
					movingRotorStartingPosSelected = 1
					movingRotorStartRotationInitialValue = 0
				} else if xPos > 430 && xPos < 470 {
					rightMouseClicked = true
					movingRotorStartingPosSelected = 2
					movingRotorStartRotationInitialValue = 0
				}
			}
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton(ebiten.MouseButtonRight)) == true {
			rightMouseClicked = false
		}

		if rightMouseClicked == true {
			xPos, yPos := ebiten.CursorPosition()
			xPos++
			sensitivity := 30
			rotorsRotationAmounts[movingRotorStartingPosSelected] += int(math.Floor(float64(((yPos)-220)/sensitivity))) - movingRotorStartRotationInitialValue
			if rotorsRotationAmounts[movingRotorStartingPosSelected] < 0 {
				rotorsRotationAmounts[movingRotorStartingPosSelected] = 26 + rotorsRotationAmounts[movingRotorStartingPosSelected]
			}
			movingRotorStartRotationInitialValue = int(math.Floor(float64(((yPos) - 220) / sensitivity)))
		}

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

func encrypt(key string) string {
	var output string
	output = plugBoardFunc(key)
	output = sendThroughRotors(output)
	output = plugBoardFunc(output)
	return output
}

func plugBoardFunc(inputedLetter string) string {
	scrambles := plugBoard
	var output string = inputedLetter
	if scrambles[inputedLetter] != " " {
		output = scrambles[inputedLetter].(string)
	}
	return output
}

func sendThroughRotors(input string) string {
	var output string = input
	for i := 0; i < 3; i++ {
		rotorOnNbm, err := strconv.Atoi((rotorNbms[i]))
		if err == nil {
			output = rotorsgui[rotorOnNbm-1][output].(string)
		}
		//rotateRotors()
		if i == 2 {
			output = reflectorgui[output].(string)
			//rotateRotors()
			for i := 2; i > -1; i-- {
				rotorOnNbm, err := strconv.Atoi((rotorNbms[i]))
				if err == nil {
					output = rotorsgui[rotorOnNbm-1][output].(string)
				}
				//rotateRotors()
			}
		}
	}
	return output
}

var rotorsgui = [5]map[interface{}]interface{}{
	map[interface{}]interface{}{
		"a": "e",
		"b": "k",
		"c": "m",
		"d": "f",
		"e": "l",
		"f": "g",
		"g": "d",
		"h": "q",
		"i": "v",
		"j": "z",
		"k": "n",
		"l": "t",
		"m": "o",
		"n": "w",
		"o": "y",
		"p": "h",
		"q": "x",
		"r": "u",
		"s": "s",
		"t": "p",
		"u": "a",
		"v": "i",
		"w": "b",
		"x": "r",
		"y": "c",
		"z": "j",
	},
	map[interface{}]interface{}{
		"a": "a",
		"b": "j",
		"c": "d",
		"d": "k",
		"e": "s",
		"f": "i",
		"g": "r",
		"h": "u",
		"i": "x",
		"j": "b",
		"k": "l",
		"l": "h",
		"m": "w",
		"n": "t",
		"o": "m",
		"p": "c",
		"q": "q",
		"r": "g",
		"s": "z",
		"t": "n",
		"u": "p",
		"v": "y",
		"w": "f",
		"x": "v",
		"y": "o",
		"z": "e",
	},
	map[interface{}]interface{}{
		"a": "b",
		"b": "d",
		"c": "f",
		"d": "h",
		"e": "j",
		"f": "l",
		"g": "c",
		"h": "p",
		"i": "r",
		"j": "t",
		"k": "x",
		"l": "v",
		"m": "z",
		"n": "n",
		"o": "y",
		"p": "e",
		"q": "i",
		"r": "w",
		"s": "g",
		"t": "a",
		"u": "k",
		"v": "m",
		"w": "u",
		"x": "s",
		"y": "q",
		"z": "o",
	},
	map[interface{}]interface{}{
		"a": "e",
		"b": "s",
		"c": "o",
		"d": "v",
		"e": "p",
		"f": "z",
		"g": "j",
		"h": "a",
		"i": "y",
		"j": "q",
		"k": "u",
		"l": "i",
		"m": "r",
		"n": "h",
		"o": "x",
		"p": "l",
		"q": "n",
		"r": "f",
		"s": "t",
		"t": "g",
		"u": "k",
		"v": "d",
		"w": "c",
		"x": "m",
		"y": "w",
		"z": "b",
	},
	map[interface{}]interface{}{
		"a": "v",
		"b": "z",
		"c": "b",
		"d": "r",
		"e": "g",
		"f": "i",
		"g": "t",
		"h": "y",
		"i": "u",
		"j": "p",
		"k": "s",
		"l": "d",
		"m": "n",
		"n": "h",
		"o": "l",
		"p": "x",
		"q": "a",
		"r": "w",
		"s": "m",
		"t": "j",
		"u": "q",
		"v": "o",
		"w": "f",
		"x": "e",
		"y": "c",
		"z": "k",
	},
}

var reflectorgui = map[interface{}]interface{}{
	"a": "e",
	"b": "j",
	"c": "m",
	"d": "z",
	"e": "a",
	"f": "l",
	"g": "y",
	"h": "x",
	"i": "v",
	"j": "b",
	"k": "w",
	"l": "f",
	"m": "c",
	"n": "r",
	"o": "q",
	"p": "u",
	"q": "o",
	"r": "n",
	"s": "t",
	"t": "s",
	"u": "p",
	"v": "i",
	"w": "k",
	"x": "h",
	"y": "g",
	"z": "d",
}

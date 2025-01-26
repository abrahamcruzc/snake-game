package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	dirUp    = Point{0, -1}
	dirDown  = Point{0, 1}
	dirLeft  = Point{-1, 0}
	dirRight = Point{1, 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	gameSpeed    = time.Second / 6
	screenWidth  = 640
	screenHeight = 480
	gridSize     = 20
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}

	currentDir := g.direction
	if ebiten.IsKeyPressed(ebiten.KeyW) && currentDir != dirDown {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyS) && currentDir != dirUp {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyA) && currentDir != dirRight {
		g.direction = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyD) && currentDir != dirLeft {
		g.direction = dirRight
	}

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()

	head := g.snake[0]
	newHead := Point{
		x: head.x + g.direction.x,
		y: head.y + g.direction.y,
	}

	// Check wall collision
	if newHead.x < 0 || newHead.x >= screenWidth/gridSize ||
		newHead.y < 0 || newHead.y >= screenHeight/gridSize {
		g.gameOver = true
		return nil
	}

	// Check self-collision
	for _, segment := range g.snake[1:] {
		if newHead == segment {
			g.gameOver = true
			return nil
		}
	}

	if newHead == g.food {
		g.snake = append([]Point{newHead}, g.snake...)
		g.spawnFood()
	} else {
		g.snake = append([]Point{newHead}, g.snake[:len(g.snake)-1]...)
	}

	return nil
}

func (g *Game) spawnFood() {
	for {
		g.food = Point{
			x: rand.Intn(screenWidth / gridSize),
			y: rand.Intn(screenHeight / gridSize),
		}
		valid := true
		for _, segment := range g.snake {
			if g.food == segment {
				valid = false
				break
			}
		}
		if valid {
			break
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw snake
	for _, p := range g.snake {
		vector.DrawFilledRect(
			screen,
			float32(p.x*gridSize),
			float32(p.y*gridSize),
			gridSize, gridSize,
			color.White,
			true,
		)
	}

	// Draw food
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*gridSize),
		float32(g.food.y*gridSize),
		gridSize, gridSize,
		color.RGBA{255, 0, 0, 255},
		true,
	)

	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		textStr := "Game Over!"
		w, h := text.Measure(textStr, face, 48)

		op := &text.DrawOptions{}
		op.GeoM.Translate(
			float64(screenWidth/2)-float64(w)/2,
			float64(screenHeight/2)-float64(h)/2,
		)

		// Draw the text
		text.Draw(screen, textStr, face, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	rand.Seed(time.Now().UnixNano())

	g := &Game{
		snake: []Point{{
			x: screenWidth / gridSize / 2,
			y: screenHeight / gridSize / 2,
		}},
		direction: dirRight,
	}

	g.spawnFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
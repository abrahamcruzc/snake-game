package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/text/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	dirUp = Point{0, -1}
	dirDown = Point{0, 1}
	dirLeft = Point{-1, 0}
	dirRight = Point{1, 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	gameSpeed = time.Second / 6
	screenWidth = 640
	screenHeight = 480
	gridSize = 20
)

type Point struct {
	x, y int
}

type Game struct {
	snake []Point
	direction Point
	lastUpdate time.Time
	food Point
	gameOver bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}
	
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.direction = dirLeft
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.direction = dirRight
	}
	
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()
	
	g.updateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) updateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]
	newHead := Point{
		head.x + direction.x,
		head.y + direction.y,
	}
	
	if newHead == g.food {
		*snake = append(
			[]Point{newHead},
			*snake...
		)
		
		g.spawnFood()
	} else {
		*snake = append(
			[]Point{newHead}, 
			(*snake)[:len(*snake)-1]...
		)
	}
}

func (g *Game) spawnFood() {
	g.food = Point{
		rand.Intn(screenWidth / gridSize), 
		rand.Intn(screenWidth / gridSize),
	}
	
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.snake {
		vector.DrawFilledRect(
			screen,
			float32(p.x * gridSize),
			float32(p.y * gridSize),
			gridSize, gridSize,
			color.White,
			true,
		)
	}
	
	vector.DrawFilledRect(
		screen,
		float32(p.food.x * gridSize),
		float32(p.food.y * gridSize),
		gridSize, gridSize,
		color.RGBA{255, 0, 0, 255},
		true,
	)
	
	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size: 48,
		}
		
		w, h := text.Measure(
			"Game Over!",
			face,
			face.Size,
		)
		
		op := &text.DrawOptions{}
		op.GeoM.Translate(screenWidth/2-w/2, screenHeight/2-h/2)
		text.Draw(
			screen, 
			"Game Over", 
			face, 
			op,
		)
		
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	s, err := text.NewGoTextFaceSource(
		bytes.NewReader(
			fonts.MPlus1pRegular_ttf,
		),
	); if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
	
	
	g := &Game{
		[]Point{{
			screenWidth / gridSize / 2, 
			screenHeight / gridSize / 2,
		}},
		Point{1, 0},
	}
	
	g.spawnFood()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
    err = ebiten.RunGame(g); if err != nil {
        log.Fatalf("Error fatal al ejecutar el juego: %v", err)
    }
}

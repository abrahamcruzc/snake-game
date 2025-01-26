package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
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
}

func (g *Game) Update() error {
	g.updateSnake(&g.snake, g.direction)
	return nil
}

func (g *Game) updateSnake(snake *[]Point, direction Point) {
	head := (*snake)[0]
	newHead := Point{
		head.x + direction.x,
		head.y + direction.y,
	}

	*snake = append([]Point{newHead}, (*snake)[:len(*snake)-1]...)

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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{
		[]Point{
			{screenWidth / gridSize / 2, screenHeight / gridSize / 2,},
			{screenWidth / gridSize / 2 - 1, screenHeight / gridSize / 2,},
		},
		Point{1, 0},

	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Snake")
    err := ebiten.RunGame(g); if err != nil {
        log.Fatalf("Error fatal al ejecutar el juego: %v", err)
    }
}

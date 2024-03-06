package main

import (
	"github.com/JoelOtter/termloop"
)

var spaceshipImg = `  █  
█████`

var enemyImg = `
█  █
 ██ `

func main() {
	game := termloop.NewGame()
	game.Screen().SetFps(60)

	startScreen(game)

	game.Start()
}

type Title struct {
	*termloop.Text
	game *termloop.Game
}

func (t *Title) Draw(s *termloop.Screen) {
	w, h := s.Size()
	t.SetPosition(w/2-len(t.Text.Text())/2, h/2)
	t.Text.Draw(s)
}

func (t *Title) Tick(evt termloop.Event) {
	if evt.Type != termloop.EventKey {
		return
	}

	gameScreen(t.game)
}

type Bullet struct {
	*termloop.Rectangle
	game *termloop.Game
}

func (b *Bullet) Tick(evt termloop.Event) {
	x, y := b.Position()
	b.SetPosition(x, y-1)

	if y < 5 {
		b.game.Screen().Level().RemoveEntity(b)
	}
}

type Spaceship struct {
	*termloop.Entity
	game         *termloop.Game
	canShoot     bool
	shootTimeout float64
}

func (s *Spaceship) Collide(collision termloop.Physical) {
	startScreen(s.game)
}

func (s *Spaceship) Tick(evt termloop.Event) {
	s.shootTimeout -= s.game.Screen().TimeDelta()
	if s.shootTimeout < 0 {
		s.canShoot = true
		s.shootTimeout = 0.5
	}

	if evt.Type != termloop.EventKey {
		return
	}

	x, y := s.Position()
	switch evt.Key {

	case termloop.KeyArrowLeft:
		if x > 5 {
			s.SetPosition(x-1, y)
		}

	case termloop.KeyArrowRight:
		w, _ := s.game.Screen().Size()
		if x < w-10 {
			s.SetPosition(x+1, y)
		}

	case termloop.KeySpace:
		if !s.canShoot {
			return
		}

		r := termloop.NewRectangle(x+2, y-1, 1, 1, termloop.ColorRed)
		b := &Bullet{Rectangle: r, game: s.game}
		s.game.Screen().Level().AddEntity(b)
		s.canShoot = false
	}
}

type EnemyGroup struct {
	changeTick int
	cnt        int
}

type Enemy struct {
	*termloop.Entity
	game  *termloop.Game
	dir   int
	group *EnemyGroup
	tick  int
}

func (e *Enemy) Collide(collision termloop.Physical) {
	if b, ok := collision.(*Bullet); ok {
		e.game.Screen().Level().RemoveEntity(b)
		e.game.Screen().Level().RemoveEntity(e)

		e.group.cnt--
		if e.group.cnt == 0 {
			startScreen(e.game)
		}
	}
}

func (e *Enemy) Tick(evt termloop.Event) {
	x, y := e.Position()
	if e.tick%2 == 0 {
		e.SetPosition(x+e.dir, y)
	}

	e.tick++

	if e.tick == e.group.changeTick {
		e.SetPosition(x, y+2)
		e.dir *= -1
		return
	}

	w, _ := e.game.Screen().Size()
	leftMargin, rightMargin := 10, w-15
	if x <= leftMargin || x >= rightMargin {
		if e.group.changeTick < e.tick-2 {
			e.group.changeTick = e.tick + 1
		}
	}
}

func gameScreen(game *termloop.Game) {
	lvl := termloop.NewBaseLevel(termloop.Cell{
		Bg: termloop.ColorBlack,
		Fg: termloop.ColorWhite,
	})
	game.Screen().SetLevel(lvl)

	w, h := game.Screen().Size()

	scanv := termloop.CanvasFromString(spaceshipImg)
	e := termloop.NewEntityFromCanvas(w/2, h-3, scanv)
	e.Fill(&termloop.Cell{Bg: termloop.ColorBlack, Fg: termloop.ColorBlue})
	ss := &Spaceship{Entity: e, game: game, canShoot: true}
	lvl.AddEntity(ss)

	group := &EnemyGroup{cnt: 2 * 5}
	for i := range 2 {
		for j := range 5 {

			ecanv := termloop.CanvasFromString(enemyImg)
			e := termloop.NewEntityFromCanvas(10+10*j, 1+4*i, ecanv)

			enemy := &Enemy{
				Entity: e,
				game:   game,
				dir:    1,
				group:  group,
			}
			lvl.AddEntity(enemy)
		}
	}
}

func startScreen(game *termloop.Game) {
	lvl := termloop.NewBaseLevel(termloop.Cell{
		Bg: termloop.ColorBlack,
		Fg: termloop.ColorWhite,
	})
	game.Screen().SetLevel(lvl)

	txt := termloop.NewText(0, 0, "Press any key to play", termloop.ColorWhite, termloop.ColorBlack)
	title := &Title{Text: txt, game: game}
	lvl.AddEntity(title)
}

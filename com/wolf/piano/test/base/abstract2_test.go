package base

import (
	"fmt"
	"testing"
)

type IGame interface {
	Name() string
}

// 抽象类
type Game struct{}

// 将 game IGame 传了进来,便可以调用”子类”的方法来获取名字。从而间接地实现了在公共方法中调用不同”子类”的实现的抽象方法
func (g *Game) play(game IGame) {
	fmt.Printf(fmt.Sprintf("%s is awesome!", game.Name()))
}

// 子类
type Dota struct {
	Game
}

func (d *Dota) Name() string {
	return "Dota"
}

type LOL struct {
	Game
}

func (l *LOL) Name() string {
	return "LOL"
}

func TestAbstractBase1(t *testing.T) {
	dota := &Dota{}
	dota.play(dota)
	lol := &LOL{}
	// 公共方法，但是个别实现依赖参数
	lol.play(lol)
}

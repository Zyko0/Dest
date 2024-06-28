package ui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/Zyko0/Alapae/assets"
	"github.com/Zyko0/Alapae/core/building"
	"github.com/Zyko0/Alapae/core/hand"
	"github.com/Zyko0/Alapae/logic"
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
)

var emptyImg = ebiten.NewImage(3, 3)

type Stats struct {
	layout *ui.Layout

	title       *uiex.Label
	stats       *uiex.Label
	right       *ui.Grid
	left        *ui.Grid
	all         *ui.Grid
	rightCurses *ui.Grid
	leftCurses  *ui.Grid
	allCurses   *ui.Grid

	descPic   *uiex.Picture
	descTitle *uiex.Label
	descText  *uiex.Label

	hovered *building.Mod
}

func newSubItemGrid(s *Stats) *ui.Grid {
	g := ui.NewGrid(5, 4).WithOptions(
		opt.Grid.Options(
			opt.RGB(10, 10, 10),
			opt.Shape(ui.ShapeBox),
			opt.Rounding(15),
			opt.Padding(2),
		),
	)
	for y := 0; y < 4; y++ {
		for x := 0; x < 5; x++ {
			g.Add(x, y, 1, 1, uiex.NewPicture(emptyImg).
				WithOptions(
					opt.Picture.Options(
						opt.EventStyle(ui.EventOptions{
							ui.Default: opt.Border(0, color.White),
							ui.Hover: func(i ui.Item) {
								if i.Data() != nil {
									i.SetBorderWidth(2)
									s.hovered = i.Data().(*building.Mod)
								}
							},
							ui.PressHover:   opt.DoEvent(ui.Hover),
							ui.JustPress:    opt.DoEvent(ui.Hover),
							ui.ReleaseHover: opt.DoEvent(ui.Hover),
						}),
					),
					opt.Picture.Image.Options(
						opt.Image.FillContainer(true),
					),
				),
			)
		}
	}
	// TODO: prepare picture containers
	return g
}

func newItemGrid(subGrid0, subGrid1 *ui.Grid) *ui.Grid {
	g := ui.NewGrid(1, 2)
	g.Add(0, 0, 1, 1, subGrid0)
	g.Add(0, 1, 1, 1, subGrid1)
	g.WithOptions(opt.Grid.Options(
		opt.Padding(2),
	))

	return g
}

func NewStats() *Stats {
	s := &Stats{}
	s.layout = ui.NewLayout(32, 18, image.Rectangle{})
	s.layout.SetDimensions(
		logic.ScreenWidth*8/10,
		logic.ScreenHeight*8/10,
	)
	s.layout.Grid().WithOptions(opt.Grid.Options(
		opt.RGB(0, 0, 0),
		opt.Rounding(32),
	))
	s.title = uiex.NewLabel("Game over").WithOptions(
		opt.Label.Text(
			opt.Text.AlignCenterX(),
			opt.Text.AlignTop(),
			opt.Text.Color(color.White),
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(48),
		),
	)
	columns := ui.NewGrid(4, 1).WithOptions(opt.Grid.Options(
		opt.Filling(ui.ColorFillingNone),
	))

	// Item grids
	s.left = newSubItemGrid(s)
	s.all = newSubItemGrid(s)
	s.right = newSubItemGrid(s)
	s.leftCurses = newSubItemGrid(s)
	s.allCurses = newSubItemGrid(s)
	s.rightCurses = newSubItemGrid(s)
	// Parents
	leftBlock := newItemGrid(s.left, s.leftCurses)
	allBlock := newItemGrid(s.all, s.allCurses)
	rightBlock := newItemGrid(s.right, s.rightCurses)

	s.stats = uiex.NewLabel("").WithOptions(
		opt.Label.RichText(
			opt.RichText.Source(assets.FontSource),
			opt.RichText.Size(24),
			opt.RichText.LineSpacing(24),
			opt.RichText.Align(ui.AlignMin, ui.AlignMin),
			opt.RichText.PaddingLeft(20),
		),
	)

	columns.Add(0, 0, 1, 1, leftBlock)
	columns.Add(1, 0, 1, 1, allBlock)
	columns.Add(2, 0, 1, 1, rightBlock)
	columns.Add(3, 0, 1, 1, s.stats)
	// TODO: stats block

	// Description box
	s.descPic = uiex.NewPicture(emptyImg).WithOptions(
		opt.Picture.Image.Options(
			opt.Image.FillContainer(true),
		),
	)
	s.descTitle = uiex.NewLabel("").WithOptions(
		opt.Label.Text(
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(32),
			opt.Text.Align(ui.AlignMin, ui.AlignCenter),
		),
	)
	s.descText = uiex.NewLabel("").WithOptions(
		opt.Label.Text(
			opt.Text.Source(assets.FontSource),
			opt.Text.Size(24),
			opt.Text.Align(ui.AlignMin, ui.AlignCenter),
		),
	)

	s.layout.Grid().Add(
		12, 0, 8, 2, s.title,
	)
	s.layout.Grid().Add(
		1, 2, 30, 12, columns,
	)
	s.layout.Grid().Add(
		1, 15, 2, 2, s.descPic,
	)
	s.layout.Grid().Add(
		4, 15, 8, 1, s.descTitle,
	)
	s.layout.Grid().Add(
		4, 16, 26, 1, s.descText,
	)

	return s
}

func (s *Stats) setItems(g *ui.Grid, mods []*building.Mod) {
	var i int
	g.ForEach(func(item ui.Item) {
		pic := item.(*uiex.Picture)
		if i < len(mods) {
			pic.SetData(mods[i])
			pic.Image().SetImage(
				assets.ItemSheetImage.SubImage(
					mods[i].SourceRect(),
				).(*ebiten.Image),
			)
		} else {
			pic.SetData(nil)
			pic.Image().SetImage(emptyImg)
		}
		i++
	})
}

type StatsContext struct {
	Title string
	Build *building.Core
}

func appendStatName(rt *uiex.RichText, name string) {
	rt.PushColorFg(color.RGBA{255, 200, 0, 255})
	rt.PushBold()
	rt.Append("\n" + name + "\n\n")
	rt.Pop()
	rt.Pop()
}

func appendStat(rt *uiex.RichText, name, value string) {
	rt.PushColorFg(color.White)
	rt.Append(name)
	rt.Pop()
	rt.PushColorFg(color.RGBA{0, 255, 0, 255})
	rt.Append(value + "\n")
	rt.Pop()
}

func (s *Stats) Update(ctx *StatsContext) {
	s.hovered = nil
	// Update items
	s.setItems(s.left, ctx.Build.Hand(hand.Left).Bonuses)
	s.setItems(s.leftCurses, ctx.Build.Hand(hand.Left).Curses)
	s.setItems(s.right, ctx.Build.Hand(hand.Right).Bonuses)
	s.setItems(s.rightCurses, ctx.Build.Hand(hand.Right).Curses)
	s.setItems(s.all, ctx.Build.Bonuses)
	s.setItems(s.allCurses, ctx.Build.Curses)
	// Update layout
	offset := image.Pt(
		logic.ScreenWidth*1/10,
		logic.ScreenHeight*1/10,
	)
	s.layout.Update(offset, ui.GetInputState())

	// Statistics
	rt := s.stats.Text().(*uiex.RichText)
	rt.SetText("")
	rt.Reset()
	appendStatName(rt, "Global")
	appendStat(rt, "Attack speed:", fmt.Sprintf("%+.1f%%", float64(ctx.Build.AttackSpeedStacks)*12.5))
	appendStat(rt, "Luck:", fmt.Sprintf("%+.1f%%", ctx.Build.Luck*100))
	appendStat(rt, "Heal per stage:", fmt.Sprintf("%+.0f", ctx.Build.HealthPerStage))
	appendStatName(rt, "Left")
	appendStat(rt, "Damage:", fmt.Sprintf("%+.1f", ctx.Build.Hand(hand.Left).Damage))
	appendStat(rt, "Critical chance:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Left).CritChance*100))
	appendStat(rt, "Critical damage:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Left).CritDamage*100))
	appendStat(rt, "Accuracy:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Left).Accuracy*100))
	appendStat(rt, "Projectile speed:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Left).ProjectileSpeed*100))
	appendStatName(rt, "Right")
	appendStat(rt, "Damage:", fmt.Sprintf("%+.1f", ctx.Build.Hand(hand.Right).Damage))
	appendStat(rt, "Critical chance:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Right).CritChance*100))
	appendStat(rt, "Critical damage:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Right).CritDamage*100))
	appendStat(rt, "Accuracy:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Right).Accuracy*100))
	appendStat(rt, "Projectile speed:", fmt.Sprintf("%+.0f%%", ctx.Build.Hand(hand.Right).ProjectileSpeed*100))

	// Update description box if necessary
	if s.hovered != nil {
		s.descPic.Image().SetImage(
			assets.ItemSheetImage.SubImage(
				s.hovered.SourceRect(),
			).(*ebiten.Image),
		)
		s.descTitle.Text().SetText(s.hovered.Name())
		s.descText.Text().SetText(s.hovered.Description())
	} else {
		s.descPic.Image().SetImage(emptyImg)
		s.descTitle.Text().SetText("")
		s.descText.Text().SetText("")
	}
}

func (s *Stats) Draw(screen *ebiten.Image) {
	s.layout.Draw(screen)
}

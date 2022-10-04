package goeinstein

import (
	"io"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type GameBackground struct {
	Widget
}

func NewGameBackground() *GameBackground {
	return &GameBackground{}
}

func (g *GameBackground) Draw() {
	// draw background
	DrawWallpaper("rain.bmp")

	// draw title
	tile := LoadImage("title.bmp")
	screen.Draw(8, 10, tile)
	tile.Free()

	titleFont := NewFont("nova.ttf", 28)
	titleFont.DrawSurface(screen.GetSurface(), 20, 20, 255, 255, 0, true, msg("einsteinPuzzle"))

	screen.AddRegionToUpdate(0, 0, screen.GetWidth(), screen.GetHeight())
}

type ToggleHintCommand struct {
	verHints *VertHints
	horHints *HorHints
}

var _ Command = (*ToggleHintCommand)(nil)

func NewToggleHintCommand(v *VertHints, h *HorHints) *ToggleHintCommand {
	t := &ToggleHintCommand{}
	t.verHints = v
	t.horHints = h
	return t
}

func (t *ToggleHintCommand) DoAction() {
	t.verHints.ToggleExcluded()
	t.horHints.ToggleExcluded()
}

type Watch struct {
	Widget

	lastRun    uint64
	elapsed    uint64
	stoped     bool
	lastUpdate uint64
	font       *Font
}

var _ TimerHandler = (*Watch)(nil)

func (w *Watch) GetElapsed() int { return int(w.elapsed) }

func NewWatch() *Watch {
	w := &Watch{}
	w.lastRun = 0
	w.elapsed = 0
	w.lastUpdate = 0
	w.Stop()
	w.font = NewFont("luximb.ttf", 16)
	return w
}

func NewWatchStream(stream io.Reader) *Watch {
	w := &Watch{}
	w.elapsed = uint64(ReadInt(stream))
	w.lastUpdate = 0
	w.Stop()
	w.font = NewFont("luximb.ttf", 16)
	return w
}

func (w *Watch) Close() {
	w.font.Close()
}

func (w *Watch) OnTimer() {
	if w.stoped {
		return
	}

	now := sdl.GetTicks64()
	w.elapsed += now - w.lastRun
	w.lastRun = now

	seconds := w.elapsed / 1000
	if seconds != w.lastUpdate {
		w.Draw()
	}
}

func (w *Watch) Stop() {
	w.stoped = true
}

func (w *Watch) Start() {
	w.stoped = false
	w.lastRun = sdl.GetTicks64()
}

func (w *Watch) Draw() {
	time := w.elapsed / 1000
	s := SecToStr(time)

	var x int32 = 700
	var y int32 = 24
	ww, h := w.font.GetSize(s)
	rect := &sdl.Rect{x - 2, y - 2, ww + 4, h + 4}
	SDL_FillRect(screen.GetSurface(), rect, sdl.MapRGB(screen.GetSurface().Format, 0, 0, 255))
	w.font.Draw(x, y, 255, 255, 255, true, s)
	screen.AddRegionToUpdate(x-2, y-2, ww+4, h+4)

	w.lastUpdate = time
}

func (w *Watch) Save(stream io.Writer) {
	WriteInt(stream, int(w.elapsed))
}

func (w *Watch) Reset() {
	w.elapsed = 0
	w.lastUpdate = 0
	w.lastRun = sdl.GetTicks64()
}

type PauseGameCommand struct {
	gameArea   *Area
	watch      *Watch
	background AreaWidgeter
}

var _ Command = (*PauseGameCommand)(nil)

func NewPauseGameCommand(a *Area, w *Watch, bg AreaWidgeter) *PauseGameCommand {
	p := &PauseGameCommand{}
	p.gameArea = a
	p.watch = w
	p.background = bg
	return p
}

func (p *PauseGameCommand) DoAction() {
	p.watch.Stop()
	area := NewArea()
	area.AddManaged(p.background, false)
	font := NewFont("laudcn2.ttf", 16)
	area.Add(NewWindowFrame(280, 275, 240, 50, "greenpattern.bmp", 6))
	area.Add(NewLabelAligh(font, 280, 275, 240, 50, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("paused")))
	area.Add(NewAnyKeyAccelDefault())
	area.Run()
	sound.Play("click.wav")
	p.gameArea.UpdateMouse()
	p.gameArea.Draw()
	p.watch.Start()
}

type WinCommand struct {
	gameArea *Area
	watch    *Watch
	game     *Game
}

var _ Command = (*WinCommand)(nil)

func NewWinCommand(a *Area, w *Watch, g *Game) *WinCommand {
	c := &WinCommand{}
	c.gameArea = a
	c.watch = w
	c.game = g
	return c
}

func (w *WinCommand) DoAction() {
	sound.Play("applause.wav")
	w.watch.Stop()
	font := NewFont("laudcn2.ttf", 20)
	ShowMessageWindow(w.gameArea, "marble1.bmp", 500, 70, font, 255, 0, 0, msg("won"))
	w.gameArea.Draw()
	scores := NewTopScores()
	defer scores.Close()
	score := w.watch.GetElapsed() / 1000
	pos := -1
	if !w.game.IsHinted() {
		if !scores.IsFull() || (score < scores.GetMaxScore()) {
			name := EnterNameDialog(w.gameArea)
			pos = scores.Add(name, score)
		}
	}
	ShowScoresWindowHighlight(w.gameArea, scores, pos)
	w.gameArea.FinishEventLoop()
}

type OkDlgCommand struct {
	res  *bool
	area *Area
}

var _ Command = (*OkDlgCommand)(nil)

func NewOkDlgCommand(a *Area, r *bool) *OkDlgCommand {
	c := &OkDlgCommand{}
	c.res = r
	c.area = a
	return c
}

func (o *OkDlgCommand) DoAction() {
	*o.res = true
	o.area.FinishEventLoop()
}

type FailCommand struct {
	gameArea *Area
	game     *Game
}

var _ Command = (*FailCommand)(nil)

func NewFailCommand(a *Area, g *Game) *FailCommand {
	f := &FailCommand{}
	f.gameArea = a
	f.game = g
	return f
}

func (f *FailCommand) DoAction() {
	sound.Play("glasbk2.wav")
	var restart bool
	var newGame bool
	font := NewFont("laudcn2.ttf", 24)
	btnFont := NewFont("laudcn2.ttf", 14)
	area := NewArea()
	area.Add(f.gameArea)
	area.Add(NewWindowFrame(220, 240, 360, 140, "redpattern.bmp", 6))
	area.Add(NewLabelAligh(font, 250, 230, 300, 100, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("loose")))
	newGameCmd := NewOkDlgCommand(area, &newGame)
	area.Add(NewButtonText(250, 340, 90, 25, btnFont, 255, 255, 0, "redpattern.bmp", msg("startNew"), newGameCmd))
	restartCmd := NewOkDlgCommand(area, &restart)
	area.Add(NewButtonText(350, 340, 90, 25, btnFont, 255, 255, 0, "redpattern.bmp", msg("tryAgain"), restartCmd))
	exitCmd := NewExitCommand(area)
	area.Add(NewButtonText(450, 340, 90, 25, btnFont, 255, 255, 0, "redpattern.bmp", msg("exit"), exitCmd))
	area.Run()
	if restart || newGame {
		if newGame {
			f.game.NewGame()
		} else {
			f.game.Restart()
		}
		f.gameArea.Draw()
		f.gameArea.UpdateMouse()
	} else {
		f.gameArea.FinishEventLoop()
	}
}

type CheatAccel struct {
	Widget

	command Command
	typed   string
	cheat   string
}

func NewCheatAccel(s string, cmd Command) *CheatAccel {
	c := &CheatAccel{}
	c.cheat = s
	c.command = cmd
	return c
}

func (c *CheatAccel) OnKeyDown(key sdl.Keycode, ch sdl.Scancode) bool {
	if key >= sdl.K_a && key <= sdl.K_z {
		s := string('a' + rune(key) - sdl.K_a)
		c.typed += s
		if len(c.typed) == len(c.cheat) {
			if c.command != nil && c.typed == c.cheat {
				c.command.DoAction()
			}
		} else {
			pos := len(c.typed) - 1
			if c.typed[pos] == c.cheat[pos] {
				return false
			}
		}
	}
	if len(c.typed) > 0 {
		c.typed = ""
	}
	return false
}

type CheatCommand struct {
	gameArea *Area
}

var _ Command = (*CheatCommand)(nil)

func NewCheatCommand(a *Area) *CheatCommand {
	c := &CheatCommand{}
	c.gameArea = a
	return c
}

func (c *CheatCommand) DoAction() {
	font := NewFont("nova.ttf", 30)
	ShowMessageWindow(c.gameArea, "darkpattern.bmp", 500, 100, font, 255, 255, 255, msg("iddqd"))
	c.gameArea.Draw()
}

type SaveGameCommand struct {
	gameArea   *Area
	watch      *Watch
	background AreaWidgeter
	game       *Game
}

var _ Command = (*SaveGameCommand)(nil)

func NewSaveGameCommand(a *Area, w *Watch, bg AreaWidgeter, g *Game) *SaveGameCommand {
	s := &SaveGameCommand{}
	s.gameArea = a
	s.watch = w
	s.background = bg
	s.game = g
	return s
}

func (s *SaveGameCommand) DoAction() {
	s.watch.Stop()

	area := NewArea()
	area.AddManaged(s.background, false)
	SaveGame(area, s.game)

	s.gameArea.UpdateMouse()
	s.gameArea.Draw()
	s.watch.Start()
}

type GameOptionsCommand struct {
	gameArea *Area
}

var _ Command = (*GameOptionsCommand)(nil)

func NewGameOptionsCommand(a *Area) *GameOptionsCommand {
	g := &GameOptionsCommand{}
	g.gameArea = a
	return g
}

func (g *GameOptionsCommand) DoAction() {
	ShowOptionsWindow(g.gameArea)
	g.gameArea.UpdateMouse()
	g.gameArea.Draw()
}

type HelpCommand struct {
	gameArea   *Area
	watch      *Watch
	background AreaWidgeter
}

var _ Command = (*HelpCommand)(nil)

func NewHelpCommand(a *Area, w *Watch, b AreaWidgeter) *HelpCommand {
	h := &HelpCommand{}
	h.gameArea = a
	h.watch = w
	h.background = b
	return h
}

func (h *HelpCommand) DoAction() {
	h.watch.Stop()
	area := NewArea()
	area.AddManaged(h.background, false)
	area.Draw()
	ShowDescription(area)
	h.gameArea.UpdateMouse()
	h.gameArea.Draw()
	h.watch.Start()
}

type hintsExcluder struct {
	verHints *VertHints
	horHints *HorHints
}

var _ RuleExcluder = (*hintsExcluder)(nil)

func NewHintsExcluder(vh *VertHints, hh *HorHints) RuleExcluder {
	return &hintsExcluder{vh, hh}
}

func (h *hintsExcluder) ExcludeRule(r Ruler) {
	h.verHints.ExcludeRule(r)
	h.horHints.ExcludeRule(r)
}

type ruleHinter struct {
	rules        *Rules
	ruleExcluder RuleExcluder
}

var _ Hinter = (*ruleHinter)(nil)

func NewRuleHinter(rs *Rules, re RuleExcluder) Hinter {
	return &ruleHinter{rs, re}
}

func (r *ruleHinter) AutoHint(pos *Possibilities) {
	r.rules.ApplyHints(pos, r.ruleExcluder)
}

var game *Game

type Game struct {
	solvedPuzzle      SolvedPuzzle
	rules             Rules
	possibilities     *Possibilities
	verHints          *VertHints
	horHints          *HorHints
	iconSet           *IconSet
	puzzle            *Puzzle
	watch             *Watch
	hinted            bool
	savedSolvedPuzzle SolvedPuzzle
	savedRules        Rules
}

func (g *Game) GetSolvedPuzzle() SolvedPuzzle    { return g.solvedPuzzle }
func (g *Game) GetRules() Rules                  { return g.rules }
func (g *Game) GetPossibilities() *Possibilities { return g.possibilities }
func (g *Game) GetVerHints() *VertHints          { return g.verHints }
func (g *Game) GetHorHints() *HorHints           { return g.horHints }
func (g *Game) IsHinted() bool                   { return g.hinted }
func (g *Game) SetHinted()                       { g.hinted = true }

func NewGame() *Game {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	return NewGameRand(rand)
}

func NewGameRand(rand *rand.Rand) *Game {
	g := &Game{
		iconSet: NewIconSet(),
	}
	g.GenPuzzle(rand)

	g.verHints = NewVertHints(g.iconSet, &g.rules)
	g.horHints = NewHorHints(g.iconSet, &g.rules)
	excluder := NewHintsExcluder(g.verHints, g.horHints)

	g.possibilities = NewPossibilities()
	OpenInitial(g.possibilities, &g.rules, excluder)

	hinter := NewRuleHinter(&g.rules, excluder)
	g.puzzle = NewPuzzle(g.iconSet, &g.solvedPuzzle, g.possibilities, hinter)
	g.watch = NewWatch()
	return g
}

func NewGameStream(stream io.Reader) *Game {
	g := &Game{
		iconSet: NewIconSet(),
	}
	g.PleaseWait()

	LoadPuzzle(&g.solvedPuzzle, stream)
	LoadRules(&g.rules, stream)
	g.savedSolvedPuzzle = g.solvedPuzzle
	g.savedRules = g.rules[:]
	g.possibilities = NewPossibilitiesStream(stream)
	g.verHints = NewVertHintsStream(g.iconSet, &g.rules, stream)
	g.horHints = NewHorHintsStream(g.iconSet, &g.rules, stream)
	excluder := NewHintsExcluder(g.verHints, g.horHints)
	hinter := NewRuleHinter(&g.rules, excluder)
	g.puzzle = NewPuzzle(g.iconSet, &g.solvedPuzzle, g.possibilities, hinter)
	g.watch = NewWatchStream(stream)
	g.hinted = true
	return g
}

func (g *Game) Close() {
	g.watch.Close()
	// g.possibilities.Close()
	g.verHints.Close()
	g.horHints.Close()
	g.puzzle.Close()
	g.DeleteRules()
}

func (g *Game) Save(stream io.Writer) {
	SavePuzzle(&g.solvedPuzzle, stream)
	SaveRules(&g.rules, stream)
	g.possibilities.Save(stream)
	g.verHints.Save(stream)
	g.horHints.Save(stream)
	g.watch.Save(stream)
}

func (g *Game) DeleteRules() {
	g.rules = nil
}

func (g *Game) PleaseWait() {
	DrawWallpaper("rain.bmp")
	window := NewWindowFrame(230, 260, 340, 80, "greenpattern.bmp", 6)
	window.Draw()
	font := NewFont("laudcn2.ttf", 16)
	label := NewLabelAligh(font, 280, 275, 240, 50, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, msg("loading"))
	label.Draw()
	screen.AddRegionToUpdate(0, 0, screen.GetWidth(), screen.GetHeight())
	screen.Flush()
}

func (g *Game) GenPuzzle(rand *rand.Rand) {
	g.PleaseWait()

	var horRules, verRules int
	for {
		if len(g.rules) > 0 {
			g.DeleteRules()
		}
		GenPuzzle(&g.solvedPuzzle, &g.rules, rand)
		GetHintsQty(&g.rules, &verRules, &horRules)
		if horRules <= 24 && verRules <= 15 {
			break
		}
	}

	g.savedSolvedPuzzle = g.solvedPuzzle
	g.savedRules = g.rules[:]

	g.hinted = options.AutoHints.value
}

func (g *Game) ResetVisuals() {
	g.possibilities.Reset()
	g.puzzle.Reset()
	g.verHints.Reset(&g.rules)
	g.horHints.Reset(&g.rules)
	OpenInitial(g.possibilities, &g.rules, NewHintsExcluder(g.verHints, g.horHints))
	g.watch.Reset()
}

func (g *Game) NewGame() {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	g.NewGameRand(rand)
}

func (g *Game) NewGameRand(rand *rand.Rand) {
	g.GenPuzzle(rand)
	g.ResetVisuals()
}

func (g *Game) Restart() {
	g.solvedPuzzle = g.savedSolvedPuzzle
	g.rules = g.savedRules[:]

	g.ResetVisuals()
	g.hinted = true
}

func (g *Game) Run() {
	area := NewArea()
	btnFont := NewFont("laudcn2.ttf", 14)

	area.SetTimer(300, g.watch)

	background := NewGameBackground()
	area.Add(background)
	cheatCmd := NewCheatCommand(area)
	area.Add(NewCheatAccel("iddqd", cheatCmd))
	winCmd := NewWinCommand(area, g.watch, g)
	area.Add(NewKeyAccel(sdl.K_F8, winCmd))
	failCmd := NewFailCommand(area, g)
	g.puzzle.SetCommand(winCmd, failCmd)
	area.AddManaged(g.puzzle, false)
	area.AddManaged(g.verHints, false)
	area.AddManaged(g.horHints, false)

	BUTTON := func(x, y int32, text string, cmd Command) {
		area.Add(NewButtonTextBevel(x, y, 94, 30, btnFont, 255, 255, 0, "btn.bmp", msg(text), false, cmd))
	}

	pauseGameCmd := NewPauseGameCommand(area, g.watch, background)
	BUTTON(12, 400, "pause", pauseGameCmd)
	toggleHintsCmd := NewToggleHintCommand(g.verHints, g.horHints)
	BUTTON(119, 400, "switch", toggleHintsCmd)
	saveCmd := NewSaveGameCommand(area, g.watch, background, g)
	BUTTON(12, 440, "save", saveCmd)
	optionsCmd := NewGameOptionsCommand(area)
	BUTTON(119, 440, "options", optionsCmd)
	exitGameCmd := NewExitCommand(area)
	BUTTON(226, 400, "exit", exitGameCmd)
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitGameCmd))
	helpCmd := NewHelpCommand(area, g.watch, background)
	BUTTON(226, 440, "help", helpCmd)
	area.AddManaged(g.watch, false)

	g.watch.Start()
	area.Run()
}

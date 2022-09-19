package goeinstein

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
)

//nolint:golint,nosnakecase,stylecheck
const MAX_SLOTS = 10

type SavedGame struct {
	fileName string
	exists   bool
	name     string
}

func NewSavedGameFile(s string) *SavedGame {
	sg := &SavedGame{
		fileName: s,
	}
	sg.exists = false

	bs, err := os.ReadFile(sg.fileName)
	if err != nil {
		panic(fmt.Errorf("read saved file (filename: %q): %w", sg.fileName, err))
	}
	sg.name = ReadString(bytes.NewReader(bs))
	sg.exists = true
	return sg
}

func NewSavedGame(s *SavedGame) *SavedGame {
	sg := &SavedGame{
		fileName: s.fileName,
		name:     s.name,
	}
	sg.exists = s.exists
	return sg
}

func (s *SavedGame) GetFileName() string { return s.fileName }
func (s *SavedGame) IsExists() bool      { return s.exists }

func (s *SavedGame) GetName() string {
	if s.exists {
		return s.name
	}
	return msg("empty")
}

type OkCommand struct {
	area *Area
	ok   *bool
}

var _ Command = (*OkCommand)(nil)

func NewOkCommand(a *Area, o *bool) *OkCommand {
	c := &OkCommand{
		area: a,
	}
	c.ok = o
	return c
}

func (o *OkCommand) DoAction() {
	*o.ok = true
	o.area.FinishEventLoop()
}

type SaveCommand struct {
	savedGame   *SavedGame
	parentArea  *Area
	saved       *bool
	font        *Font
	defaultName string
	game        *Game
}

var _ Command = (*SaveCommand)(nil)

func NewSaveCommand(sg *SavedGame, f *Font, area *Area, s *bool, dflt string, g *Game) *SaveCommand {
	c := &SaveCommand{
		savedGame:   sg,
		defaultName: dflt,
	}
	c.parentArea = area
	c.saved = s
	c.font = f
	c.game = g
	return c
}

func (s *SaveCommand) DoAction() {
	area := NewArea()
	area.AddManaged(s.parentArea, false)
	area.Add(NewWindow(170, 280, 460, 100, "blue.bmp"))
	var name string
	if s.savedGame.IsExists() {
		name = s.savedGame.GetName()
	} else {
		name = s.defaultName
	}
	area.Add(NewLabel(s.font, 180, 300, 255, 255, 0, msg("enterGame")))
	area.Add(NewInputField(340, 300, 280, 26, "blue.bmp", &name, 20, 255, 255, 0, s.font))
	exitCmd := NewExitCommand(area)
	okCmd := NewOkCommand(area, s.saved)
	area.Add(NewButtonText(310, 340, 80, 25, s.font, 255, 255, 0, "blue.bmp", msg("ok"), okCmd))
	area.Add(NewButtonText(400, 340, 80, 25, s.font, 255, 255, 0, "blue.bmp", msg("cancel"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Add(NewKeyAccel(sdl.K_RETURN, okCmd))
	area.Run()

	if *s.saved {
		*s.saved = false
		stream, err := os.OpenFile(s.savedGame.GetFileName(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			ShowMessageWindow(area, "redpattern.bmp", 300, 80, s.font, 255, 255, 255, msg("saveError"))
			panic(fmt.Errorf("open file to save game (filename: %q): %w", s.savedGame.GetFileName(), err))
		}
		WriteString(stream, name)
		s.game.Save(stream)
		err = stream.Close()
		if err != nil {
			ShowMessageWindow(area, "redpattern.bmp", 300, 80, s.font, 255, 255, 255, msg("saveError"))
			panic(fmt.Errorf("close file to save game (filename: %q): %w", s.savedGame.GetFileName(), err))
		}
		*s.saved = true
		s.parentArea.FinishEventLoop()
	} else {
		s.parentArea.UpdateMouse()
		s.parentArea.Draw()
	}
}

func GetSavesPath() string {
	path := "./einstein/save"
	EnsureDirExists(path)
	return path
}

func ShowListWindow(list []*SavedGame, commands []Command, title string, area *Area, font *Font) {
	titleFont := NewFont("nova.ttf", 26)

	area.Add(NewWindow(250, 90, 300, 420, "blue.bmp"))
	area.Add(NewLabelAligh(titleFont, 250, 95, 300, 40, ALIGN_CENTER, ALIGN_MIDDLE, 255, 255, 0, title))
	exitCmd := NewExitCommand(area)
	area.Add(NewButtonText(360, 470, 80, 25, font, 255, 255, 0, "blue.bmp", msg("close"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))

	pos := int32(150)
	no := 0
	for _, game := range list {
		area.Add(NewButtonText(260, pos, 280, 25, font, 255, 255, 255, "blue.bmp", game.GetName(), commands[no]))
		no++
		pos += 30
	}

	area.Run()
}

func SaveGame(parentArea *Area, game *Game) bool {
	path := GetSavesPath()

	area := NewArea()
	area.AddManaged(parentArea, false)
	font := NewFont("laudcn2.ttf", 14)
	saved := false

	var list []*SavedGame
	commands := make([]Command, MAX_SLOTS)
	for i := 0; i < MAX_SLOTS; i++ {
		sg := NewSavedGameFile(path + "/" + ToString(i) + ".sav")
		list = append(list, sg)
		commands[i] = NewSaveCommand(sg, font, area, &saved, "game "+ToString(i+1), game)
	}

	ShowListWindow(list, commands, msg("saveGame"), area, font)

	return saved
}

type LoadCommand struct {
	savedGame  *SavedGame
	parentArea *Area
	// saved       *bool
	font *Font
	// defaultName string
	game *Game
}

var _ Command = (*LoadCommand)(nil)

func NewLoadCommand(sg *SavedGame, f *Font, area *Area, g *Game) *LoadCommand {
	l := &LoadCommand{
		savedGame: sg,
	}
	l.parentArea = area
	l.font = f
	l.game = g
	return l
}

func (l *LoadCommand) DoAction() {
	bs, err := os.ReadFile(l.savedGame.GetFileName())
	if err != nil {
		panic(fmt.Errorf("read all file (filename: %q): %w", l.savedGame.GetFileName(), err))
	}
	stream := bytes.NewReader(bs)
	ReadString(stream)
	g := NewGameStream(stream)
	l.game = g

	l.parentArea.FinishEventLoop()
}

func LoadGame(parentArea *Area) *Game {
	path := GetSavesPath()

	area := NewArea()
	area.AddManaged(parentArea, false)
	font := NewFont("laudcn2.ttf", 14)

	var newGame *Game

	var list []*SavedGame
	var commands [MAX_SLOTS]Command
	for i := 0; i < MAX_SLOTS; i++ {
		sg := NewSavedGameFile(filepath.Join(path, ToString(i)+".sav"))
		list = append(list, sg)
		if sg.IsExists() {
			commands[i] = NewLoadCommand(list[i], font, area, newGame)
		} else {
			commands[i] = nil
		}
	}

	ShowListWindow(list, commands[:], msg("loadGame"), area, font)

	return newGame
}

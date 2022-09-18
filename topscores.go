package goeinstein

import "github.com/veandco/go-sdl2/sdl"

//nolint:golint,stylecheck
const MAX_SCORES = 10

type TopScoreEntry struct {
	name  string
	score int
}

type TopScores struct {
	scores  []*TopScoreEntry
	modifed bool
}

func (t *TopScores) IsFull() bool { return len(t.scores) >= MAX_SCORES }

func NewTopScores() *TopScores {
	t := &TopScores{}
	storage := GetStorage()

	for i := 0; i < MAX_SCORES; i++ {
		score := storage.GetInt("top_score_"+ToString(i), -1)
		if score < 0 {
			break
		}
		name := storage.GetString("top_name_"+ToString(i), "")
		t.Add(name, score)
	}

	t.modifed = false
	return t
}

func (t *TopScores) Close() {
	t.Save()
}

func (t *TopScores) Add(name string, score int) int {
	if score >= t.GetMaxScore() || len(t.scores) < 1 {
		if !t.IsFull() {
			e := &TopScoreEntry{name, score}
			t.scores = append(t.scores, e)
			t.modifed = true
			return len(t.scores) - 1
		}
		return -1
	}

	var pos int
	for i, e := range t.scores {
		if e.score > score {
			ne := &TopScoreEntry{name, score}
			t.scores = append(t.scores[:i], append([]*TopScoreEntry{ne}, t.scores[i:]...)...)
			t.modifed = true
			break
		}
		pos++
	}

	for len(t.scores) > MAX_SCORES {
		t.modifed = true
		t.scores = t.scores[:len(t.scores)-1]
	}

	if t.modifed {
		return pos
	}
	return -1
}

func (t *TopScores) Save() {
	if !t.modifed {
		return
	}

	storage := GetStorage()
	var no int

	for _, e := range t.scores {
		storage.SetString("top_name_"+ToString(no), e.name)
		storage.SetInt("top_score_"+ToString(no), e.score)
		no++
	}

	storage.Flush()
	t.modifed = false
}

func (t *TopScores) GetScores() []*TopScoreEntry { return t.scores }

func (t *TopScores) GetMaxScore() int {
	if len(t.scores) < 1 {
		return -1
	}
	return t.scores[len(t.scores)-1].score
}

type ScoresWindow struct {
	Window
}

func NewScoresWindow(x, y int32, scores *TopScores, highlight int) *ScoresWindow {
	sw := &ScoresWindow{}
	sw.Window = *NewWindow(x, y, 320, 350, "blue.bmp")

	titleFont := NewFont("nova.ttf", 26)
	entryFont := NewFont("laudcn2.ttf", 14)
	timeFont := NewFont("luximb.ttf", 14)

	txt := msg("topScores")
	w := titleFont.GetWidth(txt)
	titleFont.DrawSurface(sw.background, (320-w)/2, 15, 255, 255, 0, true, txt)

	list := scores.GetScores()
	no := 1
	pos := int32(70)
	for _, e := range list {
		s := ToString(no) + "."
		w := entryFont.GetWidth(s)
		var c uint8
		if (no - 1) == highlight {
			c = 0
		} else {
			c = 255
		}
		entryFont.DrawSurface(sw.background, 30-w, pos, 255, 255, c, true, s)
		rect := &sdl.Rect{40, pos - 20, 180, 40}
		sw.background.SetClipRect(rect)
		entryFont.DrawSurface(sw.background, 40, pos, 255, 255, c, true, e.name)
		sw.background.SetClipRect(nil)
		s = SecToStr(uint64(e.score))
		w = timeFont.GetWidth(s)
		timeFont.DrawSurface(sw.background, 305-w, pos, 255, 255, c, true, s)
		pos += 20
		no++
	}
	return sw
}

func ShowScoresWindow(parentArea *Area, scores *TopScores) {
	ShowScoresWindowHighlight(parentArea, scores, -1)
}

func ShowScoresWindowHighlight(parentArea *Area, scores *TopScores, highlight int) {
	area := NewArea()

	font := NewFont("laudcn2.ttf", 16)
	area.Add(parentArea)
	area.Add(NewScoresWindow(240, 125, scores, highlight))
	exitCmd := NewExitCommand(area)
	area.Add(NewButtonText(348, 430, 90, 25, font, 255, 255, 0, "blue.bmp", msg("ok"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Run()
}

func EnterNameDialog(parentArea *Area) string {
	area := NewArea()

	font := NewFont("laudcn2.ttf", 16)
	area.Add(parentArea)
	area.Add(NewWindow(170, 280, 460, 100, "blue.bmp"))
	storage := GetStorage()
	name := storage.GetString("lastName", msg("anonymous"))
	area.Add(NewLabel(font, 180, 300, 255, 255, 0, msg("enterName")))
	area.Add(NewInputField(350, 300, 270, 26, "blue.bmp", &name, 20, 255, 255, 0, font))
	exitCmd := NewExitCommand(area)
	area.Add(NewButtonText(348, 340, 90, 25, font, 255, 255, 0, "blue.bmp", msg("ok"), exitCmd))
	area.Add(NewKeyAccel(sdl.K_ESCAPE, exitCmd))
	area.Add(NewKeyAccel(sdl.K_RETURN, exitCmd))
	area.Run()
	storage.SetString("lastName", name)
	return name
}

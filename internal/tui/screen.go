package tui

type Screen int

const (
	ScreenUnlock Screen = iota
	ScreenList
	ScreenDetail
	ScreenEdit
	ScreenPasswordChange
	ScreenHelp
)

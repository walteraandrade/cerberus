package tui

import (
	"fmt"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/walteraandrade/cerberus/internal/clipboard"
	"github.com/walteraandrade/cerberus/internal/config"
	"github.com/walteraandrade/cerberus/internal/crypto"
	"github.com/walteraandrade/cerberus/internal/export"
	"github.com/walteraandrade/cerberus/internal/storage"
	"github.com/walteraandrade/cerberus/internal/vault"
)

type App struct {
	cfg     *config.Config
	vault   *vault.Vault
	screen  Screen
	prevScr Screen
	unlock  UnlockModel
	list    ListModel
	detail  DetailModel
	edit    EditModel
	pwChange PasswordChangeModel
	help    HelpModel
	confirm *ConfirmModel
	status  string
	width   int
	height  int

	password  []byte
	vaultPath string
	idleTimer time.Time
}

func NewApp(cfg *config.Config) App {
	vaultPath := cfg.VaultPath()
	exists := storage.Exists(vaultPath)

	return App{
		cfg:       cfg,
		screen:    ScreenUnlock,
		unlock:    NewUnlockModel(exists),
		vaultPath: vaultPath,
		idleTimer: time.Now(),
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.unlock.Init(),
		a.tickIdle(),
	)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		a.idleTimer = time.Now()

	case statusMsg:
		a.status = string(msg)
		return a, nil

	case clearStatusMsg:
		a.status = ""
		return a, nil

	case idleTickMsg:
		return a.checkIdle()
	}

	if a.confirm != nil {
		return a.updateConfirm(msg)
	}

	switch a.screen {
	case ScreenUnlock:
		return a.updateUnlock(msg)
	case ScreenList:
		return a.updateList(msg)
	case ScreenDetail:
		return a.updateDetail(msg)
	case ScreenEdit:
		return a.updateEdit(msg)
	case ScreenPasswordChange:
		return a.updatePasswordChange(msg)
	case ScreenHelp:
		return a.updateHelp(msg)
	}

	return a, nil
}

func (a App) updateUnlock(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case UnlockMsg:
		return a.handleUnlock(msg)
	default:
		var cmd tea.Cmd
		a.unlock, cmd = a.unlock.Update(msg)
		return a, cmd
	}
}

func (a App) handleUnlock(msg UnlockMsg) (tea.Model, tea.Cmd) {
	pw := []byte(msg.Password)
	params := a.kdfParams()

	if msg.Create {
		v := vault.New()
		data, err := vault.Marshal(v)
		if err != nil {
			a.unlock, _ = a.unlock.Update(UnlockErrMsg{Err: err.Error()})
			return a, nil
		}
		if err := storage.CreateVault(a.vaultPath, pw, data, params); err != nil {
			a.unlock, _ = a.unlock.Update(UnlockErrMsg{Err: err.Error()})
			return a, nil
		}
		a.vault = v
		a.password = pw
	} else {
		plaintext, err := storage.OpenVault(a.vaultPath, pw, params)
		if err != nil {
			a.unlock, _ = a.unlock.Update(UnlockErrMsg{Err: "wrong password"})
			return a, nil
		}
		v, err := vault.Unmarshal(plaintext)
		if err != nil {
			a.unlock, _ = a.unlock.Update(UnlockErrMsg{Err: "vault corrupted"})
			return a, nil
		}
		a.vault = v
		a.password = pw
	}

	a.screen = ScreenList
	a.idleTimer = time.Now()
	a.list = NewListModel(a.vault.Entries)
	return a, nil
}

func (a App) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SelectEntryMsg:
		a.screen = ScreenDetail
		a.detail = NewDetailModel(msg.Entry)
		return a, nil
	case AddEntryMsg:
		a.screen = ScreenEdit
		a.edit = NewEditModel(nil)
		return a, a.edit.Init()
	case DeleteEntryMsg:
		cm := NewConfirmModel(
			fmt.Sprintf("Delete '%s'?", msg.Entry.Title),
			doDeleteMsg{ID: msg.Entry.ID},
		)
		a.confirm = &cm
		return a, nil
	case CopyPasswordMsg:
		return a, a.copyPassword(msg.Password)
	case exportMsg:
		return a.handleExport()
	case passwordChangeMsg:
		a.screen = ScreenPasswordChange
		a.pwChange = NewPasswordChangeModel()
		return a, a.pwChange.Init()
	case helpMsg:
		a.prevScr = a.screen
		a.screen = ScreenHelp
		a.help = NewHelpModel(a.prevScr)
		return a, nil
	default:
		var cmd tea.Cmd
		a.list, cmd = a.list.Update(msg)
		return a, cmd
	}
}

func (a App) updateDetail(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BackMsg:
		a.screen = ScreenList
		a.list.SetEntries(a.vault.Entries)
		return a, nil
	case EditEntryMsg:
		a.screen = ScreenEdit
		e := msg.Entry
		a.edit = NewEditModel(&e)
		return a, a.edit.Init()
	case DeleteEntryMsg:
		cm := NewConfirmModel(
			fmt.Sprintf("Delete '%s'?", msg.Entry.Title),
			doDeleteMsg{ID: msg.Entry.ID},
		)
		a.confirm = &cm
		return a, nil
	case CopyPasswordMsg:
		return a, a.copyPassword(msg.Password)
	default:
		var cmd tea.Cmd
		a.detail, cmd = a.detail.Update(msg)
		return a, cmd
	}
}

func (a App) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BackMsg:
		a.screen = ScreenList
		a.list.SetEntries(a.vault.Entries)
		return a, nil
	case SaveEntryMsg:
		return a.handleSave(msg.Entry)
	default:
		var cmd tea.Cmd
		a.edit, cmd = a.edit.Update(msg)
		return a, cmd
	}
}

func (a App) updatePasswordChange(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case BackMsg:
		a.screen = ScreenList
		return a, nil
	case PasswordChangedMsg:
		return a.handlePasswordChange(msg)
	default:
		var cmd tea.Cmd
		a.pwChange, cmd = a.pwChange.Update(msg)
		return a, cmd
	}
}

func (a App) handlePasswordChange(msg PasswordChangedMsg) (tea.Model, tea.Cmd) {
	params := a.kdfParams()

	// Verify old password by trying to open vault
	_, err := storage.OpenVault(a.vaultPath, []byte(msg.OldPassword), params)
	if err != nil {
		a.pwChange.err = "current password is wrong"
		return a, nil
	}

	// Re-save with new password
	data, err := vault.Marshal(a.vault)
	if err != nil {
		a.pwChange.err = err.Error()
		return a, nil
	}

	if err := storage.CreateVault(a.vaultPath, []byte(msg.NewPassword), data, params); err != nil {
		a.pwChange.err = err.Error()
		return a, nil
	}

	crypto.ZeroBytes(a.password)
	a.password = []byte(msg.NewPassword)
	a.screen = ScreenList
	return a, setStatus("Password changed")
}

func (a App) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case BackMsg:
		a.screen = a.prevScr
		return a, nil
	default:
		var cmd tea.Cmd
		a.help, cmd = a.help.Update(msg)
		return a, cmd
	}
}

func (a App) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConfirmYesMsg:
		a.confirm = nil
		if del, ok := msg.Inner.(doDeleteMsg); ok {
			return a.handleDelete(del.ID)
		}
		return a, nil
	case ConfirmNoMsg:
		a.confirm = nil
		return a, nil
	default:
		cm, cmd := a.confirm.Update(msg)
		a.confirm = &cm
		return a, cmd
	}
}

func (a App) handleSave(entry vault.Entry) (tea.Model, tea.Cmd) {
	existing := a.vault.FindEntry(entry.ID)
	if existing != nil {
		entry.CreatedAt = existing.CreatedAt
		a.vault.RemoveEntry(entry.ID)
	}
	a.vault.AddEntry(entry)

	if err := a.saveVault(); err != nil {
		a.status = fmt.Sprintf("save failed: %v", err)
		return a, nil
	}

	a.screen = ScreenList
	a.list.SetEntries(a.vault.Entries)
	return a, setStatus("Saved")
}

func (a App) handleDelete(id string) (tea.Model, tea.Cmd) {
	a.vault.RemoveEntry(id)
	if err := a.saveVault(); err != nil {
		a.status = fmt.Sprintf("delete failed: %v", err)
		return a, nil
	}

	a.screen = ScreenList
	a.list.SetEntries(a.vault.Entries)
	return a, setStatus("Deleted")
}

func (a App) handleExport() (tea.Model, tea.Cmd) {
	path := filepath.Join(a.cfg.DataDir, fmt.Sprintf("cerberus-export-%s.json",
		time.Now().Format("20060102-150405")))

	if err := export.ToFile(path, export.JSONExporter{}, a.vault); err != nil {
		return a, setStatus(fmt.Sprintf("export failed: %v", err))
	}
	return a, setStatus(fmt.Sprintf("Exported to %s", path))
}

func (a App) saveVault() error {
	data, err := vault.Marshal(a.vault)
	if err != nil {
		return err
	}
	return storage.CreateVault(a.vaultPath, a.password, data, a.kdfParams())
}

func (a App) kdfParams() crypto.KDFParams {
	return crypto.KDFParams{
		Memory:      a.cfg.Argon2.Memory,
		Iterations:  a.cfg.Argon2.Iterations,
		Parallelism: a.cfg.Argon2.Parallelism,
		SaltLen:     a.cfg.Argon2.SaltLen,
		KeyLen:      a.cfg.Argon2.KeyLen,
	}
}

func (a App) copyPassword(pw string) tea.Cmd {
	return func() tea.Msg {
		timeout := time.Duration(a.cfg.ClipboardTimeout) * time.Second
		_, err := clipboard.CopyWithAutoClear(pw, timeout)
		if err != nil {
			return statusMsg(fmt.Sprintf("copy failed: %v", err))
		}
		return statusMsg(fmt.Sprintf("Copied! Clearing in %ds", a.cfg.ClipboardTimeout))
	}
}

// Idle lock
type idleTickMsg struct{}

func (a App) tickIdle() tea.Cmd {
	return tea.Tick(30*time.Second, func(time.Time) tea.Msg { return idleTickMsg{} })
}

func (a App) checkIdle() (tea.Model, tea.Cmd) {
	if a.screen == ScreenUnlock || a.cfg.LockTimeout <= 0 {
		return a, a.tickIdle()
	}

	elapsed := time.Since(a.idleTimer)
	if elapsed >= time.Duration(a.cfg.LockTimeout)*time.Second {
		return a.lockVault()
	}
	return a, a.tickIdle()
}

func (a App) lockVault() (tea.Model, tea.Cmd) {
	crypto.ZeroBytes(a.password)
	a.password = nil
	a.vault = nil
	a.screen = ScreenUnlock
	a.unlock = NewUnlockModel(storage.Exists(a.vaultPath))
	return a, tea.Batch(a.unlock.Init(), setStatus("Locked due to inactivity"))
}

func (a App) View() string {
	var view string

	switch a.screen {
	case ScreenUnlock:
		view = a.unlock.View()
	case ScreenList:
		view = a.list.View()
	case ScreenDetail:
		view = a.detail.View()
	case ScreenEdit:
		view = a.edit.View()
	case ScreenPasswordChange:
		view = a.pwChange.View()
	case ScreenHelp:
		view = a.help.View()
	}

	if a.confirm != nil {
		view += "\n" + a.confirm.View()
	}

	if a.status != "" {
		view += "\n\n" + a.status
	}

	return view
}

type doDeleteMsg struct{ ID string }
type exportMsg struct{}
type passwordChangeMsg struct{}
type helpMsg struct{}
type statusMsg string
type clearStatusMsg struct{}

func setStatus(s string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return statusMsg(s) },
		tea.Tick(3*time.Second, func(time.Time) tea.Msg { return clearStatusMsg{} }),
	)
}

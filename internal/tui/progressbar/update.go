package progressbar

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		w := msg.Width - padding*2 - 4
		if w >= maxWidth {
			w = maxWidth
		}
		m.width = w
		return m, nil

	case progressErrMsg:
		m.err = msg.err
		return m, tea.Quit

	case ProgressMsg:
		var cmds []tea.Cmd

		fileName, ratio := msg.FileName, msg.Ratio
		if pb, ok := m.Pbs[fileName]; ok {
			cmds = append(cmds, pb.Progress.SetPercent(ratio))
		}

		return m, tea.Batch(cmds...)

	case progress.FrameMsg:
		var cmds []tea.Cmd

		for _, pb := range m.Pbs {
			progressModel, cmd := pb.Progress.Update(msg)
			pb.Progress = progressModel.(progress.Model)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)

	default:
		return m, nil
	}
}

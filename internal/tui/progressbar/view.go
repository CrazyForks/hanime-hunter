package progressbar

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/acgtools/hanime-hunter/pkg/util"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	pbStyle     = lipgloss.NewStyle().MaxHeight(1).Render
	statusStyle = lipgloss.NewStyle()
)

const (
	mib = float64(1 << 20)
)

const (
	DownloadingStatus = "Downloading"
	DownloadingColor  = "#FF6600"

	MergingStatus = "Merging"
	MergingColor  = "#FFCC66"

	CompleteStatus = "Complete"
	CompleteColor  = "#00FF00"

	ErrStatus = "Error"
	ErrColor  = "#CC0000"

	RetryStatus = "Retry"
	RetryColor  = "#CC66CC"
)

func (m *Model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}

	if m.width == 0 {
		return ""
	}

	var sb strings.Builder
	pad := strings.Repeat(" ", padding)

	if m.pbs == nil {
		m.pbs = make([]*ProgressBar, 0)
	}
	if len(m.pbs) != len(m.Pbs) {
		m.pbs = pbMapToSortedSlice(m.Pbs, m.width)
	}

	for _, pb := range m.pbs {
		var stats string
		if pb.Pc != nil {
			stats = fmt.Sprintf("%d/%d", pb.Pc.Downloaded.Load(), pb.Pc.Total)
		} else {
			stats = getDLStatus(pb.Pw)
		}

		status := renderPbStatus(pb.Status)

		bar := lipgloss.JoinHorizontal(lipgloss.Top, pad, pb.Progress.View(),
			pad, stats,
			pad, pb.FileName,
			pad, status)

		sb.WriteString("\n")
		sb.WriteString(pbStyle(bar))
	}

	sb.WriteString("\n\n\n")
	sb.WriteString(helpStyle("Press ctrl+c to quit\n\n"))

	return sb.String()
}

func renderPbStatus(s string) string {
	switch s {
	case DownloadingStatus:
		return statusStyle.Foreground(lipgloss.Color(DownloadingColor)).Render(s)
	case MergingStatus:
		return statusStyle.Foreground(lipgloss.Color(MergingColor)).Render(s)
	case CompleteStatus:
		return statusStyle.Foreground(lipgloss.Color(CompleteColor)).Render(s)
	case RetryStatus:
		return statusStyle.Foreground(lipgloss.Color(RetryColor)).Render(s)
	case ErrStatus:
		return statusStyle.Foreground(lipgloss.Color(ErrColor)).Render(s)
	default:
		return ""
	}
}

func pbMapToSortedSlice(m map[string]*ProgressBar, w int) []*ProgressBar {
	res := make([]*ProgressBar, 0, len(m))
	for _, v := range m {
		v.Progress.Width = w
		res = append(res, v)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].FileName < res[j].FileName
	})

	return res
}

func getDLStatus(pw *ProgressWriter) string {
	d := time.Duration(pw.DLTime) * time.Second
	return fmt.Sprintf("%.2f MiB/%.2f MiB %s %s", float64(pw.Downloaded)/mib, float64(pw.Total)/mib, util.FormatSize(pw.Speed), d.String())
}

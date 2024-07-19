package select_pkg

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type itemDelegate struct{}

var (
	docStyle          = lipgloss.NewStyle().Margin(1, 2)
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func Run() ([]string, error) {
	ssoFolders, err := listPackages()
	if err != nil {
		return []string{}, fmt.Errorf("listing packages: %w", err)
	}

	var items []list.Item

	for _, str := range ssoFolders {
		items = append(items, item{title: str})
	}

	m := model{
		list: list.New(items, itemDelegate{}, 0, 0),
	}

	m.list.Title = "Packages"

	p := tea.NewProgram(m, tea.WithAltScreen())

	genericModel, err := p.Run()
	if err != nil {
		return []string{}, fmt.Errorf("running program: %w", err)
	}

	resultModel, ok := genericModel.(model)
	if !ok {
		return []string{}, fmt.Errorf("unable to cast resultModel to model")
	}
	fmt.Printf("Resulting choice is: %s\n", resultModel.choice)

	return []string{}, nil
}

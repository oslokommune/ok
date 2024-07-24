package charming

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: Create `ok aws` command
// TODO: Create `ok aws ecs-exec` command
// TODO: Make `ok do ecs-exec` list clusters
// TODO: Make `ok do ecs-exec` list clusters and select one to execute
// TODO:

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)

type item struct {
	title, arn string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.arn }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:

		switch keypress := msg.String(); keypress {

		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.Title())
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Not hungry? That's cool.")
	}
	return docStyle.Render(m.list.View())
}

func listClusters() ([]string, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Printf("error loading configuration: %v\n", err)
		return nil, err
	}

	ecsSvc := ecs.NewFromConfig(cfg)

	result, err := ecsSvc.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		return nil, fmt.Errorf("listing clusters: %w", err)
	}

	var clusters []string
	clusters = append(clusters, result.ClusterArns...)

	return clusters, nil
}

func basename(arn string) string {
	parts := strings.Split(arn, "/")
	return parts[len(parts)-1]
}

func Hehe() {
	hehe, _ := listClusters()

	var items []list.Item

	for _, str := range hehe {
		items = append(items, item{title: basename(str), arn: str})
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Clusters"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

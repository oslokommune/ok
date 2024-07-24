package interactive

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oslokommune/ok/pkg/pkg/common"
)

type Result struct {
	Choice  string
	Aborted bool
}

func Run(pkgManifestFilename string) (Result, error) {
	listItems, err := getListItems(pkgManifestFilename)
	if err != nil {
		return Result{}, fmt.Errorf("getting items: %w", err)
	}

	m := model{
		list: list.New(listItems, itemDelegate{}, 0, 0),
	}

	m.list.Title = "Select package to install:"

	return run(m)
}

func getListItems(pkgManifestFilename string) ([]list.Item, error) {
	manifest, err := common.LoadPackageManifest(pkgManifestFilename)
	if err != nil {
		return nil, fmt.Errorf("loading package manifest: %w", err)
	}

	var items []list.Item

	for _, p := range manifest.Packages {
		items = append(items, item{
			outputFolder: p.OutputFolder,
			ref:          p.Ref,
		})
	}
	return items, nil
}

func run(m model) (Result, error) {
	p := tea.NewProgram(m, tea.WithAltScreen())

	genericModel, err := p.Run()
	if err != nil {
		return Result{}, fmt.Errorf("running program: %w", err)
	}

	resultModel, ok := genericModel.(model)
	if !ok {
		return Result{}, fmt.Errorf("unable to cast resultModel to model")
	}

	result := Result{
		Choice:  resultModel.choice,
		Aborted: resultModel.aborted,
	}

	return result, nil
}

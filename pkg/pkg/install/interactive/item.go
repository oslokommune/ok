package interactive

type item struct {
	outputFolder string
	ref          string
}

func (i item) OutputFolder() string { return i.outputFolder }
func (i item) Ref() string          { return i.ref }
func (i item) FilterValue() string  { return i.outputFolder }

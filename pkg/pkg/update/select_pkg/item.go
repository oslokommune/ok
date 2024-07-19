package select_pkg

type item struct {
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "lol" }
func (i item) FilterValue() string { return i.title }

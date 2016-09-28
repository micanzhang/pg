package pg

type Action string

const (
	GenAction    Action = "gen"
	ListAction   Action = "list"
	NewAction    Action = "new"
	UpdateAction Action = "update"
	RemoveAction Action = "remove"
	InfoAction   Action = "info"
)

func (a Action) Run() {
	switch a {
	case GenAction:
	case ListAction:
	case NewAction:
	case UpdateAction:
	case RemoveAction:
	case InfoAction:
	default:
	}
}

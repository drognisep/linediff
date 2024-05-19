package linediff

type Tag int

const (
	Same Tag = iota
	Added
	Removed
)

type DiffSet struct {
	segments []string
	tags     []Tag
}

func (s *DiffSet) String() string {
	return ""
}

func Diff(a, b string) *DiffSet {
	return nil
}

func DiffSplit(a, b string, split Splitter) *DiffSet {
	return nil
}

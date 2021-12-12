package graph

type Graph interface {
	Size() int
	Idx(string) int
	Name(int) string
	MaybeCreate(string)
	Neighbours(int) []int
	Connect(string, string, bool)
	ChannelIterator() <-chan ThisCouldHaveBeenAvoidedIfGoLetYouImplementRangeOverCustomTypes
}

type ThisCouldHaveBeenAvoidedIfGoLetYouImplementRangeOverCustomTypes struct {
	int
	string
}

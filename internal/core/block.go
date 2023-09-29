package core

type Block interface {
}

type ExpandBlock struct {
	Name string
}

type ContentBlock struct {
	Content string
}

type UseBlock struct {
	Name string
}

type ForBlock struct {
	Name     string
	Iterable string
	Content  string
}

type OutOnlyBlock struct {
}

type RawBlock struct {
	Data []byte
}

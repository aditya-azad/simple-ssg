package core

type Block interface {
	expand()
}

type ExpandBlock struct {
	name string
}

type ContentBlock struct {
	content string
}

type UseBlock struct {
	name string
}

type ForBlock struct {
	name     string
	iterable string
	content  string
}

type OutOnlyBlock struct {
}

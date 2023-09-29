package core

type Block interface {
	expand()
}

type TemplateBlock struct {
	name  string
	props map[string]string
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

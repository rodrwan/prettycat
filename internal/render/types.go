package render

type Kind string

const (
	KindMarkdown Kind = "markdown"
	KindCode     Kind = "code"
	KindPlain    Kind = "plain"
)

type Options struct {
	Color bool
	Width int
}

type Doc struct {
	Title string
	Body  string
	Kind  Kind
}

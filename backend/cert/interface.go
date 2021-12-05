package cert

type Generator interface {
	Generate() error
}

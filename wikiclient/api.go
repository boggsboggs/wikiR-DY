package wikiclient

type Client interface {
	GetAllLinksFromPage(title string) ([]string, error)
	GetAllLinksToPage(title string) ([]string, error)
}

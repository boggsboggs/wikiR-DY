package wikiclient

type Client interface {
	GetAllLinksInPage(title string) ([]string, error)
	GetAllLinksInURL(url string) ([]string, error)
}

package dto

type SaveLinkRequest struct {
	URL string `json:"url"`
}

type SaveLinkResponse struct {
	ShortLink string `json:"short_link"`
	Error     string `json:"error"`
}

type RedirectResponse struct {
	OriginalLink string `json:"original_url"`
	Error        string `json:"error"`
}

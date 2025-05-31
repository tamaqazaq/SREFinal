package social

import (
	"fmt"
	"net/url"
)

func GetFacebookShareLink(targetURL string, description string) string {
	baseURL := "https://www.facebook.com/sharer/sharer.php?u="
	return fmt.Sprintf("%s%s&quote=%s", baseURL, url.QueryEscape(targetURL), url.QueryEscape(description))
}

func GetTwitterShareLink(targetURL string, text string) string {
	baseURL := "https://twitter.com/intent/tweet?url="
	return fmt.Sprintf("%s%s&text=%s", baseURL, url.QueryEscape(targetURL), url.QueryEscape(text))
}
func GetLinkedInShareLink(targetURL, title, summary string) string {
	baseURL := "https://www.linkedin.com/shareArticle?mini=true&url="
	return fmt.Sprintf("%s%s&title=%s&summary=%s&source=%s",
		baseURL,
		url.QueryEscape(targetURL),
		url.QueryEscape(title),
		url.QueryEscape(summary))
}

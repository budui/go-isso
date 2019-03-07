package service

import (
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type MDConverter struct {
	r      blackfriday.Renderer
	e      blackfriday.Extensions
	policy *bluemonday.Policy
}

func NewMDConverter(markup conf.Markup) *MDConverter {
	HTMLFlags := blackfriday.CommonHTMLFlags
	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: HTMLFlags,
	})
	e := blackfriday.Strikethrough | blackfriday.Autolink | blackfriday.FencedCode | blackfriday.SpaceHeadings
	p := bluemonday.NewPolicy()
	// href for <a> and align for <table>
	p.AllowAttrs("align").Matching(bluemonday.Paragraph).Globally()
	p.AllowAttrs("href").Matching(bluemonday.Paragraph).Globally()
	for _, attr := range markup.AdditionalAllowedAttributes {
		p.AllowAttrs(attr).Matching(bluemonday.Paragraph).Globally()
	}

	p.AllowElements("a", "p", "hr", "br", "ol", "ul", "li", "pre", "code", "blockquote",
		"del", "ins", "strong", "em", "h1", "h2", "h3", "h4", "h5", "h6", "table", "thead", "tbody", "th", "td")
	p.AllowElements(markup.AdditionalAllowedElements...)

	return &MDConverter{r: r, e: e, policy: p}
}

func (mdc *MDConverter) Run(input string) string {
	unsafe := blackfriday.Run([]byte(input), blackfriday.WithNoExtensions(), blackfriday.WithExtensions(mdc.e),
		blackfriday.WithRenderer(mdc.r))
	return string(mdc.policy.SanitizeBytes(unsafe))
}

package golang

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

// MetaImport represents the parsed <meta name="go-import"
// content="prefix vcs reporoot" /> tags from HTML files.
type MetaImport struct {
	Prefix, VCS, RepoRoot string
}

func (m *MetaImport) OrgRepo() (string, string) {
	repoRoot := strings.TrimSuffix(m.RepoRoot, ".git")

	parts := strings.Split(repoRoot, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2], parts[len(parts)-1]
	}
	panic(fmt.Errorf("unknown repo root: %s", m.RepoRoot))
}

func metaContent(doc *html.Node, name string) (string, error) {
	var meta *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "meta" {
			for _, attr := range node.Attr {
				if attr.Key == "name" && attr.Val == name {
					meta = node
					return
				}
			}

		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if meta != nil {
		for _, attr := range meta.Attr {
			if attr.Key == "content" {
				return attr.Val, nil
			}
		}
	}
	return "", fmt.Errorf("missing <meta name=%s> in the node tree", name)
}

func GetMetaImport(i string) (*MetaImport, error) {
	url := fmt.Sprintf("https://%s?go-get=1", i)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	content, err := metaContent(doc, "go-import")
	if err != nil {
		return nil, err
	}

	f := strings.Fields(content)

	return &MetaImport{
		Prefix:   f[0],
		VCS:      f[1],
		RepoRoot: f[2],
	}, nil
}

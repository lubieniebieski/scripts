package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Link represents a link along with its reference number
type Link struct {
	Name        string
	URL         string
	ReferenceNo int
	ID          string
}

func (l *Link) IsFootnote() bool {
	return strings.HasPrefix(l.ID, "^")
}

// MarkdownConverter converts inline links to reference links in markdown files
type MarkdownConverter struct {
	FileName string
	Links    []Link
}

func (c *MarkdownConverter) extractFootnotesFromBuffer(content []byte) {
	footnoteRegex := regexp.MustCompile(`\[(\^\d+)\]:\s(.+)`)
	matches := footnoteRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink("", string(match[2]), string(match[1]))
	}
}
func (c *MarkdownConverter) extractMarkdownLinksFromBuffer(content []byte) {
	inlineLinkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	matches := inlineLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink(string(match[1]), string(match[2]), "")
	}

	refLinkRegex := regexp.MustCompile(`\[(.*?)\]\[(.*?)\]`)
	matches = refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink(string(match[1]), "", string(match[2]))
	}
	c.extractReferenceLinksFromBuffer(content)
	c.extractFootnotesFromBuffer(content)
}
func (c *MarkdownConverter) extractReferenceLinksFromBuffer(content []byte) {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)`)
	matches := refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		matchID := string(match[1])
		matchURL := string(match[2])

		for i := range c.Links {
			if c.Links[i].ID == matchID {
				c.Links[i].URL = matchURL
				break
			}
		}
	}
}
func (c *MarkdownConverter) addLink(name string, url string, ID string) {
	link := Link{Name: name, URL: url, ID: ID}
	if !link.IsFootnote() {
		link.ReferenceNo = len(c.Links) + 1
	}
	c.Links = append(c.Links, link)
}

func (c *MarkdownConverter) extractLinks() {
	content, err := os.ReadFile(c.FileName)
	if err != nil {
		panic(err)
	}
	c.extractMarkdownLinksFromBuffer(content)
}

func (c *MarkdownConverter) updateBuffer(buffer []byte) []byte {
	c.extractMarkdownLinksFromBuffer(buffer)
	for _, link := range c.Links {
		if link.IsFootnote() {
			footnoteRef := "\n" + link.AsReference()
			buffer = bytes.ReplaceAll(buffer, []byte(footnoteRef), []byte(""))
		} else {
			linkRef := fmt.Sprintf("[%s][%d]", link.Name, link.ReferenceNo)
			linkRegex := regexp.MustCompile(fmt.Sprintf(`\[%s\]\(.*?\)`, link.Name))
			buffer = linkRegex.ReplaceAll(buffer, []byte(linkRef))
		}
	}

	if len(c.referencesList()) == 0 {
		return buffer
	}
	buffer = append(buffer, []byte("\n\n")...)
	buffer = append(buffer, []byte(strings.Join(c.referencesList(), "\n"))...)
	return buffer
}

func (c *MarkdownConverter) updateFile() {
	content, err := os.ReadFile(c.FileName)
	if err != nil {
		panic(err)
	}
	updatedContent := c.updateBuffer(content)
	err = os.WriteFile(c.FileName, updatedContent, 0644)
	if err != nil {
		panic(err)
	}
}

func (c *MarkdownConverter) referencesList() []string {
	var result []string
	for _, link := range c.Links {
		result = append(result, link.AsReference())
	}
	return result
}

func (l *Link) AsReference() string {
	var ref string
	if l.ID != "" {
		ref = l.ID
	} else {
		ref = fmt.Sprint(l.ReferenceNo)
	}
	return fmt.Sprintf("[%s]: %s", ref, l.URL)
}

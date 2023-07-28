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

func (l *Link) IsReference() bool {
	match, _ := regexp.MatchString(`^\D+$`, l.ID)
	return match && !l.IsFootnote()
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

	refLinkRegex := regexp.MustCompile(`\[([^\]]*)?\]\[(\w+)\]`)
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
	if url != "" {
		for _, link := range c.Links {
			if link.URL == url {
				return
			}
		}
	}

	if ID != "" {
		for _, link := range c.Links {
			if link.ID == ID {
				return
			}
		}
	}

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

func removeLineContainingString(buffer []byte, str string) []byte {
	lines := bytes.Split(buffer, []byte("\n"))
	var newLines [][]byte
	for _, line := range lines {
		if !bytes.Contains(line, []byte(str)) {
			newLines = append(newLines, line)
		}
	}
	return bytes.Join(newLines, []byte("\n"))
}

func (c *MarkdownConverter) updateBuffer(buffer []byte) []byte {
	c.extractMarkdownLinksFromBuffer(buffer)

	for _, link := range c.Links {
		if link.IsFootnote() || link.IsReference() {
			buffer = removeLineContainingString(buffer, link.AsReference())
		} else {
			linkRef := fmt.Sprintf("[%d]", link.ReferenceNo)
			linkRegex := regexp.MustCompile(fmt.Sprintf(`\(%s\)`, link.URL))
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := os.Args[1]
	converter := MarkdownConverter{FileName: filename}
	content, _ := os.ReadFile(filename)
	converter.extractLinksFromReferences(content)
	content = converter.clearReferences(content)
	newContent := converter.updateBuffer(content)

	err := os.WriteFile(filename, newContent, 0644)

	if err != nil {
		fmt.Printf("Error updating file: %v\n", err)
		return
	}

	fmt.Printf("File %s updated successfully!\n", filename)
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

func (l *Link) AsMarkdownLink() string {
	if l.URL == "" {
		return l.Name
	}
	if l.Name == "" {
		return fmt.Sprintf("<%s>", l.URL)
	}
	return fmt.Sprintf("[%s](%s)", l.Name, l.URL)
}

func (c *MarkdownConverter) extractLinksFromReferences(content []byte) {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)`)
	matches := refLinkRegex.FindAllSubmatch(content, -1)

	for _, match := range matches {
		c.addLink(string(""), string(match[2]), string(match[1]))
	}
}

func (c *MarkdownConverter) clearReferences(content []byte) []byte {
	refLinkRegex := regexp.MustCompile(`\[(.*?)\]:\s(.+)\n`)
	return refLinkRegex.ReplaceAll(content, []byte(""))
}

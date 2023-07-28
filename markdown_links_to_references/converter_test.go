package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestAddLink(t *testing.T) {
	converter := MarkdownConverter{}
	converter.addLink("Google", "https://www.google.com", "ref")

	expectedLinks := []Link{
		{Name: "Google", URL: "https://www.google.com", ReferenceNo: 1, ID: "ref"},
	}

	if len(converter.Links) != len(expectedLinks) {
		t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
	}

	for i, link := range converter.Links {
		if link.Name != expectedLinks[i].Name {
			t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
		}
		if link.URL != expectedLinks[i].URL {
			t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
		}
		if link.ReferenceNo != expectedLinks[i].ReferenceNo {
			t.Errorf("Expected link reference number '%d', but got '%d'", expectedLinks[i].ReferenceNo, link.ReferenceNo)
		}
		if link.ID != expectedLinks[i].ID {
			t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
		}
	}
}
func TestExtractMarkdownLinksFromBuffer(t *testing.T) {
	t.Run("works with inline links", func(t *testing.T) {
		content := []byte(`
        [Google](https://www.google.com) fdafd
        [GitHub][1]
        [Wikipedia][ref] fdsf ds
        [Example page][Example]
        [Invalid Link]
        [1]: https://github.com
        [ref]: https://www.wikipedia.org
        [Example]: https://example.com
    `)

		expectedLinks := []Link{
			{Name: "Google", URL: "https://www.google.com"},
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "Wikipedia", URL: "https://www.wikipedia.org", ID: "ref"},
			{Name: "Example page", URL: "https://example.com", ID: "Example"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(content)

		if len(converter.Links) != len(expectedLinks) {
			t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
		}

		for i, link := range converter.Links {
			if link.Name != expectedLinks[i].Name {
				t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
			}
			if link.URL != expectedLinks[i].URL {
				t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
			}
			if link.ID != expectedLinks[i].ID {
				t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
			}
		}
	})
	t.Run("works with footnotes too", func(t *testing.T) {
		mixedContent := []byte(`
        [Google](https://www.google.com)
        [GitHub][1]
		footnote example[^1]
        [1]: https://github.com
        [^1]: some footnote
    `)
		expectedLinks := []Link{
			{Name: "Google", URL: "https://www.google.com"},
			{Name: "GitHub", URL: "https://github.com", ID: "1"},
			{Name: "", URL: "some footnote", ID: "^1"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(mixedContent)

		if len(converter.Links) != len(expectedLinks) {
			t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
		}

		for i, link := range converter.Links {
			if link.Name != expectedLinks[i].Name {
				t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
			}
			if link.URL != expectedLinks[i].URL {
				t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
			}
			if link.ID != expectedLinks[i].ID {
				t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
			}
		}
	})
	t.Run("works with custom refs", func(t *testing.T) {
		mixedContent := []byte(`
        [test][test1] and [another test][test2]
				[test1]: https://www.google.com
				[test2]: https://github.com
    `)
		expectedLinks := []Link{
			{Name: "test", URL: "https://www.google.com", ID: "test1"},
			{Name: "another test", URL: "https://github.com", ID: "test2"},
		}

		converter := MarkdownConverter{}
		converter.extractMarkdownLinksFromBuffer(mixedContent)

		if len(converter.Links) != len(expectedLinks) {
			t.Errorf("Expected %d links, but got %d", len(expectedLinks), len(converter.Links))
		}

		for i, link := range converter.Links {
			if link.Name != expectedLinks[i].Name {
				t.Errorf("Expected link name '%s', but got '%s'", expectedLinks[i].Name, link.Name)
			}
			if link.URL != expectedLinks[i].URL {
				t.Errorf("Expected link URL '%s', but got '%s'", expectedLinks[i].URL, link.URL)
			}
			if link.ID != expectedLinks[i].ID {
				t.Errorf("Expected link ID '%s', but got '%s'", expectedLinks[i].ID, link.ID)
			}
		}
	})
}

func TestUpdateBuffer(t *testing.T) {
	t.Run("first example", func(t *testing.T) {
		content := []byte(`# Test Markdown File
This is a test file for the MarkdownConverter.
Here is a link to [Google](https://www.google.com).
Here is another link to [GitHub](https://github.com).
And here is a third link to [Wikipedia](https://www.wikipedia.org).
There are some footnotes[^1], too.
[^2]: Some footnote.`)

		converter := MarkdownConverter{}
		updatedContent := converter.updateBuffer(content)

		expectedContent := []byte(`# Test Markdown File
This is a test file for the MarkdownConverter.
Here is a link to [Google][1].
Here is another link to [GitHub][2].
And here is a third link to [Wikipedia][3].
There are some footnotes[^1], too.

[1]: https://www.google.com
[2]: https://github.com
[3]: https://www.wikipedia.org
[^2]: Some footnote.`)

		if !bytes.Equal(updatedContent, expectedContent) {
			t.Errorf("Expected content to be %s, but got %s", expectedContent, updatedContent)
		}
	})
	t.Run("second example", func(t *testing.T) {
		content := []byte(`[normal link](https://www.google.com), [ref link][ref], something[^1].
[^1]: Some footnote.
[ref]: https://www.wikipedia.org`)

		converter := MarkdownConverter{}
		updatedContent := converter.updateBuffer(content)

		expectedContent := []byte(`[normal link][1], [ref link][ref], something[^1].

[1]: https://www.google.com
[ref]: https://www.wikipedia.org
[^1]: Some footnote.`)

		if !bytes.Equal(updatedContent, expectedContent) {
			t.Errorf("Expected content to be %s, but got %s", expectedContent, updatedContent)
		}
	})

	t.Run("third example", func(t *testing.T) {
		content := []byte(`I watched [Extraction (2020) - IMDb](https://www.imdb.com/title/tt8936646/) and [Extraction II (2023) - IMDb](https://www.imdb.com/title/tt12263384/) -- I'm a simple man and I enjoyed them!`)

		converter := MarkdownConverter{}
		updatedContent := converter.updateBuffer(content)

		expectedContent := []byte(`I watched [Extraction (2020) - IMDb][1] and [Extraction II (2023) - IMDb][2] -- I'm a simple man and I enjoyed them!

[1]: https://www.imdb.com/title/tt8936646/
[2]: https://www.imdb.com/title/tt12263384/`)

		if !bytes.Equal(updatedContent, expectedContent) {
			t.Errorf("Expected content to be %s, but got %s", expectedContent, updatedContent)
		}
	})
}
func TestLinkAsReference(t *testing.T) {
	link := Link{Name: "Google", URL: "https://www.google.com", ReferenceNo: 1}
	expected := "[1]: https://www.google.com"
	if link.AsReference() != expected {
		t.Errorf("Link.AsReference() = %s; expected %s", link.AsReference(), expected)
	}

	link = Link{Name: "Wikipedia", URL: "https://www.wikipedia.org", ReferenceNo: 2}
	expected = "[2]: https://www.wikipedia.org"
	if link.AsReference() != expected {
		t.Errorf("Link.AsReference() = %s; expected %s", link.AsReference(), expected)
	}

	link = Link{URL: "it's a footnote", ID: "^1"}
	expected = "[^1]: it's a footnote"
	if link.AsReference() != expected {
		t.Errorf("Link.AsReference() = %s; expected %s", link.AsReference(), expected)
	}
}

func TestMain(t *testing.T) {
	// Create temporary test file
	tmpfile, err := ioutil.TempFile("", "testfile.*.md")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test content to file
	testContent := []byte("This is a [test link](https://example.com). But [this][1] already exists.\n\n[1]: https://github.com")
	if _, err := tmpfile.Write(testContent); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	// Close file before running main function
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temporary file: %v", err)
	}

	// Run main function with test file
	os.Args = []string{"main.go", tmpfile.Name()}
	main()

	// Read updated content from file
	updatedContent, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Error reading updated file: %v", err)
	}

	// Check if content was updated correctly
	expectedContent := []byte("This is a [test link][2]. But [this][1] already exists.\n\n[1]: https://github.com\n[2]: https://example.com")
	if !bytes.Equal(updatedContent, expectedContent) {
		t.Fatalf("Content was not updated correctly. Expected:\n%s\nGot:\n%s", expectedContent, updatedContent)
	}
}

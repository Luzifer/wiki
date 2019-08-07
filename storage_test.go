package main

import "testing"

func TestStoredFileParse(t *testing.T) {
	var (
		err   error
		file  string
		sFile *storedFile
	)

	// Case: Proper file with header

	file = `
---
key: value
---

# Header

content
`

	sFile, err = storedFileFromString(file)
	if err != nil {
		t.Fatalf("Parsing of proper file errored: %s", err)
	}

	if sFile.Content != "# Header\n\ncontent" {
		t.Errorf("Content did not match expectation: %q", sFile.Content)
	}

	if len(sFile.Meta) != 1 || sFile.GetMetaString("key") != "value" {
		t.Errorf("Metadata did not match expectation: %#v", sFile.Meta)
	}

	// Case: No header

	file = "# Header\n\ncontent"

	sFile, err = storedFileFromString(file)
	if err != nil {
		t.Fatalf("Parsing of proper file errored: %s", err)
	}

	if sFile.Content != "# Header\n\ncontent" {
		t.Errorf("Content did not match expectation: %q", sFile.Content)
	}

	if len(sFile.Meta) != 0 {
		t.Errorf("Metadata did not match expectation: %#v", sFile.Meta)
	}
}

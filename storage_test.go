package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	require.NoError(t, err)

	assert.Equal(t, "# Header\n\ncontent", sFile.Content)
	assert.Equal(t, map[string]any{"key": "value"}, sFile.Meta)

	// Case: No header

	file = "# Header\n\ncontent"

	sFile, err = storedFileFromString(file)
	require.NoError(t, err)

	assert.Equal(t, "# Header\n\ncontent", sFile.Content)
	assert.Len(t, sFile.Meta, 0)
}

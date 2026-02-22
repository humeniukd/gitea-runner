package functions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	_, err := Format("Hello {0}, you have {1} new messages", "Alice", 5)
	require.NoError(t, err)
	// fmt.Println(s) // Hello Alice, you have 5 new messages
}

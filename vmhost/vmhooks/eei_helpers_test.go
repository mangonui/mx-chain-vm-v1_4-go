package vmhooks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetESDTRoles_WellFormedBufferReturnsExpectedRoles(t *testing.T) {
	t.Parallel()

	role := []byte("ESDTRoleLocalMint")
	buf := []byte{'\n', byte(len(role))}
	buf = append(buf, role...)

	got := getESDTRoles(buf)
	require.NotEqual(t, int64(0), got)
}

func TestGetESDTRoles_MalformedBufferFailsClosed(t *testing.T) {
	t.Parallel()

	role := []byte("ESDTRoleLocalMint")
	buf := []byte{'\n', byte(len(role))}
	buf = append(buf, role...)

	wellFormedResult := getESDTRoles(buf)
	require.NotEqual(t, int64(0), wellFormedResult)

	malformed := append(buf, '\n', 200)

	got := getESDTRoles(malformed)
	require.Equal(t, int64(0), got)
}

func TestGetESDTRoles_BufferEndingOnDelimiterFailsClosed(t *testing.T) {
	t.Parallel()

	got := getESDTRoles([]byte{'\n'})
	require.Equal(t, int64(0), got)
}

func TestGetESDTRoles_EmptyBufferReturnsZero(t *testing.T) {
	t.Parallel()

	got := getESDTRoles([]byte{})
	require.Equal(t, int64(0), got)
}

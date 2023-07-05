package onesecmail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenRandomMailbox(t *testing.T) {
	got, err := DefaultClient.GenRandomMailbox(2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(got))
	got, err = DefaultClient.GenRandomMailbox(3)
	require.NoError(t, err)
	assert.Equal(t, 3, len(got))
	got, err = DefaultClient.GenRandomMailbox(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(got))
	got, err = DefaultClient.GenRandomMailbox(0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(got))
	got, err = DefaultClient.GenRandomMailbox(-1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(got))
}

func TestGetDomainList(t *testing.T) {
	got, err := DefaultClient.GetDomainList()
	require.NoError(t, err)
	want := []string{
		"1secmail.com",
		"1secmail.org",
		"1secmail.net",
		"kzccv.com",
		"qiott.com",
		"wuuvo.com",
		"icznn.com",
		"ezztt.com",
	}
	assert.Equal(t, want, got)
}

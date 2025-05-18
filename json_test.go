package main

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ReadJSON_Success(t *testing.T) {
	type testcase struct {
		filename string
	}
	tests := []testcase{
		{
			filename: "testdata/top-level-array.json",
		},
		{
			filename: "testdata/top-level-object.json",
		},
	}

	testfn := func(t *testing.T, tt testcase) {
		t.Helper()
		f, err := os.Open(tt.filename)
		require.NoError(t, err)

		entries, err := ReadJSON(f)
		require.NoError(t, err)
		require.Len(t, entries, 2)
		// The object-based input produces entries in random order (due to
		// map iteration), so sort by DN first.
		slices.SortFunc(entries, func(a, b *Entry) int { return slices.Compare(a.DN, b.DN) })
		require.Equal(t, NewDN("o=example,dc=example,dc=com"), entries[0].DN)
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) { testfn(t, tt) })
	}
}

func Test_ReadJSON_Failure(t *testing.T) {
	type testcase struct {
		filename string
	}
	tests := []testcase{
		{
			filename: "testdata/top-level-non-aggregate.json",
		},
		{
			filename: "testdata/non-aggregate.json",
		},
		{
			filename: "testdata/no-dn.json",
		},
		{
			filename: "testdata/no-object-class.json",
		},
		{
			filename: "testdata/multiple-dn.json",
		},
		{
			filename: "testdata/empty.json",
		},
		{
			filename: "testdata/invalid-value.json",
		},
	}

	testfn := func(t *testing.T, tt testcase) {
		t.Helper()
		f, err := os.Open(tt.filename)
		require.NoError(t, err)

		_, err = ReadJSON(f)
		require.Error(t, err)
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) { testfn(t, tt) })
	}
}

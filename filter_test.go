package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/samf/racewalk/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilters(t *testing.T) {
	var initFilter filter
	dirpath := setup()
	defer cleanup(dirpath)

	gitpath := filepath.Join(dirpath, ".git")
	err := os.Mkdir(gitpath, 0777)
	require.NoError(t, err)

	mercpath := filepath.Join(dirpath, ".hg")
	err = os.Mkdir(mercpath, 0777)
	require.NoError(t, err)

	otherpath := filepath.Join(dirpath, "other")
	err = os.Mkdir(otherpath, 0777)
	require.NoError(t, err)

	t.Run("noGit", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		f := initFilter.noGit()

		fnode, err := racewalk.MakeFileNode(gitpath)
		require.NoError(err)
		assert.True(f(*fnode))

		fnode, err = racewalk.MakeFileNode(mercpath)
		require.NoError(err)
		assert.False(f(*fnode))
	})

	t.Run("noMerc", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		f := initFilter.noHG()

		fnode, err := racewalk.MakeFileNode(mercpath)
		require.NoError(err)
		assert.True(f(*fnode))

		fnode, err = racewalk.MakeFileNode(gitpath)
		require.NoError(err)
		assert.False(f(*fnode))
	})

	t.Run("chaining", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		f := initFilter.noHG().noGit()

		fnode, err := racewalk.MakeFileNode(mercpath)
		require.NoError(err)
		assert.True(f(*fnode))

		fnode, err = racewalk.MakeFileNode(gitpath)
		require.NoError(err)
		assert.True(f(*fnode))

		fnode, err = racewalk.MakeFileNode(otherpath)
		require.NoError(err)
		assert.False(f(*fnode))
	})
}

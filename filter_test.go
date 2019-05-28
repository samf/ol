package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	t.Run("oneGit", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		dirname := fmt.Sprintf("git-%v-%v", time.Now().Unix(), os.Getpid())
		err := os.Mkdir(dirname, 0777)
		require.NoError(err)
		defer func() {
			os.RemoveAll(dirname)
		}()

		f := initFilter.oneGit()

		fnode, err := racewalk.MakeFileNode(dirname)
		require.NoError(err)
		assert.False(f(*fnode))

		err = os.Mkdir(filepath.Join(dirname, ".git"), 0777)
		require.NoError(err)

		fnode, err = racewalk.MakeFileNode(dirname)
		require.NoError(err)
		assert.True(f(*fnode))
	})

	t.Run("oneMerc", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		dirname := fmt.Sprintf("hg-%v-%v", time.Now().Unix(), os.Getpid())
		err := os.Mkdir(dirname, 0777)
		require.NoError(err)
		defer func() {
			os.RemoveAll(dirname)
		}()

		f := initFilter.oneHG()

		fnode, err := racewalk.MakeFileNode(dirname)
		require.NoError(err)
		assert.False(f(*fnode))

		err = os.Mkdir(filepath.Join(dirname, ".hg"), 0777)
		require.NoError(err)

		fnode, err = racewalk.MakeFileNode(dirname)
		require.NoError(err)
		assert.True(f(*fnode))
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

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapStorageInt8(t *testing.T) {
	m := MapStorage[int8]{}

	var val int8
	var ok bool
	var err error

	val, ok, err = m.Get("one")
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Zero(t, val)

	m.Set("zero", 0)
	m.Set("one", 1)
	m.Set("twelve", 12)

	val, ok, err = m.Get("zero")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 0, val)

	val, ok, err = m.Get("one")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 1, val)

	val, ok, err = m.Get("twelve")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 12, val)

	keys, err := m.Keys()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(keys))

	val, ok, err = m.Get("eighty") // not in map
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.EqualValues(t, 0, val)

	val, ok, err = m.Delete("twelve")
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 12, val)

	val, ok, err = m.Get("twelve") // now should be deleted
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.EqualValues(t, 0, val)
}

package engine

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testTroopID1 = 12345
	testTroopID2 = 23456
	testTroopID3 = 34567

	testTickNumber1 = 1
	testTickNumber2 = 2
	testTickNumber3 = 3
)

var testPos1 = Pos{
	X: 1,
	Y: 2,
}
var testPos2 = Pos{
	X: 2,
	Y: 3,
}
var testPos3 = Pos{
	X: 3,
	Y: 4,
}

func Test_AddTroopData(t *testing.T) {
	tc := NewTroopCache()

	tc.AddTroopData(testTickNumber1, testTroopID1, testPos1)
	tc.AddTroopData(testTickNumber1, testTroopID3, testPos2)
	tc.AddTroopData(testTickNumber2, testTroopID2, testPos2)
	tc.AddTroopData(testTickNumber3, testTroopID3, testPos3)

	require.Contains(t, tc.cache, testTroopID1)
	assert.Contains(t, tc.cache[testTroopID1].cache, testTickNumber1)

	latestUpdateFinder1 := tc.cache[testTroopID1]

	res1, found := latestUpdateFinder1.getLatestUpdate(testTickNumber1)
	require.True(t, found)
	assert.True(t, cmp.Equal(res1, testPos1))

	require.Contains(t, tc.cache, testTroopID1)
	assert.Contains(t, tc.cache[testTroopID1].cache, testTickNumber1)

	latestUpdateFinder2 := tc.cache[testTroopID2]

	res2, found := latestUpdateFinder2.getLatestUpdate(testTickNumber3)
	require.True(t, found)
	assert.True(t, cmp.Equal(res2, testPos2))

	latestUpdateFinder3 := tc.cache[testTroopID3]

	res31, found := latestUpdateFinder3.getLatestUpdate(testTickNumber1)
	require.True(t, found)
	assert.True(t, cmp.Equal(res31, testPos2))

	res32, found := latestUpdateFinder3.getLatestUpdate(testTickNumber3)
	require.True(t, found)
	assert.True(t, cmp.Equal(res32, testPos3))
}

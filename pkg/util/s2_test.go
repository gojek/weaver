package util

import (
	"testing"

	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestS2Len(t *testing.T) {
	list := s2List{
		s2.CellID(23458045904904),
		s2.CellID(23458222224904),
		s2.CellID(23888885904904),
		s2.CellID(23458999999904),
	}
	assert.Equal(t, list.Len(), 4)
}

func TestS2Less(t *testing.T) {
	list := s2List{
		s2.CellID(4542091330435678208),
		s2.CellID(4542051748017078272),
	}
	assert.True(t, list.Less(0, 1))
}

func TestS2Swap(t *testing.T) {
	first := s2.CellID(4542091330435678208)
	second := s2.CellID(4542051748017078272)
	list := s2List{
		first,
		second,
	}
	list.Swap(0, 1)
	assert.Equal(t, list[0], second)
	assert.Equal(t, list[1], first)
}

func TestContainsOverlappingS2IDsTrueCase(t *testing.T) {
	list := []uint64{
		4542051748017078272,
		4542091330435678208,
	}
	assert.True(t, ContainsOverlappingS2IDs(list))
}

func TestContainsOverlappingS2IDsFalseCase(t *testing.T) {
	list := []uint64{
		4542051748017078272,
		1504976331727699968,
	}
	assert.False(t, ContainsOverlappingS2IDs(list))
}

func TestToS2List(t *testing.T) {
	idList := []uint64{
		4542051748017078272,
		4542091330435678208,
	}
	expectedS2List := s2List{
		s2.CellID(idList[0]),
		s2.CellID(idList[1]),
	}
	assert.Equal(t, expectedS2List, toS2List(idList))
}

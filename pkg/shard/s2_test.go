package shard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	s2ShardConfigWithoutDefault = json.RawMessage(`{
		"backends": {
			"3344472479136481280": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"3346530764903677952": { "backend_name": "jkt-b", "backend": "http://jkt.b.local"}
		},
		"shard_key_separator": ","
	}`)
	s2ShardConfig = json.RawMessage(`{
		"backends": {
			"3344472479136481280": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"3346530764903677952": { "backend_name": "jkt-b", "backend": "http://jkt.b.local"},
			"default": {"backend_name": "jkt-c", "backend": "http://jkt.c.local"}
		},
		"shard_key_separator": ","
	}`)
	s2ShardConfigBad = json.RawMessage(`{
		"backends": {
			"123454321": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"123459876": "bad-data"
		},
		"shard_key_separator": ","
	}`)

	s2ShardConfigInvalidS2ID = json.RawMessage(`{
		"backends": {
			"45435hb344hj3b": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"}
		},
		"shard_key_separator": ","
	}`)
	s2ShardConfigOverlapping = json.RawMessage(`{
		"backends": {
			"3344474311656055761": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"3344474311656055760": { "backend_name": "jkt-b", "backend": "http://jkt.b.local"},
			"3344474311656055744": { "backend_name": "jkt-c", "backend": "http://jkt.c.local"},
			"3344474311656055552": { "backend_name": "jkt-d", "backend": "http://jkt.d.local"},
			"3344474311656055808": { "backend_name": "jkt-e", "backend": "http://jkt.e.local"},
			"3573054715985026048": { "backend_name": "surbaya-a", "backend": "http://surbaya.a.local"}
		},
		"shard_key_separator": ","
	}`)

	s2SmartIDShardConfig = json.RawMessage(`{
		"backends": {
			"3344472479136481280": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"3346530764903677952": { "backend_name": "jkt-b", "backend": "http://jkt.b.local"},
			"default": {"backend_name": "jkt-c", "backend": "http://jkt.c.local"}
		},
		"shard_key_separator": "-",
		"shard_key_position":2
	}`)

	s2SmartIDShardConfigBad = json.RawMessage(`{
		"backends": {
			"3344472479136481280": { "backend_name": "jkt-a", "backend": "http://jkt.a.local"},
			"3346530764903677952": { "backend_name": "jkt-b", "backend": "http://jkt.b.local"},
			"default": {"backend_name": "jkt-c", "backend": "http://jkt.c.local"}
		},
		"shard_key_position":2
	}`)
)

func TestNewS2StrategySuccess(t *testing.T) {
	strategy, err := NewS2Strategy(s2ShardConfig)
	assert.NotNil(t, strategy)
	assert.Nil(t, err)
}

func TestNewS2StrategyFailure(t *testing.T) {
	sharder, err := NewS2Strategy(s2ShardConfigBad)
	assert.Nil(t, sharder)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "json: cannot unmarshal string")
}

func TestNewS2StrategyFailureWithMissingKeySeparator(t *testing.T) {
	sharder, err := NewS2Strategy(s2SmartIDShardConfigBad)
	assert.Nil(t, sharder)
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  missing required config: shard_key_separator", err.Error())
}

func TestNewS2StrategyFailureWithOverlappingShards(t *testing.T) {
	sharder, err := NewS2Strategy(s2ShardConfigOverlapping)
	assert.Nil(t, sharder)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "[error]  Overlapping S2 IDs found in backends:")
}

func TestNewS2StrategyFailureWithInvalidS2ID(t *testing.T) {
	sharder, err := NewS2Strategy(s2ShardConfigInvalidS2ID)
	assert.Nil(t, sharder)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "[error] Bad S2 ID found in backends:")
}

func TestS2StrategyS2IDSuccess(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: ",", shardKeyPosition: -1}
	actualS2ID, err := strategy.s2ID("-6.1751,106.865")
	expectedS2ID := uint64(3344474311656055761)
	assert.Nil(t, err)
	assert.Equal(t, expectedS2ID, actualS2ID)

	actualS2ID, err = strategy.s2ID("-6.1751,110.865")
	expectedS2ID = uint64(3346531974082111711)
	assert.Nil(t, err)
	assert.Equal(t, expectedS2ID, actualS2ID)
}

func TestS2StrategyS2IDFailureForInvalidLat(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: ",", shardKeyPosition: -1}
	_, err := strategy.s2ID("-qwerty6.1751,32.865")
	assert.NotNil(t, err)
}

func TestS2StrategyS2IDFailureForInvalidLng(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: ",", shardKeyPosition: -1}
	_, err := strategy.s2ID("-6.1751,qwerty1232.865")
	assert.NotNil(t, err)
}

func TestS2StrategyS2IDFailureForInvalidLatLngObject(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: ",", shardKeyPosition: -1}
	_, err := strategy.s2ID("1116.1751,1232.865")
	assert.NotNil(t, err)
}

func TestS2StrategyS2IDFailureForInvalidAlphanumeric(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: ",", shardKeyPosition: -1}
	_, err := strategy.s2ID("-6.17511232.865")
	assert.NotNil(t, err)

	_, err = strategy.s2ID("abc,1232.865")
	assert.NotNil(t, err)
}

func TestS2StrategyS2IDWithSmartIDSuccess(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: "-", shardKeyPosition: 2}
	actualS2ID, err := strategy.s2ID("v1-foo-3344474403281829888")
	expectedS2ID := uint64(3344474403281829888)
	assert.Nil(t, err)
	assert.Equal(t, expectedS2ID, actualS2ID)

	actualS2ID, err = strategy.s2ID("v1-foo-3346532139293212672")
	expectedS2ID = uint64(3346532139293212672)
	assert.Nil(t, err)
	assert.Equal(t, expectedS2ID, actualS2ID)
}

func TestS2StrategyS2IDFailureForInvalidS2ID(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: "-", shardKeyPosition: 2}
	_, err := strategy.s2ID("v1-foo-bar")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  failed to parse s2id", err.Error())
}

func TestS2StrategyS2IDFailureForInvalidSmartID(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: "-", shardKeyPosition: 2}
	_, err := strategy.s2ID("booyeah")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  failed to get location from smart-id", err.Error())
}

func TestS2StrategyS2IDFailureForInvalidSeparatorConfig(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: "&", shardKeyPosition: 2}
	_, err := strategy.s2ID("v1-foo-3344474403281829888")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  failed to get location from smart-id", err.Error())
}

func TestS2StrategyS2IDFailureForInvalidPositionConfig(t *testing.T) {
	strategy := S2Strategy{shardKeySeparator: "-", shardKeyPosition: 3}
	_, err := strategy.s2ID("v1-foo-3344474403281829888")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  failed to get location from smart-id", err.Error())

	strategy = S2Strategy{shardKeySeparator: "-", shardKeyPosition: 4}
	_, err = strategy.s2ID("v1-foo-3344474403281829888")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  failed to get location from smart-id", err.Error())
}

func TestS2StrategyShardSuccess(t *testing.T) {
	strategy, _ := NewS2Strategy(s2ShardConfig)
	backend, err := strategy.Shard("-6.1751,106.865")
	expectedBackend := "jkt-a"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)

	backend, err = strategy.Shard("-6.1751,110.865")
	expectedBackend = "jkt-b"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)
}

func TestS2StrategySmartIDShardSuccess(t *testing.T) {
	strategy, _ := NewS2Strategy(s2SmartIDShardConfig)
	backend, err := strategy.Shard("v1-foo-3344474403281829888")
	expectedBackend := "jkt-a"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)

	backend, err = strategy.Shard("v1-foo-3346532139293212672")
	expectedBackend = "jkt-b"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)

	backend, err = strategy.Shard("v1-foo-2534")
	expectedBackend = "jkt-c"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)
}

func TestS2StrategyShardSuccessForDefaultBackend(t *testing.T) {
	strategy, _ := NewS2Strategy(s2ShardConfig)
	backend, err := strategy.Shard("-34.1751,106.865")
	expectedBackend := "jkt-c"
	assert.Nil(t, err)
	assert.Equal(t, expectedBackend, backend.Name)
}

func TestS2StrategyShardFailure(t *testing.T) {
	strategy, _ := NewS2Strategy(s2ShardConfig)
	backendForWrongLatLng, err := strategy.Shard("-126.1751,906.865")
	assert.NotNil(t, err)
	assert.Nil(t, backendForWrongLatLng)

	backendForWrongInputNum, err := strategy.Shard("-126.1751865")
	assert.NotNil(t, err)
	assert.Nil(t, backendForWrongInputNum)

	backendForWrongInputAlpha, err := strategy.Shard("abc,xyz")
	assert.NotNil(t, err)
	assert.Nil(t, backendForWrongInputAlpha)

	strategy, _ = NewS2Strategy(s2ShardConfigWithoutDefault)
	noBackendForCorrectInput, err := strategy.Shard("-6.1751,126.865")
	assert.NotNil(t, err)
	assert.Equal(t, "[error]  fail to find backend", err.Error())
	assert.Nil(t, noBackendForCorrectInput)
}

func TestGeneratesCustomError(t *testing.T) {
	ce := CustomError{ExitMessage: "Error for custom error"}
	err := ce.Error()
	assert.Equal(t, err, "[error]  Error for custom error")
}

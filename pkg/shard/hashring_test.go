package shard_test

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/gojektech/weaver/pkg/shard"

	"github.com/stretchr/testify/assert"
)

func TestNewHashringStrategy(t *testing.T) {
	shardConfig := json.RawMessage(`{
        "totalVirtualBackends": 1000,
		"backends": {
			"0-250": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"251-500": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"501-725": { "backend_name": "foobar3", "backend": "http://shard02.local"},
            "726-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
	    }
	}`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, err)
	assert.NotNil(t, hashRingStrategy)
}

func TestShouldFailToCreateWhenWrongBackends(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500,
		"backends": "foo"
	}`)
	expectedErr := fmt.Errorf("json: cannot unmarshal string into Go struct field HashRingStrategyConfig.backends of type map[string]shard.BackendDefinition")
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr.Error(), err.Error())
	assert.Nil(t, hashRingStrategy)
}

func TestShouldFailToCreateWhenNoBackends(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500
	}`)
	expectedErr := fmt.Errorf("No Shard Backends Specified Or Specified Incorrectly")
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, hashRingStrategy)
}

func TestShouldFailToCreateWhenNoBackendURL(t *testing.T) {
	shardConfig := json.RawMessage(`{
        "totalVirtualBackends": 1000,
		"backends": {
			"0-999": { "timeout": 100, "backend_name": "foobar1", "backend": "ht$tp://shard00.local"}
	    }
	}`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Contains(t, err.Error(), "first path segment in URL cannot contain colon")
	assert.Nil(t, hashRingStrategy)
}

func TestShouldFailToCreateWhenTotalVirtualBackendsIsIncorrect(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": "foo",
		"backends": {
			"0-10" : { "backend_name": "foo"}
		}
    }`)
	expectedErr := fmt.Errorf("json: cannot unmarshal string into Go struct field HashRingStrategyConfig.totalVirtualBackends of type int")
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr.Error(), err.Error())
	assert.Nil(t, hashRingStrategy)
}

func TestShouldFailToCreateWhenBackendURLIsMissing(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"backends": {
			"foo" : { "backend_name": "foo"}
		}
    }`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Contains(t, err.Error(), "missing backend url in shard config:")
	assert.Nil(t, hashRingStrategy)
}

func TestShouldDefaultTotalVirtualBackendsWhenValueMissing(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"backends": {
			"0-999" : { "backend_name": "foo", "backend": "http://backend01"}
		}
	}`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, err)
	assert.NotNil(t, hashRingStrategy)
}

func TestShouldFailToCreateWithIncorrectRange(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"backends": {
			"999-0" : { "backend_name": "foo", "backend": "http://blah"}
		}
	}`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, hashRingStrategy)
	assert.Contains(t, err.Error(), "Invalid range key 999-0 for backends")
}

func TestShouldFailToCreateWithIncorrectRangeSpec(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"backends": {
			"999-999-0" : { "backend_name": "foo", "backend": "http://blah"}
		}
	}`)
	hashRingStrategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, hashRingStrategy)
	assert.Contains(t, err.Error(), "Invalid range key format:")
}

func TestShouldFailToCreateHashRingOutOfBounds(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-500": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("Shard is out of bounds Max %d found %d", 499, 500)
	_, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr, err)
}

func TestShouldFailToCreateHashRingOnOverlap(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"249-499": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("Overlap seen in range key %d", 249)
	_, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr, err)
}

func TestShouldFailToCreateHashRingForMissingValuesInTheRangeInMiddle(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500,
		"backends": {
			"0-248": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("Shard is missing coverage for %d", 249)
	_, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr, err)
}

func TestShouldFailToCreateHashRingForMissingValuesInTheRangeAtStart(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 500,
		"backends": {
			"1-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("Shard is missing coverage for %d", 0)
	_, err := shard.NewHashRingStrategy(shardConfig)
	assert.Equal(t, expectedErr, err)
}

func TestShouldCheckBackendConfiguration(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": "foo",
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("json: cannot unmarshal string into Go struct field HashRingStrategyConfig.backends of type shard.BackendDefinition")
	strategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, strategy)
	assert.Equal(t, expectedErr.Error(), err.Error())

}

func TestShouldCheckBackendConfigurationForBackendName(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": {"foo": "bar"},
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	strategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, strategy)
	assert.Contains(t, err.Error(), "missing backend name in shard config:")
}

func TestShouldCheckBackendConfigurationForBackendUrl(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": {"foo": "bar", "backend_name": "foo"},
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	strategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, strategy)
	assert.Contains(t, err.Error(), "missing backend url in shard config:")
}

func TestShouldCheckBackendConfigurationForTimeout(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": {"backend": "http://foo", "backend_name": "foo", "timeout": "abc"},
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)
	expectedErr := fmt.Errorf("json: cannot unmarshal string into Go struct field BackendDefinition.timeout of type float64")
	strategy, err := shard.NewHashRingStrategy(shardConfig)
	assert.Nil(t, strategy)
	assert.Equal(t, expectedErr.Error(), err.Error())
}

func TestShouldShardConsistently(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)

	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	expectedBackend := "foobar2"
	backend, err := strategy.Shard("1")
	assert.Nil(t, err)
	assert.NotNil(t, backend)
	assert.Equal(t, expectedBackend, backend.Name, "Should return foobar2 for key 1")
	backend, err = strategy.Shard("1")
	assert.Equal(t, expectedBackend, backend.Name, "Should return foobar2 for key 1")

}

func TestShouldShardConsistentlyOverALargeRange(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 10,
		"backends": {
			"0-4": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"5-9": { "backend_name": "foobar2", "backend": "http://shard01.local"}
		}
	}`)

	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	shardList := []string{}
	for i := 0; i < 10000; i++ {
		backend, err := strategy.Shard(strconv.Itoa(i))
		assert.Nil(t, err, "Failed to Shard for key %d", i)
		if err != nil {
			t.Log(err)
			return
		}
		shardList = append(shardList, backend.Name)
	}

	for i := 0; i < 10000; i++ {
		backend, err := strategy.Shard(strconv.Itoa(i))
		assert.Nil(t, err, "Failed to Re - Shard for key %d", i)
		if err != nil {
			t.Log(err)
			return
		}
		assert.Equal(t, shardList[i], backend.Name, "Sharded inconsistently for key %d %s -> %s", i, shardList[i], backend.Name)
	}
}

func TestShouldShardConsistentlyAcrossRuns(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 1000,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"500-749": { "backend_name": "foobar3", "backend": "http://shard02.local"},
			"750-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
		}
	}`)

	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	shardList := []string{}
	for i := 0; i < 10000; i++ {
		backend, err := strategy.Shard(strconv.Itoa(i))
		assert.Nil(t, err, "Failed to Shard for key %d", i)
		if err != nil {
			t.Log(err)
			return
		}
		shardList = append(shardList, backend.Name)
	}

	strategy2, _ := shard.NewHashRingStrategy(shardConfig)
	for i := 0; i < 10000; i++ {
		backend, err := strategy2.Shard(strconv.Itoa(i))
		assert.Nil(t, err, "Failed to Re - Shard for key %d", i)
		if err != nil {
			t.Log(err)
			return
		}
		assert.Equal(t, shardList[i], backend.Name, "Sharded inconsistently for key %d %s -> %s", i, shardList[i], backend.Name)
	}
}

func TestShouldShardUniformally(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 1000,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"500-749": { "backend_name": "foobar3", "backend": "http://shard02.local"},
			"750-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
		}
	}`)

	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	shardDistribution := map[string]int{}
	numKeys := 10000000
	for i := 0; i < numKeys; i++ {
		backend, err := strategy.Shard(strconv.Itoa(i))
		assert.Nil(t, err, "Failed to Shard for key %d", i)
		if err != nil {
			t.Log(err)
			return
		}
		shardDistribution[backend.Name] = shardDistribution[backend.Name] + 1
	}
	mean := float64(0)
	for _, v := range shardDistribution {
		mean += float64(v)
	}
	mean = mean / 4

	sd := float64(0)

	for _, v := range shardDistribution {
		sd += math.Pow(float64(v)-mean, 2)
	}

	sd = (math.Sqrt(sd/4) / float64(numKeys)) * float64(100)
	assert.True(t, (sd < float64(2.5)), "Standard Deviation should be less than 2.5% -> %f", sd)

	t.Log("Standard Deviation:", sd)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tHeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	return m
}

func TestShouldShardWithoutLeakingMemory(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 1000,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"500-749": { "backend_name": "foobar3", "backend": "http://shard02.local"},
			"750-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
		}
	}`)

	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	numKeys := 1000
	PrintMemUsage()
	strategy.Shard("1")

	for j := 0; j < 1000; j++ {
		for i := 0; i < numKeys; i++ {
			_, _ = strategy.Shard(strconv.Itoa(i))
		}
	}

	PrintMemUsage()
}

func TestToMeasureTimeForSharding(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 1000,
		"backends": {
			"0-249": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"250-499": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"500-749": { "backend_name": "foobar3", "backend": "http://shard02.local"},
			"750-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
		}
	}`)
	strategy, _ := shard.NewHashRingStrategy(shardConfig)
	numKeys := 1000
	start := time.Now()
	for j := 0; j < 1000; j++ {
		for i := 0; i < numKeys; i++ {
			_, _ = strategy.Shard(strconv.Itoa(i))
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Elapsed Time %v", elapsed)
}

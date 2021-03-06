package accounting

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/artpar/rclone/fstest/testy"
	"github.com/stretchr/testify/assert"
)

func TestStatsGroupOperations(t *testing.T) {

	t.Run("empty group returns nil", func(t *testing.T) {
		t.Parallel()
		sg := newStatsGroups()
		sg.get("invalid-group")
	})

	t.Run("set assigns stats to group", func(t *testing.T) {
		t.Parallel()
		stats := NewStats()
		sg := newStatsGroups()
		sg.set("test", stats)
		sg.set("test1", stats)
		if len(sg.m) != len(sg.names()) || len(sg.m) != 2 {
			t.Fatalf("Expected two stats got %d, %d", len(sg.m), len(sg.order))
		}
	})

	t.Run("get returns correct group", func(t *testing.T) {
		t.Parallel()
		stats := NewStats()
		sg := newStatsGroups()
		sg.set("test", stats)
		sg.set("test1", stats)
		got := sg.get("test")
		if got != stats {
			t.Fatal("get returns incorrect stats")
		}
	})

	t.Run("sum returns correct values", func(t *testing.T) {
		t.Parallel()
		stats1 := NewStats()
		stats1.bytes = 5
		stats1.errors = 6
		stats1.oldDuration = time.Second
		stats1.oldTimeRanges = []timeRange{{time.Now(), time.Now().Add(time.Second)}}
		stats2 := NewStats()
		stats2.bytes = 10
		stats2.errors = 12
		stats2.oldDuration = 2 * time.Second
		stats2.oldTimeRanges = []timeRange{{time.Now(), time.Now().Add(2 * time.Second)}}
		sg := newStatsGroups()
		sg.set("test1", stats1)
		sg.set("test2", stats2)
		sum := sg.sum()
		assert.Equal(t, stats1.bytes+stats2.bytes, sum.bytes)
		assert.Equal(t, stats1.errors+stats2.errors, sum.errors)
		assert.Equal(t, stats1.oldDuration+stats2.oldDuration, sum.oldDuration)
		// dict can iterate in either order
		a := timeRanges{stats1.oldTimeRanges[0], stats2.oldTimeRanges[0]}
		b := timeRanges{stats2.oldTimeRanges[0], stats1.oldTimeRanges[0]}
		if !assert.ObjectsAreEqual(a, sum.oldTimeRanges) {
			assert.Equal(t, b, sum.oldTimeRanges)
		}
	})

	t.Run("delete removes stats", func(t *testing.T) {
		t.Parallel()
		stats := NewStats()
		sg := newStatsGroups()
		sg.set("test", stats)
		sg.set("test1", stats)
		sg.delete("test1")
		if sg.get("test1") != nil {
			t.Fatal("stats not deleted")
		}
		if len(sg.m) != len(sg.names()) || len(sg.m) != 1 {
			t.Fatalf("Expected two stats got %d, %d", len(sg.m), len(sg.order))
		}
	})

	t.Run("memory is reclaimed", func(t *testing.T) {
		testy.SkipUnreliable(t)
		var (
			count      = 1000
			start, end runtime.MemStats
			sg         = newStatsGroups()
		)

		runtime.GC()
		runtime.ReadMemStats(&start)

		for i := 0; i < count; i++ {
			sg.set(fmt.Sprintf("test-%d", i), NewStats())
		}

		for i := 0; i < count; i++ {
			sg.delete(fmt.Sprintf("test-%d", i))
		}

		runtime.GC()
		runtime.ReadMemStats(&end)

		t.Log(fmt.Sprintf("%+v\n%+v", start, end))
		diff := percentDiff(start.HeapObjects, end.HeapObjects)
		if diff > 1 || diff < 0 {
			t.Errorf("HeapObjects = %d, expected %d", end.HeapObjects, start.HeapObjects)
		}
	})
}

func percentDiff(start, end uint64) uint64 {
	return (start - end) * 100 / start
}

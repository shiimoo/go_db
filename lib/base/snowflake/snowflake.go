package snowflake

import (
	"sync"
	"time"
)

const (
	epoch             int64 = 1704038400000 // 设置起始时间戳，例如2024-01-01 00:00:00 UTC
	timeBitLength     uint8 = 41            // 时间戳占用的位数
	workerIDBitLength uint8 = 5             // 工作机器ID占用的位数
	sequenceBitLength uint8 = 12            // 序列号占用的位数

	maxWorkerID   int64 = -1 ^ (-1 << workerIDBitLength) // 工作机器ID的最大值
	maxSequence   int64 = -1 ^ (-1 << sequenceBitLength) // 序列号的最大值
	timeShift     uint8 = workerIDBitLength + sequenceBitLength
	workerIDShift uint8 = sequenceBitLength
	twepoch       int64 = epoch
)

type snowflake struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64
	sequence  int64
}

func checkWorkID(workerID int64) {
	if workerID < 0 || workerID > maxWorkerID {
		panic("worker ID out of range")
	}
}
func newSnowflake(workerID int64) *snowflake {
	checkWorkID(workerID)
	return &snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}
}

func (s *snowflake) Gen() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli() - twepoch
	if now < s.timestamp {
		panic("clock is moving backwards")
	}

	if now == s.timestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixMilli() - twepoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	id := ((now << timeShift) | (s.workerID << workerIDShift) | (s.sequence))
	return id
}

var unique *snowflake

func getUnique() *snowflake {
	if unique == nil {
		unique = newSnowflake(0)
	}
	return unique
}

func SetWorkerID(workerID int64) {
	checkWorkID(workerID)
	getUnique().workerID = workerID
}

func Gen() int64 {
	return getUnique().Gen()
}

func GenUint() uint {
	return uint(Gen())
}

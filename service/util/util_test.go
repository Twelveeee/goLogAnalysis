package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIpToInt(t *testing.T) {
	t.Run("test ip to int", func(t *testing.T) {
		assert.Equal(t, uint32(0), IpToInt("0.0.0.0"))
		assert.Equal(t, uint32(1000), IpToInt("0.0.3.232"))
		assert.Equal(t, uint32(4294967295), IpToInt("255.255.255.255"))
	})
}

func TestConverToTimestamp(t *testing.T) {
	t.Run("test sting To Timestamp(int)", func(t *testing.T) {
		assert.Equal(t, int64(1699632868), ConverToTimestamp("11/Nov/2023:00:14:28 +0800"))
	})
}

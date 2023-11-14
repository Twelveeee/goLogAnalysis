package entity

import (
	"testing"
)

func TestOssLogCreate(t *testing.T) {
	testInitDB()
	m1 := OssLog{}
	m1.Create()
}

package net

import (
	"github.com/sta-golang/go-lib-utils/log"
	"testing"
)

func TestLocalIP(t *testing.T) {
	log.FrameworkLogger.Debug(LocalIP())
}

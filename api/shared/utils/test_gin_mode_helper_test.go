package utils

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var testSetModeOnce sync.Once

func testSetGinMode() {
	testSetModeOnce.Do(func() {
		gin.SetMode(gin.TestMode)
	})
}


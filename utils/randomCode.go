package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateResetCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%04d", r.Intn(10000))
}

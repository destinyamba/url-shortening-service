package url_shortening

import (
	"fmt"
	"math/rand"
	"time"
)

const ch = "abcdefghijklmnopqrstuvwxyz0123456789"

func CodeGenerator() string {
	rand.NewSource(time.Now().UnixNano())
	code := make([]byte, 6)
	for i := 0; i < 6; i++ {
		code[i] = ch[rand.Intn(len(ch))]
	}
	fmt.Println(code)
	return string(code)
}

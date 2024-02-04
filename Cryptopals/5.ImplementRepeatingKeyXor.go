package main

import (
	"encoding/hex"
	"fmt"
)

var (
	text []byte = []byte(`Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`)

	enc []byte = []byte(`ICE`)
)

func ImplementRepeatingKeyXor(check string) {
	ans := make([]byte, len(text))

	tl := len(text)
	el := len(enc)
	ll := tl/el + 1

	for i := 0; i < ll; i++ {
		for r := 0; r < el; r++ {
			if i*el+r < tl {
				ans[i*el+r] = text[i*el+r] ^ enc[r]
				// fmt.Printf("%d %c %c %c\n", i*el+r, text[i*el+r], ans[i*el+r], enc[r])
			}
		}
	}

	fmt.Println(hex.EncodeToString(ans))

	if check == hex.EncodeToString(ans) {
		fmt.Println("ok")
	}
}

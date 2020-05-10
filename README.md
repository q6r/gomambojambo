# gomambojambo

Golang AST Obfuscation. The idea is to give it a package path, and it will
obfuscate the AST and write it back. This is just a toy project to experiment
with Golang AST.

Currently it does the following :

- Randomization of function names, and function calls
- For loops converted to goto with tags
- Strings obfuscated/encrypted using AES

Example given the following source code :

```
package main

import (
	"fmt"
)

func numberList(a int) []int {
	m := []int{}
	for i := 0; i < a; i++ {
		for j := 0; j < i; j++ {
			m = append(m, i + j)
		}
	}
	return m
}

func main() {
	v := numberList(10)
	sum := 0
	for i := 0; i < len(v); i++ {
		sum += v[i]
	}
	fmt.Printf("sum = %#v\n", sum)
}

```

Will be obfuscated to

```
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func ijolxvfogwskcjuvdgffytwzmukemrsm(a int) []int {
	m := []int{}
	{
		goto LOOP_INIT_mmohmb
	LOOP_INIT_mmohmb:
		;
		i := 0
		goto LOOP_COND_excfyj
	LOOP_COND_excfyj:
		if i < a {
			goto LOOP_BODY_gxiifn
		} else {
			goto LOOP_END_nxxdqv
		}
	LOOP_BODY_gxiifn:
		{
			{
				goto LOOP_INIT_nnnjjc
			LOOP_INIT_nnnjjc:
				;
				j := 0
				goto LOOP_COND_huckdb
			LOOP_COND_huckdb:
				if j < i {
					goto LOOP_BODY_cssbty
				} else {
					goto LOOP_END_rcwqup
				}
			LOOP_BODY_cssbty:
				{
					m = append(m, i+j)
					j++
					goto LOOP_COND_huckdb

				}
			LOOP_END_rcwqup:
				{
				}
			}
			i++
			goto LOOP_COND_excfyj

		}
	LOOP_END_nxxdqv:
		{
		}
	}
	return m
}

func Lenirdftvvrgwesqohuiucnfhyaehjkj(s string) string {
	key, _ := hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
	ciphertext, _ := hex.DecodeString(s)
	nonce, _ := hex.DecodeString("010101010101010101010101")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}

func main() {
	v := ijolxvfogwskcjuvdgffytwzmukemrsm(10)
	sum := 0
	{
		goto LOOP_INIT_gmayxh
	LOOP_INIT_gmayxh:
		;
		i := 0
		goto LOOP_COND_bzkbgj
	LOOP_COND_bzkbgj:
		if i < len(v) {
			goto LOOP_BODY_mezrxh
		} else {
			goto LOOP_END_oxyyzr
		}
	LOOP_BODY_mezrxh:
		{
			sum += v[i]
			i++
			goto LOOP_COND_bzkbgj
		}
	LOOP_END_oxyyzr:
		{
		}
	}
	fmt.Printf(Lenirdftvvrgwesqohuiucnfhyaehjkj("c10cd7ba54d720d09cf961814f1b2581161731150c24cf37a7b0b3"), sum)
}
```

run with `./gomambojambo -calls -loops -strings -writechanges -srcpath mycode/` and `-h` for help.

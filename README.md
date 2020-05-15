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

func feezjurstezwlipfhqotmvrfljwgclpj(a int) []int {
	m := []int{}
	{
		goto LOOP_INIT_fzjcph
	LOOP_INIT_fzjcph:
		;
		i := 0
		goto LOOP_COND_nlvrpf
	LOOP_COND_nlvrpf:
		if i < a {
			goto LOOP_BODY_mzoord
		} else {
			goto LOOP_END_kndsac
		}
	LOOP_BODY_mzoord:
		{
			{
				goto LOOP_INIT_icketl
			LOOP_INIT_icketl:
				;
				j := 0
				goto LOOP_COND_gzqnmc
			LOOP_COND_gzqnmc:
				if j < i {
					goto LOOP_BODY_udzevk
				} else {
					goto LOOP_END_pmovmj
				}
			LOOP_BODY_udzevk:
				{
					m = append(m, i+j)
					j++
					goto LOOP_COND_gzqnmc

				}
			LOOP_END_pmovmj:
				{
				}
			}
			i++
			goto LOOP_COND_nlvrpf

		}
	LOOP_END_kndsac:
		{
		}
	}
	return m
}

func main() {
	v := feezjurstezwlipfhqotmvrfljwgclpj(10)
	sum := 0
	{
		goto LOOP_INIT_eiubku
	LOOP_INIT_eiubku:
		;
		i := 0
		goto LOOP_COND_htnuhc
	LOOP_COND_htnuhc:
		if i < len(v) {
			goto LOOP_BODY_ptqcvi
		} else {
			goto LOOP_END_igzdqm
		}
	LOOP_BODY_ptqcvi:
		{
			sum += v[i]
			i++
			goto LOOP_COND_htnuhc

		}
	LOOP_END_igzdqm:
		{
		}
	}
	fmt.Printf((func(s string) string {
		k, _ := hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
		ct, _ := hex.DecodeString(s)
		n, _ := hex.DecodeString("010101010101010101010101")
		b, _ := aes.NewCipher(k)
		g, _ := cipher.NewGCM(b)
		pt, _ := g.Open(nil, n, ct, nil)
		return string(pt)
	})("c10cd7ba54d720d09caf9b82f4bbc522f575d548fe59730764cb"), sum)
}
```

run with `./gomambojambo -calls -loops -strings -writechanges -srcpath mycode/` and `-h` for help.

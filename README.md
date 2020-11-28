# gomambojambo

AST-based obfuscation. The idea is to give it a package path, and it will obfuscate it.

- Randomization of function names, and function calls
- For loops converted to goto with tags
- Strings obfuscated/encrypted using AES
- Adding some deadcode to functions

Given the following source code :

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

Will be obfuscated to :


```

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func xhrandgrmayhrxnwbvksrglqqgzxtbwc(a int) []int {
	(func() {
		zXXX := int64(6)
		sXXX := float64(4)
		{
			goto LOOP_INIT_fskwhl
		LOOP_INIT_fskwhl:
			;
			iXXX := 8
			goto LOOP_COND_ybznsy
		LOOP_COND_ybznsy:
			if iXXX < 15 {
				goto LOOP_BODY_oecdzt
			} else {
				goto LOOP_END_jcqmbc
			}
		LOOP_BODY_oecdzt:
			{
				{
					goto LOOP_INIT_ydlebm
				LOOP_INIT_ydlebm:
					;
					jXXX := iXXX
					goto LOOP_COND_zpxcxd
				LOOP_COND_zpxcxd:
					if jXXX <
						15 {
						goto LOOP_BODY_ibbxcs
					} else {
						goto LOOP_END_dhpwmq
					}
				LOOP_BODY_ibbxcs:
					{
						{
							goto LOOP_INIT_gqlrlm
						LOOP_INIT_gqlrlm:
							;

							zXXX := jXXX
							goto LOOP_COND_zbipfd
						LOOP_COND_zbipfd:
							if zXXX < 15 {
								goto LOOP_BODY_uwodsb
							} else {
								goto LOOP_END_kwjqls
							}
						LOOP_BODY_uwodsb:
							{
								sXXX = (float64(iXXX+jXXX) *
									float64(zXXX)) /
									float64(iXXX)
								zXXX++
								goto LOOP_COND_zbipfd

							}
						LOOP_END_kwjqls:
							{
							}
						}
						jXXX++
						goto LOOP_COND_zpxcxd

					}
				LOOP_END_dhpwmq:
					{
					}
				}
				iXXX++
				goto LOOP_COND_ybznsy

			}
		LOOP_END_jcqmbc:
			{
			}
		}
		if sXXX == float64(zXXX) {
		}
	})()

	m := []int{}
	{
		goto LOOP_INIT_rtnfaj
	LOOP_INIT_rtnfaj:
		;
		i := 0
		goto LOOP_COND_rdmjei
	LOOP_COND_rdmjei:
		if i < a {
			goto LOOP_BODY_ggnsph
		} else {
			goto LOOP_END_hokgpq
		}
	LOOP_BODY_ggnsph:
		{
			{
				goto LOOP_INIT_fkrzea
			LOOP_INIT_fkrzea:
				;
				j := 0
				goto LOOP_COND_wlwypq
			LOOP_COND_wlwypq:
				if j < i {
					goto LOOP_BODY_pmzmun
				} else {
					goto LOOP_END_vbqtik
				}
			LOOP_BODY_pmzmun:
				{
					m = append(m, i+j)
					j++
					goto LOOP_COND_wlwypq

				}
			LOOP_END_vbqtik:
				{
				}
			}
			i++
			goto LOOP_COND_rdmjei

		}
	LOOP_END_hokgpq:
		{
		}
	}
	return m
}
func Yrqeewbyhgvzpcnktemimipehaipmukk(
	s string) string {
	key, _ :=
		hex.DecodeString("0101010101010101010101010101010101010101010101010101010101010101")
	ciphertext,

		_ := hex.
		DecodeString(s)
	nonce, _ := hex.DecodeString("010101010101010101010101")
	block, err :=
		aes.
			NewCipher(key)
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
	(func() {
		zXXX := int64(6)
		sXXX := float64(4)
		{
			goto LOOP_INIT_fskwhl
		LOOP_INIT_fskwhl:
			;
			iXXX := 8
			goto LOOP_COND_ybznsy
		LOOP_COND_ybznsy:
			if iXXX < 15 {
				goto LOOP_BODY_oecdzt
			} else {
				goto LOOP_END_jcqmbc
			}
		LOOP_BODY_oecdzt:
			{
				{
					goto LOOP_INIT_ydlebm
				LOOP_INIT_ydlebm:
					;
					jXXX := iXXX
					goto LOOP_COND_zpxcxd
				LOOP_COND_zpxcxd:
					if jXXX <
						15 {
						goto LOOP_BODY_ibbxcs
					} else {
						goto LOOP_END_dhpwmq
					}
				LOOP_BODY_ibbxcs:
					{
						{
							goto LOOP_INIT_gqlrlm
						LOOP_INIT_gqlrlm:
							;

							zXXX := jXXX
							goto LOOP_COND_zbipfd
						LOOP_COND_zbipfd:
							if zXXX < 15 {
								goto LOOP_BODY_uwodsb
							} else {
								goto LOOP_END_kwjqls
							}
						LOOP_BODY_uwodsb:
							{
								sXXX = (float64(iXXX+jXXX) *
									float64(zXXX)) /
									float64(iXXX)
								zXXX++
								goto LOOP_COND_zbipfd

							}
						LOOP_END_kwjqls:
							{
							}
						}
						jXXX++
						goto LOOP_COND_zpxcxd

					}
				LOOP_END_dhpwmq:
					{
					}
				}
				iXXX++
				goto LOOP_COND_ybznsy

			}
		LOOP_END_jcqmbc:
			{
			}
		}
		if sXXX == float64(zXXX) {
		}
	})()

	v := xhrandgrmayhrxnwbvksrglqqgzxtbwc(10)
	sum := 0
	{
		goto LOOP_INIT_dbcywi
	LOOP_INIT_dbcywi:
		;
		i := 0
		goto LOOP_COND_vxrvig
	LOOP_COND_vxrvig:
		if i < len(v) {
			goto LOOP_BODY_wmptnj
		} else {
			goto LOOP_END_bscqrq
		}
	LOOP_BODY_wmptnj:
		{
			sum += v[i]
			i++
			goto LOOP_COND_vxrvig

		}
	LOOP_END_bscqrq:
		{
		}
	}
	fmt.Printf(Yrqeewbyhgvzpcnktemimipehaipmukk("c10cd7ba54d720d09caf9b82f4bbc522f575d548fe59730764cb"), sum)
}
```

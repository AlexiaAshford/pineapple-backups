package Encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"github.com/gookit/color"
	"github.com/tidwall/sjson"
	"math/big"
	mathRand "math/rand"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type Encrypt struct {
}

// CreateSecretKey 创建随机16位16进制随机字符串
func CreateSecretKey() []byte {
	rnd := mathRand.New(mathRand.NewSource(time.Now().Unix()))
	temp := strconv.FormatUint(rnd.Uint64(), 16)

	// 以下是强转换，如果不想使用unsafe包的话可以注释这部分然后使用被注释的return
	sh := (*reflect.StringHeader)(unsafe.Pointer(&temp))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
	//return []byte(temp)
}

// RSAEncrypt 将传入参数进行RSA加密
func RSAEncrypt(text []byte) string {
	n := new(big.Int)
	n, _ = n.SetString("00C1E3934D1614465B33053E7F48EE4EC87B14B95EF88947713D25EECBFF7E74C7977D02DC1D9451F79DD5D1C10C29ACB6A9B4D6FB7D0A0279B6719E1772565F09AF627715919221AEF91899CAE08C0D686D748B20A3603BE2318CA6BC2B59706592A9219D0BF05C9F65023A21D2330807252AE0066D59CEEFA5F2748EA80BAB81", 16) // 这边网上的base是10，不知道会不会有影响
	PublicKey := rsa.PublicKey{
		N: n,
		E: 65537,
	}
	byteText, err := rsa.EncryptPKCS1v15(rand.Reader, &PublicKey, text)
	if err != nil {
		color.Errorln("RSAEncrypt Error: ", err)
		panic(err)
	}
	return hex.EncodeToString(byteText)
}

// AESEncrypt 对传入参数进行AES加密
func AESEncrypt(text []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		color.Errorln("NewCipher Error: ", err)
		panic(err)
	}
	blockSize := block.BlockSize()

	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	origData := append(text, padText...)
	iv := []byte("0000000000000000")
	blocMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	blocMode.CryptBlocks(encrypted, origData)

	return encrypted
}

// BytesToString 将[]byte经过解密转换为string
func BytesToString(array []byte) string {
	s := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789()"

	var o func(int) string
	var e func(int, int) int

	o = func(t int) string {
		if t < 0 || t > len(s) {
			return "."
		} else {
			return string(s[t])
		}
	}

	e = func(t int, e int) int {
		n := 0
		for r := 24; r >= 0; r-- {
			if 1 == (e >> r & 1) {
				n = (n << 1) + (t >> r & 1)
			}
		}
		return n
	}

	n := ""
	r := ""
	for a := 0; a < len(array); a += 3 {
		if a+2 < len(array) {
			u := (int(array[a]) << 16) + (int(array[a+1]) << 8) + int(array[a+2])
			n += o(e(u, 7274496)) + o(e(u, 9483264)) + o(e(u, 19220)) + o(e(u, 235))
		} else {
			c := len(array) % 3
			if c == 2 {
				u := (int(array[a]) << 16) + (int(array[a+1]) << 8)
				n += o(e(u, 7274496)) + o(e(u, 9483264)) + o(e(u, 19220))
				r = "."
			} else if c == 1 {
				u := int(array[a]) << 16
				n += o(e(u, 7274496)) + o(e(u, 9483264))
				r = ".."
			}
		}
	}
	return n + r
}

// GetRequestW 获取网络请求的w
func GetRequestW(gt, challenge, aa, passtime string, c int64) string {
	hash := Hash(gt + challenge + passtime)
	text := `{"lang":"zh-cn","userresponse":"","passtime":"","imgload":"","aa":"","ep":"","rp":""}`
	t1, _ := sjson.Set(text, "userresponse", CalUserResponse(c, challenge))
	t1, _ = sjson.Set(t1, "passtime", passtime)
	t1, _ = sjson.Set(t1, "imgload", RandInt64(100, 800))
	t1, _ = sjson.Set(t1, "aa", aa)
	t1, _ = sjson.Set(t1, "ep", GetEP(gt, challenge))
	t1, _ = sjson.Set(t1, "rp", hash)
	secKey := CreateSecretKey()
	encSecKey := RSAEncrypt(secKey)
	encTextByte := AESEncrypt([]byte(t1), secKey)
	encText := BytesToString(encTextByte)
	return encText + encSecKey
}

func CalUserResponse(c int64, challenge string) string {
	mathRand.Seed(time.Now().Unix())
	n := challenge[32:]
	var r []int32
	for _, v := range n {
		if v < 57 {
			r = append(r, v-87)
		} else {
			r = append(r, v-48)
		}
	}
	m := 36*r[0] + r[1]
	a := c + int64(m)
	var e map[string]bool
	var u [][]string
	challenge = challenge[:32]
	lenC := len(challenge)
	j := 0
	for i := 0; i < lenC; i++ {
		for k := range e {
			if string(challenge[i]) == k {
			} else {
				e[string(challenge[i])] = true
				u[j] = append(u[j], string(challenge[i]))
				j += 1
				if j == 5 {
					j = 0
				}
			}
		}
	}
	h := a
	d := 4
	p := ""
	g := []int64{1, 2, 5, 10, 50}
	for h > 0 {
		if h-g[d] >= 0 {
			f := int(mathRand.Float64() * float64(len(u[d])))
			p += u[d][f]
			h -= g[d]
		} else {
			u = append(u[:d-1], u[d])
			g = append(g[:d-1], g[d])
			d -= 1
		}
	}
	return p
}

func GetEP(gt, challenge string) map[string]interface{} {
	hash := Hash(gt + challenge)
	a := time.Now().UnixNano() / 1e6
	f := a + RandInt64(2, 8)
	b := a + RandInt64(50, 80)
	l := a + RandInt64(3, 9)
	m := l + RandInt64(30, 50)
	n := m + RandInt64(1, 5)
	o := n + RandInt64(10, 50)
	p := o + RandInt64(70, 90)
	r := p + RandInt64(10, 100)
	s := r + RandInt64(1, 2)
	tm := map[string]int64{
		"a": a, "b": b, "c": b, "d": 0, "e": 0, "f": f, "g": f, "h": f, "i": f,
		"j": f, "k": 0, "l": l, "m": m, "n": n, "o": o, "p": p, "q": p, "r": a,
		"s": s, "t": s, "u": s,
	}
	return map[string]interface{}{
		"v":  "7.6.0",
		"f":  hash,
		"me": true,
		"te": false,
		"tm": tm,
	}
}

func CalT(t [][]int64) [][]int64 {
	var i [][]int64
	var o, e, n, r int64 = 0, 0, 0, 0
	for j := 0; j < len(t); j++ {
		e = t[j+1][0] - t[j][0]
		n = t[j+1][1] - t[j][1]
		r = t[j+1][2] - t[j][2]
		if e != 0 || n != 0 || r != 0 {
			if e == 0 && n == 0 {
				o += r
			} else {
				i = append(i, []int64{e, n, r + o})
				o = 0
			}
		}
	}
	if o != 0 {
		i = append(i, []int64{e, n, r})
	}
	return i
}

func CalE(t []int64) string {
	e := [][2]int64{{1, 0}, {2, 0}, {1, -1}, {1, 1}, {0, 1}, {0, -1}, {3, 0}, {2, -1}, {2, 1}}
	s := "stuvwxyz~"
	for i := 0; i < 9; i++ {
		if t[0] == e[i][0] && t[1] == e[i][1] {
			return string(s[i])
		}
	}
	return ""
}

func CalF(t [][]int64) string {
	var i, r, o = "", "", ""
	trace := CalT(t)
	for _, v := range trace {
		e := CalE(v)
		if e != "" {
			i += e
		} else {
			r += FunN(v[0])
			i += FunN(v[1])
		}
		o += FunN(v[2])
	}
	return r + "!!" + i + "!!" + o
}

// CalAA 与TT是相同算法
func CalAA(trace [][]int64, e []int64, n string) string {
	f := CalF(trace)
	var s = e[0]
	var a = e[2]
	var u = e[4]
	for j := 0; j < 4; j++ {
		r := n[j : j+2]
		c, _ := strconv.ParseInt(r, 16, 32)
		cInt32 := int32(c)
		m := string(cInt32)
		l := (s*c*c + a*c + u) % int64(len(f))
		f = f[:l] + m + f[l:]
	}
	return f
}

func FunN(t int64) string {
	e := "()*,-./0123456789:?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqr"
	var n int64 = 65
	r := ""
	var i int64 = 0
	if t < 0 {
		i = -t
	} else {
		i = t
	}
	o := i / n

	if o >= n {
		o = n - 1
	}
	if o != 0 {
		r = string(e[o])
	}
	s := ""
	if t < 0 {
		s += "!"
	}
	if r != "" {
		s += "$"
	}
	return s + r + string(e[i%n])
}

// Hash 对传入字符串进行md5加密
func Hash(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// RandInt64 返回输入区间的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return mathRand.Int63n(max-min) + min
}

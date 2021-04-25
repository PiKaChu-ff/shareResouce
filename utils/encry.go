package utils

import (
	"bytes"
	//"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"time"
	//"time"
)

const version = byte(0x00)   //定义版本号，一个字节
const addressChecksumLen = 4 //定义checksum长度为四个字节

var Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func GetRandomString() string {
	l := 20
	bytes := []byte(Alphabet)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	if inst == 0 {
		inst = 0
	}
	return inst
}

// 字节数组转 Base58,加密
func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, Alphabet[mod.Int64()])
	}

	ReverseBytes(result)
	for b := range input {
		if b == 0x00 {
			result = append([]byte{Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result
}

// Base58转字节数组，解密
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0

	for b := range input {
		if b == 0x00 {
			zeroBytes++
		}
	}

	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

// 字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func IsValidForAdress(adress []byte) bool {
	//将地址进行base58反编码，生成的其实是version+Pub Key hash+ checksum这25个字节
	version_public_checksumBytes := Base58Decode(adress)

	//[25-4:],就是21个字节往后的数（22,23,24,25一共4个字节）
	checkSumBytes := version_public_checksumBytes[len(version_public_checksumBytes)-addressChecksumLen:]
	//[:25-4],就是前21个字节（1～21,一共21个字节）
	version_ripemd160 := version_public_checksumBytes[:len(version_public_checksumBytes)-addressChecksumLen]
	//取version+public+checksum的字节数组的前21个字节进行两次256哈希运算，取结果值的前4个字节
	checkBytes := CheckSum(version_ripemd160)
	//将checksum比较，如果一致则说明地址有效，返回true
	if bytes.Compare(checkSumBytes, checkBytes) == 0 {
		return true
	}

	return false
}

func GetAddress(PublicKey string) string {

	//调用Ripemd160Hash返回160位的Pub Key hash
	ripemd160Hash := Ripemd160Hash([]byte(PublicKey))

	//将version+Pub Key hash
	version_ripemd160Hash := append([]byte{version}, ripemd160Hash...)

	//调用CheckSum方法返回前四个字节的checksum
	checkSumBytes := CheckSum(version_ripemd160Hash)

	//将version+Pub Key hash+ checksum生成25个字节
	bytes := append(version_ripemd160Hash, checkSumBytes...)

	//将这25个字节进行base58编码并返回
	return string(Base58Encode(bytes))
}

//取前4个字节
func CheckSum(payload []byte) []byte {
	//这里传入的payload其实是version+Pub Key hash，对其进行两次256运算
	hash1 := sha256.Sum256(payload)

	hash2 := sha256.Sum256(hash1[:])

	return hash2[:addressChecksumLen] //返回前四个字节，为CheckSum值
}

func Ripemd160Hash(publicKey []byte) []byte {

	//将传入的公钥进行256运算，返回256位hash值
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//将上面的256位hash值进行160运算，返回160位的hash值
	ripemd160 := New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil) //返回Pub Key hash
}

func ExtAbstract(s string) string {
	t := sha1.New()
	io.WriteString(t, s)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func AddByChar(paras []string) (retStr string) {
	retStr = ""

	maxLen := 0
	for _, v := range paras {
		if v == "" {
			return
		}
		length := strings.Count(v, "") - 1
		if length > maxLen {
			maxLen = length
		}
	}
	var tmp byte
	tmp = 0
	for i := 0; i < maxLen; i++ {
		for _, v := range paras {
			if i < (strings.Count(v, "") - 1) {
				tmp += v[i]
			}
		}
		retStr += string(tmp)
	}

	return
}

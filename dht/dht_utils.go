package dht

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"math/big"
)

func Distance(a, b []byte, bits int) *big.Int {
	//Be sure that a<b or distance will be infinite
	//distance 0,1 = 1
	//distance 1,0 = 255
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	var a_int, b_int big.Int
	(&a_int).SetBytes(a)
	(&b_int).SetBytes(b)

	var dist big.Int
	(&dist).Sub(&b_int, &a_int)

	(&dist).Mod(&dist, &ring)
	return &dist
}

func Between(Sid1, Sid2, Skey string) bool {
	//convert string to []byte
	id1 := []byte(Sid1)
	id2 := []byte(Sid2)
	key := []byte(Skey)
	// 0 if a==b, -1 if a < b, and +1 if a > b
	if bytes.Compare(key, id1) == 0 { // key == id1
		return true
	}

	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id1) == 1 && bytes.Compare(key, id2) == -1 { // id1 < key && key < id2
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // id1 > key || key < id2
			return true
		} else {
			return false
		}
	}
}

// (n + 2^(k-1)) mod (2^m)
func CalcFinger(n []byte, k int, m int) (string, []byte) {
	fmt.Println("calulcating result = (n+2^(k-1)) mod (2^m)")

	// convert the n to a bigint
	nBigInt := big.Int{}
	nBigInt.SetBytes(n)

	fmt.Printf("n            %s\n", nBigInt.String())

	fmt.Printf("k            %d\n", k)

	fmt.Printf("m            %d\n", m)

	// get the right addend, i.e. 2^(k-1)
	two := big.NewInt(2)
	addend := big.Int{}
	addend.Exp(two, big.NewInt(int64(k-1)), nil)

	fmt.Printf("2^(k-1)      %s\n", addend.String())

	// calculate sum
	sum := big.Int{}
	sum.Add(&nBigInt, &addend)

	fmt.Printf("(n+2^(k-1))  %s\n", sum.String())

	// calculate 2^m
	ceil := big.Int{}
	ceil.Exp(two, big.NewInt(int64(m)), nil)

	fmt.Printf("2^m          %s\n", ceil.String())

	// apply the mod
	result := big.Int{}
	result.Mod(&sum, &ceil)

	fmt.Printf("finger       %s\n", result.String())

	resultBytes := result.Bytes()
	resultHex := fmt.Sprintf("%x", resultBytes)

	fmt.Printf("finger (hex) %s\n", resultHex)

	return resultHex, resultBytes
}

func GenerateNodeId() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return Sha1hash(u.String())
}

func Sha1hash(str string) string {
	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(str))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

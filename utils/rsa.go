package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

/**
 * 功能：获取RSA公钥长度
 * 参数：public
 * 返回：成功则返回 RSA 公钥长度，失败返回 error 错误信息
 */
func RsaEncrypt(pubKey []byte, origData []byte) ([]byte, error) {
	if pubKey == nil {
		return nil, errors.New("public key empty")
	}
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub := pubInterface.(*rsa.PublicKey)
	//fmt.Println(pub)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

/*
   获取RSA私钥长度
   PriKey
   成功返回 RSA 私钥长度，失败返回error
*/
func RsaDecrypt(priKey []byte, ciphertext []byte) ([]byte, error) {
	if priKey == nil {
		return nil, errors.New("private key empty")
	}
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

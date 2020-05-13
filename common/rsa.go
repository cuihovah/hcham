package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type RSA struct {
	Publickey  []byte
	Privatekey []byte
}

func RSAKeyGen(bits int) error {
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	privatefile, err := os.Create("privateKey.pem")
	defer privatefile.Close()
	err = pem.Encode(privatefile, block)
	if err != nil {
		return err
	}
	publickey := &privatekey.PublicKey
	derpkix, err := x509.MarshalPKIXPublicKey(publickey)
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derpkix,
	}
	if err != nil {
		return err
	}
	publickfile, err := os.Create("publicKey.pem")
	defer publickfile.Close()
	err = pem.Encode(publickfile, block)
	if err != nil {
		return err
	}
	return nil
}

func (r *RSA) RSAEncrypt(orgidata []byte) ([]byte, error) {
	block, _ := pem.Decode(r.Publickey)
	if block == nil {
		return nil, errors.New("public key is bad")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, orgidata) //加密
}

func (r *RSA) RSADecrypt(cipertext []byte) ([]byte, error) {
	block, _ := pem.Decode(r.Privatekey)
	if block == nil {
		return nil, errors.New("public key is bad")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipertext)
}

func NewRSA() *RSA {
	r := &RSA{}
	RSAKeyGen(1024)
	r.Publickey = FileLoad("publicKey.pem")
	r.Privatekey = FileLoad("privateKey.pem")
	return r
}

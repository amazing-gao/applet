package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type (
	// WechatCrypto 微信加密
	// examples:
	// wcp := wcrypto.New("your open appid", "your token", "your aes key")
	// wcp.Encrypt("your messge")
	// wcp.Decrypt("messge from wechat")
	WechatCrypto struct {
		appID          []byte
		token          string
		encodingAESKey []byte
		iv             []byte
	}

	// UserInfo 用户信息
	UserInfo struct {
		OpenID    string `json:"openId"`
		NickName  string `json:"nickName"`
		Gender    int    `json:"gender"`
		City      string `json:"city"`
		Province  string `json:"province"`
		Country   string `json:"country"`
		AvatarURL string `json:"avatarUrl"`
		UnionID   string `json:"unionId"`
		Watermark struct {
			AppID     string `json:"appid"`
			Timestamp int    `json:"timestamp"`
		} `json:"watermark"`
	}
)

// NewWechatCrypto 新建一个微信加密、解密工具
func NewWechatCrypto(appID, token, encodingAESKey string) *WechatCrypto {
	fmt.Println("Applet.NewWechatCrypto", token, encodingAESKey, appID)

	r, _ := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	return &WechatCrypto{
		token:          token,
		appID:          []byte(appID),
		encodingAESKey: []byte(r),
		iv:             ([]byte(r))[0:16],
	}
}

// Encrypt 加密
// 输入明文消息
// 输出密文消息
func (wc *WechatCrypto) Encrypt(text string) string {
	token := make([]byte, 16)
	rand.Read(token)
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(len(text)))
	msgBytes := bytes.Join([][]byte{token, b, []byte(text), wc.appID}, []byte(""))
	// aes
	block, _ := aes.NewCipher(wc.encodingAESKey)
	B := cipher.NewCBCEncrypter(block, wc.iv)
	encoded := encode(msgBytes)
	encrypted := make([]byte, len(encoded))
	B.CryptBlocks(encrypted, encoded)
	return base64.StdEncoding.EncodeToString(encrypted)
}

// Decrypt 解密
// 输入密文消息
// 输出明文消息
func (wc *WechatCrypto) Decrypt(text string) string {
	block, error := aes.NewCipher(wc.encodingAESKey)
	if error != nil {
		panic(error)
	}
	B := cipher.NewCBCDecrypter(block, wc.iv)
	dst, error := base64.StdEncoding.DecodeString(text)

	s := make([]byte, len(dst))
	B.CryptBlocks(s, dst)
	deciphered := decode(s)
	msg := deciphered[16:]
	length := binary.BigEndian.Uint32(msg[0:4])
	return string(msg[4 : 4+length])
}

// CheckSignature 校验微信消息是否合法
func (wc *WechatCrypto) CheckSignature(timestamp, nonce, signature string) bool {
	return checkSignature([]string{timestamp, nonce}, wc.token, signature)
}

// CheckMsgSignature 开发者计算签名
func (wc *WechatCrypto) CheckMsgSignature(timestamp, nonce, msgEncrypt, signature string) bool {
	return checkSignature([]string{timestamp, nonce, msgEncrypt}, wc.token, signature)
}

// DecryptUserInfo 解密用户信息
func (wc *WechatCrypto) DecryptUserInfo(sessionKey, encryptedData, rawData, iv, signature string) (user *UserInfo, err error) {
	user = &UserInfo{}
	err = nil

	// 校验数据是否合法
	if fmt.Sprintf("%x", sha1.Sum([]byte(rawData+sessionKey))) != signature {
		err = fmt.Errorf("%s", "invalid data")
		return
	}

	sessionKeyBase64, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return
	}

	encryptedDataBase64, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return
	}

	ivBase64, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return
	}

	block, err := aes.NewCipher([]byte(sessionKeyBase64))
	if err != nil {
		return
	}

	dst := make([]byte, len(encryptedDataBase64))
	decrypter := cipher.NewCBCDecrypter(block, ivBase64)
	decrypter.CryptBlocks(dst, encryptedDataBase64)

	err = json.Unmarshal(pkcs7UnPadding(dst), user)
	if err != nil {
		return
	}

	return
}

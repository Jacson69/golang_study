package main

import (
	"fmt"
	"math/rand"
	"time"
)

//	func generateCiphertext(length int) (string, error) {
//		max := big.NewInt(10)
//		ciphertext := ""
//
//		for i := 0; i < length; i++ {
//			randomNum, err := rand.Int(rand.Reader, max)
//			if err != nil {
//				return "", err
//			}
//
//			ciphertext += randomNum.String()
//		}
//		tmp := fmt.Sprintf("%x", md5.Sum([]byte(ciphertext)))
//		return tmp, nil
//		//return ciphertext[:length], nil
//	}
func getGmtIso8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func main() {
	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(hashedPassword))

	// 生成随机字节序列
	//randomBytes := make([]byte, 12)
	//_, err := rand.Read(randomBytes)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// 将字节序列编码为Base64字符串
	//base64Str := base64.StdEncoding.EncodeToString(randomBytes)
	//fmt.Println(base64Str)
	//// 生成18位的密文
	//ciphertext := base64Str[:16]
	//
	//fmt.Println("生成的密文:", ciphertext)

	//ciphertext, err := generateCiphertext(18)
	//if err != nil {
	//	fmt.Println("生成密文时发生错误:", err)
	//	return
	//}
	//
	//fmt.Println("生成的密文:", ciphertext)

	//authorizationCode, err := generateAuthorizationCode(18)
	//if err != nil {
	//	fmt.Println("生成授权码时发生错误:", err)
	//	return
	//}
	//
	//fmt.Println("生成的授权码:", authorizationCode)

	//authorizationCode, err := generateAuthorizationCode()
	//if err != nil {
	//	fmt.Println("生成授权码时发生错误:", err)
	//	return
	//}
	//
	//fmt.Println("生成的授权码:", authorizationCode)
	//
	//valid := validateAuthorizationCode(authorizationCode)
	//if valid {
	//	fmt.Println("授权码有效")
	//} else {
	//	fmt.Println("授权码无效")
	//}
	//qrCode := generatePaymentQRCode(authorizationCode)
	//fmt.Println("支付二维码:", qrCode)

	//Generweima()
	// 设置随机种子
	seed := time.Now().UnixNano() ^ int64(rand.Intn(999999))
	rand.Seed(seed)
	randomNumber := rand.Int63n(999999999999999999)

	// 将随机数转为18位的字符串，并在前面补零
	randomNumberStr := fmt.Sprintf("%018d", randomNumber)
	fmt.Println(randomNumberStr)
	//sprintf := fmt.Sprintf("%d-%d", randomNumber, 20)
	//fmt.Println(sprintf)
}

func Generweima() {
	var user_id int64 = 160
	type ConfigStruct struct {
		Expire int64 `json:"expire"`
		Id     int64 `json:"id"`
	}
	now := time.Now().Unix()
	var expire_time int64 = 30
	expire_end := now + expire_time
	fmt.Println(expire_end)
	//var tokenExpire = getGmtIso8601(expire_end)
	var config ConfigStruct
	config.Expire = expire_end
	config.Id = user_id
	//result, err := json.Marshal(config)
	//if err != nil {
	//	panic(err)
	//}
	var key = "secret"
	r := fmt.Sprintf("%d-%d", expire_end, user_id)
	encrypted := encryptString(r, key)

	fmt.Println("Encrypted:", encrypted)

	decrypted := decryptString(encrypted, key)

	fmt.Println("Decrypted:", decrypted)
	//cbc := utils.AesEncryptCBC(result, []byte(key))
	//toString := base64.URLEncoding.EncodeToString(cbc)
	//fmt.Println(toString)
	//result1 := string(utils.AesDecryptCBC(cbc, []byte(key)))
	//fmt.Println(result1)

	//str := "abcdd"
	//key := "1234567891011121"
	//
	//cbc1 := utils.AesEncryptCBC([]byte(str), []byte(key))
	//
	//result := string(utils.AesDecryptCBC(cbc1, []byte(key)))
	//fmt.Println(result)
	//err = qrcode.WriteFile(string(cbc), qrcode.Medium, 256, "qr1.png")
	//if err != nil {
	//	fmt.Println("生成二维码错误 error:", err)
	//}

}

func encryptString(input, key string) string {
	encrypted := ""
	for i := 0; i < len(input); i++ {
		char := input[i] ^ key[i%len(key)] // 使用异或操作进行加密
		encrypted += string(char)
	}
	//if len(encrypted) > 20 {
	//	//encrypted = encrypted[:20]
	//	fmt.Println(len(encrypted))
	//	//return "", fmt.Errorf("Encrypted string exceeds maximum length of 20 characters")
	//}
	return encrypted
}

func decryptString(encrypted, key string) string {
	decrypted := ""
	for i := 0; i < len(encrypted); i++ {
		char := encrypted[i] ^ key[i%len(key)] // 使用异或操作进行解密
		decrypted += string(char)
	}
	return decrypted
}

//func generateAuthorizationCode(length int) (string, error) {
//	randomBytes := make([]byte, (length+3)/4*3) // 生成足够长度的随机字节序列
//	_, err := rand.Read(randomBytes)
//	if err != nil {
//		return "", err
//	}
//
//	// 将字节序列编码为Base64字符串
//	authorizationCode := base64.URLEncoding.EncodeToString(randomBytes)
//
//	return authorizationCode[:length], nil
//}

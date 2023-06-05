package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/skip2/go-qrcode"
	"image/png"
	"os"
)

// 生成条形码图片
func generateBarcodeCode(userID string) (string, error) {
	// 使用Code 128生成条形码
	code, err := code128.Encode(userID)
	if err != nil {
		return "", err
	}

	// 创建一个要输出数据的文件
	file, _ := os.Create("qr3.png")
	defer file.Close()

	// 将条形码编码为字符串
	barcodeCode, err := barcode.Scale(code, 800, 100)
	if err != nil {
		return "", err
	}
	// 将code128的条形码编码为png图片
	png.Encode(file, barcodeCode)

	if err != nil {
		return "", err
	}
	return barcodeCode.Content(), nil
}

// 生成二维码图片
func QrCode(userID string) error {
	err := qrcode.WriteFile(userID, qrcode.Medium, 256, "qr1.png")
	if err != nil {
		fmt.Println("生成二维码错误 error:", err)
		return err
	}
	return nil
}

func main() {
	encrypted := "4310514082179637"
	code, err := GenerateQRCode(encrypted)
	if err != nil {
		fmt.Println("Error generating barcode code:", err)
		return
	}
	base64Str := base64.StdEncoding.EncodeToString(code)
	// 输出Base64字符串
	fmt.Println("Base64:", base64Str)
	c2, err := GenerateBarcode(encrypted)
	if err != nil {
		fmt.Println("Error generating barcode code:", err)
		return
	}
	base64Str2 := base64.StdEncoding.EncodeToString(c2)
	// 输出Base64字符串
	fmt.Println("Base64:", base64Str2)
}

// 生成条形码，生成byte，不生成图片
func GenerateBarcode(text string) ([]byte, error) {
	// 使用Code 128生成条形码
	code, err := code128.Encode(text)
	if err != nil {
		return nil, err
	}

	// 将条形码编码为字符串
	barcodeCode, err := barcode.Scale(code, 400, 80)
	if err != nil {
		return nil, err
	}

	// 创建缓冲区
	buf := new(bytes.Buffer)
	// 将图像数据写入缓冲区
	if err := png.Encode(buf, barcodeCode); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 生成二维码，输出byte，不生成图片
func GenerateQRCode(text string) ([]byte, error) {
	// 生成二维码图片
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// 将图片数据编码为PNG格式
	var pngBuf []byte
	pngBuf, err = qr.PNG(256)
	if err != nil {
		return nil, err
	}

	return pngBuf, nil
}

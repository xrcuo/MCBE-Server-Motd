/*
 * @Author: NyanCatda
 * @Date: 2021-12-26 21:23:59
 * @LastEditTime: 2022-01-03 16:29:15
 * @LastEditors: NyanCatda
 * @Description: Java服务器状态图片生成
 * @FilePath: \MotdBE\StatusImg\StatusImgJava.go
 */
package StatusImg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/BlackBEDevelopment/MCBE-Server-Motd/MotdBEAPI"
	"github.com/golang/freetype"
)

func ServerStatusImgJava(Host string) *bytes.Buffer {
	//获取服务器信息
	ServerData, err := MotdBEAPI.MotdJava(Host)
	if err != nil {
		fmt.Println(err)
	}
	if ServerData.Status == "offline" {
		offlineImgFile, err := os.Open("StatusImg/background.png")
		if err != nil {
			fmt.Println(err)
		}
		offlineImg, err := png.Decode(offlineImgFile)
		if err != nil {
			fmt.Println(err)
		}
		//将图片写入Buffer
		Buff := bytes.NewBuffer(nil)
		err = png.Encode(Buff, offlineImg)
		if err != nil {
			fmt.Println(err)
		}
		return Buff
	}

	//读取背景图片
	backgroundFile, err := os.Open("StatusImg/background.png")
	if err != nil {
		fmt.Println(err)
	}
	backgroundImg, err := png.Decode(backgroundFile)
	if err != nil {
		fmt.Println(err)
	}

	//转换类型
	img := backgroundImg.(*image.NRGBA)

	//读取字体数据
	fontBytes, err := ioutil.ReadFile("StatusImg/unifont-12.1.04.ttf")
	if err != nil {
		fmt.Println(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
	}

	//设置标题字体
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(72)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(30)
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(白色)
	f.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255}))
	pt := freetype.Pt(20, 30+int(f.PointToFixed(26))>>8)
	f.DrawString("MOTD: "+RemoveColorCode(ServerData.Motd), pt)

	//设置服务器图标
	ServerFaviconBase64 := strings.Replace(ServerData.Favicon, "data:image/png;base64,", "", -1)
	Favicon, _ := base64.StdEncoding.DecodeString(ServerFaviconBase64)
	FaviconBuffer := bytes.NewBuffer(Favicon)
	FaviconImg, _, err := image.Decode(FaviconBuffer)
	if err == nil {
		draw.Draw(img, img.Bounds(), FaviconImg, image.Pt(-550, -55), draw.Over)
	}

	//设置内容字体
	f = freetype.NewContext()
	//设置分辨率
	f.SetDPI(72)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(30)
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(白色)
	f.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))
	pt = freetype.Pt(20, 75+int(f.PointToFixed(26))>>8)
	f.DrawString("协议版本: "+strconv.Itoa(ServerData.Agreement), pt)
	pt = freetype.Pt(20, 125+int(f.PointToFixed(26))>>8)
	f.DrawString("游戏版本: "+ServerData.Version, pt)
	pt = freetype.Pt(20, 175+int(f.PointToFixed(26))>>8)
	f.DrawString("在线人数: "+strconv.Itoa(ServerData.Online)+"/"+strconv.Itoa(ServerData.Max), pt)
	pt = freetype.Pt(20, 325+int(f.PointToFixed(26))>>8)
	f.DrawString("连接延迟: "+strconv.FormatInt(ServerData.Delay, 10), pt)

	//将图片写入Buffer
	Buff := bytes.NewBuffer(nil)
	err = png.Encode(Buff, img)
	if err != nil {
		fmt.Println(err)
	}
	return Buff
}

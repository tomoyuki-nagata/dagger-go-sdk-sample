package repository

import (
	"fmt"
	"io"
	"net/http"
)

// 皇居の住所を返却します。
func Get() (string, error) {
	rsp, err := http.Get("https://zipcloud.ibsnet.co.jp/api/search?zipcode=100-0001")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 関数を抜ける際に必ずresponseをcloseするようにdeferでcloseを呼ぶ
	defer rsp.Body.Close()

	// レスポンスを取得し出力
	body, _ := io.ReadAll(rsp.Body)
	return string(body), err
}

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Ttsign struct {
	Params  string
	Data    string
	Cookies string
}

func (t *Ttsign) hash(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *Ttsign) getBaseString() string {
	baseStr := t.hash(t.Params)
	if t.Data != "" {
		baseStr += t.hash(t.Data)
	} else {
		baseStr += strings.Repeat("0", 32)
	}
	if t.Cookies != "" {
		baseStr += t.hash(t.Cookies)
	} else {
		baseStr += strings.Repeat("0", 32)
	}
	return baseStr
}

func (t *Ttsign) GetValue() map[string]string {
	return t.encrypt(t.getBaseString())
}

func (t *Ttsign) encrypt(data string) map[string]string {
	unix := time.Now().Unix()
	len := 0x14
	key := []int{0xDF, 0x77, 0xB9, 0x40, 0xB9, 0x9B, 0x84, 0x83, 0xD1, 0xB9, 0xCB, 0xD1, 0xF7, 0xC2, 0xB9, 0x85, 0xC3, 0xD0, 0xFB, 0xC3}
	paramList := []int{}
	for i := 0; i < 12; i += 4 {
		temp := data[8*i : 8*(i+1)]
		for j := 0; j < 4; j++ {
			h, _ := strconv.ParseInt(temp[j*2:(j+1)*2], 16, 32)
			paramList = append(paramList, int(h))
		}
	}
	paramList = append(paramList, []int{0x0, 0x6, 0xB, 0x1C}...)
	h, _ := strconv.ParseInt(fmt.Sprintf("%x", int(unix)), 16, 32)
	paramList = append(paramList, int((h&0xFF000000)>>24), int((h&0x00FF0000)>>16), int((h&0x0000FF00)>>8), int((h&0x000000FF)>>0))
	eorResultList := []int{}
	for i, v := range paramList {
		eorResultList = append(eorResultList, v^key[i])
	}
	for i := 0; i < len; i++ {
		c := t.reverse(eorResultList[i])
		d := eorResultList[(i+1)%len]
		e := c ^ d
		f := t.rbitAlgorithm(e)
		h := ((f ^ 0xFFFFFFFF) ^ len) & 0xFF
		eorResultList[i] = h
	}
	result := ""
	for _, v := range eorResultList {
		result += t.hexString(v)
	}
	return map[string]string{
		"X-SS-REQ-TICKET": strconv.FormatInt(unix*1000, 10),
		"X-Khronos":       strconv.FormatInt(unix, 10),
		"X-Gorgon":        "0404b0d30000" + result,
	}
}

func (t *Ttsign) rbitAlgorithm(num int) int {
	result := ""
	tmpString := fmt.Sprintf("%b", num)
	for len(tmpString) < 8 {
		tmpString = "0" + tmpString
	}
	for i := 0; i < 8; i++ {
		result += string(tmpString[7-i])
	}
	r, _ := strconv.ParseInt(result, 2, 32)
	return int(r)
}

func (t *Ttsign) hexString(num int) string {
	tmpString := fmt.Sprintf("%x", num)
	if len(tmpString) < 2 {
		tmpString = "0" + tmpString
	}
	return tmpString
}

func (t *Ttsign) reverse(num int) int {
	tmpString := t.hexString(num)
	r, _ := strconv.ParseInt(tmpString[1:]+tmpString[:1], 16, 32)
	return int(r)
}

// func main() {
// 	url := "https://api19-va.tiktokv.com/aweme/v1/user/profile/other/?sec_user_id=MS4wLjABAAAAjq0XOponqLGAoqP7qHkmRoqVm_IzcfFt43KQdouIWXqdYSgxxy2Ao4GX7ifhFKNv&address_book_access=1&from=0&os_api=25&device_type=ASUS_Z01QD&ssmix=a&manifest_version_code=2021704040&dpi=240&carrier_region=DE&uoo=0&region=DE&app_name=musical_ly&version_name=17.4.4&timezone_offset=3600&ts=1683035921&ab_version=17.4.4&residence=DE&pass-route=1&cpu_support64=true&pass-region=1&current_region=DE&storage_type=1&ac2=wifi&appTheme=light&app_type=normal&ac=wifi&host_abi=armeabi-v7a&channel=googleplay&update_version_code=2021704040&_rticket=1683035921475&device_platform=android&iid=7228561267690456837&build_number=17.4.4&locale=de-DE&op_region=DE&version_code=170404&timezone_name=Europe%2FBerlin&cdid=90c40bcb-3180-4f7d-86fa-a5d1e5364edc&openudid=d8e48a4690fded4b&sys_region=DE&device_id=7191957835704632838&app_language=de&resolution=720*1280&device_brand=Asus&language=de&os_version=7.1.2&aid=1233"
// 	signed := ttsign{
// 		params:  strings.Split(url, "?")[1],
// 		data:    "",
// 		cookies: "",
// 	}
// 	x := signed.getValue()
// 	b, _ := json.Marshal(x)
// 	fmt.Println(string(b))
// }

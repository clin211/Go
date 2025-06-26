package wechat

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	redirectOauthURL       = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	webAppRedirectOauthURL = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	webAccessTokenURL      = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL  = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL            = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	checkAccessTokenURL    = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
	accessTokenURL         = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	getTicketURL           = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// 获取redirectOauthURL链接。其在微信中跳转后可以获取code
// 因为微信转码会造成部分链接参数丢失的情况，使用urlEncode对链接进行处理
func RedirectOauthUrl(appID, redirectUrl string) string {
	if appID == "" || redirectUrl == "" {
		return ""
	}

	// url encode
	v := url.Values{}
	v.Add("redirectUrl", redirectUrl) // 添加map
	encodeUrl := v.Encode()
	encodeUrl = strings.TrimLeft(encodeUrl, "redirectUrl=") //去掉url中多余的字符串
	urlStr := fmt.Sprintf(redirectOauthURL, appID, encodeUrl, "snsapi_userinfo", "")
	return urlStr
}

// 网页授权access_token
type ResWebAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

// 获取网页授权access_token
func GetWebAccessToken(appID, secret, code string) (res ResWebAccessToken, err error) {
	urlStr := fmt.Sprintf(webAccessTokenURL, appID, secret, code)
	body, err := util.HttpGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	if res.Errcode != 0 {
		err = fmt.Errorf("GetWebAccessToken error : errcode=%v , errmsg=%v", res.Errcode, res.Errmsg)
		return
	}
	return
}

type WxUserInfo struct {
	ID         int        `json:"id"`
	Openid     string     `json:"openid"`
	Nickname   string     `json:"nickname"`
	Headimgurl string     `json:"headimgurl"`
	Sex        int        `json:"sex"`
	Province   string     `json:"province"`
	City       string     `json:"city"`
	Country    string     `json:"country"`
	Name       string     `json:"name"`
	Mobile     string     `json:"mobile"`
	Address    string     `json:"address"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Errcode    int        `json:"errcode"`
	Errmsg     string     `json:"errmsg"`
}

func GetUserInfo(accessToken, openID string) (res WxUserInfo, err error) {
	urlStr := fmt.Sprintf(userInfoURL, accessToken, openID)
	body, err := util.HttpGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	if res.Errcode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", res.Errcode, res.Errmsg)
		return
	}
	return
}

// 普通access_token
type ResAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

// 获取普通access_token
func GetAccessToken(appID, secret string) (res ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, appID, secret)
	body, err := util.HttpGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	if res.Errcode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", res.Errcode, res.Errmsg)
		return
	}
	return
}

type resTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
}

func GetTicket(accessToken string) (res resTicket, err error) {
	urlStr := fmt.Sprintf(getTicketURL, accessToken)
	body, err := util.HttpGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	if res.Errcode != 0 {
		err = fmt.Errorf("getTicket Error : errcode=%d , errmsg=%s", res.Errcode, res.Errmsg)
		return
	}
	return
}

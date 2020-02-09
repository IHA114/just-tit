package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const KeezmoviesApiURL = "http://www.keezmovies.com/wapi/"
const KeezmoviesApiTimeout = 5

type KeezmoviesEmbedCode map[string]interface{}
type KeezmoviesSingleVideo map[string]interface{}

type KeezmoviesController struct {
	beego.Controller
}

func (c *KeezmoviesController) Get() {
	aux := strings.Replace(c.Ctx.Request.URL.Path, ".html", "", -1)
	str := strings.Split(aux, "/")
	videoID := str[2]

	redirect := "https://www.keezmovies.com/video/title-" + videoID + "?utm_source=just-tit.com&utm_medium=embed&utm_campaign=hubtraffic_dsmatilla"

	BaseDomain := "https://just-tit.com"

	type TemplateData = map[string]interface{}

	c.Data["ID"] = videoID
	c.Data["Domain"] = BaseDomain

	videocode := KeezmoviesGetVideoByID(videoID)
	_, ok := videocode["video"]
	if !ok { c.Redirect(redirect, 307) }
	video := videocode["video"].(map[string]interface{})
	embedcode := KeezmoviesGetVideoEmbedCode(videoID)
	if embedcode["error"] != nil {
		log.Println("[KEEZMOVIES][GET]",embedcode["error"])
		c.Redirect(redirect, 307)
	}
	embed := embedcode["video"].(map[string]interface{})

	str2, _ := base64.StdEncoding.DecodeString(fmt.Sprintf("%s", embed["embed_code"]))
	c.Data["Embed"] = template.HTML(fmt.Sprintf("%+v", html.UnescapeString(string(str2))))
	c.Data["PageTitle"] = fmt.Sprintf("%s", video["title"])
	c.Data["PageMetaDesc"] = fmt.Sprintf("%s", video["title"])
	c.Data["Thumb"] = fmt.Sprintf("%s", video["image_url"])
	c.Data["Url"] = fmt.Sprintf(BaseDomain+"/keezmovies/%s.html", videoID)
	c.Data["Width"] = "650"
	c.Data["Height"] = "550"
	c.Data["KeezmoviesVideo"] = video

	if c.Data["PageTitle"] == "" {
		c.Redirect(redirect, 307)
	}

	if c.GetString("tp") == "true" {
		c.TplName = "video/player.html"
	} else {
		c.Data["Result"] = doSearch(fmt.Sprintf("%s", fmt.Sprintf("%s", video["title"])))
		c.TplName = "index.html"
	}
}

func KeezmoviesGetVideoByID(ID string) KeezmoviesSingleVideo {
	timeout := time.Duration(KeezmoviesApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, _ := client.Get(fmt.Sprintf(KeezmoviesApiURL+"getVideoById?output=json&video_id=%s", ID))
	b, _ := ioutil.ReadAll(resp.Body)
	var result KeezmoviesSingleVideo
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[KEEZMOVIES][GETVIDEOBYID]",err)
	}
	return result

}

func KeezmoviesGetVideoEmbedCode(ID string) KeezmoviesEmbedCode {
	timeout := time.Duration(KeezmoviesApiTimeout * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, _ := client.Get(fmt.Sprintf(KeezmoviesApiURL+"getVideoEmbedCode?output=json&video_id=%s", ID))
	b, _ := ioutil.ReadAll(resp.Body)
	var result KeezmoviesEmbedCode
	err := json.Unmarshal(b, &result)
	if err != nil {
		log.Println("[KEEZMOVIES][GETVIDEOEMBEDCODE]",err)
	}
	return result
}

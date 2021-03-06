package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	HOME_FORMAT = `
	<h1>
	    <a href='%s'>%s</a>
	</h1>
	<h2>
	    %s
	</h2>
	<br />
	`
	// HomeFormat ...
	HomeFormat string = `
   [%s](%s)
   > %s
`
)

type HomeInterface interface {
	GetData() string
}

type Home struct {
}

type Jue struct {
}

type info struct {
	Success int `json:"Success"`
	Result  []a
}

type a struct {
	Newsid      int    `json:"newsid"`
	Title       string `json:"title"`
	Orderdate   string `json:"orderdate"`
	Description string `json:"description"`
	Isad        bool   `json:"isad"`
	WapNewsUrl  string `json:"WapNewsUrl"`
	NewsTips    []Ad   `json:"NewsTips"`
	//Newsid        int         `json:"newsid"`
	// V             string      `json:"v"`
	// Postdate      string      `json:"postdate"`
	// Image         string      `json:"image"`
	// Slink         string      `json:"slink"`
	// Hitcount      int         `json:"hitcount"`
	// Commentcount  int         `json:"commentcount"`
	// Cid           int         `json:"cid"`
	// Url           string      `json:"url"`
	// Live          int         `json:"live"`
	// Lapinid       int         `json:"lapinid"`
	// Forbidcomment string      `json:"forbidcomment"`
	// Imagelist     interface{} `json:"imagelist"`
	// C             string      `json:"c"`
	// Client        string      `json:"client"`
	// Sid           int         `json:"sid"`
	// PostDateStr   string      `json:"PostDateStr"`
	// HitCountStr   string      `json:"HitCountStr"`
	// NewsTips      interface{} `json:"NewsTips"`
}

type Ad struct {
	TipClass string
	TipName  string
}

func (h *Home) GetData() (z string, err error) {
	t := strconv.FormatInt(time.Now().Unix(), 10) + "000"
	z = ""
	d, _ := base64.StdEncoding.DecodeString(homeURL)
	m := string(d)
	str := m + homePath
	param := url.Values{}
	u, _ := url.Parse(str)
	param.Set("Tag", "")
	param.Set("ot", t)
	param.Set("page", "0")
	u.RawQuery = param.Encode()
	uPath := u.String()
	log.Println(uPath)
	resp, err := http.Get(uPath)
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(data))
	info := info{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		panic(err)
	}
	for _, v := range info.Result {
		z += fmt.Sprintf(HomeFormat, v.Title, v.WapNewsUrl, v.Description)
	}
	return z, err
}

func (h *Home) GetOneData(open bool) (string, error) {
	haveMore := true
	z := ""
	s, _ := time.LoadLocation("Asia/Shanghai")
	t := strconv.FormatInt(time.Now().In(s).Unix(), 10) + "000"
	d, _ := base64.StdEncoding.DecodeString(homeURL)
	str := string(d) + homePath
	rt := 0
	for haveMore {
		haveMore = false
		param := url.Values{}
		u, _ := url.Parse(str)
		param.Set("Tag", "")
		param.Set("ot", t)
		param.Set("page", "0")
		u.RawQuery = param.Encode()
		uPath := u.String()
		log.Println(uPath)

		req, err := http.NewRequest(http.MethodGet, uPath, nil)
		req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.79 Safari/537.36")
		client := &http.Client{}
		resp, _ := client.Do(req)

		//resp, err := http.Get(uPath)
		if resp.StatusCode != 200 {
			rt++
			continue
		}
		if resp.StatusCode != 200 && rt >= 3 {
			break
		}
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(data))
		info := info{}
		err = json.Unmarshal(data, &info)
		if err != nil {
			log.Println(err.Error())
		}
		for _, v := range info.Result {
			if v.Newsid == 1 {
				continue
			}
			now, err := time.ParseInLocation("2006-01-02T15:04:05", v.Orderdate, s)
			if err != nil {
				now, err = time.ParseInLocation(time.RFC3339, v.Orderdate, s)
				if err != nil {
					continue
				}
			}
			if now.Format("2006-01-02") != time.Now().In(s).Format("2006-01-02") {
				continue
			}
			if len(v.NewsTips) > 0 {
				if v.NewsTips[0].TipName == "广告" {
					continue
				}
			}
			t = strconv.FormatInt(now.Unix(), 10) + "000"
			// if !open {
			// 	v.WapNewsUrl = ""
			// }
			// z += fmt.Sprintf(HOME_FORMAT, v.WapNewsUrl, v.Title, v.Description)
			z += fmt.Sprintf(HomeFormat, v.Title, v.WapNewsUrl, v.Description)
			haveMore = true
			time.Sleep(time.Second)
		}
	}
	return z, nil
}

func (j *Jue) GetOneData(open bool) (string, error) {
	haveMore := true
	z := ""
	s, _ := time.LoadLocation("Asia/Shanghai")
	d, _ := base64.StdEncoding.DecodeString(jueURL)
	wap, _ := base64.StdEncoding.DecodeString(jueWap)
	str := string(d) + juePath
	before := ""
	for haveMore {
		haveMore = false
		param := url.Values{}
		u, _ := url.Parse(str)
		param.Set("uid", "")
		param.Set("device_id", "")
		param.Set("token", "")
		param.Set("src", "web")
		param.Set("before", before)
		param.Set("limit", "30")
		u.RawQuery = param.Encode()
		uPath := u.String()
		log.Println(uPath)
		resp, err := http.Get(uPath)
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(data))
		info := JueResult{}
		err = json.Unmarshal(data, &info)
		if err != nil {
			panic(err)
		}
		for _, v := range info.D.List {
			now, err := time.ParseInLocation("2006-01-02T15:04:05Z", v.CreatedAt, s)
			if err != nil {
				panic(err)
			}
			if now.Format("2006-01-02") != time.Now().In(s).Format("2006-01-02") {
				continue
			}
			//wap_url := ""
			if open {
				_ = string(wap) + v.ObjectId
			}

			//z += fmt.Sprintf("<a href='%s'><h2>%s</h2></a><br />", wap_url, v.Content)
			z += fmt.Sprintf("<h2>%s %s</h2><br />", v.Content, v.Url)
			for _, vv := range v.Pictures {
				z += fmt.Sprintf("<img src='%s' width='600' height='auto'/>", vv)
			}
			before = v.CreatedAt
			haveMore = true
		}
	}
	return z, nil
}

func (j *Jue) GetData() (string, error) {
	d, _ := base64.StdEncoding.DecodeString(jueURL)
	m := string(d)
	str := m + juePath
	param := url.Values{}
	u, _ := url.Parse(str)
	param.Set("uid", "")
	param.Set("device_id", "")
	param.Set("token", "")
	param.Set("src", "web")
	param.Set("before", "")
	param.Set("limit", "30")
	u.RawQuery = param.Encode()
	uPath := u.String()
	log.Println(uPath)
	resp, err := http.Get(uPath)
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(data))
	info := JueResult{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		panic(err)
	}
	z := ""
	for _, v := range info.D.List {
		z += fmt.Sprintf("<h2>%s</h2><br />", v.Content)
		for _, vv := range v.Pictures {
			z += fmt.Sprintf("<img src='%s' width='600' height='auto'/>", vv)
		}
	}
	return z, err
}

type JueResult struct {
	S int     `json:"s"`
	M string  `json:"m"`
	D JueList `json:"d"`
}

type JueList struct {
	Total int         `json:"total"`
	List  []JueObject `json:"list"`
}

type JueObject struct {
	Uid       string   `json:"uid"`
	Content   string   `json:"content"`
	Pictures  []string `json:"pictures"`
	CreatedAt string   `json:"createdAt"`
	ObjectId  string   `json:"objectId"`
	Url       string   `json:"url"`
}

func NewHome() *Home {
	return &Home{}
}

package connection

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fgv50/tools"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/twmb/murmur3"
)

type resps struct {
	url        string
	body       string
	header     map[string][]string
	server     string
	statuscode int
	length     int
	title      string
	jsurl      []string
	favhash    string
}

func rndua() string {
	ua := []string{"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2226.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1",
		"Mozilla/5.0 (Windows NT 6.3; rv:36.0) Gecko/20100101 Firefox/36.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10; rv:33.0) Gecko/20100101 Firefox/33.0",
		"Mozilla/5.0 (X11; Linux i586; rv:31.0) Gecko/20100101 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:31.0) Gecko/20130401 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (compatible, MSIE 11, Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows; Intel Windows) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.67"}
	n := rand.Intn(13) + 1
	return ua[n]
}

func gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}

func getfavicon(httpbody string, turl string) string {
	faviconpaths := xegexpjs(`href="(.*?favicon....)"`, httpbody)
	var faviconpath string
	u, err := url.Parse(turl)
	if err != nil {
		panic(err)
	}
	turl = u.Scheme + "://" + u.Host
	if len(faviconpaths) > 0 {
		fav := faviconpaths[0][1]
		if fav[:2] == "//" {
			faviconpath = "http:" + fav
		} else {
			if fav[:4] == "http" {
				faviconpath = fav
			} else {
				faviconpath = turl + "/" + fav
			}

		}
	} else {
		faviconpath = turl + "/favicon.ico"
	}
	return favicohash(faviconpath)
}

func Httprequest(url string, cycle int, chsize int) (*resps, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	cookie := &http.Cookie{
		Name:  "rememberMe",
		Value: "me",
	}
	req.AddCookie(cookie)
	req.Header.Set("Accept", "*/*;q=0.8")
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", rndua())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	httpbody := string(result)
	httpbody = tools.ToUtf8(httpbody, contentType)
	title := gettitle(httpbody)
	httpheader := resp.Header
	var server string
	capital, ok := httpheader["Server"]
	if ok {
		server = capital[0]
	} else {
		Powered, ok := httpheader["X-Powered-By"]
		if ok {
			server = Powered[0]
		} else {
			server = "None"
		}
	}
	var jsurl []string
	if cycle < chsize {
		jsurl = Jsjump(httpbody, url)
	} else {
		jsurl = []string{""}
	}
	favhash := getfavicon(httpbody, url)
	s := resps{url, httpbody, resp.Header, server, resp.StatusCode, len(httpbody), title, jsurl, favhash}
	return &s, nil
}


func xegexpjs(reg string, resp string) (reslut1 [][]string) {
	reg1 := regexp.MustCompile(reg)
	if reg1 == nil {
		log.Println("regexp err")
		return nil
	}
	result1 := reg1.FindAllStringSubmatch(resp, -1)
	return result1
}

func Mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	_, err := h32.Write([]byte(raw))
	if err == nil {
		return fmt.Sprintf("%d", int32(h32.Sum32()))
	} else {
		//log.Println("favicon Mmh3Hash32 error:", err)
		return "0"
	}
}

func StandBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()

}

func favicohash(host string) string {
	timeout := time.Duration(8 * time.Second)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   timeout,
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse /* 不进入重定向 */
		},
	}
	resp, err := client.Get(host)
	if err != nil {
		//log.Println("favicon client error:", err)
		return "0"
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println("favicon file read error: ", err)
			return "0"
		}
		return Mmh3Hash32(StandBase64(body))
	} else {
		return "0"
	}
}

func Jsjump(str string, url string) []string {
	regs := []string{`(window|top)\.location\.href = "(.*?)"`, `redirectUrl = "(.*?)"`, `<meta.*?http-equiv=.*?refresh.*?url=(.*?)>`}
	var results []string
	for _, reg := range regs {
		result1 := xegexpjs(reg, str)
		//fmt.Println(result1)
		if result1 != nil {
			if len(result1) > 0 {
				for _, m := range result1 {
					s := len(m)
					if strings.Contains(m[s-1], "http") {
						continue
					} else {
						str2 := strings.Trim(m[s-1], "/")
						str2 = strings.ReplaceAll(str2, "../", "/")
						if str2[:1] == "/" {
							results = append(results, url+str2)
						} else {
							results = append(results, url+"/"+str2)
						}
					}
				}
			} else {
				continue
			}
		}
	}
	return results
}

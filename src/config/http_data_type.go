package config
import (
	"net/url"
	"net/http"
	"log"
	"time"
	"fmt"
	"strings"
	"regexp"
)

type Method string
const (
	POST Method = "POST"
	GET Method = "GET"
)
type Scheme string
const (
	HTTP Scheme = "http://"
	HTTPS Scheme = "https://"
)
// TODO: setup proxy service
type Link struct{
	Protocol Scheme

	URL *url.URL

	// this link retry times for http request
	RetryTimes int

	// "GET" or "POST"
	Type Method

	// form data for "POST" request
	FormData url.Values

	Header http.Header
	Cookies []*http.Cookie

	// Dialer timeout and client timeout.
	Timeout time.Duration

	// whether this Link is available. such as http://javascript_xxx is not available.
	Available bool
	UnavailableAddr string
}
// note: only handle http or https, others will be setup unavailable. such as mailto:xxx.com
func NewLink(addr string) *Link{
	var err error
	link := new(Link)
	link.RetryTimes = 3
	link.Protocol = HTTPS
	link.Available = true
	reg := regexp.MustCompile("\\s+")
	addr = reg.ReplaceAllString(addr, "")
	if addr == ""{
		log.Printf("url is \"\"\n", addr)
		return &Link{Available: false, UnavailableAddr: "no addr"}
	}
	n := strings.Index(addr, string(HTTPS))
	if n >= 0{
		addr = addr[n:]
	}else{
		n = strings.Index(addr, string(HTTP))
		if n >= 0{
			addr = addr[n:]
		}else{

			items := strings.Split(addr, ":")
			tmp := false
			if len(items) == 2{
				tmp, _ = regexp.MatchString("\\d+", items[1])
				if !tmp{
					return &Link{Available: false, UnavailableAddr: addr}
				}
			}
			if len(items) == 1 || tmp{
				addr = fmt.Sprintf("%s%s", string(link.Protocol), addr)
			}
			if len(items) > 2{
				return &Link{Available: false, UnavailableAddr: addr}
			}
		}
	}
	link.URL, err = url.Parse(addr)
	if err != nil{
		log.Printf("parse url %s failed!\n", addr)
		return &Link{Available: false, UnavailableAddr: addr}
	}
	link.Type = GET
	link.FormData = make(url.Values)
	link.Header = make(http.Header)
	// link.AddHeader("User-Agent", "xxx")

	link.Cookies = make([]*http.Cookie, 0, 16)
	link.Timeout = 30 * time.Second
	return link
}
func (this *Link) AddHeader(key, value string){
	this.Header.Add(key, value)
}
func (this *Link) AddCookie(cookie *http.Cookie){
	this.Cookies = append(this.Cookies, cookie)
}
func (this *Link) AddReqParam(key, value string){
	this.FormData.Add(key, value)
}

type RawData struct{
	LINK *Link
	Data string

	// if LINK is not available, the RawData should be the same.
	Available bool
}
func NewRawData(link *Link, data string) *RawData{
	if link == nil || !link.Available{
		return &RawData{Available: false, LINK: link}
	}
	rd := new(RawData)
	rd.Data = data
	rd.LINK = link
	rd.Available = true
	return rd
}

type Text struct{
	LINK *Link
	Data string

	// if LINK is not available, the Text should be the same.
	Available bool
}
func NewText(link *Link, data string) *Text{
	if link == nil || !link.Available{
		return &Text{Available: false, LINK: link}
	}
	txt := new(Text)
	txt.Data = data
	txt.LINK = link
	txt.Available = true
	return txt
}

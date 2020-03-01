package sina
import (
	"config"
	"fmt"
	"github.com/antchfx/htmlquery"
	"slog"
	"strings"
	"spider"
)
/***************************************layer-0 resolve strategy*************************************/
type RP0 struct{
	spider.ResolverPart
}
func NewRP0(name string, rawData *config.RawData) config.IResolverPart{
	d := new(RP0)
	d.InitObject(name, rawData)
	return d
}
// default option is resolving the page and getting all links in this page.
func (this *RP0) GenerateLinks() ([]*config.Link, error){
	links := make([]*config.Link, 0, 128)
	doc, err := htmlquery.Parse(strings.NewReader(this.GetRawData().Data))
	if err != nil{
		slog.Error(this.GetLogName(), fmt.Sprintf("Resolve url %s faild.\n%s", this.GetRawData().LINK.URL.String(), err))
		return links, err
	}
	for _, e := range htmlquery.Find(doc, "//a[@href]"){
		href := htmlquery.SelectAttr(e, "href")

		// TODO:maybe need to config the generated link, such as set cookie, header and so on. 
		// what's more, the link maybe a relative path, need to transfer by infomation of this.GetRawData().LINK.
		if strings.HasPrefix(href, "//"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[2:])
		}
		if strings.HasPrefix(href, "/"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[1:])
		}
		if href != ""{
			link := config.NewLink(href)
			links = append(links, link)
		}
	}
	return links, nil
}
// default option is return html data directly, not the pure text.
func (this *RP0) GenerateText() (*config.Text, error){
	text := config.NewText(this.GetRawData().LINK, this.GetRawData().Data)
	return text, nil
}

/***************************************layer-1 resolve strategy*************************************/
type RP1 struct{
	spider.ResolverPart
}
func NewRP1(name string, rawData *config.RawData) config.IResolverPart{
	d := new(RP1)
	d.InitObject(name, rawData)
	return d
}
// default option is resolving the page and getting all links in this page.
func (this *RP1) GenerateLinks() ([]*config.Link, error){
	links := make([]*config.Link, 0, 128)
	doc, err := htmlquery.Parse(strings.NewReader(this.GetRawData().Data))
	if err != nil{
		slog.Error(this.GetLogName(), fmt.Sprintf("Resolve url %s faild.\n%s", this.GetRawData().LINK.URL.String(), err))
		return links, err
	}
	for _, e := range htmlquery.Find(doc, "//a[@href]"){
		href := htmlquery.SelectAttr(e, "href")

		// TODO:maybe need to config the generated link, such as set cookie, header and so on. 
		// what's more, the link maybe a relative path, need to transfer by infomation of this.GetRawData().LINK.
		if strings.HasPrefix(href, "//"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[2:])
		}
		if strings.HasPrefix(href, "/"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[1:])
		}
		if href != ""{
			link := config.NewLink(href)
			links = append(links, link)
		}
	}
	return links, nil
}
// default option is return html data directly, not the pure text.
func (this *RP1) GenerateText() (*config.Text, error){
	text := config.NewText(this.GetRawData().LINK, this.GetRawData().Data)
	return text, nil
}

/***************************************layer-2 resolve strategy*************************************/
type RP2 struct{
	spider.ResolverPart
}
func NewRP2(name string, rawData *config.RawData) config.IResolverPart{
	d := new(RP2)
	d.InitObject(name, rawData)
	return d
}
// default option is resolving the page and getting all links in this page.
func (this *RP2) GenerateLinks() ([]*config.Link, error){
	links := make([]*config.Link, 0, 128)
	doc, err := htmlquery.Parse(strings.NewReader(this.GetRawData().Data))
	if err != nil{
		slog.Error(this.GetLogName(), fmt.Sprintf("Resolve url %s faild.\n%s", this.GetRawData().LINK.URL.String(), err))
		return links, err
	}
	for _, e := range htmlquery.Find(doc, "//a[@href]"){
		href := htmlquery.SelectAttr(e, "href")

		// TODO:maybe need to config the generated link, such as set cookie, header and so on. 
		// what's more, the link maybe a relative path, need to transfer by infomation of this.GetRawData().LINK.
		if strings.HasPrefix(href, "//"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[2:])
		}
		if strings.HasPrefix(href, "/"){
			href = fmt.Sprintf("%s://%s/%s", this.GetRawData().LINK.URL.Scheme, this.GetRawData().LINK.URL.Host, href[1:])
		}
		if href != ""{
			link := config.NewLink(href)
			links = append(links, link)
		}
	}
	return links, nil
}
// default option is return html data directly, not the pure text.
func (this *RP2) GenerateText() (*config.Text, error){
	text := config.NewText(this.GetRawData().LINK, this.GetRawData().Data)
	return text, nil
}

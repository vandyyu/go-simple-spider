package spider
import (
	"config"
	"fmt"
	"github.com/antchfx/htmlquery"
	"slog"
	"strings"
)
type DefaultResolverPart struct{
	ResolverPart
}
func NewDefaultResolverPart(name string, rawData *config.RawData) config.IResolverPart{
	d := new(DefaultResolverPart)
	d.InitObject(name, rawData)
	return d
}
// default option is resolving the page and getting all links in this page.
func (this *DefaultResolverPart) GenerateLinks() ([]*config.Link, error){
	links := make([]*config.Link, 0, 128)
	doc, err := htmlquery.Parse(strings.NewReader(this.rawData.Data))
	if err != nil{
		slog.Error(this.GetLogName(), fmt.Sprintf("Resolve url %s faild.\n%s", this.rawData.LINK.URL.String(), err))
		return links, err
	}
	for _, e := range htmlquery.Find(doc, "//a[@href]"){
		href := htmlquery.SelectAttr(e, "href")

		// TODO:maybe need to config the generated link, such as set cookie, header and so on. 
		// what's more, the link maybe a relative path, need to transfer by infomation of this.rawData.LINK.
		if strings.HasPrefix(href, "//"){
			href = fmt.Sprintf("%s://%s/%s", this.rawData.LINK.URL.Scheme, this.rawData.LINK.URL.Host, href[2:])
		}
		if strings.HasPrefix(href, "/"){
			href = fmt.Sprintf("%s://%s/%s", this.rawData.LINK.URL.Scheme, this.rawData.LINK.URL.Host, href[1:])
		}
		if href != ""{
			link := config.NewLink(href)
			links = append(links, link)
		}
	}
	return links, nil
}
// default option is return html data directly, not the pure text.
func (this *DefaultResolverPart) GenerateText() (*config.Text, error){
	text := config.NewText(this.rawData.LINK, this.rawData.Data)
	return text, nil
}

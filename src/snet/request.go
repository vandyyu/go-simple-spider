package snet
import (
	"config"
	 "net"
	 "net/http"
	 "slog"
	 "fmt"
	 "io"
	 "io/ioutil"
	 "strings"
	 "time"
 )

// const MAX_NUM_CONN_SELENIUM int = 10
// TODO: TCP pool to be implemented. connTCP should be a global object, and can be reused many times.
// TODO: default request type is "GET", "POST" method is not implemented.
// TODO: many other request info doesn't transfer to selenium, such as header, cookies and so on. now is noly transfer url, such as "https://www.baidu.com".
// TODO: set timeout.
func RequestSelenium(logName string, link *config.Link) string{
	slog.Info(logName, fmt.Sprintf("Start request %s by selenium.", link.URL.String()))
	data := make([]byte, 0, 1 * 1024 * 1024)
	for retryIndex := 0; retryIndex < link.RetryTimes; retryIndex++{
		conn, err := net.Dial("tcp", ":9999")
		if err != nil{
			slog.Error(logName, fmt.Sprintf("Connet tcp :9999 in localhost failed! %dth retrying Sleeping 2 secs.", retryIndex+1), err)
			time.Sleep(2 * time.Second)
			if conn != nil{
				conn.Close()
			}
			continue
		}
		conn.Write([]byte(link.URL.String() + "\n"))

		data = data[:0]
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil && err == io.EOF{
				break
			}else if err != nil{
				msg := fmt.Sprintf("Request url %s from selenium, stream doesn't read competely.", link.URL.String())
				slog.Error(logName, msg, err)
				data = data[:0]
				break
			}else{
				data = append(data, buf[:n]...)
			}
		}
		if conn != nil{
			conn.Close()
		}
		if string(data) != ""{
			break
		}else{
			slog.Info(logName, fmt.Sprintf("Request url %s by selenium failed, the %dth retrying. Sleeping 1 secs", link.URL.String(), retryIndex+1))
			time.Sleep(1 * time.Second)
		}

	}
	if string(data) == ""{
		slog.Warning(logName, fmt.Sprintf("Request url %s by selenium still failed after retrying %d times.", link.URL.String(), link.RetryTimes))
	}else{
		msg := fmt.Sprintf("End request %s by selenium. Data length is: %d bytes.", link.URL.String(), len(data))
		slog.Info(logName, msg)
	}
	return string(data)
}

// this method use golang's net package, "GET" and "POST" are included.
// just simple encapsulation
// TODO: add proxy service
func RequestNative(logName string, link *config.Link) string{
	slog.Info(logName, fmt.Sprintf("Start request %s by native http.", link.URL.String()))
	var (
		ret string
		client *http.Client
		req *http.Request
		resp *http.Response
		err error
	)
	ret = ""
	if link.Timeout > 0{
		transport := &http.Transport{
			 Dial: (&net.Dialer{
			       Timeout: link.Timeout,
			 }).Dial,
			 MaxIdleConns: 0,
		}
		client = &http.Client{
			Timeout: link.Timeout,
			Transport: transport,
		}
	}else{
		client = &http.Client{}
	}
	for retryIndex := 0; retryIndex < link.RetryTimes; retryIndex++{
		if link.Type == config.GET{
			link.URL.RawQuery = link.FormData.Encode()
			reqURL := link.URL.String()

			req, err = http.NewRequest("GET", reqURL, nil)
		}
		if link.Type == config.POST{
			req, err = http.NewRequest("POST", link.URL.String(), strings.NewReader(link.FormData.Encode()))
		}
		req.Header = link.Header
		for _, c := range(link.Cookies){
			req.AddCookie(c)
		}
		resp, err = client.Do(req)
		if err != nil{
			msg := fmt.Sprintf("Request url %s by native, request http failed!", link.URL.String())
			slog.Error(logName, msg, err)
			ret = ""
			break
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			msg := fmt.Sprintf("Request url %s by native, read data failed.", link.URL.String())
			slog.Error(logName, msg, err)
			ret = ""
			break
		}
		ret = string(body)
		if ret != "" && resp.StatusCode == 200{
			break
		}else{
			slog.Info(logName, fmt.Sprintf("Request url %s by native http failed, the %dth retrying. Sleeping 1 secs.", link.URL.String(), retryIndex+1))
			time.Sleep(1 * time.Second)
		}

	}
	if ret == ""{
		slog.Warning(logName, fmt.Sprintf("Request url %s by native http still failed after retrying %d times.", link.URL.String(), link.RetryTimes))
	}else{
		slog.Info(logName, fmt.Sprintf("End request %s by native http. Data length is: %d bytes.", link.URL.String(), len(ret)))
	}
	if resp != nil{
		resp.Body.Close()
	}
	return ret
}

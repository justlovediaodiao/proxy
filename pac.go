package proxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func isValid(line string) bool {
	return line != "" && strings.Index(line, "!") != 0 && strings.Index(line, "[") != 0
}

func toValidList(list []string) []string {
	var result = make([]string, 0, len(list))
	for _, line := range list {
		if isValid(line) {
			result = append(result, line)
		}
	}
	return result
}

func merge(gfw string, proxyURL string) (string, error) {
	abp, err := os.ReadFile("resources/abp.js")
	if err != nil {
		return "", err
	}
	rule, err := os.ReadFile("resources/user-rule.txt")
	if err != nil {
		return "", err
	}
	var ruleList = strings.Split(string(rule), "\n")
	var gfwList = strings.Split(gfw, "\n")
	gfwText, err := json.MarshalIndent(toValidList(gfwList), "", "    ")
	if err != nil {
		return "", err
	}
	ruleText, err := json.MarshalIndent(toValidList(ruleList), "", "    ")
	if err != nil {
		return "", err
	}
	var text = string(abp)
	text = strings.ReplaceAll(text, "__USERRULES__", string(ruleText))
	text = strings.ReplaceAll(text, "__RULES__", string(gfwText))
	text = strings.ReplaceAll(text, "__PROXY__", proxyURL)
	return text, nil
}

func getGfwList(c *Config) ([]byte, error) {
	// if failed, use proxy
	res, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	if err != nil {
		if c.Protocol == "socks5" || c.Protocol == "http" {
			var proxyURL = fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
			var transport = http.DefaultTransport.(*http.Transport)
			transport.Proxy = func(*http.Request) (*url.URL, error) {
				return url.Parse(proxyURL)
			}
			res, err = http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
			transport.Proxy = nil
		}
	}

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	content, err = base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return nil, err
	}
	return content, nil
}

// UpdatePAC update pac file.
func UpdatePAC(c *Config) error {
	content, err := getGfwList(c)
	if err != nil {
		fmt.Println("get online gfwlist failed, use local instead")
		content, err = os.ReadFile("resources/gfwlist.txt")
		if err != nil {
			return err
		}
	} else {
		err = os.WriteFile("resources/gfwlist.txt", content, 0644)
		if err != nil {
			return err
		}
	}
	abp, err := merge(string(content), c.ProxyURL)
	if err != nil {
		return err
	}
	err = os.WriteFile("resources/pac.js", []byte(abp), 0644)
	if err != nil {
		return err
	}
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("resources/pac.js")
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
	} else {
		w.Header().Set("Content-Type", "text/plain; chart=utf-8")
		w.Write(content)
	}
}

// StartServer start a http server serve for pac file.
func StartServer(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(handle))
}

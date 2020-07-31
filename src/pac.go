package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
	abp, err := ioutil.ReadFile("resources/abp.js")
	if err != nil {
		return "", err
	}
	rule, err := ioutil.ReadFile("resources/user-rule.txt")
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

func getGfwList(p *Proxy) ([]byte, error) {
	// if failed, use proxy
	res, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	if err != nil {
		if p.Protocol == "socks5" || p.Protocol == "http" {
			var proxyURL = fmt.Sprintf("%s://%s:%d", p.Protocol, p.Host, p.Port)
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
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	content, err = base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return nil, err
	}
	return content, nil
}

func UpdatePAC(p *Proxy) error {
	content, err := getGfwList(p)
	if err != nil {
		fmt.Println("get online gfwlist failed, use local instead")
		content, err = ioutil.ReadFile("resources/gfwlist.txt")
		if err != nil {
			return err
		}
	} else {
		err = ioutil.WriteFile("resources/gfwlist.txt", content, 0644)
		if err != nil {
			return err
		}
	}
	abp, err := merge(string(content), p.URL())
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("resources/pac.js", []byte(abp), 0644)
	if err != nil {
		return err
	}
	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("resources/pac.js")
	if err != nil {
		w.WriteHeader(500)
		w.Write(nil)
	} else {
		w.Header().Set("Content-Type", "text/plain; chart=utf-8")
		w.Write(content)
	}
}

func StartServer(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(handle))
}

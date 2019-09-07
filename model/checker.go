package model

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func TCPCheck(params *healthCheckParams) error {
	d := net.Dialer{Timeout: params.Timeout}
	l4 := fmt.Sprintf("%s:%d", params.Addr, params.Port)
	conn, err := d.Dial("tcp", l4)

	if err != nil {
		return err
	}

	defer conn.Close()
	return nil
}

func HTTPCheck(params *healthCheckParams, protocol string) error {
	url := fmt.Sprintf("%s://%s:%d/%s", protocol, params.Addr, params.Port, params.Path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if params.HostName != "" {
		req.Header.Set("Host", params.HostName)
	}

	client := &http.Client{Timeout: params.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if params.SearchWord != "" && 0 > strings.Index(string(body), params.SearchWord) {
		return fmt.Errorf("the search word does not match the response body")
	}

	return nil
}

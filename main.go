package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"sync"
	"time"
)

//func main() {
//	jar, _ := cookiejar.New(nil)
//	c := &http.Client{
//		Jar: jar,
//	}
//
//	req, _ := http.NewRequest(http.MethodGet, "https://w.wallhaven.cc/full/g8/wallhaven-g866qq.jpg", nil)
//
//	req.Header.Set("Range", "bytes=0-9,10-11")
//
//	resp, err := c.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	for k, v := range resp.Header {
//		fmt.Printf("k: %s, v: %v\n", k, v)
//	}
//
//	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
//	if err != nil {
//		panic(err)
//	}
//
//	boundary := params["boundary"]
//	reader := multipart.NewReader(resp.Body, boundary)
//	part, err := reader.NextPart()
//	if err != nil {
//		panic(err)
//	}
//
//	data, _ := ioutil.ReadAll(part)
//	fmt.Printf("len: %d, content: %v\n", len(data), data)
//}

func main() {
	length := ContentLen(downloaderURL)
	fmt.Printf("length: %d\n", length)

	partLen := 1024 * 1024 // 1MB

	wg := &sync.WaitGroup{}
	startTime := time.Now()
	for start, end := 0, 0; start < length; start += partLen {
		if length-start >= partLen {
			end = start + partLen - 1
		} else {
			end = length - 1
		}
		fmt.Printf("start: %d, end: %d\n", start, end)

		wg.Add(1)
		go func(start, end int) {
			DownloadPart(start, end)
			wg.Done()
		}(start, end)
	}

	wg.Wait()
	endTime := time.Now()

	fmt.Printf("dur: %f\n", endTime.Sub(startTime).Seconds())
	fmt.Printf("speed: %fKB/s\n", float64(length)/endTime.Sub(startTime).Seconds()/1024)

	fmt.Printf("--------------")
	time.Sleep(time.Second)
	Download()
}

//const downloaderURL = "https://w.wallhaven.cc/full/g8/wallhaven-g866qq.jpg"
const downloaderURL = "http://ppe.oss-cn-shenzhen.aliyuncs.com/collections/182/2/full_res.jpg"

var (
	client = newClient()
)

func newReq(start, end int) *http.Request {
	req, err := http.NewRequest(http.MethodGet, downloaderURL, nil)
	if err != nil {
		panic(err)
	}

	if start != -1 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	}
	return req
}

func newClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	tp := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxConnsPerHost:       10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c := &http.Client{
		Jar:       jar,
		Transport: tp,
	}
	return c
}

func DownloadPart(start, end int) []byte {
	req := newReq(start, end)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func readPart(resp *http.Response) []byte {
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}

	boundary := params["boundary"]
	reader := multipart.NewReader(resp.Body, boundary)
	part, err := reader.NextPart()
	if err != nil {
		panic(err)
	}

	data, _ := ioutil.ReadAll(part)
	return data
}

func ContentLen(url string) int {
	resp, err := client.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	v := resp.Header.Get("Content-Length")
	res, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	for k, v := range resp.Header {
		fmt.Printf("k: %s, v: %v\n", k, v)
	}
	return res
}

func Download() {
	c := newClient()
	resp, err := c.Get(downloaderURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	start := time.Now()
	data, _ := ioutil.ReadAll(resp.Body)
	end := time.Now()
	fmt.Printf("%s\n", data[:50])
	fmt.Printf("dur: %f\n", end.Sub(start).Seconds())
	fmt.Printf("speed: %fKB/s", float64(len(data))/end.Sub(start).Seconds()/1024)
}

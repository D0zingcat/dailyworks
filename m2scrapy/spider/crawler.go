package spider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/d0zingcat/go-logger/logger"
)

var pagesCount int
var failedUrls []string
var mu *sync.Mutex = &sync.Mutex{}

func init() {
	logger.SetRollingFile(".", "spider.log", 10, 50, logger.MB)
	logger.SetLevel(logger.DEBUG)
	htmlBytes, err := reqPage(HOME_URL)
	if err != nil {
		logger.Error("Get total page count failed!")
		panic(err)
	}
	re := regexp.MustCompile(`<span aria-current='page' class='page-numbers current'>(\d+)</span>`)
	pagesMatch := re.FindAllStringSubmatch(string(htmlBytes), -1)
	if len(pagesMatch) > 0 && len(pagesMatch[0]) > 1 {
		page := pagesMatch[0][1]
		pagesCount, err = strconv.Atoi(page)
		if err != nil {
			logger.Error("Can not convert page number")
			panic(err)
		}
	}
}

func Process(n int, dir string) {
	count := pagesCount
	flag := make([]int, pagesCount+1)
	ch := make(chan int)
	i := 1
	for ; i <= pagesCount-n; i += n {
		go dispatch(i, i+n, ch, dir)
	}
	go dispatch(i, pagesCount+1, ch, dir)
	for count > 0 {
		flag[<-ch] = 1
		count--
	}
	logger.Info("Fail to get these urls: ", failedUrls)
}

func dispatch(start, end int, ch chan int, dir string) {
	for i := start; i < end; i++ {
		dynUrl := fmt.Sprintf(TEMPLATE_URL, i)
		content, err := reqPage(dynUrl)
		if err != nil {
			logger.Error("Req ", dynUrl, " error!")
		}
		content = strings.Replace(content, "\r\n", "", -1)
		content = strings.Replace(content, "\r", "", -1)
		content = strings.Replace(content, "\n", "", -1)
		re := regexp.MustCompile(`<li class="comment byuser(.*?</li>)`)
		comments := re.FindAllString(content, -1)
		for _, item := range comments {
			re := regexp.MustCompile(`<img src="(.+?)".*?/>`)
			imgs := re.FindAllStringSubmatch(item, -1)
			err := storePic(imgs[0][1], dir, strconv.Itoa(i))
			if err != nil {
				// logger.Error(err)
				continue
			}
		}
		ch <- i
	}
}

func storePic(url, location, prefix string) error {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		err = os.Mkdir(location, 0744)
		if err != nil {
			logger.Error("create dir failed!")
			return fmt.Errorf("Dir create fail")
		}
	}
	ss := strings.Split(url, "/")
	filename := ss[len(ss)-1]
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Fail to request the pic: ", url)
		failedUrls = conAppendSlice(failedUrls, url)
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Fail to read pic response: ", url)
		failedUrls = conAppendSlice(failedUrls, url)
		return err
	}

	err = ioutil.WriteFile(filepath.Join(location, prefix+"-"+filename), bodyBytes, 0744)
	if err != nil {
		logger.Error("Store pic failed: ", url)
		failedUrls = conAppendSlice(failedUrls, url)
		return err
	}
	return nil
}

func conAppendSlice(s []string, e string) []string {
	mu.Lock()
	s = append(s, e)
	mu.Unlock()
	return s
}
func reqPage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Fail to request the page")
		return "", err
	}
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Fail to read response")
		return "", err
	}
	return string(htmlBytes), nil
}

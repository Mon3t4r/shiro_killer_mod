package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func GetCommandArgs() {
	flag.StringVar(&UserAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36", "User-Agent")
	flag.StringVar(&UrlFile, "f", "", "url文件 例如:url.txt")
	flag.StringVar(&Url, "u", "", "单个url目标 例如 http://hackforfun.com")
	flag.StringVar(&Header, "header", "", "数据包中添加header信息,如有多个用,连接 如 Abc:123,BCD:456")
	flag.StringVar(&Method, "m", "GET", "Request Method")
	flag.StringVar(&PostContent, "content", "", "POST Method Content")
	flag.IntVar(&Timeout, "timeout", 3, "Request timeout time(s)")
	flag.IntVar(&Interval, "interval", 0, "Each request interval time(s)")
	flag.StringVar(&HttpProxy, "proxy", "", "Set up http proxy e.g. http://127.0.0.1:8080")
	flag.StringVar(&SKey, "k", "", "指定密钥文件")
	flag.StringVar(&AesMode, "mode", "", "Specify CBC or GCM encryption mode")
	flag.IntVar(&t, "t", 50, "Number of goroutines")
	flag.StringVar(&CheckContent, "chk", "rO0ABXNyADJvcmcuYXBhY2hlLnNoaXJvLnN1YmplY3QuU2ltcGxlUHJpbmNpcGFsQ29sbGVjdGlvbqh/WCXGowhKAwABTAAPcmVhbG1QcmluY2lwYWxzdAAPTGphdmEvdXRpbC9NYXA7eHBwdwEAeA==", "Check Content")
	flag.StringVar(&NRemeberMe, "rm", "rememberMe", "Name of rememberMe")
	flag.StringVar(&OutPutfile, "o", "", "out filename")
	flag.Parse()
}

var outchan = make(chan string)

func StartTask(TargetUrl string) {

	if !ShiroCheck(TargetUrl) {
		_, result := KeyCheck(TargetUrl)
		outchan <- fmt.Sprintln("\n", TargetUrl, ": \n", result)
	} else {
		outchan <- fmt.Sprintln(TargetUrl, ": ", "Shiro not exist!")
	}
}
func ShiroCheck(TargetUrl string) bool {
	ok, _ := HttpRequest("wotaifu", TargetUrl)
	return ok
}

const (
	BarLen = 100
)

func KeyCheck(TargetUrl string) (bool, string) {
	Content, _ := base64.StdEncoding.DecodeString(CheckContent)
	isFind, Result := false, ""

	var builder strings.Builder // create a string builder

	Keylen := len(ShiroKeys)
	isFind = false
	for i := range ShiroKeys {
		percent := float64(i+1) / float64(Keylen)
		// 计算当前需要打印多少个等号
		eqLen := int(percent * float64(BarLen))
		builder.Reset()                                         // reset the builder
		builder.WriteString("\r[")                              // write the left bracket
		builder.WriteString(strings.Repeat("=", eqLen))         // write the equal signs
		builder.WriteString(strings.Repeat(" ", BarLen-eqLen))  // write the spaces
		builder.WriteString("] ")                               // write the right bracket and a space
		builder.WriteString(fmt.Sprintf("%.2f%%", percent*100)) // write the percentage
		fmt.Print(builder.String())                             // print the builder's string
		time.Sleep(time.Duration(Interval) * time.Second)
		isFind, Result = FindTheKey(ShiroKeys[i], Content, TargetUrl)
		if isFind {
			break
		}
	}
	return isFind, Result
}
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
func main() {
	GetCommandArgs()
	if UrlFile == "" && Url == "" {
		flag.Usage()
		fmt.Println("[Error] UrlFile (-f) or Url -u must be specified.")
		os.Exit(1)
	}
	if SKey != "" { //是否读取密钥本
		KeyF, err := os.Open(SKey)
		if err != nil {
			panic(err)
		}
		defer KeyF.Close()
		krd := bufio.NewReader(KeyF)
		for {
			UnFormatted, _, err := krd.ReadLine()
			if err == io.EOF {
				break
			}
			if len(UnFormatted) > 0 {
				ShiroKeys = append(ShiroKeys, string(UnFormatted))
			}

		}
		ShiroKeys = RemoveDuplicatesAndEmpty(ShiroKeys) //key去重
	}
	//work goroutines

	var workurl = make(chan string)
	for i := 0; i < t; i++ {
		go func() {
			for url := range workurl {
				StartTask(url)
				wg.Done()
			}
		}()
	}
	//out goroutines
	var outf *os.File
	var err error
	if OutPutfile != "" {
		outf, err = os.OpenFile(OutPutfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("[Error] OutFile (-o) must be specified.")
			os.Exit(1)
		}
		defer outf.Close()
	}
	go func() {
		for outStr := range outchan {
			fmt.Print(outStr)
			if outf != nil {
				outf.WriteString(outStr)
				outf.Sync()
			}
		}
	}()

	if UrlFile != "" {

		UrlF, err := os.Open(UrlFile)
		if err != nil {
			panic(err)
		}
		defer UrlF.Close()
		rd := bufio.NewReader(UrlF)
		startTime := time.Now()

		for {
			UnFormatted, _, err := rd.ReadLine()
			if err == io.EOF {
				break
			}
			TargetUrl := string(UnFormatted)
			if !strings.Contains(TargetUrl, "http://") && !strings.Contains(TargetUrl, "https://") {
				TargetUrl = "https://" + TargetUrl
			}
			wg.Add(1)
			workurl <- TargetUrl

		}
		wg.Wait()
		endTime := time.Since(startTime)
		fmt.Println("\nDone! Time used:", int(endTime.Minutes()), "m", int(endTime.Seconds())%60, "s")
	}
	if Url != "" {
		startTime := time.Now()
		TargetUrl := string(Url)
		if !strings.Contains(TargetUrl, "http://") && !strings.Contains(TargetUrl, "https://") {
			TargetUrl = "https://" + TargetUrl
		}
		wg.Add(1)
		workurl <- TargetUrl
		wg.Wait()
		endTime := time.Since(startTime)
		fmt.Println("\nDone! Time used:", int(endTime.Minutes()), "m", int(endTime.Seconds())%60, "s")
	}
}

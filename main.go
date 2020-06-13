package main

import (
	"fmt"
	"bufio"
	"flag"
	"os"
	"sync"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
	"strings"
	"io/ioutil"
	"regexp"
	)

func init() {
	flag.Usage = func() {
		h := []string{
			"",
			"Urlive (Check url is live *HTTP status code \"200 ok\" only)",
			"",
			"By : viloid [Sec7or - Surabaya Hacker Link]",
			"",
			"Basic Usage :",
			" ▶ echo http://domain.com/path/file.extparam=value | urlive",
			" ▶ cat listurls.txt | urlive -c 50",
			"",
			"Options :",
			"  -H, --header        Add Header to the request",
			"  -c, --concurrency   Increase concurrency level (*default 20)",			
			"  -x, --proxy         Add HTTP proxy",
			"  -m, --match         Add match specific string (*Sensitive Case)",
			"  -o, --output        Output to file",
			"",
			"",
		}
		fmt.Fprintf(os.Stderr, strings.Join(h, "\n"))
	}
}

func main() {

	var headers headerArgs
	flag.Var(&headers, "header", "")
	flag.Var(&headers, "H", "")

	var concurrency int
	flag.IntVar(&concurrency, "concurrency", 20, "Concurrency level")
	flag.IntVar(&concurrency, "c", 20, "Concurrency level")

	var proxy string
	flag.StringVar(&proxy, "proxy", "", "")
	flag.StringVar(&proxy, "x", "", "")

	var match string
	flag.StringVar(&match, "match", "", "")
	flag.StringVar(&match, "m", "", "")

	var outputFile string
	flag.StringVar(&outputFile, "output", "", "")
	flag.StringVar(&outputFile, "o", "", "")

	flag.Parse()

	client := newClient(proxy)	
	jobs := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			for u := range jobs {

				req, err := http.NewRequest("GET", u, nil)
				if err != nil {
					return
				}
				if headers == nil {
					req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; urlive/0.1; +https://github.com/vsec7/urlive)")
				}
				
				// add headers to the request
				for _, h := range headers {
					parts := strings.SplitN(h, ":", 2)

					if len(parts) != 2 {
						continue
					}
					req.Header.Set(parts[0], parts[1])
				}

				// send the request
				resp, err := client.Do(req)
				if err != nil {
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {

					if match != "" {
						data, _ := ioutil.ReadAll(resp.Body)
						var re, _ = regexp.Compile(match)
						
						if re.MatchString(string(data)) == true {
							if outputFile != "" {
								file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
								if err != nil {
									fmt.Printf("[!] Failed Creating File: %s", err)
								}
								buf := bufio.NewWriter(file)
								buf.WriteString(u+"\n")
								buf.Flush()
								file.Close()
							}
							fmt.Printf("%s\n", u)
						}						
					} else {
						if outputFile != "" {
							file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
							if err != nil {
								fmt.Printf("[!] Failed Creating File: %s", err)
							}
							buf := bufio.NewWriter(file)
							buf.WriteString(u+"\n")
							buf.Flush()
							file.Close()
						}
						fmt.Printf("%s\n", u)
					}
				}
			}
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		jobs <- sc.Text()
	}
	close(jobs)
	wg.Wait()
}


func newClient(proxy string) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:		30,
		IdleConnTimeout:	time.Second,
		TLSClientConfig:	&tls.Config{InsecureSkipVerify: true},
		DialContext:		(&net.Dialer{
			Timeout:	time.Second * 5,
		}).DialContext,
	}

	if proxy != "" {
		if p, err := url.Parse(proxy); err == nil {
			tr.Proxy = http.ProxyURL(p)
		}
	}

	re := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &http.Client{
		Transport:		tr,
		CheckRedirect: 	re,
		Timeout:		time.Second * 5,
	}
}

type headerArgs []string

func (h *headerArgs) Set(val string) error {
	*h = append(*h, val)
	return nil
}

func (h headerArgs) String() string {
	return "string"
}
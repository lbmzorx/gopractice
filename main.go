// HttpGet project main.go
package main

import (
	"encoding/json"
	//	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func num(ch chan int) {
	num := <-ch
	num = num + 1
	ch <- num
}

func Count(ch chan int) {
	ch <- 1
}

type RespRes struct {
	status int32  `json:"status"`
	msg    string `json:"msg"`
}

func main() {

	var connection = 100
	//	var method = "GET"
	var user_id string
	var url = "http://www.foxdouatt.com/grab-product/grab?grab_id=33&user_id="

	timestamp := time.Now()

	chs := make([]chan int, connection)

	chnum := make(chan int, 2)

	chnum <- 0

	//	timeout := flag.Int("o", 5, "-o N/*")
	//	t := time.Duration(*timeout) * time.Second
	//	Client := http.Client{Timeout: t}*/

	for i := 0; i < connection; i++ {
		fmt.Println("start...", i)
		chs[i] = make(chan int)
		go func(ch chan int, step int, chnum chan int) {

			user_id = strconv.FormatInt(timestamp.Unix(), 10)

			//			req, _ := http.NewRequest(method, url+user_id+string(step), nil)

			//			resp, err := Client.Do(req)

			fmt.Println(url + user_id + strconv.Itoa(step+1))
			resp, err := http.Get(url + user_id + strconv.Itoa(step+1))
			if err != nil {
				fmt.Println(step, "|Failed.", err)
				Count(ch)
				return
			}
			defer resp.Body.Close()

			body, err1 := ioutil.ReadAll(resp.Body)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			fmt.Println(step, '|', string(body))
			if resp.StatusCode != 200 {

				fmt.Println("请求失败")
			}

			//			stb := &RespRes{}
			var dat map[string]interface{}
			err = json.Unmarshal([]byte(body), &dat)
			if err != nil {
				fmt.Println("Unmarshal faild")
			} else {
				fmt.Println("status:", dat["status"], "msg", dat["msg"])
				if dat["status"] == true {
					num(chnum)
				}
			}

			Count(ch)

		}(chs[i], i, chnum)
	}
	for _, ch := range chs {
		<-ch
	}

	countnum := <-chnum

	tend := time.Now()

	fmt.Println("count:", countnum)
	fmt.Println("user time:", tend.Sub(timestamp))

	//	fmt.Println("Hello World!")
}

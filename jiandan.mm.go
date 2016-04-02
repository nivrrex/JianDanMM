package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func get_mm_url(url string, user_agent string) (html string, err error) {
	client := new(http.Client)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log := fmt.Sprintf("Http Get url Error : %s   %s\n", err, url)
		fmt.Print(log)
		return
	}
	req.Header.Add("User-Agent", user_agent)
	resp, err := client.Do(req)
	if err != nil {
		log := fmt.Sprintf("Http Get url Error : %s   %s\n", err, url)
		fmt.Print(log)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	html = string(body)

	return html, err
}

func get_file_name(url string) string {
	re := strings.Replace(url, "https://", "", -1)
	re = strings.Replace(re, "http://", "", -1)
	re = strings.Replace(re, "/", "-", -1)
	re = strings.Replace(re, "?", "", -1)
	if re != "" {
		return re
	} else {
		return ""
	}
}

func get_scale(scale float64) (level string) {
	if scale >= 10 {
		level = "very"
	} else if scale >= 5 {
		level = "high"
	} else if scale >= 1.5 {
		level = "good"
	} else if scale >= 0.75 {
		level = "normal"
	} else if scale < 0.75 {
		level = "bad"
	}
	return level
}

func get_jpg(url string, file_name string, id string, support string, unsupport string, rank string) (log string) {
	resp, err := http.Get(url)
	if err != nil {
		log = fmt.Sprintf("Http get jpg Error:%s %s %s\n", id, err, url)
		fmt.Print(log)
		// handle error
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	f, err := os.Create(file_name)
	if err != nil {
		log = fmt.Sprintf("Write jpg Error:%s %s %s\n", id, err, url)
		fmt.Print(log)
		return
	}
	f.WriteString(string(body))
	f.Close()

	log = fmt.Sprintf("Write OK ... %s %s %s %s %s %s\n", id, url, file_name, support, unsupport, rank)
	fmt.Print(log)

	return log
}

type ConfigJson struct {
	StartPage string
	EndPage   string
}

func main() {
	var configJson ConfigJson
	execDirAbsPath, _ := os.Getwd()
	f_config, err := os.Open(execDirAbsPath + "\\config.json")
	if err != nil {
		fmt.Printf("Reading config.json is error ... \n%s", err)
		// handle error
	}
	input, err := ioutil.ReadAll(f_config)
	if err != nil {
		fmt.Printf("Reading config.json is error ... \n%s", err)
		// handle error
	}
	err = json.Unmarshal(input, &configJson)
	if err != nil {
		fmt.Printf("Reading config.json is error ... \n%s", err)
		// handle error
	}
	f_config.Close()

	startpage, _ := strconv.Atoi(configJson.StartPage)
	endpage, _ := strconv.Atoi(configJson.EndPage)

	url_hash := make(map[string]map[string]string)
	f_log, err := os.Create("log_result.log")
	if err != nil {
		// handle error
	}
	log.SetOutput(f_log)

	for i := startpage; i < endpage+1; i++ {
		body, _ := get_mm_url("http://jandan.net/ooxx/page-"+strconv.Itoa(i), "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.6.2000 Chrome/30.0.1599.101 Safari/537.37"+strconv.Itoa(i))

		fmt.Println()
		log_str := "reading " + "http://jandan.net/ooxx/page-" + strconv.Itoa(i)
		fmt.Println(log_str)
		log.Print(log_str)
		fmt.Println()

		reg := `(?m)comment-(\d*?)">\d*?<\/a><\/span><p>.*?\s*<img src="([a-zA-z]+:\/\/[^\s]+?)" \/><\/p>`
		regCom, _ := regexp.Compile(reg)
		search := regCom.FindAllStringSubmatch(body, -1)
		if search != nil {
			//fmt.Println("%s",search)
			for _, value := range search {
				id := value[1]
				url := value[2]
				file_name := get_file_name(url)
				_, ok := url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["URL"] = url
				url_hash[id]["FileName"] = file_name
				//fmt.Println(id,"   ",value[2],"   ",file_name)
			}
		} else {
			//fmt.Println("%v",line)
		}

		reg = `(?m)comment-(\d*?)">\d*?<\/a><\/span><p>.*?\s*<img src="[a-zA-z]+:\/\/[^\s]+?" org_src="([a-zA-z]+:\/\/[^\s]+?)" o`
		regCom, _ = regexp.Compile(reg)
		search = regCom.FindAllStringSubmatch(body, -1)
		if search != nil {
			//fmt.Println("%s",search)
			for _, value := range search {
				id := value[1]
				url := value[2]
				file_name := get_file_name(url)
				_, ok := url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["URL"] = url
				url_hash[id]["FileName"] = file_name
				//fmt.Println(id,"   ",value[2],"   ",file_name)
			}
		} else {
			//fmt.Println("%v",line)
		}

		reg = `<div class="vote" id="vote-(\d*?)">.*?<span id="cos_support-\d*?">(\d*?)<\/span>.*?<span id="cos_unsupport-\d*?">(\d*?)<\/span>`
		regCom, _ = regexp.Compile(reg)
		search = regCom.FindAllStringSubmatch(body, -1)
		if search != nil {
			//fmt.Println("%s",search)
			for _, value := range search {
				id := value[1]
				support := value[2]
				unsupport := value[3]
				support_i, _ := strconv.ParseFloat(support, 64)
				unsupport_i, _ := strconv.ParseFloat(unsupport, 64)
				scale := float64(support_i / unsupport_i)
				rank := get_scale(scale)

				_, ok := url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["Support"] = support
				url_hash[id]["UnSupport"] = unsupport
				url_hash[id]["Scale"] = strconv.FormatFloat(scale, 'f', 3, 64)
				url_hash[id]["Rank"] = rank
				url_hash[id]["FileName"] = rank + "---" + url_hash[id]["FileName"]

				//write the picture
				log_str := get_jpg(url_hash[id]["URL"], url_hash[id]["FileName"], id, support, unsupport, rank)
				log.Print(log_str)

				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				time.Sleep(time.Duration(100+r.Intn(200)) * time.Millisecond)
				//fmt.Println(id,"   ",support,"   ",unsupport)
			}
		} else {
			//fmt.Println("%v",line)
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(1800+r.Intn(1500)) * time.Millisecond)
	}
	f_log.Close()

	//for k, v := range url_hash {
	//     for i, j := range v {
	//			fmt.Println(k, " " , i , " " , j)
	//          fmt.Println(k, " " , i , " " , j)
	//     }
	//    fmt.Println()
	//}

}

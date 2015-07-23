package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"regexp"
	"strings"
	"strconv"
	"os"
	"time"
	"math/rand"
	"log"
)


func get_mm_url(url string,user_agent string) (html string,err error){
	client := new(http.Client)
	req, err := http.NewRequest("GET", url , nil)	
	if err != nil {
		// handle error
	}
	req.Header.Add("User-Agent", user_agent)
	
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)	

	html = string(body)
	return html,err
}

func get_file_name(url string)string{
    re := strings.Replace(url,"https://","",-1)
    re =  strings.Replace(re,"http://","",-1)
    re =  strings.Replace(re,"/","-",-1)
    re =  strings.Replace(re,"?","",-1)
    if re != "" {
        return re
	}else{
        return ""
    }
}

func get_scale(scale float64)(level string){
    if scale >= 10 {
        level = "very"
	}else if scale >= 5 {
        level = "high"
	}else if scale >= 1.5 {
        level = "good"
	}else if scale >= 0.75 {
        level = "normal"
	}else if scale < 0.75 {
        level = "bad"
	}
    return level
}

func get_jpg(url string , file_name string,id string)(log string){
	resp, err := http.Get(url)
	if err != nil {
		log = fmt.Sprintf("Http Error:%s %s %s\n", id , err , url)
		fmt.Print("Http Error:",id,err,url)
		// handle error
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	
    f,err:=os.Create(file_name)
    if err != nil {
		log = fmt.Sprintf("Write Error:%s %s %s\n" , id , err , url)
		fmt.Print(log)
		return
    }
    f.WriteString(string(body))
	log = fmt.Sprintf("Write OK ... %s %s\n" , id , url)
	fmt.Print(log)

	defer resp.Body.Close()
    defer f.Close()
	return log
}


func main() {

	url_hash := make(map[string]map[string]string)
	fo, err := os.Create("log_result.log")
    if err != nil {
		// handle error
    }
    log.SetOutput(fo)

	for i := 900; i < 1475; i++ {
		body , _ := get_mm_url("http://jandan.net/ooxx/page-" + strconv.Itoa(i),"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.6.2000 Chrome/30.0.1599.101 Safari/537.37" + strconv.Itoa(i))

		fmt.Println()
		log_str := "reading " + "http://jandan.net/ooxx/page-" + strconv.Itoa(i)
		fmt.Println(log_str)
		log.Print(log_str)
		fmt.Println()
		
		reg := `(?m)comment-(\d*?)">\d*?</a></span><p>.*?\s*<img src="([a-zA-z]+://[^\s]+?)" /></p>`
		regCom, _ := regexp.Compile(reg)
		search := regCom.FindAllStringSubmatch(body, -1)
		if search != nil {
			//fmt.Println("%s",search)
			for _, value := range search {
				id := value[1]
				url := value[2]
				file_name := get_file_name(url)
				_,ok :=  url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["URL"] = url
				url_hash[id]["FileName"] = file_name
				//fmt.Println(id,"   ",value[2],"   ",file_name)
			}
		} else {
			//fmt.Printf("%v",line)
		}

		reg = `(?m)comment-(\d*?)">\d*?</a></span><p>.*?\s*<img src="[a-zA-z]+://[^\s]+?" org_src="([a-zA-z]+://[^\s]+?)" o`
		regCom, _ = regexp.Compile(reg)
		search = regCom.FindAllStringSubmatch(body, -1)
		if search != nil {
			//fmt.Println("%s",search)
			for _, value := range search {
				id := value[1]
				url := value[2]
				file_name := get_file_name(url)
				_,ok :=  url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["URL"] = url
				url_hash[id]["FileName"] = file_name
				//fmt.Println(id,"   ",value[2],"   ",file_name)
			}
		} else {
			//fmt.Printf("%v",line)
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
				support_i,_ := strconv.ParseFloat(support,64)
				unsupport_i,_ := strconv.ParseFloat(unsupport,64)
				scale :=  float64(support_i/unsupport_i)
				rank := get_scale(scale)

				_,ok :=  url_hash[id]
				if !ok {
					url_hash[id] = make(map[string]string)
				}
				url_hash[id]["Support"] = support
				url_hash[id]["UnSupport"] = unsupport
				url_hash[id]["Scale"]  = strconv.FormatFloat(scale, 'f' , 3 ,64)
				url_hash[id]["Rank"]  = rank
				url_hash[id]["FileName"]  = rank + "---" + url_hash[id]["FileName"]
				
				//write the picture
				log_str := get_jpg(url_hash[id]["URL"],url_hash[id]["FileName"],id)
				log.Print(log_str)

				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				time.Sleep(time.Duration(100 + r.Intn(200)) * time.Millisecond )
				//fmt.Println(id,"   ",support,"   ",unsupport)
			}
		} else {
			//fmt.Printf("%v",line)
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(1800 + r.Intn(1500)) * time.Millisecond)
	}
	
	//for k, v := range url_hash {
	//     for i, j := range v {
	//			fmt.Print(k, " " , i , " " , j)
	//     }
	//	  fmt.Println()
	//}
	defer fo.Close()

}
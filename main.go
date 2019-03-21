// 码云 webhook

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Conf struct {
	ListenPort string //监听端口
	Repos      []Repo
}

type Repo struct {
	RepoName string //仓库名
	WorkPath string //项目路径
	Secret   string //webhook 密钥
	Command  string //go 程序需要的额外命令：1、删除旧程序2、编译新程序3、重启程序
}

//码云推送的post数据
type PostBody struct {
	Password string
}

var conf Conf

//匹配的仓库
var matchRepo Repo

func init() {
	content, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalln("解析json错误：", err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	url := strings.Trim(r.URL.Path, "/")
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, item := range conf.Repos {
		if item.RepoName == url {
			matchRepo = item
			break
		}
	}
	if matchRepo.RepoName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//检查密码
	var postBody PostBody
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &postBody)
	if err != nil {
		_, _ = fmt.Fprintln(w, "解析post数据错误：", err)
		return
	}

	if postBody.Password != matchRepo.Secret {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//切换到代码仓库跟路径
	err = os.Chdir(matchRepo.WorkPath)
	if err != nil {
		_, _ = fmt.Fprintln(w, "切换路径错误：", err)
		return
	}
	//git pull
	out, err := exec.Command("bash", "-c", "git pull 2>&1").Output()
	if err != nil {
		_, _ = fmt.Fprintln(w, "git pull 执行错误："+err.Error())
		return
	}
	_, _ = fmt.Fprintln(w, string(out))

	///--------------- php 语言仓库到此结束---------------
	if matchRepo.Command == "" {
		return
	}

	//go 程序需要的额外命令：1、删除旧程序2、编译新程序3、重启程序
	out, err = exec.Command("bash", "-c", matchRepo.Command).Output()
	if err != nil {
		_, _ = fmt.Fprintln(w, "bash执行错误："+err.Error())
		return
	}
}

func main() {
	http.HandleFunc("/", handle)

	log.Fatal(http.ListenAndServe(conf.ListenPort, nil))
}

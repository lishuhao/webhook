// 码云 webhook

package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
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
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var conf Conf
var WorkPath = ""
var Secret = ""
var MatchRepo = false //是否匹配

func init() {
	content, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(content, &conf)
}

func main() {
	http.HandleFunc("/", handle)

	log.Fatal(http.ListenAndServe(conf.ListenPort, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	url := strings.Trim(r.URL.Path, "/")
	if url == "" {
		fmt.Fprintln(w, "404")
		return
	}

	for _, item := range conf.Repos {
		if item.RepoName == url {
			MatchRepo = true
			WorkPath = item.WorkPath
			Secret = item.Secret
		}
		break
	}
	if !MatchRepo {
		fmt.Fprintln(w, "404")
		return
	}

	//检查密码
	body, _ := ioutil.ReadAll(r.Body)
	pwd := jsoniter.Get(body, "password").ToString()
	if pwd != Secret {
		fmt.Fprintln(w, "密码错误")
		return
	}

	cmd := fmt.Sprintf("cd %s && git pull 2>&1", WorkPath)
	//cmd := fmt.Sprintf("whoami", WorkPath)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Fprintln(w, cmd+" 执行错误："+err.Error())
	}
	fmt.Fprintln(w, string(out))
}

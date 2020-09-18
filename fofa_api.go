package main

/*
data   :  2020.9.17
version:  1.0.0
author :  BurnyMcDull
*/

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/bitly/go-simplejson"
	"github.com/gookit/color"
)

type FofaAccount struct {
	Email string `json:"email"`
	User  string `json:"username"`
	Fcoin int    `json:"fcoin"`
}

type Flag struct {
	Sql     string // fofa_sql
	Size    int    // default number in one page
	History bool   // history
}

func Init() *Flag {
	cfg, _ := goconfig.LoadConfigFile("conf.ini")
	size, _ := cfg.Int("fofa_config", "defaultsize")
	ishistory := cfg.MustBool("fofa_config", "history")
	fofa_flag := Flag{}
	flag.StringVar(&fofa_flag.Sql, "q", "", "查询语句 (支持domain,host,ip,header,body,title，运算符支持== = != =~)")
	flag.IntVar(&fofa_flag.Size, "s", size, "每页数量")
	flag.BoolVar(&fofa_flag.History, "f", ishistory, "是否获取历史数据 （default "+strconv.FormatBool(ishistory)+")")
	return &fofa_flag
}

//验证用户身份
func verifyUser(mail string, key string) bool {
	fofaaccount := FofaAccount{}
	resp, err := http.Get("https://fofa.so/api/v1/info/my?email=" + mail + "&key=" + key)
	if err != nil {
		color.Red.Println("[Error]cant verify user")
		return false
	} else {
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				color.Red.Println("[Error]cant get the response\n")
				return false
			}
			//fmt.Printf(string(body) + "\n")
			err_json := json.Unmarshal([]byte(body), &fofaaccount)
			if err_json != nil {
				color.Red.Println("[Error]cant verify user")
			}
			color.Green.Println("[Success] verify account")
			fmt.Printf("Email:%v \nUser:%v \nFcoon:%v \n", fofaaccount.Email, fofaaccount.User, fofaaccount.Fcoin)
			return true
		} else {
			color.Red.Println("[Error]check the mail and key value ")
			return false
		}
	}
}

//验证语句
func verifySql(sql string) bool {
	if sql == "" {
		return true
	}
	return false
}

func writedata(f *os.File, data []interface{}, len_fileds int) {
	var writestr string
	for i := 0; i < len(data); i++ {
		var writestr_data string
		for m := 0; m < len_fileds-1; m++ {
			writestr_data = writestr_data + data[i].([]interface{})[m].(string) + ","
		}
		writestr = writestr_data + data[i].([]interface{})[len_fileds-1].(string) + "\n"
		_, err := f.Write([]byte(writestr))
		if err != nil {
			color.Red.Println("[Error]cant write file")
			return
		}

	}
	color.Green.Println("[Success]write file:" + f.Name() + " success")

	//fmt.Printf("%v", data[0].([]interface{})[0])
	return
}

func dump_fofa_data(sql string, mail string, key string, size int, fields string, full bool) {
	sql_input := []byte(sql)
	qbase64 := base64.StdEncoding.EncodeToString(sql_input)
	resp, err := http.Get("http://fofa.so/api/v1/search/all?email=" + mail + "&key=" + key + "&fields=" + fields + "&size=" + strconv.Itoa(size) + "&page=1&qbase64=" + qbase64 + "&full=" + strconv.FormatBool(full))
	if err != nil {
		color.Red.Println("[Error]cant get the response\n")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red.Println("[Error]cant get the response\n")
		return
	}
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		color.Red.Println("[Error]json error")
		return
	}
	total, err := js.Get("results").Array()
	if err != nil {
		color.Red.Println("[Error]cant find results，maybe dont have enough fcoin")
		return
	}
	fileds_arr := strings.Split(fields, ",")
	f, err := os.Create("./file/" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv")
	writestr := fileds_arr[0]
	for i := 1; i < len(fileds_arr); i++ {
		writestr = writestr + "," + fileds_arr[i]
	}
	writestr = writestr + "\n"
	_, err = f.Write([]byte(writestr))
	writedata(f, total, len(fileds_arr))
	f.Close()
	//fmt.Printf("%v", total[0].([]interface{}))

}

func main() {
	fofa_flag := &Flag{}
	fofa_flag = Init()
	flag.Parse()
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	mail, err := cfg.GetValue("fofa_user", "mail")
	key, err := cfg.GetValue("fofa_user", "key")
	size, _ := cfg.Int("fofa_config", "defaultsize")
	ishistory := cfg.MustBool("fofa_config", "history")
	fields, err := cfg.GetValue("fofa_config", "fields")

	if err != nil {
		color.Red.Println("[Error]wrong conf.ini")
	}
	if verifySql(fofa_flag.Sql) {
		flag.Usage()
		return
	}
	if verifyUser(mail, key) {
		//fmt.Println(fofa_flag.Sql)
		dump_fofa_data(fofa_flag.Sql, mail, key, size, fields, ishistory)
	} else {
		return
	}
}

# fofa_dump_golang

### V1.0.0



## 参考

### 重写python项目 fofa_dump，利用golang 进行重写，完成查询并写入csv文件的功能，批量化查询功能将在需要的时候继续去写。



## Usage

编译：

```go
go build fofa_api
```

```go
Usage of fofa_api:
  -f	是否获取历史数据 （default false)
  -q string
    	查询语句 (支持domain,host,ip,header,body,title，运算符支持== = != =~)
  -s int
    	每页数量 (default 10)
```

### 注意查看config.ini 文件，所有配置信息均在内。
# 用Creek以SQL方式分析日志文件

## 场景假设
我们有一个文本文件需要分析，里面有结构化、或者非结构化的数据。并且我们没有可用的BI分析工具可用，系统自带的文本工具又不能满足需求。

进一步假设，我们线上服务器有一个很大的日志文件，每一行包含了API请求信息，包括时间戳，用时，和API。格式如下：
```
2019-12-24 11:13:56.007  INFO 231375 [http-nio-8009-exec-177] --- c.b.b.p.w.logging.BceRequestIdFilter     : [681b1e43-cc08-42d6-aebe-4cc26a08558d][] [status:200,time:36ms] PUT /v1/api1
```
为了分析线上API的耗时情况，我们希望从日志中找到如下答案：

 1. 每个API每分钟的最大耗时
 2. 找到api2耗时超过100ms的所有时间点
 3. 每个API每分钟的请求量，和平均耗时

# 方案
定义一个Creek的数据源读取日志文件，分析用SQL表达，通过一个脚本动态调用creek的API将SQL转化成一个流式作业可执行文件，并且运行使之产出期望的数据。

Creek支持以grok pattern的方式匹配文本，提取感兴趣部分。我们通过定义一个grok pattern把时间戳、用时，和API分别提取出来，定义成流式作业source的输入schema。

每当用户输入一个SQL分析语句，将这个SQL和前面定义的schema构造成一个Creek的作业定义，并且请求[Creek的API](https://cloud.baidu.com/doc/RE/s/bk4di6p4c)动态生成一个可执行的流式作业程序，并且自动允许，输入为文件，输出到标准输出。

该方案不需要你安装任何软件(curl除外)，只需要有网络即可。

## 作业定义模板
准备一个作业定义的模板，定义好源的schema等信息，**文件地址**和**sql字段**用占位符表示，后面动态替换。

针对前面提到的日志格式，准备如下作业定义模板(这个模板的定义是一次性的，后面可以复用)，取名api_time_usage.json:
```
{
    "sources": [{
        "schema": {
            "format": "TXT_BY_GROK",
            "formatAttr": {
                "pattern": "%{TIMESTAMP_ISO8601:ts}.*,time:%{INT:timeuse}ms] %{WORD:method:string} %{NOTSPACE:api:string}"
            },
            "fields": [{
                "name": "ts",
                "type": "SQL_TIMESTAMP"
            }, {
                "name": "method",
                "type": "STRING"
            }, {
                "name": "api",
                "type": "STRING"
            }, {
                "name": "timeuse",
                "type": "DOUBLE"
            }]
        },
        "watermark": 0,
        "name": "t",
        "eventTime": "ts",
        "type": "FILE",
        "attr": {
            "input": "FILE_PLACE_HOLDER"
        }
    }],
    "sink": {
        "schema": {
            "format": "CSV"
        },
        "name": "mysink",
        "type": "STDOUT",
        "attr": {}
    },
    "name": "demojob",
    "timeType": "EVENTTIME",
    "sql": "INSERT INTO mysink SQL_PLACE_HOLDER "
}

```
FILE_PLACE_HOLDER和SQL_PLACE_HOLDER为占位符，运行时替换。

## 执行脚本
我们需要一个脚本来完成整个过程：作业定义的替换、可执行文件的生成、运行，取名为fsql.sh:
```
#!/bin/bash
creekGenUrl='http://creek.baidubce.com/v1/creek/generate?flag=exe&architecture=linux-amd64'
if [ "$1" == "" ] || [ "$2" == ""  ] || [ "$3" == "" ];then
    echo "usage: fsql.sh <job.json> <file> <sql>"
    exit 0
fi
 
if [[ ! -f "$1" ]];then
  echo "oops! file does not found: $1"
fi
 
json=$(cat "$1")
json="${json/FILE_PLACE_HOLDER/$2}"
json="${json/SQL_PLACE_HOLDER/$3}"

code=$(curl --compressed -l -k -H 'Content-Type:application/json;charset=utf-8' -o creek -w %{http_code} -X POST -d "$json" $creekGenUrl)
if [[ $code -eq "200" ]];then
         chmod 775 creek
else
    echo "failed to generate creek, http code:$code"    
fi

./creek
```
脚本假设你的运行环境为linux-amd64。如果不是如此，请修改第二行的architecture参数。

## 执行
有了作业定义模板和fsql.sh，就可以开始执行了。为了演示，我们这里以[示例日志文件sample_api.log](https://github.com/ubyyj/creek_ext/blob/master/jobs/sample_api.log)为例。

### 为了统计每个API每分钟的最大耗时，我们可以这样运行:
```
./fsql.sh api_time_usage.json sample_api.log "SELECT api, max(timeuse) FROM t GROUP BY TUMBLE(rowtime, INTERVAL '1' MINUTE), api"
```
运行结果如下:
```
/v1/api1,163
/v1/api1,177
/v1/api1,152
/v1/api2,36
/v1/api2,64
/v1/api2,133
/v1/api2,42
```

### 找到api2耗时超过100ms的所有时间点，执行:
```
./fsql.sh ../jobs/api_time_usage.json ../jobs/sample_api.log "SELECT CAST(ts AS VARCHAR), api FROM t WHERE api='/v1/api2' AND timeuse > 100 "
```
CAST(ts AS VARCHAR)将时间戳转化成容易辨识的字符串。运行结果如下:
```
2019-12-24 11:33:25,/v1/api2
```

### 每个API每分钟的请求量，和平均耗时，执行:
```
./fsql.sh ../jobs/api_time_usage.json ../jobs/sample_api.log "SELECT CAST(TUMBLE_START(rowtime, INTERVAL '1' MINUTE) AS VARCHAR), api, COUNT(*), AVG(timeuse) FROM t GROUP BY TUMBLE(rowtime, INTERVAL '1' MINUTE), api"
```
运行结果如下:
```
2019-12-24 11:13:00,/v1/api1,16,90.0625
2019-12-24 11:14:00,/v1/api1,14,134.7142857143
2019-12-24 11:18:00,/v1/api1,7,113
2019-12-24 11:13:00,/v1/api2,6,33.6666666667
2019-12-24 11:14:00,/v1/api2,6,60.5
2019-12-24 11:33:00,/v1/api2,6,52.5
2019-12-24 11:51:00,/v1/api2,6,38.3333333333
```
## 总结
当我们定义好一个文本文件的schema后，借助Creek的即时构建，零依赖的特性，可以像分析DB一样分析你的文本文件。通过标准的SQL语法对文件进行各种即席分析，满足各种业务需求。


# balloons-websocket

用于构建实时应用程序的基础架构和API，balloons提供了最好的基础架构和API，以大规模地提供实时体验。向最终用户提供快速稳定的实时消息。让我们处理实时消息传递的复杂性，以便您可以专注于代码。

balloons的实时API向开发人员公开了整个balloons基础架构，从而可以轻松地以任何规模提供实时功能。无需建立和运营基础架构，而将精力集中在交付真正重要的功能上。

Simple API = happy developer :)

https://github.com/fanpei0121/balloons-websocket

## 提供的功能
1. 实现了```/httpPushMessage``` http协议推送消息
2. 实现了```/wsPushMessage```   ws协议推送消息
3. 实现了```/readMessage```     订阅频道


## 运行
本项目使用Go Mod管理依赖。
```shell
go run main.go
```


## 依赖
不依赖任何第三方服务，可以直接运行在windows, linux, macos 平台

## 使用

**只有同一个app_key并且同一个channel下才会收到消息， 订阅的readChannel必须和推送消息的writeChannel相同**

**sign生成规则：**

1.声明一个字符串，内容为"000"+当前系统时间戳

2.将Secret key（Secret key 为 .env文件的SECRET_KEY 配置参数）作为秘钥，用AES算法对字符串进行加密

3.使用Base64对加密结果进行编码，结果就是sign

4.可以请求http://127.0.0.1:3000/getSign 获取示例sign

**app_key: 应用标识**

**writeChannel：推送消息的目标通道**

**readChannel： 订阅消息的目标通道**


### 功能
1. javascript ws方式推送消息
```shell
let sign = encodeURI("VURxMi9KZk9LVW5Dd1pOeTg5cldnVjZyUS95Q3lwMHp4ZXljSE9adnN5cz0=");
let ws = new WebSocket("ws://127.0.0.1:3000/wsPushMessage?app_key=my_app&sign="+sign+"&writeChannel=channel01");
ws.send("message");
```
2. http推送消息
* 地址：http://127.0.0.1:3000/httpPushMessage
* 类型：POST
* 状态码：200
* 简介：无
* 请求接口格式：

```
├─ app_key: String (服务标识)
├─ sign: String (签名)
├─ writeChannel: String (数据推送的通道)
└─ data: String (数据)

```

* 返回接口格式：

```
├─ code: Number 
├─ data: String 
└─ msg: String 

```
3. javascript 订阅频道
```
let sign = encodeURI("VURxMi9KZk9LVW5Dd1pOeTg5cldnVjZyUS95Q3lwMHp4ZXljSE9adnN5cz0=");
let ws = new WebSocket("ws://127.0.0.1:3000/readMessage?app_key=my_app&sign="+sign+"&readChannel=channel01");
ws.onopen = function() {
    console.log("client：打开连接");
};
ws.onmessage = function(e) {
    console.log('接受到消息:' + e.data);
};
ws.onclose = function(params) {
    console.log("client：关闭连接");
};
```
# boe-backend

### 文件上传

文件上传需要先使用文件名(也可以是路径)，向预签名接口请求一个上传链接，随后向链接发起put请求，请求体为待上传的文件。

#### 文件预签名

Get `boe.vinf.top:8888/file/presign`

**request**

```json
{
    "path": "file_path"
}
```

例子file_path: `img/test.png`

**response**

```json
{
    "code": 200,
    "data": "http://boe.vinf.top:8888/asset/1/img/test.png?X-Amz-Algorithm=...",
    "message": "success"
}
```

#### 上传文件

Put`http://boe.vinf.top:8888/asset/1/img/test.png?X-Amz-Algorithm=...`

**request**

file

**response**

200



#### url说明

资源前缀`/asset/`

在后端会自动获取用户id，拼接到文件路径中，防止用户覆盖其他用户的文件

`http://boe.vinf.top:8888/asset/1/img/test.png`

所以前端可以根据base_url`http://boe.vinf.top:8888/asset/`、uid、文件路径来拼接出最后的文件访问id。

或者，直接解析预签名返回的url，去除后面的get参数后，获得访问url



### websocket json 约定

**device->backend**

#### 设备上线

设备发送hello信息

```json
{
    "type":"hello",
    "mac":"xxx"
}
```

后端返回hi，并尝试通过mac注册设备

```json
{
    "type":"hi",
    "msg":"success"
}
```

紧接着后端发送`deviceInfo`，向客户端询问设备详细

```json
{
    "type":"deviceInfo",
}
```

#### 同步设备信息

设备收到`deviceInfo`，返回`deviceInfo`

```json
{
    "type":"deviceInfo",
    "info":{
    }
}
```

后端收到`deviceInfo`会将info中信息插入/更新到数据库

#### 心跳保持

设备定时向后端发送ping信息，包括设备运行时间，当前运行的planId

```json
{
    "type":"ping",
    "runningTime": 1234,
    "planId":1
}
```

设备端收到`ping`信息，返回`pong`，并记录ping中信息

```json
{
    "type":"pong",
}
```

当设备异常掉线，后端会将ping中包含的实时性数据保存到数据库，并从内存中删除设备相关信息

#### 计划下发

**backend->device**

当用户发布计划时，后端向设备发送`planList`计划列表，删除时发送`deletePlan`

```json
{
    "type":"planList",
    "plan":[]
}
{
    "type":"deletePlan",
    "planIds":[
        "1","2"
    ]
}
```


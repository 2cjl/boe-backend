# boe-backend

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

#### 控制

##### 截图

**backend->device**

```json
{
    "type":"screenshot"
}
```

**device->backend**

```json
{
    "type":"screenshot",
    "data":"image in base64 format"
}
```

##### 亮度

**backend->device**

```json
{
    "type":"brightness"
    "data": 0.8
}
```


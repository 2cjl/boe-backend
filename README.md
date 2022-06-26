# boe-backend

#### ws json 约定

**device->backend**

```json
{
    "type":"hello",
    "mac":"xxx"
}
{
    "type":"ping",
    "runningTime": 1234,
    "planId":1
}
{
    "type":"deviceInfo",
    "info":{
    }
}
```

**backend->device**

```json
{
    "type":"hi",
    "msg":"succ"
}
{
    "type":"pong",
}
{
    "type":"syncPlan"
}
{
    "type":"deletePlan",
    "planIds":[
        "1","2"
    ]
}
```


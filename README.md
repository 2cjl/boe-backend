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
{
    "type":"syncPlan"
}
```

**backend->device**

```json
{
    "type":"hi",
    "msg":"ok"
}
{
    "type":"pong",
}
{
    "type":"planList",
    "plan":[
    	{},
		{},
    ]
}
{
    "type":"deletePlan",
    "planIds":[
        "1","2"
    ]
}
```


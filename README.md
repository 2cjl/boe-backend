# boe-backend

ws json 约定

```json
device->backend
{
    "type":"hello",
    "mac":"xxx"
}
{
    "type":"ping",
    "running_time": 1234,
    "plan_id":1
}
{
    "type":"device_info",
    "info":{
    }
}
{
    "type":"sync_plan"
}

backend->device
{
    "type":"hi",
    "msg":"ok"
}
{
    "type":"pong",
}
{
    "type":"plan_list",
    "plan":[
    	{},
		{},
    ]
}
{
    "type":"delete_plan",
    "plan_id":[
        "1","2"
    ]
}

```


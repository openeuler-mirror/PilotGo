# api协议

## agent与server通讯协议

### 通用字段描述

数据格式：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|message_id|string|消息id|y||
|message_type|int|消息类型|y||
|status|int|消息状态|y||
|data|object|具体消息数据|y||

示例：

    {
        "message_id":"xxxxxx",
        "message_type":1,
        "status":0,
        "data":{
        }
    }


### 心跳
描述：由agent主动定时发送，server
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|message_type|int|消息类型，1|y||
|agent_version|string|agent版本字符串|y||
|up_time|string|xxx格式时间字符串|y||
示例：

    {
        "message_id":"xxxxxx",
        "message_type":1,
        "status":0,
        "data":{
            "agent_version":"v1.1.1",
            "agent_uptime":"2021-06-19T09:36:23+08:00"
        }
    }

### 系统信息
描述：获取agent所在机器系统信息
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|os_name|string|os名称|y||
|os_pretty_name|string|os名称完整描述|y||
|os_id|string|os id|y||
|os_version|string|os版本|y||
|os_version_id|string|os版本id|y||
|os_arch|string|os架构|y||
|kernel_version|string|kernel版本|y||
|os_uptime|string|os启动时间|y||
示例：

    {
        "message_id":"xxxxxx",
        "message_type":1,
        "status":0,
        "data":{
            "agent_version":"v1.1.1",
            "up_time":"2021-06-19T09:36:23+08:00"
        }
    }


## PilotGo server http API 协议
### 部门管理 API
#### 新增部门
描述：向部门管理树中添加新的部门
请求方法：POST
url：/machinemanager/adddepart
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|PID|string|上一级部门ID|y||
|ParentDepart|string|上一级部门名称|y||
|Depart|string|部门名称|y||

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门信息入库成功"
}
```
#### 部门信息
描述：获取全部部门信息
请求方法：GET
url：/machinemanager/departinfo
请求参数：无


响应结果：
类型：json
字段：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|label|string|部门名称|y||
|id|int|部门id|y||
|pid|int|上级部门id|y||
|children|[]MachineNode|下一级部门节点|y||

示例：
```json
{
	"code": 200,
	"data": {
		"label": "xx",
		"id": 1,
		"pid": 0,
		"children": [
			{
				"label": "",
				"id": 2,
				"pid": 1,
				"children": null
			},
			{
				"label": "xx",
				"id": 3,
				"pid": 1,
				"children": null
			}
		]
	}
}
```
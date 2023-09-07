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
#### 修改部门名称
描述：修改部门管理树中部门的名称
请求方法：POST
url：/machinemanager/updatedepart
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门ID|y||
|DepartName|string|修改的部门名称|y||

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门更新成功"
}
```

#### 删除部门信息
描述：删除部门管理树中的指定部门
请求方法：POST
url：/machinemanager/deletedepartdata
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门ID|y||

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门删除成功"
}
```

#### 返回全部部门信息
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
#### 特定部门信息
描述：获取全部部门信息
请求方法：GET


### 用户管理 API
#### 用户登录
描述：用户登录PilotGo系统  
请求方法：POST
url：/user/login
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|email|string|登录用户账号邮箱|y||
|password|string|用户密码|y||  

请求示例：

    {
        "email":"test@qq.com",
        "password":"1234"
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|departId|int|部门id|y||
|departName|string|部门名字|y||
|roleId|string|用户角色|y||
|token|string|用户身份令牌|y||
|userType|string|用户类型|y||  

示例：

    {
        "code":200,
        "data":
            {
                "departId":1,
                "departName":"麒麟",
                "roleId":"1",
                "token":"xxxxx",
                "userType":0
            },
        "msg":"登录成功!"
    }
#### 用户退出
描述：用户退出系统  
请求方法：GET
url：/user/logout
请求参数：无  
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"退出成功!"
    }
#### 查询所有用户
描述：返回数据库中所有的用户  
请求方法：GET
url：/user/searchAll
请求参数：page、size

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|page|int|页码数|y||
|size|int|每页数量|y||

请求示例：

    user/searchAll?page=1&size=10  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|departName|string|部门名字|y||
|departPId|int|部门父id|y||
|departid|string|部门id|y||
|email|string|用户账号邮箱|y||
|id|int||y||
|phone|string|手机号|y||
|role|[]string|用户角色id|y|数组|
|userType|int|用户类型|y||
|username|string|用户名|y||  

示例：

    {
        "code":200,
        "data":[
            {
                "departName":"xxx",
                "departPId":13,
                "departid":15,
                "email":"xxx",
                "id":5,
                "phone":"xxx",
                "role":["xxx"],
                "userType":3,
                "username":"xxx"
            },
        ],
        "ok": true, 
        "page": 1, 
        "size": 10, 
        "total": 1
    }
#### 用户信息高级搜索
描述：根据用户邮箱模糊查询出符合搜索要求的用户数据  
请求方法：POST
url：/user/userSearch
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|email|string|用户账号邮箱相关字符|y|模糊查询条件|  

请求示例：

    {
        "email":"xxxx"
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|departName|string|部门名字|y||
|departPId|int|部门父id|y||
|departid|string|部门id|y||
|email|string|用户账号邮箱|y||
|id|int||y||
|phone|string|手机号|y||
|role|[]string|用户角色id|y|数组|
|userType|int|用户类型|y||
|username|string|用户名|y||  

示例：

    {
        "code":200,
        "data":[
            {
                "departName":"xxx",
                "departPId":13,
                "departid":15,
                "email":"xxx",
                "id":5,
                "phone":"xxx",
                "role":["xxx"],
                "userType":3,
                "username":"xxx"
            },
        ],
        "ok": true, 
        "page": 1, 
        "size": 10, 
        "total": 1
    }

#### 添加用户
描述：添加新用户  
请求方法：POST
url：/user/register
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|username|string|用户名|y||
|password|string|密码|y||
|phone|string|手机号|y||
|email|string|用户账号邮箱|y||
|departName|string|部门名字|y||
|departPId|int|部门父id|y||
|departid|string|部门id|y||
|role|[]string|用户角色id|y|数组|
|userType|int|用户类型|y||  

请求示例:

    {
        "username":"xxx",
        "password":"xxx",
        "phone":"xxx",
        "email":"xxx",
        "departName":"xxx",
        "departPId":13,
        "departid":15,
        "role":["xxx"],
        "userType":3
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"添加用户成功!"
    }
#### 重置密码
描述：用户重置密码  
请求方法：POST
url：/user/reset
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|email|string|用户邮箱|y||

请求示例：

    {
        "email":"xxxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"密码重置成功!"
    }
#### 删除用户
描述：删除用户  
请求方法：POST
url：/user/delete
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|email|[]string|用户账号邮箱|y|数组|

请求示例:

    {
        "code":200,
        "data":null,
        "msg":"密码重置成功!"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"用户删除成功!"
    }
#### 修改用户信息
描述：修改用户信息  
请求方法：POST
url：/user/update
请求类型：json   
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Pid|int|部门父id|n||
|id|int|部门id|n||
|phone|string|手机号|n||
|departName|string|部门名字|n||

请求示例:

    {
        "Pid":2,
        "id":1,
        "phone":"xxxx",
        "departName":"xxxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"用户信息修改成功!"
    }
#### 批量添加用户
描述：批量添加新用户  
请求方法：POST
url：/user/import
请求体：.xlsx文件  
字段描述：（上传文件中各字段必须严格按照以下顺序）  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|username|string|用户名|y||
|phone|string|手机号|y||
|email|string|用户账号邮箱|y||
|DepartName|string|部门|y||
|RoleID|string|角色|y||  

响应结果：  
类型：json  
示例：

    {
        "code":200,
        "data":null,
        "msg":"批量添加用户成功!"
    }
### 日志模块管理 API
#### 删除日志
描述：删除父日志及所有子日志  
请求方法：POST
url：/agent/delete
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|ids|[]int|父日志id|y|数组|  

请求示例:

    {
        "ids":[2,3,4]
    }
响应结果：  
类型：json  
字段描述：无    
示例：

    {
        "code":200,
        "data":null,
        "msg":"日志删除成功!"
    }

#### 父日志查询
描述：查询所有的父日志  
请求方法：GET
url：/agent/log_all
请求参数：departId、page、size

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|page|int|页码数|y||
|size|int|每页数量|y||
|departId|int|部门id|y||

请求示例：

    /agent/log_all?page=1&size=10&departId=1  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|created_at|string|日志产生时间|y||
|userName|string|操作机器用户|y||
|departName|string|用户所属部门|y||
|type|string|用户操作机器类型|y|软件包安装或者服务重启等|
|status|string|操作状态|y|成功数，操作的机器总数，成功率|  

示例：

    {
        "code":200,
        "data":[
            {
                "id":1,
                "created_at":"2022-03-18T15:27:32+08:00",
                "userName":"test@qq.com",
                "departName":"xxx",
                "type":"xxx",
                "status":"2,3,0.67"
            },
        ],
        "ok":true,
        "page":1,
        "size":10,
        "total":20
    }

#### 子日志查询
描述：查询所有的子日志  
请求方法：GET
url：/agent/logs
请求参数：id

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|条目id|y||
  
请求示例：

    /agent/logs?id=1  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|logparent_id|int|父日志的id|y||
|ip|string|机器ip|y||
|code|int|状态码|y||
|object|string|操作对象|y|服务名字或者软件包名字|
|action|string|操作类型|y|软件包安装或者服务重启等|
|message|string|状态返回消息|y||  

示例：

    {
        "code":200,
        "data":[
            {
                "id":1,
                "logparent_id":1,
                "ip":"xxx.xxx.xxx.xxx",
                "code":400,
                "object":"kernel",
                "action":"软件包安装",
                "message":"获取uuid失败"
            },
        ],
        "ok":true,
        "page":1,
        "size":10,
        "total":1
    }

### 权限角色管理 API
#### 获取登录用户权限
描述：获取登录用户权限  
请求方法：POST
url：/user/permission
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|roleId|[]int|角色id|y|数组|  

请求示例:

    {
        "roleId":[2,3,4]
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|button|[]string|登录用户拥有的权限按钮|y|数组|
|menu|[]string|登录用户拥有的权限菜单|y|数组|
|userType|string|用户类型|y||  

示例：

    {
        "code":200,
        "data":
            {
                "button":["xxx","xxxx"],
                "menu":["xxx","xxxx"],
                "userType":1
            },
        "msg":"用户权限列表"
    }

#### 角色查询
描述：查询所有角色的权限  
请求方法：GET
url：/user/roles
请求参数：无  
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|button|[]string|登录用户拥有的权限按钮|y|数组|
|menu|[]string|登录用户拥有的权限菜单|y|数组|
|type|string|用户类型|y||
|role|string|用户角色|y||
|description|string|用户信息描述|y||  

示例：

    {
        "code":200,
        "data":[
            {
                "buttons":["xxx","xxx"],
                "description":"超级管理员",
                "id":1,
                "menus":["xx","xxx"],
                "role":"超级用户",
                "type":0
            },
        ],
        "ok":true,
        "page":1,
        "size":10,
        "total":1
    }

#### 获取用户角色
描述：获取所有用户的角色  
请求方法：GET
url：/user/role
请求参数：无   
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|role|string|用户角色名称|y||
|type|int|用户类型|y||
|description|string|用户角色信息描述|y||    

示例：

    {
        "code":200,
        "data":{
            "role":[
                {
                    "ID":1,
                    "role":"超级用户",
                    "type":0,
                    "description":"超级管理员"
                },
            ]
        },
        "msg":"获取用户角色"
    }

#### 变更角色权限
描述：角色权限变更  
请求方法：POST
url：/user/roleChange
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|角色id|y||
|menus|[]string|权限菜单|y|数组|
|buttonId|[]string|权限按钮|y|数组|  

请求示例:

    {
        "id":3,
        "menus":["xxx","xxx"],
        "buttonId":["xxx","xx"]
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|角色id|y||
|menus|string|权限菜单|y||
|buttonId|string|权限按钮|y||  

示例：

    {
        "code":200,
        "data":{
            "data":{
                "ID":0,
                "menus":"xxx,xxx",
                "buttonId":"xxx,xx"
            }
        },
        "msg":"角色权限变更成功"
    }

#### 添加角色
描述：添加用户角色  
请求方法：POST
url：/user/addRole
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Role|string|角色名称|y||
|Description|string|角色描述|y||  

请求示例:

    {
        "Role":"xxx",
        "Description":"xxx"
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Role|string|角色名称|y||
|Type|int|用户类型|y||
|Description|string|角色描述|y||  

示例：

    {
        "code":200,
        "data":null,
        "msg":"新增角色成功!"
    }

#### 删除用户角色
描述：删除用户角色  
请求方法：POST
url：/user/delRole
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|角色id|y||  

请求示例:

    {
        "id":2
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"角色删除成功!"
    }

#### 编辑角色信息
描述：编辑角色信息  
请求方法：POST
url：/user/roleChange
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Role|string|角色名称|n||
|Description|string|角色描述|n||  

请求示例:

    {
        "Role":"xxx",
        "Description":"xxx"
    }
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|角色id|y||
|Role|string|角色名称|y||
|Description|string|角色描述|y||  

示例：

    {
        "code":200,
        "data":{
            "data":{
                "ID":1,
                "role":"test1",
                "description":"12342"
            }
        },
        "msg":"角色信息修改成功"
    }

### Agent机器管理 API
#### agent注册列表
描述：agent注册列表  
请求方法：GET
url：/api/agent_list
请求参数：无    
响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|IP|string|agent的ip|y||
|agent_uuid|string|agent的uuid|y||
|agent_version|string|agent版本字符串|y||  

示例：

    {
        "code":200,
        "data":[
            {
            "IP": "xxx.xxx.xxx.xxx",
            "agent_uuid": "xxxx",
            "agent_version": "v0.0.1"
            },
        ]
    }

#### 机器OS信息
描述：agent系统信息  
请求方法：GET
url：/api/os_info
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/os_info?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|IP|string|agent的ip|y||
|KernelArch|string|平台架构|y||
|KernelVersion|string|内核版本|y||
|Platform|string|系统|y||
|PlatformVersion|string|系统平台版本|y||
|Uptime|string|系统启动时间|y||  

示例：

    {
        "code":200,
        "data":{
            "os_info":{
                "IP":"xxx.xxx.xxx.xxx",
                "KernelArch":"x86_64",
                "KernelVersion":"5.10.0",
                "Platform":"openeuler",
                "PlatformVersion":"21.03",
                "Uptime":"2022年 01月 11日 星期x 10:10:10 CST"
            }
        },
        "msg":"Success"
    }

#### 机器CPU信息
描述：agent CPU信息  
请求方法：GET
url：/api/cpu_info
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/cpu_info?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|CpuNum|int|agent的cpu核数|y||
|ModelName|string|机器cpu型号|y||  

示例：

    {
        "code":200,
        "data":{
            "CPU_info":{
                "CpuNum":4,
                "ModelName":"Intel(R) Core(TM) i5"
            }
        },
        "msg":"Success"
    }

#### 机器内存信息
描述：agent 内存信息  
请求方法：GET
url：/api/memory_info
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/memory_info?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|MemFree|int|agent的空闲内存|y||
|MemTotal|int|机器内存大小|y||  

示例：

    {
        "code":200,
        "data":{
            "memory_info":{
                "MemFree":1200,
                "MemTotal":1234
            }
        },
        "msg":"Success"
    }

#### 机器服务信息
描述：agent 服务列表信息  
请求方法：GET
url：/api/service_list
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/service_list?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Active|string|服务状态|y||
|Name|string|服务名称|y||  

示例：

    {
        "code":200,
        "data":{
            "service_list":[
                {
                    "Active":"active",
                    "Name":"dev-cdrom.device"
                }]
        },
        "msg":"Success"
    }

#### 机器服务状态信息
描述：agent 服务状态信息  
请求方法：GET
url：/api/service_status
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/service_status?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Active|string|服务状态|y||
|Name|string|服务名称|y||  

示例：

    {
        "code": 200,
        "data": {
            "service_status": "inactive"
        },
        "msg": "Success"
    }

#### 机器服务启动信息
描述：agent 服务启动信息  
请求方法：POST
url：/agent/service_start
请求类型：json  

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器uuid|y||
|service|string|服务名称|y||
|userName|string|操作用户名称|y||

请求示例:

    {
        "uuid":"xxx",
        "service":"xxx",
        "userName":"xxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code": 200,
        "data": {
            "service_start": null
        },
        "msg": "Success"
    }

#### 机器服务重启信息
描述：agent 服务重启信息  
请求方法：POST
url：/agent/service_restart
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器uuid|y||
|service|string|服务名称|y||
|userName|string|操作用户名称|y||

请求示例:

    {
        "uuid":"xxx",
        "service":"xxx",
        "userName":"xxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code": 200,
        "data": {
            "service_restart": null
        },
        "msg": "Success"
    }

#### 机器服务停止信息
描述：agent 服务停止信息  
请求方法：POST
url：/agent/service_stop
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器uuid|y||
|service|string|服务名称|y||
|userName|string|操作用户名称|y||

请求示例:

    {
        "uuid":"xxx",
        "service":"xxx",
        "userName":"xxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code": 200,
        "data": {
            "service_stop": null
        },
        "msg": "Success"
    }

#### 机器磁盘使用信息
描述：agent 磁盘使用信息  
请求方法：GET
url：/api/disk_use
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/disk_use?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|device|string|磁盘分区名字|y||
|fstype|string|文件类型|y||
|path|string|挂载点|y||
|total|string|磁盘总大小|y||
|used|string|已使用大小|y||
|usedPercent|string|使用率|y||  

示例：

    {
        "code":200,
        "data":{
            "disk_use":[{
                "device":"/dev/dm-0",
                "fstype":"ext2/ext3",
                "path":"/",
                "total":"16G",
                "used":"11G",
                "usedPercent":"xx%"
            }]
        },
        "msg":"Success"
    }

#### 机器注册用户信息
描述：agent 所有注册用户信息  
请求方法：GET
url：/api/user_all
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/user_all?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Username|string|用户名字|y||
|ShellType|string|shell类型|y||
|HomeDir|string|家目录|y||  

示例：

    {
        "code":200,
        "data":{
            "user_all":[{
                "HomeDir":"/root",
                "ShellType":"/bin/bash",
                "Username":"root"
            }]
        },
        "msg":"获取机器所有用户数据成功!"
    }

#### 机器当前登录用户信息
描述：agent 当前登录用户信息  
请求方法：GET
url：/api/user_info
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/user_info?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Userid|string|用户id|y||
|Groupid|string|用户组id|y||
|GroupName|string|用户组名字|y||
|HomeDir|string|家目录|y||
|Username|string|用户名|y||  

示例：

    {
        "code": 200,
        "data": {
            "user_info": {
                "GroupName": "root",
                "Groupid": "0",
                "HomeDir": "/root",
                "Userid": "0",
                "Username": "root"
            }
        },
        "msg": "获取当前登录用户信息成功!"
    }

#### 机器网卡IO信息
描述：agent 网卡IO信息  
请求方法：GET
url：/api/net_io
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/net_io?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Name|string|网卡名字|y||
|BytesRecv|string|接收字节数|y||
|BytesSent|string|发送字节数|y||
|PacketsRecv|string|接收包|y||
|PacketsSent|string|发送包|y||  

示例：

    {
        "code":200,
        "data":{
            "net_io":[{
                "BytesRecv":55745202,
                "BytesSent":91099990,
                "Name":"ens33",
                "PacketsRecv":176816,
                "PacketsSent":475332
            }]
        },
        "msg":"Success"
    }

#### 机器网卡信息
描述：agent 网卡信息  
请求方法：GET
url：/api/net_nic
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/net_nic?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Name|string|网卡名字|y||
|IPAddr|string|IP地址|y||
|MacAddr|string|机器mac地址|y||  

示例：

    {
        "code":200,
        "data":{
            "net_nic":[{
                "IPAddr":"192.168.160.134",
                "MacAddr":"00:00:00:00:00:00",
                "Name":"ens33"
            }]
        },
        "msg":"Success"
    }

#### 机器安装所有软件包
描述：agent 所有软件包  
请求方法：GET
url：/api/rpm_all
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/rpm_all?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|rpm|string|软件包名字|y||  

示例：

    {
        "code":200,
        "data":{
            "rpm_all":[
                "kernal",
                "docker"
            ]
        },
        "msg":"Success"
    }

#### 机器某个软件包详细信息
描述：获取agent某个软件包详细信息  
请求方法：GET
url：/api/rpm_info
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/rpm_info?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Name|string|软件包名字|y||
|Version|string|软件包版本|y||
|Release|string|软件包发行版|y||
|Architecture|string|架构|y||
|Summary|string|软件包说明|y||  

示例：

    {
        "code":200,
        "data":{
            "rpm_info":{
                "Architecture":"x86_64",
                "Name":"gnupg2",
                "Release":"1.oe1",
                "Summary":"xxxxx",
                "Version":"4.4"
            }
        },
        "msg":"Success"
    }

#### 机器软件包卸载
描述：agent软件包卸载  
请求方法：POST
url：/agent/rpm_remove
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|UUIDs|[]string|机器uuid|y|批量机器卸载，数组|
|RPM|string|软件包名字|y||
|UserName|string|操作机器用户|y||  

请求示例:

    {
        "uuid":"xxx,xxxx",
        "RPM":"xxx",
        "userName":"xxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"软件包卸载成功"
    }

#### 机器软件包安装
描述：agent软件包安装  
请求方法：POST
url：/agent/rpm_install
请求类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|UUIDs|[]string|机器uuid|y|批量机器安装，数组|
|RPM|string|软件包名字|y||
|UserName|string|操作机器用户|y||  

请求示例:

    {
        "uuid":"xxx,xxxx",
        "RPM":"xxx",
        "userName":"xxx"
    }
响应结果：  
类型：json  
字段描述：无  
示例：

    {
        "code":200,
        "data":null,
        "msg":"软件包安装成功"
    }

#### 机器基础信息
描述：agent 基础信息  
请求方法：GET
url：/api/os_basic
请求参数：uuid

字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器唯一标识|y||

请求示例：

    api/os_basic?uuid=xxxxxxx  

响应结果：  
类型：json  
字段描述：  
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|IP|string|机器IP|y||
|state|int|机器状态|y|离线/在线/未分配|
|depart|string|机器所属部门|y||  

示例：

    {
        "code":200,
        "data":{
            "IP":"192.168.100.10",
            "depart":"xxx",
            "state":3
        },
        "msg":"Success"
    }

### 部门管理 API
#### 新增部门
描述：向部门管理树中添加新的部门
请求方法：POST
url：/machinemanager/adddepart
请求参数：
类型：form
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|PID|string|上一级部门ID|y||
|ParentDepart|string|上一级部门名称|y||
|Depart|string|部门名称|y||

示例：

	/machinemanager/adddepart?PID=1&ParentDepart=xx&Depart=xx


响应结果：
类型：json


示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门信息入库成功"
}
```
#### 修改部门名称
描述：修改部门管理树中部门的名称
请求方法：POST
url：/machinemanager/updatedepart
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门ID|y||
|DepartName|string|修改的部门名称|y||


示例：

```json
{
    "DepartID": 1,
    "DepartName": "xx"
}
```

响应结果：
类型：json

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门更新成功"
}
```

#### 删除部门信息
描述：删除部门管理树中的指定部门
请求方法：POST
url：/machinemanager/deletedepartdata
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门ID|y||

示例：
```json
{
    "DepartID": 1
}
```

响应结果：
类型：json

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "部门删除成功"
}
```

#### 返回全部部门信息
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
|children|[]Object|下级部门|y||


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

#### 特定部门信息
描述：获取指定部门及其子部门信息
请求方法：GET
url：/machinemanager/depart
请求参数：
类型：form
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门id|y||

示例：

	/machinemanager/depart?DepartID=1


响应结果：
类型：json
字段：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|label|string|部门名称|y||
|id|int|部门id|y||
|pid|int|上级部门id|y||
|children|[]Object|下级部门|y||

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

#### 修改部门信息
描述：修改机器的所属部门
请求方法：POST
url：/machinemanager/modifydepart
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|machineid|int|机器ID|y|
|departid|int|部门ID|y|

示例：

```json
{
    "machineid":"1",
    "departid":6
}
```

响应结果：
类型：json

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "机器部门修改成功"
}
```

#### 部门机器信息
描述：获取指定部门及其子部门信息
请求方法：GET
url：/machinemanager/depart
请求参数：
类型：form
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|DepartID|int|部门id|y||
|page|int|分页的页码|y||
|size|int|分页的当页容量|y||
|ShowSelect|bool||y||

示例：

	/machinemanager/machineinfo?DepartId=8&page=1&size=10&ShowSelect=true


响应结果：
类型：json
字段：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|机器id|y||
|departid|int|部门id|y||
|departname|string|部门名称|y||


示例：

```json
{
	"code": 200,
	"data": [
		{
			"id": 6,
			"departid": 9,
			"departname": "xx",
			"ip": "192.168.160.107",
			"uuid": "7891234566",
			"cpu": "intel i5",
			"state": "1",
			"systeminfo": "centos 6"
		},
		{
			"id": 17,
			"departid": 9,
			"departname": "xx",
			"ip": "192.168.160.118",
			"uuid": "91234567817",
			"cpu": "intel i7",
			"state": "1",
			"systeminfo": "centos 7"
		},
	],
	"page": 1,
	"size": 2,
	"total": 6
}
```

#### 未分配机器资源池
描述：未分配的机器资源池
请求方法：GET
url：/machinemanager/sourcepool
请求参数：
类型：form
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|page|int|分页的页码|y||
|size|int|分页的当页容量|y||

示例：

	/machinemanager/sourcepool?page=1&size=10

响应参数：
类型：json
字段：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|机器id|y||
|departid|int|部门id|y||
|departname|string|部门名称|y||
|ip|string|机器IP|y||
|uuid|string|机器UUID|y||
|cpu|string|机器cpu信息|y||
|state|string|机器状态|y||
|systeminfo|string|机器系统信息|y||


示例：

```json
{
	"code": 200,
	"data": [
		{
			"id": 1,
			"departid": 1,
			"departname": "麒麟",
			"ip": "192.168.160.102",
			"uuid": "2345678901",
			"cpu": "intel i3 ",
			"state": "1",
			"systeminfo": "centos 7"
		},
		{
			"id": 25,
			"departid": 1,
			"departname": "麒麟",
			"ip": "192.168.160.126",
			"uuid": "89123456725",
			"cpu": "amd64",
			"state": "1",
			"systeminfo": "kylin v10"
		},
	],
	"ok": true,
	"page": 1,
	"size": 10,
	"total": 6
}
```

#### 机器删除
描述：机器删除
请求方法：GET
url：/machinemanager/deletemachinedata
请求参数：
类型：form
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|uuid|string|机器uuid|y||

示例：

	/machinemanager/deletemachinedata?uuid=a3fffcac-8626-4b4d-badc-3aaa9310a0d9


响应参数：
类型：json

示例：

```json
{
	"code": 200,
	"data": null,
	"msg": "机器删除成功"
}
```

### 批次管理
#### 批次建立
描述：批次建立
请求方法：POST
url：/batchmanager/createbatch
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Name|string|批次名|y||
|Description|string|批次描述|y||
|Manager|string|创建人|y||
|DepartID|[]string|部门ID|y||
|DepartName|[]string|部门名|y||
|Machine|[]string|机器ID|y||


示例：

```json
{  
    "Name":"new",
    "Description":"这是一条测试信息",
    "Manager":"wh",
    "DepartID":["1"],
    "DepartName":["xx","xx"],
    "Machine": ["1","2"]
}
```

响应参数：
类型：json

示例：

```json
{
	"code": 200,
	"data": null,
	"msg": "批次入库成功"
}
```

#### 修改批次信息
描述：批次修改
请求方法：POST
url：/batchmanager/updatebatch
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|BatchID|string|批次ID|y||
|BatchName|string|批次名|y||
|Descrip|string|描述|y||

示例：

```json
{
    "BatchID": "1",
    "BatchName": "kylin",
    "Descrip": "update"
}
```

响应参数：
类型：json

示例：

```json
{
	"code": 422,
	"data": null,
	"msg": "批次修改成功"
}
```

#### 删除批次
描述：批次删除
请求方法：POST
url：/batchmanager/deletebatch
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|BatchID|string|批次ID|y||


示例：

```json
{ 
    "BatchID":["1"]
}
```

响应参数：
类型：json

示例：

```json
{
	"code": 422,
	"data": null,
	"msg": "批次删除成功"
}
```

#### 批次信息
描述：批次信息返回
请求方法：GET
url：/machinemanager/batchinfo
请求参数：
类型：form
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|page|int|分页的页码|y||
|size|int|分页的当页容量|y||

示例：

	/machinemanager/sourcepool?page=1&size=10

响应参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|CreatedAt|string|批次创建时间|y||
|UpdatedAt|string|批次更新时间|y||
|ID|int|批次ID|y||
|description|string|批次描述|y||
|manager|string|创建人|y||
|DepartID|[]string|部门ID|y||
|DepartName|[]string|部门名|y||
|machinelist|[]string|机器ID|y||


示例：
```json
{
	"code": 200,
	"data": [
		{
			"ID": 17,
			"CreatedAt": "2022-04-13T04:00:45+08:00",
			"UpdatedAt": "2022-04-13T04:00:45+08:00",
			"DeletedAt": null,
			"name": "new",
			"description": "",
			"manager": "xx",
			"machinelist": "1,2",
			"Depart": "1",
			"DepartName": "xx"
		},
		{
			"ID": 16,
			"CreatedAt": "2022-04-12T13:55:57+08:00",
			"UpdatedAt": "2022-04-12T13:55:57+08:00",
			"DeletedAt": null,
			"name": "test111",
			"description": "",
			"manager": "xx",
			"machinelist": "xx",
			"Depart": "1",
			"DepartName": "xx"
		},
	],
	"ok": true,
	"page": 1,
	"size": 10,
	"total": 3
}
```

#### 批次机器信息
描述：返回该批次机器信息
请求方法：GET
url：/machinemanager/batchinfo
请求参数：
类型：json
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|page|int|分页的页码|y||
|size|int|分页的当页容量|y||
|ID|int|批次ID|y||

示例：
```json
{
    "page":1,
    "size":10,
    "ID":5
}
```

响应参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|id|int|机器id|y||
|departid|int|部门id|y||
|ip|string|机器IP|y||
|machineuuid|string|机器UUID|y||
|CPU|string|机器cpu信息|y||
|state|string|机器状态|y||
|sysinfo|string|机器系统信息|y||


示例：
```json
{
	"code": 200,
	"data": [
		{
			"id": 7,
			"departid": 10,
			"ip": "192.168.160.108",
			"machineuuid": "8912345677",
			"CPU": "intel i7",
			"state": 1,
			"sysinfo": "centos 6"
		},
		{
			"id": 18,
			"departid": 10,
			"ip": "192.168.160.119",
			"machineuuid": "12345678918",
			"CPU": "and 64",
			"state": 1,
			"sysinfo": "suse"
		},
		
	],
	"page": 0,
	"size": 0,
	"total": 4
}
```

### 监控接口
#### 范围查询
描述：对一个时间段进行监控数据查询
请求方法：POST
url：
请求参数：/prometheus/queryrange
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|machineip|string|查询的机器ip|y||
|query|int|查询对应的序号|y||
|starttime|string|开始时间|y||
|endtime|string|结束时间|y||


示例：

```json
{
    "machineip":"192.xxx.xxx.xx:9100",
    "query":6,
    "starttime":"1648620854",
    "endtime":"1648621000"
}
```

响应参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|time|string|当前时间|y||
|value|string|查询对应的序号|y||

示例：

```json
{
	"code": 200,
	"data": [
		{
			"device": "docker0",
			"label": [
				{
					"time": "2022-03-30 14:14:14",
					"value": "0"
				},
				{
					"time": "2022-03-30 14:14:24",
					"value": "0"
				},
				{
					"time": "2022-03-30 14:14:34",
					"value": "0"
				},
			]
		},
	]
}
```

#### 当前时间查询
描述：对当前时间的监控数据获取
请求方法：POST
url：/prometheus/query
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|machineip|string|查询的机器ip|y||
|query|int|查询对应的序号|y||
|time|string|时间|y||

示例：

```json
{
    "machineip":"192.xxx.xxx.xx:9100",
    "query":3,
    "time":"1648620552"
}
```

响应参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|time|string|当前时间|y||
|value|string|查询对应的序号|y||

示例：

```json
{
	"code": 200,
	"data": [
		{
			"device": "dm-0",
			"label": {
				"time": "2022-03-30 14:09:12",
				"value": "1.799640071985603"
			}
		},
		{
			"device": "dm-1",
			"label": {
				"time": "2022-03-30 14:09:12",
				"value": "0"
			}
		},
	]
}
```

#### 告警获取
描述：返回当前的所有告警信息
请求方法：GET
url：/prometheus/alert
请求参数：无


响应结果：
类型：json 
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|alertname|string|告警名称|y||
|instance|string|节点名称|y||
|job|string|工作名称|y||
|annotations|string|注释|y||
|state|string|机器状态|y||
|activeAt|string|机器活跃时间|y||

示例：

```json
{
	"code": 200,
	"data": [
		{
			"alertname": "InstanceDown",
			"instance": "192.xx.160.xx:9090",
			"job": "12345678936",
			"annotations": "Instance 192.168.160.137:9090 down",
			"state": "firing",
			"activeAt": "2022-04-12T11:11:56.02661039Z"
		}
	]
}
```

#### 告警邮件发送
描述：发送邮件
请求方法：POST
url：/prometheus/query
请求参数：
类型：json
字段描述：
|名称|类型|说明|是否必须|备注|
|-|-|-|-|-|
|Email|string|邮箱地址|y||
|alertname|string|名称|y||
|IP|string|报警IP|y||
|summary|string|报警描述|y||
|StartsAt|string|开始时间|y||
|EndsAt|string|结束时间|y||

示例:

```json
{
    "Email": ["xx@qq.com"],
    "Labels": {
        "alertname": "短信服务",    
        "IP": "192.168.1.1"   
    },
    "Annotations": {
        "summary": "短信账号全部欠费了，无法切换可用服务，发不出短信"   
    },
    "StartsAt": "2022-03-18T07:54:52.898371829Z",   
    "EndsAt": "2022-03-18T12:58:52.898371829Z"      
}
```

响应结果：
类型：json

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "success"
}
```
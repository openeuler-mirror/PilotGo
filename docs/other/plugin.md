# 插件开发描述

## plugin-sdk
该工具提供了用于开发PilotGo插件的SDK
- 实现了go-plugin的Server接口和Client接口
- 定义了插件需要实现的PilotGo接口
  1. 插件加载时执行
  ```bash
  func OnLoad() error
  ```
  2. 插件关闭时执行
  ```bash
  func OnClose() error
   ```

  3. 获取插件的基本信息
  ```bash
  func GetManifest() PluginManifest
   ```
  其中PluginManifest包括:

  |id|name| author | version|
  |:----:|:----:|:----:|:----:|
  |插件标识名|插件名|插件开发者|插件版本|

  4. 获取插件的配置，目前的提供配置是在插件中的后端服务的配置信息，例如服务端口等
  ```bash
  func GetConfiguration() []PluginConfig
   ```
    其中PluginConfig包括:

  |title|description| key | type| values |
    |:----:|:----:|:----:|:----:|:----:|
  |配置标题|配置描述|配置标识|配置类型|配置值|
  目前提供的type在plugin-sdk中是：

  |UrlValue|PortValue| ProtocolValue  |
  |:----:|:----:|:----:|
  |地址|端口|协议|

  5. 获取插件的前端资源
  ```bash
  func GetWebExtension() []WebExtension
   ```

  其中PluginConfig包括:

  |type|pathMatchRegex| source |
    |:----:|:----:|:----:|
  |资源类型|资源路径||
  目前提供的type在plugin-sdk中是：

  |CSS|JavaScript| Html  |
    |:----:|:----:|:----:|
  |css资源|javaScript资源|html资源|
- 提供基本的信息
  1. 提供插件进程与主进程连接的握手配置
    ```bash
  var HandshakeConfig = plugin.HandshakeConfig{
	      ProtocolVersion: 1,
	      MagicCookieKey: "PILOTGO_PLUGIN",
	      MagicCookieValue: "Mz1K0OGpIRs",
  }
   ```

## PilotGo server API插件接口
### 插件管理API
#### 获取插件列表
描述：会扫描存放在plugins目录下的插件

请求方法：GET

url：/plugin/list

请求参数：无

响应结果：

类型：json

字段：

|名称|类型|说明|是否必须|备注|
|:----:|:----:|:----:|:----:|:----:|
|name|string|插件名称|y||
|status|int|插件状态（1为运行中，0为关闭）|y||
|url|string|运行的插件前端访问地址|y||
示例：
```json
{
	"code": 200,
	"data":  [
			{
				"name": "",
				"status": 0,
				"url": ""
			},
			{
				"name": "xx",
				"status": 1,
				"url": "http://localhost:8888"
			}
		]
}
```
#### 启动插件
描述：启动关闭状态的插件

请求方法：POST

url：/plugin/load

请求参数：

类型：json

字段描述：

|名称|类型|说明|是否必须|备注|
|:----:|:----:|:----:|:----:|:----:|
|Name|string|插件名|y||

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "插件加载成功"
}
```
#### 关闭插件
描述：关闭运行状态的插件

请求方法：POST

url：/plugin/unload

请求参数：

类型：json

字段描述：

|名称|类型|说明|是否必须|备注|
|:----:|:----:|:----:|:----:|:----:|
|Name|string|插件名|y||

示例：
```json
{
	"code": 200,
	"data": null,
	"msg": "插件卸载成功"
}
```
## 开发一个插件
  - 需要先引入plugin-sdk和go-plugin
  - 实现plugin-sdk定义的基本接口
      ```
      type PluginInterface interface {
        OnLoad() error
        OnClose() error
        GetManifest() PluginManifest
        GetConfiguration() []PluginConfig
        GetWebExtension() []WebExtension
      }
    ```
  - 定义pluginMap以便能够被识别，插件使用方能够分配实例，这里的参数p实现了定义的接口的对象
    ```
	    var pluginMap = map[string]plugin.Plugin{
		        "pilotGo_plugin":&plugin_sdk.PilotGoPlugin{Impl: p},
	    }
    ```
- 调用serve方法用于与插件使用方进行连接，需要传入的配置是plugin_sdk定义的握手配置和pluginMap
    ```
  	 plugin.Serve(&plugin.ServeConfig{
		          HandshakeConfig:plugin_sdk.HandshakeConfig,
		          Plugins:pluginMap,
	      })
    ```
- 打包然后存放到plugins目录下即可被识别到

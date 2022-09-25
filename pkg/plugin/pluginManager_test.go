package plugin

//
//import (
//	"reflect"
//	"testing"
//)
//
//func TestRegister(t *testing.T) {
//	// 基础测试
//	plugin1 := init_plugin()
//	err := globalManager.Regist(&plugin1)
//	if err != nil {
//		t.Errorf("Register Basic Test Failed")
//	}
//
//	// 缺少信息测试1: 缺少Url, Port
//	plugin2 := Plugin{
//		Name:        "plugin1",
//		Version:     "1.1.3",
//		Description: "This is a test example",
//		Depends:     nil,
//		Protocol:    "http",
//	}
//
//	err2 := globalManager.Regist(&plugin2)
//	if err2 == nil {
//		t.Errorf("Register Insufficient Info Test 1 Failed")
//	}
//
//	// 缺少信息测试2: 缺少Version, Description, Depends
//	plugin3 := Plugin{
//		Name:     "plugin3",
//		Url:      "www.test3.com",
//		Port:     "113",
//		Protocol: "https",
//	}
//
//	err3 := globalManager.Regist(&plugin3)
//	if err3 != nil {
//		t.Errorf("Register Insufficient Info Test 2 Failed: %s", err3.Error())
//	}
//
//	// 空信息测试
//	plugin4 := Plugin{
//		Name:        "",
//		Version:     "",
//		Description: "",
//		Depends:     nil,
//		Url:         "",
//		Port:        "",
//		Protocol:    "",
//	}
//	err4 := globalManager.Regist(&plugin4)
//	if err4 == nil {
//		t.Errorf("Register Empty Test Failed")
//	}
//
//	// 信息重复测试1： 名字重复
//	plugin2.Url = "www.test1.com"
//	err5 := globalManager.Regist(&plugin2)
//	if err5 == nil {
//		t.Errorf("Register Duplication Test Failed")
//	}
//
//	// 信息重复测试2：端口重复
//	plugin2.Name = "plugin2"
//	plugin2.Port = "3389"
//	err6 := globalManager.Regist(&plugin2)
//	if err6 != nil {
//		t.Errorf("Register Duplication Test Failed: %s", err6.Error())
//	}
//}
//
//func TestRemove(t *testing.T) {
//	plugin1 := init_plugin()
//	globalManager.Regist(&plugin1)
//
//	err1 := globalManager.Remove("plugin1")
//	err2 := globalManager.Remove("plug")
//	err3 := globalManager.Remove("plugin1")
//
//	if err1 != nil {
//		t.Errorf("Remove Test Failed: %s", err1.Error())
//	}
//	if err2 == nil || err3 == nil {
//		t.Errorf("Remove Test Failed")
//	}
//}
//
//func TestGet(t *testing.T) {
//	plugin1 := init_plugin()
//
//	// 基础测试
//	globalManager.Regist(&plugin1)
//	url, port, protocal, err1 := globalManager.Get("plugin1")
//
//	if url != plugin1.Url || port != plugin1.Port || protocal != plugin1.Protocol || err1 != nil {
//		if err1 != nil {
//			t.Errorf("Get Basic Test Failed: %s", err1.Error())
//		} else {
//			t.Errorf("Get Basic Test Failed:\n Url: %s\n, Port: %s\n, Protocal: %s\n", url, port, protocal)
//		}
//	}
//
//	// 未注册名字测试
//	url2, port2, protocal2, err2 := globalManager.Get("DNE")
//	if url2 != "" || port2 != "" || protocal2 != "" || err2 == nil {
//		t.Errorf("Get Invalid Name Test 1 Failed")
//	}
//
//	// 缺少信息测试： 缺少Description, Depends
//	plugin3 := Plugin{
//		Name:     "plugin3",
//		Version:  "1.3",
//		Url:      "www.info-test.com",
//		Port:     "32",
//		Protocol: "http",
//	}
//	globalManager.Regist(&plugin3)
//	url3, port3, protocal3, err3 := globalManager.Get("plugin3")
//
//	if url3 != plugin3.Url || port3 != plugin3.Port || protocal3 != plugin3.Protocol || err3 != nil {
//		if err3 != nil {
//			t.Errorf("Get Basic Test Failed: %s", err1.Error())
//		} else {
//			t.Errorf("Get Basic Test Failed:\n Url: %s\n, Port: %s\n, Protocal: %s\n", url, port, protocal)
//		}
//	}
//}
//
//func TestGetAll(t *testing.T) {
//	// 空值测试
//	globalManager.Remove("plugin1")
//	if globalManager.GetAll() != nil {
//		t.Errorf("GetAll Empty Test Failed")
//	}
//
//	plugin1 := init_plugin()
//	plugin2 := Plugin{
//		Name:        "plugin2",
//		Version:     "1.1",
//		Description: "An example for plugin manager",
//		Url:         "www.test2.com",
//		Port:        "2100",
//		Protocol:    "https",
//	}
//
//	// 基础测试
//	globalManager.Regist(&plugin1)
//	globalManager.Regist(&plugin2)
//
//	plugins := globalManager.GetAll()
//
//	plugin_want := make([]Plugin, 0)
//	plugin_want = append(plugin_want, plugin1, plugin2)
//
//	if !reflect.DeepEqual(plugins, plugin_want) {
//		t.Errorf("GetAll Basic Test Failed：%+v\n, want: %+v\n", plugins, plugin_want)
//	}
//}
//
//func TestCheck(t *testing.T) {
//	// check1 = false
//	check1 := globalManager.Check("plug")
//
//	plugin1 := init_plugin()
//	globalManager.Regist(&plugin1)
//	// check2 = true
//	check2 := globalManager.Check("plugin1")
//
//	if check1 || !check2 {
//		t.Errorf("Check Test Failed")
//	}
//}
//
//func init_plugin() Plugin {
//	plugin := Plugin{
//		Name:        "plugin1",
//		Version:     "2.1",
//		Description: "Plugin example",
//		Depends:     nil,
//		Url:         "www.test1.com",
//		Port:        "3389",
//		Protocol:    "https",
//	}
//	return plugin
//}

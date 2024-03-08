<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-04-14 17:00:07
  LastEditTime: 2022-07-01 14:16:52
 -->
<template>
 <div class="content">
   <div class="select">
     <el-form :model="firewallForm" ref="firewallForm" :rules="rules" class="form">
        <el-form-item label="执行命令:" class="p_text" prop="linux_value">
          <el-select
            v-model="firewallForm.linux_value"
            placeholder="请选择指令">
            <el-option
              v-for="item in LinuxOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="showNext" prop="zone" label="区域位置:">
          <el-select
            v-model="firewallForm.zone"
            placeholder="请选择防火墙区域"
          >
            <el-option
              v-for="item in ZoneOptions"
              :key="item.$index"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="showNext" prop="protocol" label="网络协议:">
          <el-select
            v-model="firewallForm.protocol"
            placeholder="请选择网络协议"
          >
            <el-option
              v-for="item in protocols"
              :key="item.$index"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="showNext" prop="port" label=' 端 口 号 :'>
          <el-input
            placeholder="请输入端口号"
            v-model="firewallForm.port"
            clearable
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div class="btn">
        <el-button plain  type="primary" @click="handleClick">确定</el-button>
        <el-button plain  type="primary" @click="handlerClickReset">重置</el-button>
      </div>
   </div>
   <div class="info" v-loading="loading">
     <div class="left">
      <div class="lt" style="width:100%;height:30%;">
        <el-result v-if="isRun" icon="success" title="防火墙运行中"/>
        <el-result v-else icon="warning" title="防火墙未运行"/>
      </div>
      <div class="lc" style="width:100%;height:50%;text-align: center; color:#777">
        <p>网卡：{{defaultNic}}</p>
        <p>默认区域：{{defaultZone}}</p>
      </div>
      <div class="lb" style="width:100%;height:10%;">
        <p>更换默认区域：</p><br/>
        <el-select v-model="defaultZone" style="width:100%;"
          @change="zoneChange">
          <el-option
            v-for="item in ZoneOptions"
            :key="item.$index"
            :label="item"
            :value="item">
          </el-option>
        </el-select>
      </div>
     </div>
     <div class="right">
      <p class="tooltip1"><b style="color: #000;text-align: center; width: 8%;">区域</b>* FirewallD区域定义了绑定的网络连接、网卡以及源地址的可信程度。区域是服务、端口以及来源的组合。</p>
      <el-tabs tab-position="left" style="width:100%;height: 96%; border-top: #777;" @tab-click="handleZone" v-model="selectZone">
        <el-tab-pane v-for="zone in ZoneOptions" :key="zone.$index" :label="zone" :name="zone" style="width:100%;height:100%;">
          <el-tabs v-model="activeName" v-loading="zoneLoading" style="width:100%;height: 100%;">
            <el-tab-pane label="服务" name="service">
              <p class="tooltip">* 可以在这里定义区域中哪些服务是可信的。可连接至绑定到这个区域的连接、网卡和源的所有主机和网络及可以访问可信服务。</p><br/>
              <el-row>
                <el-checkbox-group v-model="hasServices">
                  <el-col v-for="item in services" :key="item.$index" :span="3">
                    <el-checkbox @change="checked=>handleService(checked,item)"  :label="item"/>
                  </el-col>
                </el-checkbox-group>
              </el-row>
            </el-tab-pane> 
            <el-tab-pane label="端口" name="port">
              <p class="tooltip">* 添加可让允许访问的主机或者网络访问的附加端口或者端口范围</p><br/>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="协议">端口</el-descriptions-item>
                <el-descriptions-item v-for="item in ports" :key="item.$index" :label="item.protocol">{{item.port}}</el-descriptions-item>
            </el-descriptions>
            </el-tab-pane>
            <el-tab-pane label="来源" name="source">
              <p class="tooltip">* 添加条目以便在该区域绑定源地址或范围。还可以绑定到MAC源地址，但会有所限制。端口转发及伪装不适用于MAC源绑定</p><br/>
              <span v-for="item in sources" :key="item.$index">{{item}}<el-divider/></span>
            </el-tab-pane>
          </el-tabs>
        </el-tab-pane>
      </el-tabs>    
     </div>
    
   </div>
 </div>
</template>
<script>
import { reStart, close, openPort, deleteOpenPort, 
  getConfig, changeDefaultZone, getZoneConfig, addService,
  delService, addSource, delSource} from "@/request/firewall";
export default {
  name: "FirewallInfo",
  data() {
    return {
      loading: true,
      zoneLoading: false,
      showService: false,
      activeName: 'service',
      defaultNic: '',
      defaultZone: '',
      runSatus: '',
      isRun: false,
      selectZone: '',
      hasServices: [],
      services: [],
      sources: [],
      ports: [],
      params: {},
      sourceTypes: ["IP","MAC","ipset"],
      firewallForm: {
        zone: "",
        linux_value: "",
        protocol: "",
        port: "",
      },
      rules:{
        linux_value:[{ required: true,message:'必选', trigger: "change" }],
        zone:[{ required: true, message:'请选择防火墙区域位置',trigger: "change" }],
        protocol:[{ required: true, message:'请选择协议',trigger: "change" }],
        port:[{ required: true, message:'请输入端口号',trigger: "blur" }],
      },
      linux_value: [
        "重启防火墙",
        "关闭防火墙",
        "指定区域开放端口",
        "指定区域删除开放端口",
      ],
      LinuxOptions: [
        {
          value: "1",
          label: "重启防火墙",
        },
        {
          value: "2",
          label: "关闭防火墙",
        },
        {
          value: "4",
          label: "指定区域开放端口",
        },
        {
          value: "6",
          label: "指定区域删除端口",
        },
      ],
      protocols: ["TCP","UDP"],
      ZoneOptions: [],
    }
  },
  computed: {
    showNext() {
      let value = this.firewallForm.linux_value;
      return ['4','5','6'].includes(value);
    }
  },
  mounted() {
    let obj = this.params = {uuid:this.$route.params.detail};
    getConfig(obj).then(res => {
      if(res.data.code === 200 && res.data.data.firewalld_config.status.toLowerCase() == 'running') {
        this.isRun = true;
        let {defaultZone,nic,services,zones} = res.data.data.firewalld_config;
        this.ZoneOptions = zones;
        this.defaultNic = nic;
        this.defaultZone = defaultZone;
        this.services = services;
        this.handleZone({name:defaultZone});
        this.loading = false;
      } else {
        this.isRun = false;
        this.loading = false;
      }
    })
  },
  methods: {
    handleZone(tag) {
      this.selectZone = tag.name;
      this.zoneLoading = true;
      this.activeName = "service";
      this.hasServices = [];
      this.ports = [];
      this.sources = [];
      getZoneConfig({...this.params,zone:tag.name}).then(res => {
        if(res.data.code === 200) {
          this.zoneLoading = false;
          let {ports,services,sources} = res.data.data.firewalld_zone;
          this.hasServices = services;
          this.ports = ports;
          this.sources = sources;
        }
      })
    },
    zoneChange() {
      changeDefaultZone({...this.params,zone:this.defaultZone}).then(res => {
        if(res.data.code == 200) {
          this.$message.success(res.data.msg)
        } else {
          this.$message.error(res.data.msg)
        }
      })
    },
    handleService(value,service) {
      let _this = this;
      let confirmText = value ? "请确认是否添加"+ service + "服务" : "请确认是否移除"+ service + "服务";
      _this.$confirm(confirmText, '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }).then(() => {
          _this.handleServiceConfirm(value,service)
        }).catch(() => {
          if(value) {
            let index = _this.hasServices.indexOf(value);
            _this.hasServices.splice(index,1);
          } else {
            _this.hasServices.push(service);
          }
        })
    },
    handleServiceConfirm(value,item) {
      let params = {
        uuid: this.$route.params.detail,
        zone: this.selectZone,
        service: item
      }
      value ? addService(params).then(res => this.handleResult(res))
       : delService(params).then(res => this.handleResult(res))
    },
    displayTooltip() {
      let message = `请输入ipv4或ipv6地址，格式为address[/mask].
          1.对于ipv4地址，该掩码必须为网络掩码或者一个数字.
          2.对于ipv6地址，则该掩码为一个数字.`;
      this.$prompt(message, '地址', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        this.displayTip = false;
        this.sForm.address = value;
      })
    },
    handleResult(res) {
      if(res.data.code === 200) {
        this.$message({
          type:'success',
          message: res.data.msg,
        })
      } else {
        this.$message({
          type: 'error',
          message: res.data.msg
        })
      }
    },
    handleClick() {
      let data = {
        zone:this.firewallForm.zone,
        port:this.firewallForm.port,
        protocol:this.firewallForm.protocol.toLowerCase()
      };
      this.$refs.firewallForm.validate((valid)=>{
        if (valid) {
          switch (this.firewallForm.linux_value) {
            case "":
              break;
            case "1":
              reStart(this.params).then(res => {
                this.handleResult(res);
              });
              break;
            case "2":
              close(this.params).then(res => {
                this.handleResult(res);
              });
              break;
            case "4":
              openPort({...this.params,...data}).then(res => {
                this.handleResult(res);
              });
              break;
            case "6":
              deleteOpenPort({...this.params,...data}).then(res => {
                this.handleResult(res);
              });
              break;
          }
        }
      })
    },
    handlerClickReset() {
      this.display = true;
      this.$refs.firewallForm.resetFields();
    },
  },

}
</script>
<style scoped lang="scss">
.content {
  width: 100%;
  height: 100%;
  .select {
    width: 98%;
    display: flex;
    justify-content:baseline;
    .form {
      width: 45%;
      height: 100%;
      .el-select,
      .el-input {
        width: 80%;
      }
    }
    .btn {
      width: 20%;
      height: 100%;
    }
  }
  .info {
    width: 98%;
    height: 76%;
    border: 3px solid #ddd;
    display: flex;
    justify-content: space-around;
    .left {
      width:14%;
      height:100%;
      float: left;
      display: flex;
      font-size: 16px;
      justify-content: flex-start;
      flex-direction: column;
      align-items: center;
      border-right: 3px solid #ddd;
    }
    .right {
      width:86%;
      height:100%;
      .tooltip1 {
        height: 4%;
        display: flex;
        align-items: center;
        // color: rgb(241, 139, 14);
        color: #777;
        font-size:16px;
        border-bottom: 1px solid #ddd;
      }
      .tooltip {
        // color: rgb(241, 139, 14);
        color: #777;
        font-size: 14px;
      }
    }
  }
}
</style>

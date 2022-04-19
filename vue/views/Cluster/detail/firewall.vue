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
  LastEditTime: 2022-04-19 09:56:16
 -->
<template>
 <div class="content">
   <div class="select">
     <el-form :model="firewallForm" ref="firewallForm" :rules="rules" class="form">
        <el-form-item label="执行命令:" class="p_text" prop="linux_value">
          <el-select
            v-model="firewallForm.linux_value"
            placeholder="请选择指令"
          >
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
            @keyup.enter.native="handleClick('firewallForm')"
          >
            <el-option
              v-for="item in ZoneOptions"
              :key="item.zone_value"
              :label="item.zone_label"
              :value="item.zone_value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item v-if="showNext" prop="port" label=' 端 口 号 :'>
          <el-input
            placeholder="请输入端口号"
            v-model="firewallForm.port"
            clearable
            @keyup.enter.native="handleClick('firewallForm')"
          >
          </el-input>
        </el-form-item>
      </el-form>
      <div class="btn">
        <el-button plain  type="primary" @click="handleClick">确定</el-button>
        <el-button plain  type="primary" @click="handlerClickReset">重置</el-button>
      </div>
   </div>
   <div class="info">
     <div class="detail" v-if="display">
       <p class="title">防火墙配置信息：</p>
       <div>{{ fireArea }}</div>
       <div class="config" v-html="fireConfig"></div>
     </div>
     <div class="result" v-else>
       <p class="title">执行结果：</p>
       <el-descriptions :column="1" size="medium" border>
        <el-descriptions-item label="执行动作">{{ action }}</el-descriptions-item>
        <el-descriptions-item label="结果">
          {{result+":"}}
          <p class="progress" v-show="result != ''">
            <span :style="{background: result === '成功' ? 'rgb(109, 123, 172)' : 'rgb(223, 96, 88)'}">100%</span>
          </p>
        </el-descriptions-item>
      </el-descriptions>
     </div>
   </div>
 </div>
</template>
<script>
import { reStart, close, openPort, deleteOpenPort, getConfig } from "@/request/firewall";
export default {
  name: "FirewallInfo",
  data() {
    return {
      fireArea: '',
      fireConfig: '',
      display: true,
      params: {},
      result: '',
      action: '暂无',
      firewallInfo: {
        Name: '',
      },
      firewallForm: {
        zone: "",
        linux_value: "",
        port: "",
      },
      rules:{
        linux_value:[{ required: true,message:'必选', trigger: "change" }],
        zone:[{ required: true, message:'请选择防火墙区域位置',trigger: "change" }],
        port:[{ required: true, message:'请输入端口号',trigger: "change" }],
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

      ZoneOptions: [
        {
          zone_value: "public",
          zone_label: "默认区域public",
        },
        {
          zone_value: "trusted",
          zone_label: "trusted",
        },
        {
          zone_value: "dmz",
          zone_label: "dmz",
        },
        {
          zone_value: "internal",
          zone_label: "internal",
        },
        {
          zone_value: "work",
          zone_label: "work",
        },
        {
          zone_value: "home",
          zone_label: "home",
        },
        {
          zone_value: "drop",
          zone_label: "drop",
        },
        {
          zone_value: "block",
          zone_label: "block",
        },
      ],

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
      if(res.data.code === 200) {
        let result = res.data.data.firewalld_config;
        let configDiv = '';
        this.fireArea = result && result[0];
        result.forEach((item,index) => {
          if(index == 0) return;
          configDiv += '<p>' + item + '</p>';
        })
        this.fireConfig = configDiv;
      }
    })
  },
  methods: {
    handleResult(res) {
      this.display = false;
      if(res.data.code === 200) {
        this.result = "成功";
        this.$message({
          type:'success',
          message: res.data.msg,
        })
      } else {
        this.result = "失败";
        this.$message({
          type: 'error',
          message: res.data.msg
        })
      }
    },
    handleClick() {
      this.action = this.LinuxOptions.filter(item => item.value == this.firewallForm.linux_value)[0].label;
      this.$refs.firewallForm.validate((valid)=>{
        if (valid) {
          let data = {};
          data["zone"]=this.firewallForm.zone;
          data["port"]=this.firewallForm.port;
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
  display: flex;
  flex-direction: column;
  align-items: center;
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
    height: 80%;
    overflow-y: auto;
    .detail {
      width: 100%;
      height: 100%;
      .title {
        width: 30%;
        margin: 2% 0;
      }
      .config {
        width: 100%;
        height: 100%;
        text-indent: 1em;
        display: flex;
        flex-direction: column;
        justify-content: space-evenly;
        p {
          width: 100%;
        }
      }
    }
    .result {
      width: 100%;
      height: 90%;
      .title {
        width: 30%;
        margin: 2% 0;
      }
      .progress {
        display: inline-block;
        width:74%; 
        margin-left: 2%;
        border: 1px solid rgba(11, 35, 117,.5);  
        background: #fff; 
        border-radius: 10px; 
        text-align:left;
        span {
          display: inline-block;
          text-align:center;
          color: #fff;
          width: 100%;
          border: 1px solid #fff;
          border-radius: 10px;
        }
      }
    }
  }
}
</style>

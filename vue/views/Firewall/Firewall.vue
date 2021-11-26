<template>
  <div>
    <div class="firewall-conditions">
      <span>防火墙配置</span>
      <div>
        <el-form :model="firewallForm" ref="firewallForm" :rules="rules">
          <el-row>
            <el-col :span="6">
              <el-form-item prop="ip" label="主 机 IP" class="p_text">
                <el-input
                  class="firewall-input"
                  size="mini"
                  placeholder="请输入主机ip"
                  v-model="firewallForm.ip"
                  clearable
                  @keyup.enter.native="handleClick('firewallForm')"
                >
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item prop="host_user" label="用户名" class="p_text">
                <el-input
                  class="firewall-input"
                  size="mini"
                  placeholder="请输入用户名"
                  v-model="firewallForm.host_user"
                  clearable
                  @keyup.enter.native="handleClick('firewallForm')"
                >
                </el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item prop="host_password" label="密码" class="p_text">
                <el-input
                  class="firewall-input"
                  size="mini"
                  placeholder="请输入密码"
                  v-model="firewallForm.host_password"
                  clearable
                  @keyup.enter.native="handleClick('firewallForm')"
                >
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="6">
              <el-form-item label="执行命令" class="p_text" prop="linux_value">
                <el-select
                  size="mini"
                  class="firewall-input1"
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
            </el-col>
          </el-row>
          <el-row v-if="firewallForm.linux_value==='4'||firewallForm.linux_value==='5'||firewallForm.linux_value==='6'">
            <el-col :span="6">
              <el-form-item prop="zone" label="区域位置" class="p_text">
                <el-select
                  size="mini"
                  class="firewall-input1"
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
            </el-col>
            <el-col :span="6">
              <el-form-item prop="port" label="端口号" class="p_text">
                <el-input
                  size="mini"
                  class="firewall-input2"
                  placeholder="请输入端口号"
                  v-model="firewallForm.port"
                  clearable
                  @keyup.enter.native="handleClick('firewallForm')"
                >
                </el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="4"> </el-col>
            <el-button size="mini" type="primary" @click="handleClick">确定</el-button>
            <el-button type="primary" size="mini" @click="handlerClickReset">重置</el-button>
          </el-row>
        </el-form>
      </div>
    </div>
    <div class="firewall-config">
      <div class="config">
        <el-descriptions title="防火墙配置信息" :data="tmpData">
          <el-descriptions-item label="Config">
            {{tmpData}}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
  </div>
</template>

<script>
import {
  FirewallConfig,
  FirewallStop,
  FirewallRestart,
  FirewallReload,
} from "@/request/api";
import {
  FirewallAddZonePort,
  FirewallDelZonePort,
  FirewallAddZonePortPermanent,
} from "@/request/api";

export default {
  name: "Firewall",
  data() {
    return {
      tmpData1:"",
      tmpData:"",
      firewallForm: {
        ip: "",
        host_user: "",
        host_password: "",
        zone: "",
        linux_value: "",
        port: "",
      },
      rules:{
        ip:[{ required: true, message:'主机ip不可为空', trigger: "change" }],
        host_user:[{ required: true,message:'用户名不可为空', trigger: "change" }],
        host_password:[{ required: true, message:'请输入用户密码',trigger: "change" }],
        linux_value:[{ required: true,message:'必选', trigger: "change" }],
        zone:[{ required: true, message:'请选择防火墙区域位置',trigger: "change" }],
        port:[{ required: true, message:'请输入端口号',trigger: "change" }],
      },
      linux_value: [
        "重启防火墙",
        "关闭防火墙",
        "更新防火墙",
        "指定区域开放端口",
        "指定区域永久开放端口",
        "删除开放的端口",
        "获取防火墙配置",
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
          value: "3",
          label: "更新防火墙",
        },
        {
          value: "4",
          label: "指定区域开放端口",
        },
        {
          value: "5",
          label: "指定区域永久开放端口",
        },
        {
          value: "6",
          label: "删除开放的端口",
        },
        {
          value: "7",
          label: "获取防火墙配置",
        },
      ],

      ZoneOptions: [
        {
          zone_value: "public",
          zone_label: "public",
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
    };
  },
  methods: {
    handleClick() {
      this.$refs.firewallForm.validate((valid)=>{
        if (valid) {
          let data = new FormData();
          data.append("ip", this.firewallForm.ip);
          data.append("host_user", this.firewallForm.host_user);
          data.append("host_password", this.firewallForm.host_password);
          switch (this.firewallForm.linux_value) {
            case "":
              break;
            case "1":
              FirewallRestart(data);
              this.$message({
                type:'success',
                message:'重启防火墙成功',
              })
              break;
            case "2":
              FirewallStop(data);
              this.$message({
                type:'success',
                message:'关闭防火墙成功',
              })
              break;
            case "3":
              FirewallReload(data);
              this.$message({
                type:'success',
                message:'更新防火墙成功',
              })
              break;
            case "4":
              data.append("zone", this.firewallForm.zone);
              data.append("port", this.firewallForm.port);
              FirewallAddZonePort(data);
              this.$message({
                type:'success',
                message:'开放端口成功',
              })
              break;
            case "5":
              data.append("zone", this.firewallForm.zone);
              data.append("port", this.firewallForm.port);
              FirewallAddZonePortPermanent(data);
              this.$message({
                type:'success',
                message:'永久开放端口成功',
              })
              break;
            case "6":
              data.append("zone", this.firewallForm.zone);
              data.append("port", this.firewallForm.port);
              FirewallDelZonePort(data);
              this.$message({
                type:'success',
                message:'删除端口成功',
              })
              break;
            case "7":
              const _this=this
              FirewallConfig(data).then(function (res){
                _this.tmpData = res.data.data["tmp"];
              })
              this.$message({
                type:'success',
                message:'获取防火墙配置成功',
              })
          }
        } else {
          this.$message({
            type:'error',
            message:'表单提交失败！'
          })
          return false;
        }
      })
    },
    handlerClickReset() {
      this.firewallForm.host_password = '';
      this.firewallForm.host_user = "";
      this.firewallForm.ip = "";
      this.firewallForm.linux_value = "";
      this.firewallForm.zone = "";
      this.firewallForm.port = "";
      this.$refs.firewallForm.resetFields();
    },
  },
};
</script>

<style lang="scss" scoped>
.firewall-conditions {
  background: #fff;
  /*height: 100px;*/
  line-height: 50px;
  .el-row {
    margin-left: 30px;
  }

  .el-col {
    .p_text {
      font-size: 14px;
    }
  }
}

.firewall-input {
  margin-left: 18px;
  width: 230px;
}
.firewall-input1 {
  margin-left: 10px;
  width: 230px;
}
.firewall-input2 {
  margin-left: 18px;
  width: 230px;
}
.firewall-config {
  width: 100%;
  margin: 10px auto;
  float: left;
}
.config {
  margin-top: 15px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}
</style>

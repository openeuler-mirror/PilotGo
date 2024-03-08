<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
      <el-form-item label="配置方式" prop="BOOTPROTO">
        <el-select v-model="form.BOOTPROTO" @change="handleSelect">
          <el-option
            v-for="item in type"
            :key="item.id"
            :value="item.value"
            :label="item.name"
          >
          </el-option>
        </el-select>
      </el-form-item>
        <el-form-item label="IPv4地址" prop="IPADDR" v-show="showIPv4">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.IPADDR"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4子网掩码" prop="NETMASK" v-show="showIPv4">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.NETMASK"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4网关" prop="GATEWAY" v-show="showIPv4">
           <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.GATEWAY"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4首选DNS" prop="DNS1" v-show="showIPv4">
         <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.DNS1"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4备选DNS" prop="DNS2" v-show="showIPv4">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.DNS2"
            autocomplete="off"
          ></el-input>
        </el-form-item>
      </el-form>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleSubmitForm()">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  updateNet  } from "@/request/cluster";
export default {
  props: {
    net: {
      type: Object,
    },
  },
  data() {
    return {
      uuid: '',
      showIPv4: true,
      type: [
        {
          id: 1,
          name: '手动',
          value: 'static'
        },
        {
          id: 2,
          name: '动态DHCP',
          value: 'dhcp'
        }
      ],
      form: {
        BOOTPROTO: "",
        IPADDR: "",
        NETMASK: "",
        GATEWAY: "",
        DNS1: "",
        DNS2: "",
      },
      rules: {
        BOOTPROTO: [{ 
          required: true, 
          message: "请选择分配方式",
          trigger: "blur"
        }],
        IPADDR: [{ 
          required: true, 
          trigger: "blur", 
          message:"请输入IP地址" 
        }],
        NETMASK: [{ 
          required: true, 
          trigger: "blur", 
          message:"请输入子网掩码" 
        }],
        GATEWAY: [{ 
          required: true, 
          trigger: "blur", 
          message:"请输入网关" 
        }],
        DNS1: [{ 
          required: true, 
          trigger: "blur", 
          message:"请输入首选DNS" 
        }],
      },
    }
  },
  mounted() {
    this.form.BOOTPROTO = this.net.BOOTPROTO === '手动' ? 'static' : 'dhcp';
    this.showIPv4 = this.net.BOOTPROTO === '手动' ? true : false;
    this.form.IPADDR = this.net.IPADDR;
    this.form.NETMASK = this.net.NETMASK;
    this.form.GATEWAY = this.net.GATEWAY;
    this.form.DNS1 = this.net.DNS1;
    this.form.DNS2 = this.net.DNS2;
    this.uuid = this.$route.params.detail;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSelect(type) {
      this.showIPv4 = type === 'static' ? true : false; 
    },
    handleSubmitForm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateNet({macUUID: this.uuid, ...this.form})
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click");
                this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.msg);
              }
            })
        }
      });
    },
  },
};
</script>
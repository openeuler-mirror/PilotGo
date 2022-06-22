<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
      <el-form-item label="配置方式" prop="ip">
        <el-select v-model="form.ip" @change="handleSelect">
          <el-option
            v-for="item in type"
            :key="item.id"
            :value="item.id"
            :label="item.name"
          >
          </el-option>
        </el-select>
      </el-form-item>
        <el-form-item label="IPv4地址" prop="ip" v-show="showIPv4">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.ip"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4子网前缀长度" prop="kernel" v-show="showIPv4">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.kernel"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4网关" prop="repo" v-show="showIPv4">
           <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.repo"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4首选DNS" prop="service" v-show="showIPv4">
         <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.service"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv4备选DNS" prop="kernel" v-show="showIPv4">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.kernel"
            autocomplete="off"
          ></el-input>
        </el-form-item>

        <el-form-item label="IPv6" prop="IPv6" v-show="showIPv4">
        <el-switch
          v-model="isIPv6"
          active-text="开"
          inactive-text="关"
          @change="handle6Change">
        </el-switch>
      </el-form-item>
        <el-form-item label="IPv6地址" prop="ip" v-show="showIPv6">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.ip"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv6子网前缀长度" prop="kernel" v-show="showIPv6">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.kernel"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv6网关" prop="repo" v-show="showIPv6">
           <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.repo"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv6首选DNS" prop="service" v-show="showIPv6">
         <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.service"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="IPv6备选DNS" prop="kernel" v-show="showIPv6">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.kernel"
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
import {  updateIp  } from "@/request/cluster";
export default {
  props: {
    ip: {
      type: String
    } 
  },
  data() {
    return {
      isIPv4: true,
      isIPv6: false,
      showIPv4: true,
      showIPv6: false,
      type: [
        {
          id: 1,
          name: '手动'
        },
        {
          id: 2,
          name: '动态DHCP'
        }
      ],
      form: {
        ip: "",
        repo: "",
        kernel: "",
        service: "",
      },
      rules: {
        ip: [{ 
          required: true, 
          message: "请输入IP",
          trigger: "blur"
        }],
        kernel: [{ 
          required: true, 
          trigger: "blur", 
          message:"修改后需要重启生效" 
        }],
      },
      disabled: true,
      repos: [],
      service: [],
    }
  },
  mounted() {
    this.form.ip = this.ip;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSelect(type) {
      this.showIPv4 = type === 1 ? true : false; 
    },
    handle6Change(value) {
      this.showIPv6 = value;
    },
    handleSubmitForm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateIp({ip: this.form.ip, data: this.form})
            .then((res) => {
              if (res.data.status === "success") {
                this.$emit("click");
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("修改失败，请检查输入内容");
            });
        } else {
          this.$message.error("修改失败，请检查输入内容");
          return false;
        }
      });
    },
  },
};
</script>
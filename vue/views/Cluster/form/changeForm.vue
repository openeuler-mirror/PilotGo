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
  Date: 2022-03-25 10:03:53
  LastEditTime: 2022-03-25 10:46:10
 -->
<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="IP:" prop="ip">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.ip"
            :disabled="disabled"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="原部门:" prop="formerDept">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.formerDept"
            :disabled="disabled"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="新部门:" prop="currentDept">
          <el-input
            class="ipInput"
            controls-position="right"
            :disabled="disabled"
            v-model="form.currentDept"
            autocomplete="off"
          ></el-input>
          <ky-tree
            :getData="getChildNode" 
            :showEdit="false"
            ref="tree" 
            @nodeClick="handleSelectDept">
          </ky-tree>
        </el-form-item>
      </el-form>

      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleChange">确 定</el-button>
      </div>
  </div>
</template>
<script>
import kyTree from "@/components/KyTree";
import { getChildNode, changeMacDept } from "@/request/cluster";
export default {
  name: 'ChangeForm',
  components: {
    kyTree,
  },
  props: {
    row: {
      type: Object
    }
  },
  data() {
   return {
      disabled: true,
      machineid: 0,
      departid: 0,
      form: {
        ip: '',
        formerDept: '',
        currentDept: ''
      },
      rules: {
        ip: [
          { 
            required: true, 
            message: "ip未识别到",
            trigger: "blur" 
          }],
        formerDept: [
          { 
            required: true, 
            message: "原部门未识别到",
            trigger: "blur" 
          }],
        currentDept: [
        { 
          required: true, 
          message: "请选择新部门",
          trigger: "blur" 
        }],
      }
   }
  },
  mounted() {
    this.form.ip = this.row.ip;
    this.form.formerDept = this.row.departname;
    this.departid = this.row.departid;
    this.machineid = this.row.id;
  },
  methods: {
    getChildNode,
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSelectDept(data) {
      if(data) {
        this.form.currentDept = data.label;
        this.departid = data.id;
      }
    },
    handleChange() {
      let params = {
        "machineid": this.machineid,
        "departid": this.departid,
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          changeMacDept(params).then((res) => {
            if (res.data.code === 200) {
              this.$emit("click","success");
              this.$message.success(res.data.msg);
              this.$refs.form.resetFields();
            } else {
              this.$message.error(res.data.msg);
            }
          })
          .catch((res) => {
            this.$message.error("更换部门失败");
          });
        }
      });
    },
  }
}
</script>
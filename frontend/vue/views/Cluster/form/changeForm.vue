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
  LastEditTime: 2022-04-07 14:34:23
 -->
<template>
  <div class="content">
    <div class="mac">
        <div slot="header" class="clearfix">
          <span>已选机器列表:</span>
        </div>
        <small-table
            ref="stable"
            :data="machineArr"
            :height="tHeight">
            <template v-slot:content>
            <el-table-column
              prop="ip"
              label="ip">
            </el-table-column>
            <el-table-column
              prop="departname"
              label="原部门">
            </el-table-column>
            </template>
          </small-table>
    </div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
        class="dept"
      >
        <!-- <el-form-item label="IP:" prop="ip">
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
        </el-form-item> -->
        <el-form-item label="部门:" prop="currentDept">
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
        <el-button type="primary" @click="handleContinue">继续选择</el-button>
        <el-button type="primary" @click="handleChange">确 定</el-button>
      </div>
  </div>
</template>
<script>
import kyTree from "@/components/KyTree";
import { getChildNode, changeMacDept } from "@/request/cluster";
import SmallTable from "@/components/SmallTable";
export default {
  name: 'ChangeForm',
  components: {
    kyTree,
    SmallTable
  },
  props: {
    machines: {
      type: Array,
    }
  },
  data() {
   return {
      machineArr: [],
      tHeight: 340,
      disabled: true,
      machineid: 0,
      departid: 0,
      form: {
        currentDept: ''
      },
      rules: {
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
    let keys = {};
    this.machines.forEach((item) => keys[item.ip]=item);
    for(let key in keys) {
      this.machineArr.push(keys[key])
    }
  },
  methods: {
    getChildNode,
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleContinue() {
      this.$emit("click");
    },
    handleSelectDept(data) {
      if(data) {
        this.form.currentDept = data.label;
        this.departid = data.id;
      }
    },
    handleChange() {
      let macIds = [];
      this.machineArr.forEach(item => {
        macIds.push(item.id)
      })
      let params = {
        "machineid": macIds.toString(),
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
<style scoped lang="scss">
.content {
  display: flex;
  flex-wrap: wrap;
  .mac {
      width: 50%;
      max-height: 400px;
      overflow:auto;
      .clearfix {
        font-size: 16px;
      }
      .text {
        font-size: 16px;
      }
    }
    .dept {
      width: 40%;
    }
    .dialog-footer {
      width: 100%;
    }
}
  
</style>
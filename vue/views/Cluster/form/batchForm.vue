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
  Date: 2022-02-25 16:33:46
  LastEditTime: 2022-02-25 17:39:43
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="批次名称:" prop="batchName">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.batchName"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="已选机器:" prop="mechines">
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
              prop="uuid"
              label="id">
            </el-table-column>
            <el-table-column
              prop="departid"
              label="部门">
            </el-table-column>
            </template>
          </small-table>
        </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.description"
            autocomplete="off"
          ></el-input>
        </el-form-item>
      </el-form>

      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleContinue">继续选择</el-button>
        <el-button type="primary" @click="handleConfirm">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  createBatch  } from "@/request/batch";
import SmallTable from "@/components/SmallTable";
export default {
  components: {
    SmallTable
  },
  props: {
    departInfo: {
      type: Object,
      default: function() {
        return {
          id: 1
        }
      }
    },
    machines: {
      type: Array,
    }
  },
  data() {
    return {
      machineArr: [],
      tHeight: 180,
      form: {
        batchName: "",
        description: ""
      },
      rules: {
        batchName: [{ 
          required: true, 
          message: "请填写批次名称", 
          trigger: "blur" 
        }]
      },
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
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click",{isBatch:false});
    },
    handleContinue() {
      this.$emit("click",{isBatch:true});
    },
    handleConfirm() {
      let machineuuids = [];
      let deptids = [];
      let deptNames = [];
      this.machineArr.forEach(item => {
        machineuuids.push(item.uuid);
        deptids.push(item.departid+'');
        deptNames.push(item.departname);
      })
      let data = {
        'Name': this.form.batchName, 
        'Descrip': this.form.description, 
        'Manager': this.$store.getters.userName, 
        "DepartID": [...new Set(deptids)],
        "DepartName": deptNames,
        "Machine": machineuuids || [],
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          createBatch(data)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click",{isBatch:false});
                this.$refs.form.resetFields();
                this.$message.success(res.data.msg);
              } else {
                this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("创建失败，请检查输入内容");
            });
        }
      });
    },
  },
};
</script>
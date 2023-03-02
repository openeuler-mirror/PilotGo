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
  Date: 2022-04-25 10:27:45
  LastEditTime: 2022-05-19 11:22:58
 -->
<template>
  <div class="content">
    <div class="dept panel">
      <div class="title">部门列表</div>
      <ky-tree ref="tree" :getData="getChildNode" :showSelect="showSelect" :showEdit="showChange"
        @checkClick="handleCheck" @nodeClick="handleSelectDept">
      </ky-tree>
    </div>
    <div class="info panel">
      <el-form :model="form" :rules="rules" ref="form" label-width="100px">
        <el-form-item label="批次名称:" prop="batchName">
          <el-input class="ipInput" type="text" size="medium" v-model="form.batchName" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input class="ipInput" type="text" size="medium" v-model="form.description" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <el-transfer :data="initMac" filterable filter-placeholder="请输入关键字" :titles="['备选项', '已选项']" v-model="targetMac">
        <el-button class="transfer-footer" slot="left-footer" type="primary" plain size="small"
          @click="handleReset">重置</el-button>
        <el-button class="transfer-footer" slot="right-footer" type="primary" plain size="small"
          @click="handleConfirm">创建</el-button>
      </el-transfer>
    </div>
  </div>
</template>

<script>
import kyTree from "@/components/KyTree";
import AuthButton from "@/components/AuthButton";
import { createBatch } from "@/request/batch";
import { getMacIps, getChildNode } from "@/request/cluster";
export default {
  name: "CreateBatch",
  components: {
    kyTree,
    AuthButton,
  },
  data() {
    return {
      showChange: false,
      showSelect: true,
      isNodeCheck: true,
      flag: 0,
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
      initMac: [], // 备选项
      targetMac: [], // 已选项
      machineArr: [],
      choosedDept: [],
      form: {
        batchName: "",
        description: ""
      },
      checkedNode: [],
      searchData: {
        DepartId: 1,
      },
      showSelect: true,
    };
  },
  mounted() {
  },
  methods: {
    getChildNode,
    handleReset() {
      this.flag = 0;
      this.initMac = [];
      this.choosedDept = [];
      this.$refs.tree.setCheckedKeys([])
    },
    handleSelectDept(data) {
      if (this.choosedDept.indexOf(data.id) === -1) {
        this.choosedDept.push(data.id);
      } else {
        return;
      }
      if (this.flag == 0 || this.isNodeCheck) {
        this.initMac = [];
        this.flag++;
      }
      this.isNodeCheck = false;
      getMacIps({ DepartId: data.id }).then(res => {
        if (res.data.code === 200) {
          res.data.data.forEach(item => {
            this.initMac.push({
              key: item.uuid,
              deptId: item.departid,
              label: item.ip + '-' + item.departname,
              macId: item.id,
              disabled: false
            })
          })
        }
      })

    },
    handleFilter(params) {
      // 处理节点选择状态
      if (params.checked) {
        this.initMac.push({
          key: params.data.id,
          label: params.data.label,
          disabled: false
        });
      } else {
        let delIndex = this.initMac.map((item, index) => {
          if (item.key == params.data.id) {
            return index;
          }
        }).filter(item => item >= 0)[0];
        this.initMac.splice(delIndex, 1);
      }
    },
    handleCheck(params) {
      if (this.flag === 0 || !this.isNodeCheck) {
        this.initMac = [];
        this.flag++;
      }
      this.isNodeCheck = true;
      this.handleFilter(params);
    },
    createByList() {
      let deptids = [];
      let macIds = [];
      let checkedIp = [];
      this.targetMac.forEach(uuid => {
        checkedIp.push(...this.initMac.filter(item => item.key === uuid));
      });
      checkedIp.forEach(item => {
        macIds.push(item.macId);
        deptids.push(item.deptId);
      })
      return {
        Name: this.form.batchName,
        Description: this.form.description,
        Manager: this.$store.getters.userName,
        DepartID: [...new Set(deptids)],
        Machines: [...new Set(macIds)]
      }
    },
    createByNode() {
      return {
        Name: this.form.batchName,
        Description: this.form.description,
        Manager: this.$store.getters.userName,
        deptids: this.targetMac,
      }
    },
    handleConfirm() {
      let params = this.isNodeCheck ? this.createByNode() : this.createByList();
      this.$refs.form.validate((valid) => {
        if (valid) {
          createBatch(params)
            .then((res) => {
              if (res.data.code === 200) {
                this.$refs.form.resetFields();
                this.initMac = []
                this.targetMac = []
                this.$message.success(res.data.msg);
              } else {
                this.$message.error(res.msg);
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

<style scoped lang="scss">
.content {
  width: 100%;
  display: flex;
  justify-content: space-around;

  .dept {
    height: 100%;
    width: 20%;
    display: inline-block;

    .title {
      width: 100%;
      height: 8%;
      font-weight: bold;
      border-radius: 6px 6px 0 0;
      color: #fff;
      background: rgb(45, 69, 153);
      display: flex;
      justify-content: center;
      align-items: center;
    }
  }

  .info {
    .el-form {
      width: 56%;
      height: 18%;
    }

    width: 78%;
    height: 100%;
    float: right;
  }
}
</style>

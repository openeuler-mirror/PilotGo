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
  LastEditTime: 2022-03-21 16:25:05
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="cluster">
    <router-view v-if="$route.meta.breadcrumb"></router-view>
    <div v-if="!$route.meta.breadcrumb">
    <div class="dept panel">
      <ky-tree :getData="getChildNode" ref="tree" @nodeClick="handleSelectDept"></ky-tree>
    </div>
    <div class="info panel">
      <ky-table
        class="cluster-table"
        ref="table"
        :getData="getClusters"
        :searchData="searchData"
        :treeNodes="checkedNode"
      >
        <template v-slot:table_search>
          <div>{{ departName }}</div>
        </template>
        <template v-slot:table_action>
          <auth-button name="rpm_install"  @click="handleIssue" v-show="!isBatch" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> rpm下发 </auth-button>
          <auth-button name="rpm_uninstall"  @click="handleUnInstall" v-show="!isBatch" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> rpm卸载 </auth-button>
          <auth-button name="create_batch"  @click="handleAddBatch" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 创建批次 </auth-button>
          <el-popconfirm title="确定删除所选项目吗？" @confirm="handleDelete">
            <auth-button name="cluster_delete"  slot="reference" v-show="!isBatch" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </auth-button>
          </el-popconfirm>
        </template>
        <template v-slot:table>
          <el-table-column label="ip">
            <template slot-scope="scope">
              <router-link :to="$route.path + scope.row.uuid">
                {{ scope.row.ip }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="departname" label="部门"> 
          </el-table-column>
          <el-table-column prop="cpu" label="cpu" width="130"> 
          </el-table-column>
          <el-table-column label="状态">
            <template slot-scope="scope">
              <span v-if="scope.row.state == 1">正常</span>
              <span v-if="scope.row.state == 2">离线</span>
              <span v-if="scope.row.state == 3">空闲</span>
            </template>
          </el-table-column>
           <el-table-column prop="systeminfo" label="系统信息"> 
          </el-table-column>
          <el-table-column label="防火墙配置" width="120">
            <template slot-scope="scope">
              <el-button
                size="mini"
                @click="handleFireWall(scope.row.ip)">
                <em class="el-icon-edit-outline"></em>
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button size="mini" type="primary" plain 
                @click="handleUpdateIp(scope.row.ip)"> 
                编辑 </el-button>
            </template>
          </el-table-column>
        </template>
      </ky-table>
    </div>

    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="560px"
    >
     <update-form v-if="type === 'update'" :ip='ip' @click="handleClose"></update-form>   
     <batch-form v-if="type === 'batch'" :departInfo='departInfo' :machines='machines' @click="handleClose"></batch-form>   
     <rpm-issue v-if="type === 'issue'" :acType='title' :machines='machines' @click="handleClose"></rpm-issue>   
    </el-dialog>
  </div>
  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import kyTree from "@/components/KyTree";
import AuthButton from "@/components/AuthButton";
import UpdateForm from "./form/updateForm";
import BatchForm from "./form/batchForm";
import RpmIssue from "./form/rpmIssue";
import { getClusters, deleteIp, getChildNode } from "@/request/cluster";
export default {
  name: "Cluster",
  components: {
    UpdateForm,
    BatchForm,
    RpmIssue,
    kyTable,
    kyTree,
    AuthButton,
  },
  data() {  
    return {
      title: '',
      type: '',
      ip: '',
      acType: '',
      isBatch: false,
      checkedNode: [],
      departName: '',
      departInfo: {},
      machines: [],
      display: false,
      disabled: false,
      searchData: {
        DepartId: 1,
        showSelect: true,
      },
    };
  },
  mounted() {
    this.departName = '机器列表';
    this.showSelect = ['0','1'].includes(this.$store.getters.userType) ? true : false;
  },
  methods: {
    getClusters,
    getChildNode,
    handleClose(params) {
      this.display = false;
      this.title = "";
      this.type = "";
      if(params.isBatch) {
        this.isBatch = true;
      } else {
        this.isBatch = false;
        this.machines = [];
        this.$refs.table.handleSearch();
      }
    },
    handleUpdateIp(ip) {
      this.display = true;
      this.title = "编辑IP";
      this.type = "update"; 
      this.ip = ip;
    },
    handleAddBatch() {
      this.display = true;
      this.title = "创建批次";
      this.type = "batch"; 
      this.machines = this.$refs.table.selectRow.rows;
    },
    handleIssue() {
      this.machines = [];
      this.display = true;
      this.title = "软件包下发";
      this.type = "issue"; 
      this.machines = this.$refs.table.selectRow.rows;
    },
    handleUnInstall() {
      this.machines = [];
      this.display = true;
      this.title = "软件包卸载";
      this.type = "issue"; 
      this.machines = this.$refs.table.selectRow.rows;
    },
    handleDelete() {
      let ids = this.$refs.table.selectRow.rows[0];
      deleteIp({ uuid: ids }).then((res) => {
        if (res.data.code === 200) {
          this.$refs.table.handleSearch();
          this.$message.success("删除成功");
        } else {
          this.$message.success("删除失败");
        }
      });
    },
    handleSelectDept(data) {
      if(data) {
        this.departName = data.label + '机器列表';
        this.searchData.DepartId = data.id;
        this.departInfo = data;
        this.$refs.table.handleSearch();
      }
    },
    handleNodeCheck(data) {
      this.checkedNode = [];
      if(data) {
        this.checkedNode = data.checkedKeys;
      }
    },
    handleFireWall(ip) {
      this.$router.push({
        name: 'Firewall',
        query: { ip: ip }
      })
    }
  },
};
</script>

<style scoped lang="scss">
.cluster {
  height: 94%;
  margin-top: 10px;
  .dept {
    width: 36%;
    display: inline-block;
  }
  .info {
    width: 60%;
    float: right;
  }

  .deleteHostText {
    margin-left: 10px;
    .del-host {
      color: red;
      font-weight: 600;
      font-size: 18px;
    }
  }
}
</style>

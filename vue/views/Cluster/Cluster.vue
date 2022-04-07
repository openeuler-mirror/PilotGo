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
  LastEditTime: 2022-04-07 14:50:40
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="cluster">
    <router-view v-if="$route.meta.breadcrumb"></router-view>
    <div style="width:100%;height:100%" v-if="!$route.meta.breadcrumb">
    <div class="dept panel">
      <ky-tree :getData="getChildNode" :showEdit="showChange" ref="tree" @nodeClick="handleSelectDept"></ky-tree>
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
          <auth-button name="create_batch"  @click="handleAddBatch" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 创建批次 </auth-button>
          <el-dropdown trigger="click">
            <el-button class="kylin-item-button">
              批量配置<i class="el-icon-arrow-down el-icon--right"></i>
            </el-button>
            <el-dropdown-menu slot="dropdown">
              <auth-drop @click.native="handleChange" name="rpm_install"> 变更部门 </auth-drop>
              <auth-drop @click.native="handleIssue" name="rpm_install"> rpm下发 </auth-drop>
              <auth-drop @click.native="handleUnInstall" name="rpm_uninstall"> rpm卸载 </auth-drop>
              <!-- <auth-drop @click.native="handleFireWall" name="rpm_uninstall"> 防火墙配置 </auth-drop> -->
            </el-dropdown-menu>
          </el-dropdown>
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
            <!-- <template slot-scope="scope">
              {{scope.row.departname}}
              <span v-if="showChange" title="变更部门" class="el-icon-edit-outline deptchange" @click="handleChange(scope.row)"></span>
            </template> -->
          </el-table-column>
          <el-table-column prop="cpu" label="cpu" width="220"> 
          </el-table-column>
          <el-table-column label="状态">
            <template slot-scope="scope">
              <span v-if="scope.row.state == 1">在线</span>
              <span v-if="scope.row.state == 2">离线</span>
              <span v-if="scope.row.state == 3">未分配</span>
            </template>
          </el-table-column>
           <el-table-column prop="systeminfo" label="系统信息"> 
          </el-table-column>
          <!-- <el-table-column label="防火墙配置">
            <template slot-scope="scope">
              <el-button
                size="mini"
                @click="handleFireWall(scope.row.ip)">
                <em class="el-icon-setting"></em>
              </el-button>
            </template>
          </el-table-column> -->
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button size="mini" type="primary" plain 
                @click="handleProme(scope.row.ip)">               
                监控
              </el-button>
              <!-- <el-button size="mini" type="primary" plain 
                @click="handleFireWall(scope.row.ip)">
                防火墙
              </el-button> -->
              <!-- <el-button size="mini" type="primary" plain 
                @click="handleChange(scope.row)"> 
                变更部门 </el-button> -->
            </template>
          </el-table-column>
        </template>
      </ky-table>
    </div>

    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="760px"
    >
     <update-form v-if="type === 'update'" :ip='ip' @click="handleClose"></update-form>   
     <change-form v-if="type === 'change'" :machines='batchMAcs' @click="handleClose"></change-form>   
     <batch-form v-if="type === 'batch'" :departInfo='departInfo' :machines='batchMAcs' @click="handleClose"></batch-form>   
     <rpm-issue v-if="type === 'issue'" :acType='title' :machines='machines' @click="handleClose"></rpm-issue>   
    </el-dialog>
  </div>
  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import kyTree from "@/components/KyTree";
import AuthButton from "@/components/AuthButton";
import AuthDrop from "@/components/AuthDrop";
import UpdateForm from "./form/updateForm";
import BatchForm from "./form/batchForm";
import ChangeForm from "./form/changeForm";
import RpmIssue from "./form/rpmIssue";
import { getClusters, deleteIp, getChildNode } from "@/request/cluster";
export default {
  name: "Cluster",
  components: {
    UpdateForm,
    BatchForm,
    ChangeForm,
    RpmIssue,
    kyTable,
    kyTree,
    AuthButton,
    AuthDrop,
  },
  data() {  
    return {
      title: '',
      type: '',
      ip: '',
      row: {},
      acType: '',
      showChange: false,
      isBatch: false,
      checkedNode: [],
      departName: '',
      departInfo: {},
      machines: [],
      batchMAcs: [],
      display: false,
      disabled: false,
      searchData: {
        DepartId: 1,
        showSelect: true,
      },
    };
  },
  mounted() {
    this.showChange = true;//['0','1'].includes(this.$store.getters.userType);
    getClusters({DepartId: 1}).then(res => {
      if(res.data.code === 200 && res.data.total !== 0) {
        let name = res.data.data[0].departname;
        this.departName = name === '' ? '机器列表' : name + '机器列表';
      } else {
        this.departName = "机器列表"
      }
    })
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
    handleChange() {
      this.display = true;
      this.title = "变更部门";
      this.type = "change"; 
      this.machines = this.$refs.table.selectRow.rows;
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
    },
    handleProme(ip) {
      this.$store.dispatch('setSelectIp', ip)
      this.$router.push({
        name: 'Prometheus',
        query: { ip: ip }
      })
    }
  },
  watch: {
    machines: function(newValue,oldValue) {
      this.batchMAcs = newValue.concat(oldValue)
    }
  }
};
</script>

<style scoped lang="scss">
.cluster {
  height: 94%;
  display: flex;
  margin-top: 10px;
  .dept {
    height: 100%;
    width: 20%;
    display: inline-block;
  }
  .info {
    width: 78%;
    height: 100%;
    float: right;
    .deptchange {
      cursor: pointer;
    }
    .deptchange:hover {
      color: rgb(108, 173, 228)
    }
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

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
 LastEditTime: 2023-09-08 16:27:01
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="cluster">
    <transition name="fade-transform" mode="out-in">
      <router-view v-if="$route.name == 'MacDetail' || $route.name == 'createBatch' || $route.name == 'Prometheus'">
      </router-view>
    </transition>
    <div style="width:100%;height:100%" v-if="$route.name == 'macList'">
      <div class="dept panel">
        <ky-tree :getData="getChildNode" :showEdit="showChange" ref="tree" @nodeClick="handleSelectDept"
          @refresh="handleRefresh"></ky-tree>
        <span class="sourceBtn" @click="getSourcePool">未分配资源池</span>
      </div>
      <div class="info panel">
        <ky-table class="cluster-table" ref="table" :isSource="isSource" :showSelect="showSelect" :getData="getClusters"
          :getSourceData="getSourceMac" :searchData="searchData" :treeNodes="checkedNode">
          <template v-slot:table_search>
            <div>{{ departName }}</div>
          </template>
          <template v-slot:table_action>
            <auth-button name="dept_change" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"
              @click="handleChange"> 变更部门 </auth-button>
            <auth-button name="dept_change" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"
              @click="handleDelete"> 删除 </auth-button>
          </template>
          <template v-slot:table>
            <el-table-column label="ip">
              <template slot-scope="scope">
                <span class="ipLink" @click="handleDetail(scope.row)" title="查看机器详情">
                  {{ scope.row.ip }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="departname" label="部门">
            </el-table-column>
            <el-table-column prop="cpu" label="cpu">
            </el-table-column>
            <el-table-column label="状态">
              <template slot-scope="scope">
                <state-dot :state="scope.row.state"></state-dot>
              </template>
            </el-table-column>
            <el-table-column prop="systeminfo" label="系统">
            </el-table-column>
          </template>
        </ky-table>
      </div>

      <el-dialog :title="title" top="10vh" :before-close="handleClose" :visible.sync="display" :width="dialogWidth"
        :fullscreen="isFull">
        <update-form v-if="type === 'update'" :ip='ip' @click="handleClose"></update-form>
        <change-form v-if="type === 'change'" :machines='batchMAcs' @click="handleClose"></change-form>
      </el-dialog>
    </div>
  </div>
</template>

<script>
import kyTree from "@/components/KyTree";
import AuthDrop from "@/components/AuthDrop";
import StateDot from "@/components/StateDot";
import UpdateForm from "./form/updateForm";
import ChangeForm from "./form/changeForm";
import RpmIssue from "./form/rpmIssue";
import { getClusters, deleteIp, getChildNode, getSourceMac } from "@/request/cluster";
export default {
  name: "Cluster",
  components: {
    UpdateForm,
    ChangeForm,
    RpmIssue,
    kyTree,
    AuthDrop,
    StateDot,
  },
  data() {
    return {
      title: '',
      type: '',
      ip: '',
      row: {},
      acType: '',
      isFull: false,
      dialogWidth: '760px',
      showChange: false,
      checkedNode: [],
      departName: '',
      departInfo: {},
      machines: [],
      batchMAcs: [],
      display: false,
      disabled: false,
      searchData: {
        DepartId: this.$store.getters.UserDepartId,
      },
      isSource: 1,
      showSelect: true,
    };
  },
  mounted() {
    this.showChange = [0].includes(this.$store.getters.userType);
    this.departName = this.$store.getters.tableTitle || '机器列表';
  },
  watch: {
    machines: function (newValue, oldValue) {
      this.batchMAcs = newValue.concat(oldValue)
    },
    '$route': {
      handler() {
        if (this.$route.name == 'MacDetail') {
          this.departName = "机器列表";
        }
      }
    }
  },
  methods: {
    getClusters,
    getSourceMac,
    getChildNode,
    handleClose(params) {
      this.display = false;
      this.title = "";
      this.type = "";
      this.$refs.table.handleSearch();
    },
    handleRefresh() {
      // 节点树更新
      this.$refs.table.handleSearch();
    },
    handleChange() {
      this.display = true;
      this.title = "变更部门";
      this.type = "change";
      this.dialogWidth = "760px";
      this.machines = this.$refs.table.selectRow.rows;
    },
    handleDelete() {
      let uuidArray = [];
      this.$refs.table.selectRow.rows.forEach(item => {
        uuidArray.push(item.uuid);
      })
      this.$confirm('是否确认删除该机器?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteIp({ deluuid: uuidArray }).then(res => {
          if (res.data.code === 200) {
            this.$message.success(res.data.msg);
            this.$refs.table.handleSearch();
          } else {
            this.$message.error(res.data.msg);
          }
        })
      }
      ).cache(() => {
        console.log('delete ip failed')
      });
    },
    handleUpdateIp(ip) {
      this.display = true;
      this.title = "编辑IP";
      this.type = "update";
      this.dialogWidth = "760px";
      this.ip = ip;
    },
    handleSelectDept(data) {
      if (data) {
        this.departName = data.label + '机器列表';
        this.searchData.DepartId = data.id;
        this.departInfo = data;
        this.$refs.table.handleSearch(this.searchData);
      }
    },
    handleNodeCheck(data) {
      this.checkedNode = [];
      if (data) {
        this.checkedNode = data.checkedKeys;
      }
    },
    getSourcePool() {
      this.isSource = (Math.random() + 1) * 100;
      this.departName = "未分配资源池";
    },
    handleDetail(row) {
      this.handleSelectIP(row.ip);
      this.$router.push({
        path: `/cluster/macList/${row.uuid}`,
      })
    },
    handleSelectIP(ip) {
      this.$store.dispatch('setSelectIp', ip)
      this.$store.commit('SET_IMMUTABLE', false)
    },
    handleProme(ip) {
      this.handleSelectIP(ip);
      this.$router.push({
        name: 'Prometheus',
        query: { ip: ip }
      })
    }
  },
};
</script>

<style scoped lang="scss">
.cluster {
  width: 100%;
  height: 100%;
  display: flex;

  .dept {
    height: 100%;
    width: 20%;
    display: inline-block;

    .sourceBtn {
      display: block;
      background: rgb(255, 191, 0);
      color: #fff;
      width: 80%;
      padding: 4px;
      border-radius: 6px;
      margin: 10% auto 0;
      text-align: center;
      cursor: pointer;
    }
  }

  .info {
    width: 78%;
    height: 100%;
    float: right;

    .ipLink {
      color: rgb(64, 158, 255);
      cursor: pointer;

      &:hover {
        color: rgb(242, 150, 38);
      }
    }

    .deptchange {
      cursor: pointer;
    }

    .deptchange:hover {
      color: rgb(108, 173, 228)
    }
  }

  .term {
    width: 100%;
    height: 100%;

    .term_head {
      position: relative;
      width: 100%;
      font-size: 16px;
      border: 1px solid rgb(109, 123, 172);
      border-radius: 10px 10px 0 0;
      background: rgb(109, 123, 172);
      color: #fff;
      display: flex;
      justify-content: space-between;
    }

    .termTitle {
      display: inline-block;
      width: 30%;
      padding: 0.3% 0 0 1%;
    }

    .closeChart {
      display: inline-block;
      width: 4px;
      height: 4px;
      position: absolute;
      top: 2%;
      right: 2%;
      z-index: 1;
      cursor: pointer;
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

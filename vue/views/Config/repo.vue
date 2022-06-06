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
  Date: 2022-01-17 09:41:31
  LastEditTime: 2022-06-06 17:13:43
 -->
<template>
 <div class="panel" style="height:100%">
   <transition name="fade-transform" mode="out-in">
    <router-view v-if="$route.name == 'sysctl' || $route.name == 'libconfig'"></router-view>
   </transition>
   <div class="info" v-if="$route.name == 'repo'">
      <div class="repoList" v-show="!showCompare">
        <div class="select"><span>请选择机器：</span>
        <el-autocomplete
          style="width:50%"
          class="inline-input"
          v-model="macIp"
          :fetch-suggestions="querySearch"
          placeholder="请输入ip关键字"
          @select="handleSelect"
        ></el-autocomplete></div>
        <ky-table
          :getData="getRepos"
          :searchData="searchData"
          ref="table"
        >
        <template v-slot:table_search>
          <div>repo列表</div>
        </template>
        <template v-slot:table>
          <el-table-column prop="path" label="repo路径" width="140">
          </el-table-column>
          <el-table-column prop="name" label="repo名" width="120">
            <template slot-scope="scope">
              <span title="详情" class="repoDetail" @click="handleDetail(scope.row)">{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button type="text"  size="mini" @click="handleEdit(scope.row)">编辑</el-button>   
              <el-button type="text"  size="mini" @click="handleDowlond(scope.row)">下载到库</el-button>   
              <el-button type="text"  size="mini" @click="handleCompare(scope.row)">对比</el-button>   
            </template>
          </el-table-column>
        </template>
      </ky-table>
      </div>
      <div class="rightInfo" :style="{'width':compareWidth}">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>{{rightTtile}}</span>
            <el-button v-if="showEdit" style="float: right; padding: 3px 0" type="text" @click="handleEditConfirm">确认</el-button>
            <el-button v-if="showCompare" style="float: right; padding: 3px 0" type="text" @click="handleDetail">退出</el-button>
          </div>
          <div class="edit" v-if="showEdit">
            <br/>
            <vue-editor v-model="detail" id="container" :editor-toolbar="customToolbar"/>
          </div>
          <div class="detail" v-if="showDetail">
            <pre>{{detail || '点击文件名查看详情'}}</pre>
          </div>
          <div class="compare" v-if="showCompare">
            <code-diff
              :old-string="init"
              :new-string="detail"
              file-name="test.txt"
              output-format="side-by-side"/>
          </div> 
        </el-card>
        
      </div>
      <el-dialog 
        :title="title"
        :before-close="handleClose" 
        :visible.sync="display" 
        :width="dialogWidth">
        <download-form v-if="type === 'download'" :uuid="searchData.uuid" :row="row" @click="handleClose"></download-form>
      </el-dialog>
   </div>
 </div>
</template>
<script>
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import DownloadForm from "./form/downloadForm.vue";
import { VueEditor } from "vue2-editor";
import { getallMacIps } from '@/request/cluster'
import { getRepos, getRepoDetail, updateRepo } from "@/request/config"
export default {
  name: "repoConfig",
  components: {
    kyTable,
    AuthButton,
    DownloadForm,
    VueEditor,
  },
  data() {
    return {
      title: '',
      type: '',
      macIp: '',
      row: {},
      compareWidth: '54%',
      init: 'file empty',
      detail: '',
      oldValue: 'asdsd/n',
      newValue: 'axc/nasdaasd',
      rightTtile: '文件详情',
      uuid: 'c11d8252-eac3-4c07-ac0e-99d8080d0a05',
      dialogWidth: '70%',
      display: false,
      showEdit: false,
      showDetail: true,
      showCompare: false,
      rowData: [],
      searchData: {
        uuid: 'c11d8252-eac3-4c07-ac0e-99d8080d0a05'
      },
      customToolbar: [
      ["bold", "italic", "underline"],
      [{'align': ["","center", "right","justify"]}],
    ]
    }
  },
  mounted() {
    getallMacIps().then(res => {
      this.ips = [];
      this.ipData = [];
      if(res.data.code === 200) {
        this.ips = res.data.data && res.data.data;
        this.ips.forEach(item => {
            this.ipData.push({'value':item.ip_dept, 'uuid':item.uuid, 'ip':item.ip})
          })
      }
    })
  },
  methods: {
    getRepos,
    querySearch(queryString, cb) {
      var ipData = this.ipData;
      var results = queryString ? ipData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): ipData;
      cb(results);
    },
    handleSelect(mac) {
      this.searchData.uuid = mac && mac.uuid;
      this.macIp = mac && mac.ip;
      this.$refs.table.handleSearch();
    },
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = '';
      this.$refs.table.handleSearch();
    },
    getDetail(row) {
      getRepoDetail({uuid: this.searchData.uuid,file: row.path+'/'+row.name}).then(res => {
        if(res.data.code === 200) {
          this.detail = res.data && res.data.data.file;
        }
      })
    },
    handleDetail(row) {
      this.compareWidth = '54%';
      this.showEdit = false;
      this.showDetail = true;
      this.showCompare = false;
      if(row) {
        this.rightTtile = '文件详情'+' - '+ row.name;
        this.getDetail(row);
      } else {
        this.rightTtile = '文件详情';
      }
    },
    handleCompare(row) {
      this.showEdit = false;
      this.showDetail = false;
      this.showCompare = true;
      this.rightTtile = '文件对比';
      this.compareWidth = '96%'
    },
    handleDowlond(row) {
      this.display = true;
      this.title = '文件下载到库';
      this.type = 'download';
      this.row = row;
    },
    handleEdit(row) {
      this.rightTtile = '编辑文件' + ' - ' + row.name;
      this.showEdit = true;
      this.showDetail = false;
      this.showCompare = false;
      this.row = row;
      this.getDetail(row);
    },
    handleEditConfirm() {
      let params = {
        path: this.row.path,
        uuid: this.searchData.uuid,
        name: this.row.name,
        file: this.detail.replace(/<[^>]+>/g, ''),
        ip: this.macIp,
      }
      updateRepo(params).then(res => {
        if(res.data.code === 200) {
          this.$message.success(res.data.msg)
        } else {
          this.$message.error(res.data.msg)
        }
      })

    },
  }
}
</script>
<style scoped lang="scss">
.info {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: space-evenly;
  .repoList {
    width: 44%;
    height: 100%;
    .select {
      height: 6%;
      margin: 4px;
    }
    .ky-table {
      height: 90%;
      .header {
        height: 6%;
      }
    }
  }
  .rightInfo {
    width: 54%;
    height: 100%;
    .edit {
      width: 100%;
      height: 90%;
      overflow: auto;
    }
    .detail {
      height: 100%;
      pre {
        height: 90%;
        white-space: pre-line;
        overflow-y: auto;
      }
    }
    .compare {
      width: 100%;
      height: 96%;
      overflow: auto;
    }
  }
  .repoDetail {
    cursor: pointer;
  }
  .repoDetail:hover {
    color: rgb(64, 158, 255);
  }
}
</style>

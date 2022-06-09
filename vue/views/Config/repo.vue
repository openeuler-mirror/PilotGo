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
  LastEditTime: 2022-06-09 16:00:54
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
          style="width:40%"
          class="inline-input"
          v-model="macIp"
          :fetch-suggestions="querySearch"
          placeholder="请输入ip关键字"
          @select="handleSelect"
        ></el-autocomplete></div>
        <ky-table
          :getData="getAllFiles"
          :searchData="searchData"
          ref="table"
          :isRowClick="openRowClick"
        >
        <template v-slot:table_search>
          <div>配置文件列表</div>
        </template>
        <template v-slot:table>
          <el-table-column type="expand" width="2">
            <template>
              <el-steps 
                :active="files.length"
                align-center 
                v-loading="loading" 
                element-loading-text="文件请求中"
                element-loading-spinner="el-icon-loading"
                element-loading-background="rgba(255, 255, 255, 0.8)">
                <el-step v-for="item in files" 
                  @click.native="selectRollbackFile(item)"
                  class="step"
                  :key="item.id" 
                  :description="item.name"/>
              </el-steps>
            </template>
          </el-table-column>
          <el-table-column prop="path" label="文件路径" width="220">
          </el-table-column>
          <el-table-column prop="name" label="文件名" width="160">
            <template slot-scope="scope">
              <span title="详情" class="repoDetail" @click="handleDetail(scope.row)">{{scope.row.name}}</span>
            </template>
          </el-table-column>
          <el-table-column prop="type" label="文件类型">
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-button type="text"  size="mini" @click="handleEdit(scope.row)">编辑</el-button>   
              <el-button type="text"  size="mini" @click="handleDowlond(scope.row)">下载到库</el-button>   
              <el-button type="text"  size="mini" @click="handleCompare(scope.row)">对比</el-button>   
              <el-button type="text"  size="mini" @click="getRollbackFile(scope.row)">回滚</el-button>   
            </template>
          </el-table-column>
        </template>
      </ky-table>
      </div>
      <div class="rightInfo" :style="{'width':compareWidth}">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>{{rightTtile}}</span>
             <el-popconfirm 
              v-if="showRollback"
              title="确定回滚当前文件版本？"
              cancel-button-type="default"
              confirm-button-type="danger"
              @confirm="handleRollBack">
              <el-button  slot="reference" style="float: right; padding: 3px 0" type="text">回滚</el-button>
             </el-popconfirm>
              <el-popconfirm 
              v-if="showEdit"
              title="确定使用当前文件？"
              cancel-button-type="default"
              confirm-button-type="danger"
              @confirm="handleEditConfirm">
              <el-button  slot="reference" style="float: right; padding: 3px 0" type="text">确认</el-button>
             </el-popconfirm>
            <el-button v-if="showCompare" style="float: right; padding: 3px 0" type="text" @click="handleDetail">退出</el-button>
          </div>
          <div class="editor" v-if="showEdit">
            <quill-editor
              class="ql-editor"
              :options="editorOptions"
              :content="detail"
              ref="myQuillEditor"
              @change="onEditorChange($event)">
          </quill-editor>
          </div>
          <div class="detail" v-if="showDetail || showRollback">
            <pre>{{detail || '点击文件名查看详情'}}</pre>
          </div>
          <div class="compare" v-if="showCompare">
           <compare-form :detail="detail" :files="files"/>
          </div> 
        </el-card>
        
      </div>
      <el-dialog 
        :title="title"
        :before-close="handleClose" 
        :visible.sync="display" 
        :width="dialogWidth">
        <download-form v-if="type === 'download'" :ipDept="ipDept" :macIp="macIp" :uuid="searchData.uuid" :row="row" @click="handleClose"></download-form>
      </el-dialog>
   </div>
 </div>
</template>
<script>
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import DownloadForm from "./form/downloadForm.vue";
import CompareForm from "./form/compareForm.vue";
import { quillEditor } from 'vue-quill-editor'
import 'quill/dist/quill.snow.css'
import { getallMacIps } from '@/request/cluster'
import { getAllFiles, getRepoDetail, updateRepo, getFileHistories, fileRollback } from "@/request/config"
export default {
  name: "repoConfig",
  components: {
    kyTable,
    AuthButton,
    DownloadForm,
    CompareForm,
    quillEditor,
  },
  data() {
    return {
      loading: true,
      title: '',
      type: '',
      macIp: '',
      ipDept: '',
      row: {},
      compareWidth: '46%',
      detail: '',
      content: '',
      files: [],
      historyFiles: [
        {id:1,name:'fairy-2022/4/21',description:'test file'},
        {id:2,name:'test-2022/4/21',description:'test file2'},
        {id:3,name:'test33-2022/4/21',description:'test file3'},
        {id:4,name:'test33-2022/4/21',description:'test file3'},
        {id:5,name:'test33-2022/4/21',description:'test file3'},
        {id:6,name:'test33-2022/4/21',description:'test file3'},
      ],
      rightTtile: '文件详情',
      dialogWidth: '70%',
      display: false,
      showEdit: false,
      showDetail: true,
      showCompare: false,
      showRollback: false,
      rowData: [],
      searchData: {
        uuid: 'c11d8252-eac3-4c07-ac0e-99d8080d0a05'
      },
      openRowClick: true,
      editorOptions: {
        modules:{
          toolbar:[
            ['bold', 'italic', 'underline'],
            ['code-block']
          ],
        },
      }
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
    getAllFiles,
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
      this.ipDept = mac && mac.value.split('-')[1];
      this.$refs.table.handleSearch();
    },
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = '';
      this.$refs.table.handleSearch();
    },
    getRollbackFile(row) {
      this.row = row;
      this.$refs.table.changeExpandRow(row);
      getFileHistories({uuid: this.searchData.uuid,name:row.name}).then(res=> {
        if(res.data.code === 200) {
          this.loading = false;
          this.files = res.data.data && res.data.data.oldfiles;
        }
      })
    },
    selectRollbackFile(currFile) {
      this.showEdit = false;
      this.showDetail = false;
      this.showCompare = false;
      this.showRollback = true;
      this.rightTtile = '文件详情'+' - '+ currFile.name;
      this.detail = currFile.file;
    },
    handleRollBack() {
      let params = {
        path: this.row.path,
        uuid: this.searchData.uuid,
        name: this.row.name,
        file: this.detail,
        ip: this.macIp,
        ipDept: this.ipDept,
      }
      fileRollback(params).then(res => {
        if(res.data.code === 200) {
          this.rightTtile = '文件详情';
          this.showDetail = true;
          this.showRollback = false;
          this.detail = '';
          this.$message.success(res.data.msg)
          this.handleClose();
        } else {
          this.$message.error(res.data.msg)
        }
      })
    },
    getDetail(row) {
      getRepoDetail({uuid: this.searchData.uuid,file: row.path+'/'+row.name}).then(res => {
        if(res.data.code === 200) {
          this.detail = res.data && res.data.data.file;
        }
      })
    },
    handleDetail(row) {
      this.compareWidth = '46%';
      this.showEdit = false;
      this.showDetail = true;
      this.showCompare = false;
      this.showRollback = false;
      if(row) {
        this.rightTtile = '文件详情'+' - '+ row.name;
        this.getDetail(row);
      } else {
        this.rightTtile = '文件详情';
      }
    },
    handleCompare(row) {
      this.getDetail(row);
      getFileHistories({uuid: this.searchData.uuid,name:row.name}).then(res=> {
        if(res.data.code === 200) {
          this.files = res.data.data && res.data.data.oldfiles;
          this.compareDisplay();
        }
      })
      
    },
    compareDisplay() {
      this.showEdit = false;
      this.showDetail = false;
      this.showCompare = true;
      this.showRollback = false;
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
      this.showRollback = false;
      this.row = row;
      this.getDetail(row);
    },
    //editor内容改变事件
    onEditorChange({ quill, html, text }) {
      this.content = text;
    },
    handleEditConfirm() {
      let params = {
        path: this.row.path,
        uuid: this.searchData.uuid,
        name: this.row.name,
        file: this.content,
        ip: this.macIp,
        ipDept: this.ipDept,
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
    width: 52%;
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
    width: 46%;
    height: 100%;
    .editor {
      width: 100%;
      height: 96%;
      .ql-editor {
        padding: 0;
        overflow: hidden;
      }
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
.step :hover {
  cursor: pointer;
}
</style>

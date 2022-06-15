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
  LastEditTime: 2022-06-14 14:43:15
  Description: provide agent log manager of pilotgo
 -->
<template>
  <el-container>
    <el-aside style="width: 10%">
      <div class="logo"> 
        <img src="../../assets/logo.png" alt=""> 
        <span>PilotGo</span>
      </div>
      <el-menu id="el-menu"
        :collapse="isCollapse"
        :unique-opened="true"
        @select="handleSelect"
        class="el-menu-vertical-demo"
        background-color="#fff"
        :default-active="activePanel">
        <sidebar-item :routes="routesData"></sidebar-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="height: 6%">
        <bread-crumb class="breadCrumb"></bread-crumb>
        <div class="header-function">
          <el-dropdown class="header-function__username">
            <div :title="username">
              <em class="el-icon-s-custom"></em>
              <span>{{username.length > 16 ? username.split('@')[0] : username}}</span>
            </div>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item @click.native="updatePwd">修改密码</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
          <div class="logOut" title="退出" @click="handleLogOut">
            <em class="el-icon-s-unfold"></em>
          </div>
        </div>
        <tags-view />
      </el-header>
      <el-main>
        <div class="bodyContent">
          <transition name="fade-transform" mode="out-in">
            <keep-alive :include="cachedViews">
              <router-view :key="key" />
            </keep-alive>
          </transition>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import SidebarItem from "./components/SidebarItem";
import BreadCrumb from "./components/BreadCrumb";
import TagsView from "./components/TagsView";
import Config from '../../../config/index.js'; //Config.dev.proxyTable['/'].target.split('//')[1]
export default {
  name: "Home",
  components: {
    SidebarItem,
    BreadCrumb,
    TagsView
  },
  data() {
    return {
      crumbs: [],
      isCollapse: false,
      // cachedViews: ['Batch','Overview','Prometheus']
    };
  }, 
  mounted() {
    this.create();
  },
  computed: {
    cachedViews() {
      return this.$store.getters.cachedViews
    },
    key() {
      return this.$route.path
    },
    routesData() {
        return this.$store.getters.getPaths
    },
    activePanel() {
        return this.$store.getters.activePanel
    },
    menuKey() {
      return this.$store.state.menuIndex;
    },
    clusterIp() {
      return this.$store.state.selectedClusterIp
        ? this.$store.state.selectedClusterIp
        : null;
    },
    username(){
      return this.$store.getters.userName;
    },
    breadCrumb() {
      this.crumbs = [];
      if(this.$route.params.id) {
        this.$route.meta.breadCrumb[1].name = this.$route.params.id
      }
    }
  },
  methods: {
    handleSelect(key) {
      this.$store.commit("SET_ACtiVEPANEL", key);
    },
    handleLogOut() {
      this.$confirm('此操作将注销登录, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch("logOut").then((res) => {
          this.$router.push("/login");
        });
      })
    },
    updatePwd() {
      this.$prompt('请输入邮箱', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
        inputErrorMessage: '邮箱格式不正确'
      }).then(({ value }) => {
        this.$message({
          type: 'success',
          message: '你的邮箱是: ' + value
        });
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '取消输入'
        });
      });
    },
    create() {
      let es = new EventSource(Config.dev.proxyTable['/'].target+'/event');
      es.addEventListener('message', event => {
          this.$notify({
          title: '警告',
          message: event.data,
          type: 'warning',
          duration: 0
        });
      });
      es.addEventListener('error', event => {
          if (event.readyState == EventSource.CLOSED) {
              console.log('event was closed');
          };
      });
      es.addEventListener('close', event => {
          console.log(event.type);
          es.close();
      });
    }

  },
};
</script>

<style lang="scss" scope>
.el-container {
  height: calc(100%);
  width: calc(100%);
  overflow: hidden;
  
  .el-aside {
    height: 100%;
    overflow: hidden;
    .logo {
      width: 100%;
      height: 10%;
      font-size: 3em;
      font-family: fantasy;
      color: rgb(241, 139, 14);
      display: flex;
      justify-content: space-around;
      align-items: center;
      position: relative;
      z-index: 1999;
      background: rgb(11, 35, 117);
      img {
        width: 50%;
      }
      span {
        display: inline-block;
        margin-left: -30%;
      }
    }
    .el-menu {
      width: 100%;
      height: 90%;
    }
    .aside-footer {
      position: absolute;
      bottom: 0;
      left: 0;
      margin-bottom: 80px;
      margin-left: 24px;
      padding: 10px 10px 10px 0px;
      cursor: pointer;

      span {
        margin-left: 10px;
        font-size: 14px;
      }
    }
  }
  .el-container {
    width: 88%;
    height: 100%;
    overflow: hidden;
    background: #fff;
    .el-header {
      position: relative;
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: solid 1px #e6e6e6;
      .header-logoName {
        height: 100%;
        font-size: 28px;
        color: #fff;
        float: left;
        display: flex;
        justify-content: space-evenly;
        align-items: center;
        img {
        height: 120%;
      }
      }
      .header-function {
        width: 20%;
        height: 100%;
        float: right;
        display: flex;
        justify-content: end;
        align-items: center;

        .header-function__username {
          width: 30px;
          height: 30px;
          text-align: center;
          color: rgb(241, 139, 14);
          font-size: 24px;
          border-radius: 30px;
          background-color: #fff;
          transition: width 500ms ease-in-out;
          span {
            display: none;
          }
        }
        .header-function__username:hover {
          border: 1px solid rgb(241, 139, 14);
          width: 70%;
          span {
            animation: text 1s 1;
            font-size: 16px;
            display: inline-block;
          }
          @keyframes text {
            0% {
              color: transparent;
            }
            100% {
              color: rgb(241, 139, 14);
            }
          }
        }
        .logOut {
          cursor: pointer;
          font-size: 24px;
          margin: 0 0 0 6px;
          color: rgb(241, 139, 14);
        }
        .logOut:hover {
          color: rgb(202, 205, 210);
        }
      }
    }
    .el-main {
      height: 92%;
      padding: 3.6% 0 0;
      .bodyContent {
        width: 100%;
        height: 100%;
      }
      .breadCrumb {
        width: 100%;
        height: 8%;
        .el-breadcrumb {
          width: 100%;
          height: 100%;
          display: flex;
          align-items: center;
        }
      }
    }
}
}
</style>

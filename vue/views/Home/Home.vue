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
  LastEditTime: 2022-02-25 17:40:15
  Description: provide agent log manager of pilotgo
 -->
<template>
  <el-container>
    <el-header>
      <div class="header-logoName">PilotGo运维平台</div>
      <div class="header-function">
        <el-dropdown class="header-function__username" trigger="click">
          <div>
            <em class="el-icon-s-custom"></em>
            <span>{{username}}</span>
          </div>
          <el-dropdown-menu slot="dropdown" >
            <el-dropdown-item @click.native="updatePwd">修改密码</el-dropdown-item>
            <el-dropdown-item @click.native="handleLogOut">注销登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-header>
    <el-container class="aside-main">
      <el-aside width="200px">
        <el-menu
          id="el-menu"
          :uniqueOpened="true"
          @select="handleSelect"
          class="el-menu-vertical-demo"
          background-color="#fff"
          :default-active="activePanel"
        >
          <sidebar-item :routes="routesData"></sidebar-item>
        </el-menu>
      </el-aside>
      <el-main>
        <bread-crumb></bread-crumb>
        <keep-alive>
          <router-view></router-view>
        </keep-alive>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import SidebarItem from "./components/SidebarItem";
import BreadCrumb from "./components/BreadCrumb";
export default {
  name: "Home",
  components: {
    SidebarItem,
    BreadCrumb,
  },
  data() {
    return {
      crumbs: [],
    };
  },
  computed: {
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
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消操作'
        });
      });
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
  },
};
</script>

<style lang="scss" scope>
.el-container {
  .el-header {
    background: url(../../assets/header-bg.png);
    background-size: cover;

    .header-logoName {
      font-size: 28px;
      color: #fff;
      float: left;
      line-height: 60px;
      margin-left: 8px;
    }
    .header-function {
      float: right;
      color: #ffffff;
      font-size: 28px;
      line-height: 60px;
      margin-right: 10px;
      .header-function__translate {
        margin-right: 30px;
      }

      .header-function__username {
        float: left;
        font-size: 18px;
        color: #ffffff;
        margin-left: 10px;
        span {
          font-size: 16px;
        }
      }
    }
  }

  .aside-main {
    height: calc(100% - 60px);
  }
  .el-aside {
    background: #fff;
    position: relative;
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

  .el-main {
    height: 100%;
    background: #fafafa;
    .cockpit {
      width: 100%;
      height: 100%;
    }
  }
}
</style>

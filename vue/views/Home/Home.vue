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
         <keep-alive>
            <router-view></router-view>
        </keep-alive>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import SidebarItem from "./components/SidebarItem";
export default {
  name: "Home",
  components: {
    SidebarItem,
  },

  data() {
    return {
      /* routesData: [
        {path: 'overview', name: 'overview', title: '概览', index: "1", icon_class: 'el-icon-location'},
        {path: 'cluster', name: 'cluster', title: '机器管理', index: "2", icon_class: 'el-icon-s-platform'},
        {path: 'batch', name: 'batch', title: '批次管理', index: "3", icon_class: 'el-icon-menu'},
        {path: 'plug_in', name: 'plug_in', title: '插件管理', index: "4",  icon_class: 'el-icon-document'},
        // {path: 'prometheus', name: 'prometheus', title: 'Prometheus', index: "5",  icon_class: 'el-icon-odometer'},
        // {path: 'cockpit', name: 'cockpit', title: 'cockpit', index: "6",  icon_class: 'el-icon-setting'},
        {path: 'usermanager', name: 'usermanager', title: '用户管理', index: "7",  icon_class: 'el-icon-user-solid'},
        {path: 'rolemanager', name: 'rolemanager', title: '角色管理', index: "8",  icon_class: 'el-icon-s-custom'},
        {path: 'firewall', name: 'firewall', title: '防火墙配置', index: "9",  icon_class: 'el-icon-s-home'},
      ] */
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
    }
  },
  methods: {
    // userInfo() {
    //     this. =true;
    // },
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
      console.log("修改密码")
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
    background: #fafafa;
    .cockpit {
      width: 100%;
      height: 100%;
    }
  }
}
</style>

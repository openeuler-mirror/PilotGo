<template>
  <el-container>
    <el-header>
      <div class="header-logoName">PilotGo运维平台</div>
      <div class="header-function">
        <el-dropdown class="header-function__username"  @command="commandHandler">
          <div>
            <i class="el-icon-s-custom"></i>
            <span>{{user}}</span>
          </div>
          <el-dropdown-menu slot="dropdown" >
            <el-dropdown-item command="userinfo" divided @click="open">修改密码</el-dropdown-item>
            <el-dropdown-item command="logout" divided>注销登录</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </div>
    </el-header>
    <el-container class="aside-main">
      <el-aside width="240px">
        <el-menu
          id="el-menu"
          :uniqueOpened="true"
          class="el-menu-vertical-demo"
          @select="handleSelect"
          background-color="#fff"
          :default-active="menuKey"
        >
          <el-menu-item index="1">
            <i class="el-icon-location"></i>
            <template #title>概览</template>
          </el-menu-item>
          <el-menu-item index="2">
            <i class="el-icon-menu"></i>
            <template #title>集群</template>
          </el-menu-item>
          <el-menu-item index="3">
            <i class="el-icon-document"></i>
            <template #title>插件管理</template>
          </el-menu-item>
          <el-menu-item index="4">
            <i class="el-icon-monitor"></i>
            <template #title>Prometheus</template>
          </el-menu-item>
          <el-menu-item index="5">
            <i class="el-icon-setting"></i>
            <template #title>cockpit</template>
          </el-menu-item>
          <el-menu-item index="6">
            <i class="el-icon-user-solid"></i>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="7">
            <i class="el-icon-user-solid"></i>
            <template #title>防火墙配置</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <overview v-show="menuKey == '1'"></overview>
        <cluster v-on:selectIp="selectClusterIp" v-show="menuKey == '2'"></cluster>
        <plug-in v-show="menuKey == '3'"></plug-in>
        <monitor v-show="menuKey == '4'"></monitor>
        <iframe class="cockpit" :src="clusterIp" v-show="menuKey == '5' && clusterIp != null">{{ clusterIp }}</iframe>
        <div v-show="menuKey == '5' && clusterIp == null">请选择集群IP</div>
        <user-man v-show="menuKey == '6'"></user-man>
        <firewall v-show="menuKey=='7'"></firewall>
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import Cluster from "./Cluster/Cluster";
import Monitor from "./Monitor/Monitor";
import Overview from "./Overview/Overview";
import PlugIn from "./Plug-in/Plug-in";
import UserMan from "./UserManager/UserMan";
import UserInfo from "./UserInfo";
import Cookies from "js-cookie";
import Firewall from "./Firewall/Firewall";

export default {
  name: "Home",
  components: {
    Firewall,
    Monitor,
    Overview,
    PlugIn,
    Cluster,
    UserMan,
    UserInfo,
  },

  data() {
    return {};
  },
  computed: {
    menuKey() {
      return this.$store.state.menuIndex;
    },
    clusterIp() {
      return this.$store.state.selectedClusterIp
        ? this.$store.state.selectedClusterIp
        : null;
    },
    user(){
      return Cookies.get('email')
    }
  },
  methods: {
    // userInfo() {
    //     this. =true;
    // },
    handleSelect(key) {
      this.$store.commit("menuKeySelect", key);
    },
    selectClusterIp(ip) {
      this.$store.commit("mutateSelectedClusterIp", ip);
      this.handleSelect("5");
    },
    commandHandler(cmd) {
      if (cmd == 'logout') {
        this.$confirm('此操作将注销登录, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$router.push({ path: "/login" });
        }).catch(() => {
          this.$message({
            type: 'info',
            message: '已取消操作'
          });
        });
      }else if (cmd == 'userinfo') {
        console.log("error submit!!");
        return false;
      }
    },
    open() {
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
    background: url(../assets/header-bg.png);
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

    .el-menu {
      i {
        margin-right: 12px;
      }
      .el-menu-item,
        .el-submenu__title {
          height: 46px;
          line-height: 46px;
          margin: 10px 0;
        }
      .el-menu-item.is-active {
        background: #f2f8ff !important;
        border-right: 3px solid #0076ff;
        color: #0076ff;
      }
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

  .el-main {
    background: #fafafa;
    .cockpit {
      width: 100%;
      height: 100%;
    }
  }
}
</style>

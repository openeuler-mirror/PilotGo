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
 LastEditTime: 2023-02-06 13:49:27
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div style="width:calc(100%);height:calc(100%)">
    <el-container>
      <el-aside style="width: 8%">
        <div class="logo">
          <img src="../../assets/logo.png" alt="">
          <!-- <span>PilotGo</span> -->
        </div>
        <el-menu id="el-menu" :collapse="isCollapse" :unique-opened="true" @select="handleSelect"
          class="el-menu-vertical-demo" background-color="#fff" :default-active="activePanel">
          <sidebar-item :routes="routesData"></sidebar-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-dialog title="修改密码" :visible.sync="dialogFormVisible" width="650px">
          <el-form :model="form" ref="form" label-width="80px" :rules="rules">
            <el-form-item label="邮箱">
              <el-input :value="username" autocomplete="off" :disabled="true"></el-input>
            </el-form-item>
            <el-form-item label="新密码" prop="password">
              <el-input type="password" v-model="form.password" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="确认密码" prop="passwordValid">
              <el-input type="password" v-model="form.passwordValid" autocomplete="off"></el-input>
            </el-form-item>
          </el-form>
          <div slot="footer" class="dialog-footer">
            <el-button @click="dialogFormVisible = false">取 消</el-button>
            <el-button type="primary" @click="confirmClick">确 定</el-button>
          </div>
        </el-dialog>
        <el-header style="height: 6%">
          <bread-crumb class="breadCrumb"></bread-crumb>
          <div class="header-function">
            <el-dropdown class="header-function__username">
              <div :title="username">
                <em class="el-icon-s-custom"></em>
                <span>{{ username.length > 16 ? username.split('@')[0] : username }}</span>
              </div>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item @click.native="update">修改密码</el-dropdown-item>
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
            <!--插件页-->
            <component v-for="item in iframeComponents" :key="item.name" :is="item.name" :url="item.url"
              :plugin_type="item.plugin_type" :name="item.name" :path="item.path" v-if="$route.path === item.path">
            </component>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import SidebarItem from "./components/SidebarItem";
import BreadCrumb from "./components/BreadCrumb";
import TagsView from "./components/TagsView";
import { updatePwd } from "@/request/user";
export default {
  name: "Home",
  components: {
    SidebarItem,
    BreadCrumb,
    TagsView
  },
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请输入密码'));
      } else {
        if (this.form.passwordValid !== '') {
          this.$refs.form.validateField('passwordValid');
        }
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.form.password) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    return {
      ws: null,
      crumbs: [],
      isCollapse: false,
      dialogFormVisible: false,
      // cachedViews: ['Batch','Overview','Prometheus']
      form: {
        name: '',
        password: '',
        passwordValid: ''
      },
      rules: {
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 3, max: 16, message: '最大长度16位', trigger: 'blur' },
          { validator: validatePass, trigger: 'blur' }
        ],
        passwordValid: [
          { validator: validatePass2, trigger: 'blur' }
        ]
      },

    }
  },
  mounted() {
    this.initSocket();
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
    iframeComponents() {
      return this.$store.getters.iframeComponents;
    },
    menuKey() {
      return this.$store.state.menuIndex;
    },
    clusterIp() {
      return this.$store.state.selectedClusterIp
        ? this.$store.state.selectedClusterIp
        : null;
    },
    username() {
      return this.$store.getters.userName;
    },
    breadCrumb() {
      this.crumbs = [];
      if (this.$route.params.id) {
        this.$route.meta.breadCrumb[1].name = this.$route.params.id
      }
    },
    socketUrl() {
      return (location.protocol === "http:" ? "ws" : "wss") + "://" +
        location.host +
        "/event";
    }
  },
  methods: {
    handleSelect(key) {
      this.$store.commit("SET_ACTIVE_PANEL", key);
    },
    handleLogOut() {
      this.$confirm('此操作将注销登录, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let params = {
          email: this.$store.getters.email,
          departName: this.$store.getters.departName
        };
        this.$store.dispatch("logOut", params).then((res) => {
          this.$router.push("/login");
        });
      }).catch(() => {

      })
    },
    update() {
      this.dialogFormVisible = true
      this.$refs.form.clearValidate();
      // this.$prompt('请输入邮箱', '提示', {
      //   confirmButtonText: '确定',
      //   cancelButtonText: '取消',
      //   inputPattern: /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
      //   inputErrorMessage: '邮箱格式不正确'
      // }).then(({ value }) => {
      //   this.$message({
      //     type: 'success',
      //     message: '你的邮箱是: ' + value
      //   });
      // }).catch(() => {
      //   this.$message({
      //     type: 'info',
      //     message: '取消输入'
      //   });
      // });
    },
    // socket
    initSocket() {
      this.ws = new WebSocket(this.socketUrl)
      this.openSocket()
      this.onmessage()
      this.closeSocket()
      this.errorSocket()
    },
    // 打开连接
    openSocket() {
      this.ws.onopen = () => {
        console.log('打开ws连接')
      }
    },
    // 接收消息
    onmessage() {
      this.ws.onmessage = (event) => {
        this.$notify({
          title: '警告',
          message: event.data,
          type: 'warning',
          duration: 3000
        });
      }
    },
    // 关闭连接
    closeSocket() {
      this.ws.onclose = () => {
        let _this = this;
        _this.ws.close();
        console.log('断开重连')
        setTimeout(() => {
          _this.initSocket();
        }, 15000)
      }
    },
    // 连接错误
    errorSocket() {
      this.ws.onerror = () => {
        let _this = this;
        _this.ws.close();
        this.$message.error('websoket连接失败,请刷新!')
      }
    },
    confirmClick() {
      this.$refs.form.validate(async (valid) => {
        if (valid) {
          const res = await updatePwd({
            email: this.username,
            password: this.form.password
          })
          if (res.data.code === 200) {
            this.$message.success('修改成功')
          } else {
            this.$message.error(res.data.msg)
          }
          this.dialogFormVisible = false
          this.$refs.form.resetFields();

        }
      });

    }

  },
  beforeDestroy() {
    this.ws.close();
  }
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
      display: flex;
      justify-content: space-around;
      align-items: center;
      position: relative;
      z-index: 1999;
      background-color: #fff;
      border-right: 1px solid #e6e6e6;
      border-bottom: 1px solid #e6e6e6;

      img {
        height: 90%;
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

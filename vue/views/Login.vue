<template>
  <div class="login-page">
    <div class="login-form-bg">
      <img alt="drawing of a cat" src="../assets/login-bg.png" />
      <div class="login-form">
        <div class="login-form_text">
          <p>PilotGo运维平台</p>
          <span>账户登录</span>
        </div>

        <el-form
          :model="loginForm"
          status-icon
          :rules="rules"
          ref="loginForm"
          class="form"
          v-loading="loading"
        >
          <el-form-item class="form_item" prop="email" label="邮箱">
            <el-input
              type="text"
              v-model="loginForm.email"
              class="form_item__input"
              placeholder="请输入邮箱"
            >
            </el-input>
          </el-form-item>
          <el-form-item class="form_item" prop="password" label="输入密码">
            <el-input
              type="password"
              v-model="loginForm.password"
              class="form_item__input"
              placeholder="请输入密码"
              @keyup.enter.native="submitLogin"
            >
            </el-input>
          </el-form-item>
        </el-form>
        <el-row class="form-btn">
          <el-button
            class="btn"
            type="primary"
            @click="submitLogin"
          >
            登 录
          </el-button>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script>
import { checkEmail } from "@/rules/check";
import { encrypt } from "@/utils/crypto";
export default {
  name: "Login",
  data() {
    return {
      loading: false,
      loginForm: {
        email: "",
        password: "",
      },
      rules: {
        email: [
          { 
            required: true, 
            message: "请输入邮箱", 
            trigger: "blur" 
          },
          {
            validator: checkEmail,
            message: "请输入正确的邮箱格式",
            trigger: "change"
          }
        ],
        password: [
          { 
            required: true, 
            message: "请输入密码", 
            trigger: "blur" 
          }
        ],
      },
    };
  },
  methods: {
    submitLogin() {
      this.$refs.loginForm.validate((valid) => {
        if (valid) {
            this.loading = true;
            let data = {
                username: this.loginForm.email,
                password: encrypt(this.loginForm.password, this.loginForm.email)
            }
            this.$store.dispatch("loginByEmail", data).then((res) => {
                this.loading = false;
                this.$router.push({
                    path: '/home',
                    query: {
                        page: 1,
                        per_page: 20
                    }
                })
            }).catch(error => {
                this.loading = false;
                this.$message({
                    message: error.message,
                    type: 'error'
                })
            })
        }
      });
    },
  },
};
</script>

<style lang="scss">

.login-page {
  height: 100%;
  background: url(../assets/bg.png);
  background-size: cover;
  position: relative;

  .login-form-bg {
    width: 800px;
    height: 440px;
    background-image: linear-gradient(#248ae4, #17b0f8);
    border-radius: 10px;
    position: absolute;
    top: 50%;
    left: 50%;
    margin-top: -220px;
    margin-left: -400px;

    img {
      width: 330px;
      position: absolute;
      top: 0;
      left: 0;
      margin-top: 70px;
      margin-left: 60px;
    }
    .login-form {
      position: absolute;
      top: 0;
      right: 0;
      background: #f4f4f4;
      width: 360px;
      height: inherit;
      border-radius: 10px;

      .login-form_text {
        margin: 20px auto;
        width: 260px;
        p {
          font-size: 30px;
          margin: 10px auto;
        }
        span {
          font-size: 24px;
        }
      }

      .form {
        width: 260px;
        margin: 0 auto;

        .form_item {
          .form_item__input {
            background: transparent;

            .el-input__inner {
              border: none;
              background: transparent;
              border-bottom: 1px solid #14b8fd;
              border-radius: 0;
            }
          }
        }
      }

      .form-btn {
        color: #fff;
        width: 150px;
        margin: 0 auto 5px auto;
        .btn {
          width: 150px;
          border-radius: 40px;
          margin-top: 26px;
        }
      }
      .form-register {
        color: #867d7d;
        font-size: 14px;
        width: 150px;
        float: right;
      }
    }
  }
}


</style>

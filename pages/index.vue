<template>
  <div class="container">
    <div>
      <el-form :model="loginForm" status-icon :rules="rules" label-width="100px" ref="loginForm">
        <el-form-item label="用户名" prop="userName">
          <el-input class="login-input" autocomplete="off" prefix-icon="el-icon-user" v-model="loginForm.userName" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input class="login-input" autocomplete="off" prefix-icon="el-icon-lock" v-model="loginForm.password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="login-input" @click="submitForm('loginForm')">登陆</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { setTimeout } from 'timers'

export default {
  data() {
    const checkUserName = (rule, value, callback) => {
      if (/^[0-9a-zA-Z_]{1,}$/.test(value)) {
        callback()
        return
      }
      callback(new Error('用户名必须是字母、数字或者下划线'))
    }
    return {
      loginForm: {
        userName: '',
        password: ''
      },
      rules: {
        userName: [{ validator: checkUserName, trigger: 'blur' }]
      }
    }
  },
  methods: {
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (!valid) {
          return false
        }
        this.$axios.post('/auth/login', {
          userName: this.loginForm.userName,
          password: this.loginForm.password,
          redirectUrl: this.$route.query.redirect
        }).then(res => {
          this.$notify({
            message: '登陆成功'
          })
          setTimeout(() => window.location.reload(), 1000)
        })
      });
    }
  }
}
</script>

<style>
.container {
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.login-input {
  width: 300px;
}

.title {
  font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
    'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>

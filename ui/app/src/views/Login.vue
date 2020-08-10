<template>
  <b-container fluid>
    <div class="Login">
      <b-navbar type="dark" variant="dark">
        <b-navbar-brand href="https://www.smartvisionworks.com/"></b-navbar-brand>
      </b-navbar>
      <b-container>
        <b-row align-v="center">
          <b-col></b-col>
          <b-col cols="3">
            <b-card type="dark">
              <b-form>
                <b-input v-model="loginInfo.username" label="Username" placeholder="Username" />
                <b-input v-model="loginInfo.password" label="Password" type="password" placeholder="Password" />
                <b-button @click="loginUser" @mousedown="preventDefault">Login</b-button>
              </b-form>
            </b-card>
          </b-col>
          <b-col></b-col>
        </b-row>
      </b-container>
    </div>
  </b-container>
</template>

<script>
  export default {
    name: 'Login',
    data() {
      return {
        loginInfo: {
          username: '',
          password: ''
        }
      }
    },
    methods: {
      //disables a keyboard key's default action
      preventDefault (e) {
        e.preventDefault()
      },
      loginUser () {
        
      }
    }
  }
</script>

<style scoped>
  .card {
    margin-top: 20px;
    background-color: #343a40;
    border: none;
  }
  .form-control {
    margin-bottom: 5px;
  }
</style>
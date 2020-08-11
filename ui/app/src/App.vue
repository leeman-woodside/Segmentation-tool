<template>
  <div id="app">
    <router-view />
    <vue-progress-bar></vue-progress-bar>
  </div>
</template>

<script>

export default {
  data() {
    return {
    }
  },

  beforeCreate() {
    this.$Progress.start()
  },

  mounted() {
    document.getElementById("overlay").style.display = "none";
    document.getElementById("loader").style.display = "none";
  }
}
</script>

<style>
html, body {
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
  background-color: black !important;
  background-image: url(./assets/AINode.jpg);
  background-size: 70%;
  background-position: right bottom;
  background-repeat: no-repeat;
  overflow: auto;
}
body {
  padding-bottom: 8%;
  height: 100%;
  overflow-y: scroll;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#nav {
  padding: 30px;
}
#nav a {
  font-weight: bold;
  color: #2c3e50;
}
#nav a.router-link-exact-active {
  color: #42b983;
}
</style>

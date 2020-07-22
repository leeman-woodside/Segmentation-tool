<template>
  <div id="home">
    <Tool/>
  </div>
</template>

<script>
import Tool from '../components/Tool.vue'

export default {
  name: 'App',
  components: {
    Tool
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
  background-image: url(../assets/AINode.jpg);
  background-size: 70%;
  background-position: right bottom;
  background-repeat: no-repeat;
  overflow: auto;
}
</style>
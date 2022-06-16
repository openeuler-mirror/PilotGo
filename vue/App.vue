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
  Date: 2022-01-14 10:41:07
  LastEditTime: 2022-06-14 09:41:14
 -->
<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: "App",
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
html,
body{
  padding: 0px;
  margin: 0px;
  /* height: 100%; */
  background-color: rgb(11, 35, 117);
  width: 100vw; 
  height: 100vh; 
  transform-origin: 0 0 0; /*指定缩放的基本点*/
  -moz-transform-origin: 0 0 0;
  -ms-transform-origin: 0 0 0;
  -webkit-transform-origin: 0 0 0;
  -o-transform-origin: 0 0 0;
  overflow: hidden;
}
</style>

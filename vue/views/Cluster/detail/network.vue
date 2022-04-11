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
  Date: 2022-04-11 13:07:46
  LastEditTime: 2022-04-11 16:39:14
 -->
<template>
 <div>
   haskd
 </div>
</template>
<script>
import { getNetwork, getNetNic, getNetTcp, getNetUdp } from '@/request/cluster';
export default {
  name: "NetworkInfo",
  data() {
    return {
      netData: [],
    }
  },
  mounted() {
    getNetwork({uuid:this.$route.params.detail}).then(res => {
      this.netData = res.data.data.net_io;
    })
    getNetNic({uuid:this.$route.params.detail}).then(res =>{
      
    })
    getNetTcp({uuid:this.$route.params.detail}).then(res =>{
      
    })
    getNetUdp({uuid:this.$route.params.detail}).then(res =>{
      
    })
  }
}
</script>
<style scoped lang="scss">
</style>
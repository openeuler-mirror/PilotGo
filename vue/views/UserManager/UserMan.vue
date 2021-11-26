<template>
  <div>
    <el-tabs v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="所有用户" name="first">
        <users-man>
        </users-man>
      </el-tab-pane>
      <el-tab-pane label="角色管理" name="third">
        <users-role>
        </users-role>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<script>
import UsersMan from "@/components/UserManager/UsersMan";
import UsersRole from "@/components/UserManager/UsersRole";
export default {
  name: "UserMan",
  components: {UsersRole, UsersMan},
  data() {
    return {
      activeName: 'second',
    }
  },
  mounted() {

  },
  methods: {
    handleClick(tab, event) {
      console.log(tab, event);
    }
  }
}
</script>
<style>

</style>

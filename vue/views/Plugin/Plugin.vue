<template>
  <div class="plugin" style="width:100%;height:100%;">
    <ky-table ref="table" :getData="getPlugins">
      <template v-slot:table_search>
        <div>插件列表</div>
      </template>
      <template v-slot:table_action>
        <el-button @click="handleCreate"> 添加插件 </el-button>
        <el-popconfirm title="确定删除所选项目吗？" @confirm="handleDeleteItems">
          <el-button slot="reference"> 删除 </el-button>
        </el-popconfirm>
      </template>
      <template v-slot:table>
        <el-table-column prop="name" label="名称" width="150">
        </el-table-column>
        <el-table-column prop="version" label="版本" width="150">
        </el-table-column>
        <el-table-column prop="description" label="概述" show-overflow-tooltip>
        </el-table-column>
        <el-table-column prop="enabled" label="状态" width="150">
        </el-table-column>
      </template>
    </ky-table>

    <el-dialog :title="title" :before-close="handleClose" :visible.sync="display" width="560px">
      <add-form v-if="type === 'create'" @click="handleClose"></add-form>
      <iframe v-if="type === 'cockpitIp'" :src="cockpitIp" title="插件"></iframe>
    </el-dialog>
  </div>
</template>

<script>
import { getPlugins, deletePlugins } from "@/request/plugin";
import AddForm from "./form/addForm.vue"
import _import from '../../router/_import';
export default {
  name: "Plugin",
  components: {
    AddForm,
  },
  data() {
    return {
      display: false,
      title: "",
      type: "",
      deleteDialogVisible: false, // 移除插件dialog是否显示
    };
  },

  methods: {
    getPlugins,
    handleClose(type) {
      this.display = false;
      this.title = "";
      this.type = "";
      if (type === 'success') {
        this.refresh();
      }
    },
    refresh() {
      this.$refs.table.handleSearch();
    },
    handleCreate() {
      this.display = true;
      this.title = "添加插件";
      this.type = "create";
    },
    handleDeleteItems() {
      this.$refs.table.selectRow.rows.forEach(item => {
        deletePlugins({ UUID: item.uuid }).then((res) => {
          if (res.data.code === 200) {
            this.refresh();
            this.$store.dispatch('SetDynamicRouters', []).then(() => {
              this.$store.dispatch('GenerateRoutes');
            })
            this.$message.success("删除成功");
          } else {
            this.$message.error("删除插件错误：" + res.data.msg)
          }
        });
      });
    },
  },
};
</script>

<style scoped lang="scss">
.plugin {
  width: 100%;
}
</style>

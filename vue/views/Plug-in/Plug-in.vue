<template>
  <div class="plugIn">
    <ky-table
        class="cluster-table"
        ref="table"
        :getData="getPlugins"
      >
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
          <el-table-column prop="plugin" label="名称" width="150">
          </el-table-column>
          <el-table-column prop="version" label="版本" width="150">
          </el-table-column>
          <el-table-column prop="description" label="概述" show-overflow-tooltip>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="150">
          </el-table-column>
        </template>
      </ky-table>

      <el-dialog
        :title="title"
        :before-close="handleClose" 
        :visible.sync="display" 
        width="560px"
      >
        <add-form v-if="type === 'create'" @click="handleClose"></add-form>
        <iframe v-if="type === 'cockpitIp'" :src="cockpitIp" title="插件"></iframe>
      </el-dialog>
  </div>
</template>

<script>
import { getPlugins, deletePlugins } from "@/request/plugIn";
import kyTable from "@/components/KyTable";
import AddForm from "./form/addForm.vue"
export default {
  name: "PlugIn",
  components: {
    kyTable,
    AddForm,
  },
  data() {
    return {
      cockpitIp: "",
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
      if(type === 'success') {
        this.refresh();
      }
    }, 
    refresh(){
      this.$refs.table.handleSearch();
    },
    handleCreate() {
      this.display = true;
      this.title = "添加插件";
      this.type = "create";
    },
    handleDeleteItems() {
      let names = [];
      this.$refs.table.selectRow.rows.forEach(item => {
        names.push(item.plugin);
      });
      deletePlugins({ plugin: names }).then((res) => {
        if (res.data.status === "success") {
          this.refresh();
          this.$message.success("删除成功");
        }
      });
    },
  },
};
</script>

<style scoped lang="scss">
.plugIn {
  width: 100%;
  margin-top: 10px;
}
</style>

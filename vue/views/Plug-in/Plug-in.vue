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
             <!--
        <template v-slot:table_action>
          <el-button @click="handleCreate"> 添加插件 </el-button>
          <el-popconfirm title="确定删除所选项目吗？" @confirm="handleDeleteItems">
            <el-button slot="reference"> 删除 </el-button>
          </el-popconfirm>
        </template>
        -->
        <template v-slot:table>
          <el-table-column prop="name" label="名称" width="200">
          </el-table-column>
        <!--  <el-table-column prop="version" label="版本" width="200">
          </el-table-column>
          <el-table-column prop="description" label="概述" show-overflow-tooltip>
          </el-table-column> -->
          <el-table-column prop="status" label="状态" width="150">
           <template slot-scope="scope">
           <p v-if="scope.row.status=='0'">
           停止状态
           </p>
           <p v-else-if="scope.row.status=='1'">
           运行状态
           </p>
           </template>
          </el-table-column>
           <el-table-column prop="status" label="操作" width="300">
                       <template slot-scope="scope">
                <el-button  v-if="scope.row.status=='0'" size="mini" type="primary" plain name="plugin_onload" 
                  @click="handleLoad(scope.row.name)"> 安装
                </el-button>
                                <el-button  v-if="scope.row.status=='1'" size="mini" type="primary" plain name="plugin_unload" 
                  @click="handleUnLoad(scope.row.name)"> 卸载
                </el-button>
                                                <el-button  v-if="scope.row.status=='1'" size="mini" type="primary" plain name="" 
                  @click="handleOpen(scope.row.url)">进入
                </el-button>
            </template>
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
import { getPlugins, loadPlugin, unLoadPlugin } from "@/request/plugIn";
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
    handleLoad(name) {
      loadPlugin({ Name: name }).then((res) => {
        console.log(res);
        if (res.data.code === 200) {
          this.refresh();
          this.$message.success("安装成功");
        }
      });
    },
   handleUnLoad(name) {
      unLoadPlugin({ Name: name }).then((res) => {
        console.log(res);
        if (res.data.code === 200) {
          this.refresh();
          this.$message.success("卸载成功");
        }
      });
    },
    handleOpen(url) {
      this.$router.push(
        {
          path:'/plugin-web',
          query:{
            header_title:"插件展示",
            urlPath:url,
          }
        }
      )
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

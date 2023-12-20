<template>
  <div class="plugin" style="width:100%;height:100%;">
    <ky-table ref="table" :getData="getPluginsPaged">
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
        <el-table-column prop="url" label="服务端地址" width="250">
        </el-table-column>
        <el-table-column prop="status" label="连接状态" width="150">
          <template slot-scope="scope">
            <el-icon :class="{ 'el-icon-success': scope.row.status, 'el-icon-error': !scope.row.status }"></el-icon>
            <span>{{ scope.row.status === true ? '连接' : '断开' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="lastheatbeat" label="上次心跳时间" width="250">
        </el-table-column>
        <el-table-column prop="enabled" label="启用状态" width="150">
          <template slot-scope="scope">
            <span>{{ scope.row.enabled === 0 ? '已禁用' : '已启用' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right">
          <template slot-scope="scope">
            <el-button size="mini" type="primary" plain name="default_all" @click="handlePluginState(scope.row)">
              {{ scope.row.enabled === 1 ? '禁用' : '启用' }}
            </el-button>
          </template>
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
import { getPluginsPaged, deletePlugins, unLoadPlugin } from "@/request/plugin";
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
    getPluginsPaged,
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
            this.del_plugin_tag(item.name)
          } else {
            this.$message.error("删除插件错误：" + res.data.msg)
          }
        });
      });
    },
    // 启用/停用插件
    handlePluginState(item) {
      let targetEnabled = item.enabled === 1 ? 0 : 1
      unLoadPlugin({ uuid: item.uuid, enable: targetEnabled }).then(res => {
        if (res.data.code === 200) {
          this.refresh();
          this.$store.dispatch('SetDynamicRouters', []).then(() => {
            this.$store.dispatch('GenerateRoutes').catch((err) => {
              console.log("generate router error: ", err)
            })
          }).catch((err) => {
            console.log("update dynamic router error: ", err)
          })
          this.$message.success(res.data.msg);
          if (targetEnabled === 0) {
            // 禁用插件，删除tag标签
            this.del_plugin_tag(item.name);
          }
        } else {
          this.$message.error(res.data.msg);
        }
      })
    },

    // 删除对应tag标签显示
    del_plugin_tag(tagName) {
      let openTags = this.$store.getters.visitedViews;
      let tagArr = openTags.filter(tag => tag.title === tagName);
      if (tagArr.length > 0) {
        this.$store.dispatch('tagsView/delView', tagArr[0]).then()
      }
    }
  },
};
</script>

<style scoped lang="scss">
.plugin {
  width: 100%;
}

.el-icon-success {
  color: green;
}

.el-icon-error {
  color: red;
}
</style>

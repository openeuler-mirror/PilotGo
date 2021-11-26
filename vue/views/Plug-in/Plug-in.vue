<template>
  <div class="plugIn">
    <div class="plugIn-title">
      <span class="plugIn-title__text">插件列表</span>
      <div class="plugIn-title__operate">
        <el-button size="mini" @click="addDialogVisible = true"
          >添加插件</el-button
        >
        <el-popconfirm
          title="确定删除所选项目吗？"
          @confirm="handleDeleteItems"
        >
          <el-button size="mini" slot="reference">删除</el-button>
        </el-popconfirm>
      </div>
    </div>
    <el-table
      class="plugIn-table"
      :header-cell-style="{
        'background-color': '#f6f8fd',
        'border-bottom': '2px solid #e4eaf9',
        padding: '8px 4px',
      }"
      :cell-style="rowStyle"
      :data="pluginList"
      style="width: 100%"
      @selection-change="handleSelectChange"
    >
      <el-table-column type="selection" width="55" v-model="multipleSelection">
      </el-table-column>
      <el-table-column prop="plugin" label="名称" width="150">
      </el-table-column>
      <el-table-column prop="version" label="版本" width="150">
      </el-table-column>
      <el-table-column prop="description" label="概述" show-overflow-tooltip>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="150">
      </el-table-column>
    </el-table>

    <el-row class="plugIn-footer">
      <el-pagination
        class="pagination"
        background
        :page-sizes="sizes"
        layout="total, sizes, prev, pager, next, jumper"
        prev-text="上一页"
        next-text="下一页"
        :page-size="pageSize"
        :total="pluginList.length"
        @size-change="handleChangePageSize"
        @current-change="handleCurrentChangePlugInTable"
      >
      </el-pagination>
    </el-row>

    <el-dialog title="添加插件" :visible.sync="addDialogVisible" width="460px">
      <el-form
        :model="addPluginForm"
        :rules="rules"
        ref="addPluginForm"
        label-width="100px"
      >
        <el-form-item label="主机地址:" prop="host">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addPluginForm.host"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="端口:" prop="port">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addPluginForm.port"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="协议类型:">
          <el-select v-model="addPluginForm.protocol" placeholder="请选择">
            <el-option
              v-for="item in protocols"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCancelForm('addPluginForm')"
            >取 消</el-button
          >
          <el-button type="primary" @click="handleSubmitForm()"
            >确 定</el-button
          >
        </span>
      </template>
    </el-dialog>

    <el-dialog
      title="cockpit配置"
      :visible.sync="cockpitDialogVisible"
      width="860px"
    >
      <iframe :src="cockpitIp" frameborder="0"></iframe>
    </el-dialog>
  </div>
</template>

<script>
import { getPlugins, insertPlugin, deletePlugins } from "@/request/api";
export default {
  name: "PlugIn",
  data() {
    // ip验证规则
    const validateIP = (rule, value, callback) => {
      let _this = this;
      if (!value && value == 0) {
        callback(new Error("IP不能为空！"));
      } else if (!_this.ipReg.test(value)) {
        callback(new Error("请输入正确的IP地址"));
      } else {
        callback();
      }
    };
    const validatePort = (rule, value, callback) => {
      if (!value) {
        callback(new Error("端口不能为空！"));
      } else {
        callback();
      }
    };
    const validatePlugin = (rule, value, callback) => {
      if (!value) {
        callback(new Error("插件名称不能为空！"));
      } else {
        callback();
      }
    };
    const validateProtocol = (rule, value, callback) => {
      if (!value) {
        callback(new Error("协议不能为空！"));
      } else {
        callback();
      }
    };
    return {
      cockpitIp: "",
      ipReg:
        /^([1-9]|[1-9]\d|1\d{2}|2[0-1]\d|22[0-3])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$/,
      addDialogVisible: false, // 添加插件dialog是否显示
      deleteDialogVisible: false, // 移除插件dialog是否显示
      cockpitDialogVisible: false,
      addPluginForm: {
        // 添加插件需要填的信息

        host: "",
        port: "",
        protocol: "",
      },
      rules: {
        // 添加插件需要填的信息的验证规则
        host: [{ required: true, validator: validateIP, trigger: "blur" }],
        port: [{ required: true, validator: validatePort }],
        plugin: [{ required: true, validator: validatePlugin }],
        protocol: [{ required: true, validator: validateProtocol }],
      },
      radio: { id: "", name: "" }, // table选中的row的id和name
      sizes: [14, 20, 30, 40, 50, 100], // 选择table每页显示多少条
      pageSize: 14, // 默认table每页多少条数据
      plugInTableData: [],
      tableData: [],
      pluginList: [],
      protocols: [
        { label: "http", value: "http" },
        { label: "https", value: "https" },
      ],
      multipleSelection: [],
    };
  },

  methods: {
    // 改变table行的背景颜色条纹式
    rowStyle({ rowIndex }) {
      if (rowIndex % 2 == 0) {
        return "background: #fff; padding: 8px 0;";
      } else {
        return "background: #f2f7ff; padding: 8px 0";
      }
    },
    // table每页显示条数
    currentChangePlugInTable(currentPage, pageSize) {
      let _this = this;

      // 初次改变时，排除当前页数为空的情况
      if (!currentPage) {
        currentPage = 1;
      }

      _this.tableData = _this.plugInTableData.slice(
        (currentPage - 1) * pageSize,
        currentPage * pageSize
      );
    },
    // 每页显示条目个数改变
    handleChangePageSize(size) {
      let _this = this;
      _this.pageSize = size;

      this.currentChangePlugInTable(_this.currentPage, _this.pageSize);
    },
    // 当前页改变
    handleCurrentChangePlugInTable(current) {
      let _this = this;
      _this.currentPage = current;
      _this.currentChangePlugInTable(_this.currentPage, _this.pageSize);
    },
    handleSubmitForm() {
      let _this = this;
      _this.$refs["addPluginForm"].validate((valid) => {
        if (valid) {
          insertPlugin(_this.addPluginForm)
            .then((res) => {
              console.log(res.data.status);
              if (res.data.status === "success") {
                _this.refreshList();
                _this.addDialogVisible = false;
                _this.$message.success("添加成功");
              } else {
                _this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              _this.$message.error("添加失败，请检查输入内容");
            });
        } else {
          _this.$message.error("添加失败，请检查输入内容");
          return false;
        }
      });
    },
    handleCancelForm(addPluginForm) {
      this.$refs[addPluginForm].clearValidate();
      this.addDialogVisible = false;
    },
    handleDeleteItems() {
      let _this = this;
      let names = [];
      for (let i of _this.multipleSelection) {
        names.push(i.plugin);
      }
      deletePlugins({ plugin: names }).then((res) => {
        if (res.data.status === "success") {
          _this.refreshList();
          _this.$message.success("删除成功");
        }
      });
    },
    refreshList() {
      getPlugins().then((res) => {
        let _this = this;
        _this.pluginList = res.data.data;
        if (res.data.data.length > 0) {
          for (let i of _this.pluginList) {
            if (i.plugin == "cockpit") {
              _this.$store.commit("mutateCockpitPluginServer", i.host);
            }
          }
        }
      });
    },
    handleSelectChange(val) {
      this.multipleSelection = val;
    },
  },

  mounted() {
    let _this = this;
    // 设置table每页显示条数
    this.refreshList();
    _this.tableData = _this.plugInTableData.slice(0, _this.pageSize);

    let count = 0; //唯一标识计数

    // 给数据添加唯一标识属性
    _this.plugInTableData.forEach((tableItem) => {
      tableItem.id = count;
      count++;
    });
  },
};
</script>

<style scoped lang="scss">
.plugIn {
  width: 100%;
  margin-top: 10px;
  .plugIn-title {
    height: 44px;
    background: #3e9df9;
    line-height: 44px;
    padding: 0 10px;
    .plugIn-title__text {
      font-size: 14px;
      color: #fff;
    }

    .plugIn-title__operate {
      float: right;

      .el-button {
        color: #fff;
        background: transparent;
        border: none;
      }
      .el-button:hover,
      .el-button:active {
        background: #fff;
        color: #3e9df9;
      }
    }
  }

  .plugIn-footer {
    float: right;
    margin-top: 20px;

    .total {
      float: left;
      font-size: 14px;
      color: #606266;
      line-height: 32px;
    }
    .pagination {
      float: right;
    }
  }

  .ipInput {
    width: 240px;
  }
  .deletePlugInText {
    margin-left: 20px;
  }
}
</style>

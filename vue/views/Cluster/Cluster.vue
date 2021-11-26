<template>
  <div class="cluster">
    <div class="cluster-title">
      <span class="cluster-title__text">主机列表</span>
      <div class="cluster-title__operate">
        <el-button size="mini" @click="addIPDialogVisible = true"
          >添加IP</el-button
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
      class="cluster-table"
      :header-cell-style="{
        'background-color': '#f6f8fd',
        'border-bottom': '2px solid #e4eaf9',
        padding: '8px 4px',
      }"
      :cell-style="rowStyle"
      ref="multipleTable"
      :data="clusterList"
      style="width: 100%"
      @selection-change="handleClickRow"
    >
      <el-table-column
        prop="id"
        v-model="multipleSelection"
        type="selection"
        width="55"
      >
      </el-table-column>

      <el-table-column label="IP" width="150">
        <template slot-scope="scope">
          <a @click="pluginLogin(scope.row.ip)">{{ scope.row.ip }}</a>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="150">
        <template slot-scope="scope">
          <span v-if="scope.row.system_status == 0">异常</span>
          <span v-else>正常</span>
        </template>
      </el-table-column>
      <el-table-column label="类型" width="150">
        <template slot-scope="scope">
          <span v-if="scope.row.machine_type == 0">虚拟机</span>
          <span v-else>物理机</span>
        </template>
      </el-table-column>
      <el-table-column prop="system_info" label="系统信息"> </el-table-column>
      <el-table-column prop="system_version" label="系统版本" width="150">
      </el-table-column>
      <el-table-column prop="arch" label="架构" width="150"> </el-table-column>
      <el-table-column prop="installation_time" label="安装时间">
      </el-table-column>
    </el-table>

    <el-row class="cluster-footer">
      <el-pagination
        class="pagination"
        background
        :page-sizes="sizes"
        layout="total, sizes, prev, pager, next, jumper"
        prev-text="上一页"
        next-text="下一页"
        :page-size="pageSize"
        :total="clusterTableData.length"
        @size-change="handleChangePageSize"
        @current-change="handleCurrentChangePlugInTable"
      >
      </el-pagination>
    </el-row>

    <el-dialog title="添加IP" :visible.sync="addIPDialogVisible" width="560px">
      <el-form
        :model="addIPForm"
        :rules="rules"
        ref="addIPForm"
        label-width="100px"
      >
        <el-form-item label="IP:" prop="ip">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="addIPForm.ip"
            autocomplete="off"
          ></el-input>
        </el-form-item>

        <el-form-item label="系统信息:" prop="system_info">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.system_info"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="系统版本:" prop="system_version">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.system_version"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="架构:" prop="arch">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.arch"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="安装时间:" prop="installation_time">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.installation_time"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="机器类型:" prop="machine_type">
          <el-select v-model="addIPForm.machine_type" placeholder="请选择">
            <el-option
              v-for="item in machinetypes"
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
          <el-button @click="addIPDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="handleSubmitForm('addIPForm')"
            >确 定</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import axios from "axios";
import { getClusters, insertIp, deleteIp } from "@/request/api";
export default {
  name: "Cluster",
  data() {
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
    return {
      addIPDialogVisible: false,
      delHostDialogVisible: false,
      ipReg:
        /^([1-9]|[1-9]\d|1\d{2}|2[0-1]\d|22[0-3])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$/,
      sizes: [14, 20, 30, 40, 50, 100], // 选择table每页显示多少条
      pageSize: 14, // 默认table每页多少条数据
      clusterTableData: [
        {
          ip: "",
          status: "",
          type: "",
          systemInfo: "",
          systemVersion: "",
          Architecture: "",
          installationTime: "",
        },
      ],
      tableData: [],
      multipleSelection: [],
      addIPForm: {
        ip: "",
        system_info: "",
        system_version: "",
        arch: "",
        installation_time: "",
        machine_type: 0,
      },
      rules: {
        // 添加插件需要填的信息的验证规则
        ip: [{ required: true, validator: validateIP, trigger: "blur" }],
        system_info: [{ required: true, trigger: "blur" }],
        system_version: [{ required: true, trigger: "blur" }],
        arch: [{ required: true, trigger: "blur" }],
        installation_time: [{ required: true, trigger: "blur" }],
        machine_type: [{ required: true, trigger: "blur" }],
      },
      clusterList: [
        {
          ip: "",
          status: "",
          type: "",
          systemInfo: "",
          systemVersion: "",
          Architecture: "",
          installationTime: "",
        },
      ],
      deleteIPList: [],
      machinetypes: [
        { label: "虚拟机", value: 0 },
        { label: "物理机", value: 1 },
      ],
      cockpitUrl: "",
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

      _this.tableData = _this.clusterTableData.slice(
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
    // 点击选中/取消选框
    handleClickRow(val) {
      this.multipleSelection = val;
    },

    // 当前页改变
    handleCurrentChangePlugInTable(current) {
      let _this = this;
      _this.currentPage = current;
      _this.currentChangePlugInTable(_this.currentPage, _this.pageSize);
    },
    // 删除主机界面
    handleDeleteItems() {
      let _this = this;
      let ids = [];
      for (let i of _this.multipleSelection) {
        ids.push(i.id);
      }
      deleteIp({ id: ids }).then((res) => {
        if (res.data.status === "success") {
          _this.refreshList();
          this.$message.success("删除成功");
        }
      });
    },
    pluginLogin(name) {
      axios
        .get(
          "https://" +
            this.$store.state.cockpitPluginServer +
            ":8888/plugin/cockpit/port",
          {
            params: {
              ip: name,
              port: "9090",
            },
          }
        )
        .then((res) => {
          let cockpitUrl = "https://" + name + ":" + res.data.port;
          this.$emit("selectIp", cockpitUrl);
        });
    },
    handleSubmitForm() {
      let _this = this;
      this.$refs["addIPForm"].validate((valid) => {
        if (valid) {
          insertIp(_this.addIPForm)
            .then((res) => {
              if (res.data.status === "success") {
                _this.refreshList();
                _this.addIPDialogVisible = false;
              } else {
                _this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("添加失败，请检查输入内容");
            });
        } else {
          this.$message.error("添加失败，请检查输入内容");
          return false;
        }
      });
    },
    refreshList() {
      getClusters().then((res) => {
        let _this = this;
        _this.clusterList = res.data.data;
      });
    },
  },

  mounted() {
    let _this = this;
    if (_this.clusterTableData.length > 14) {
      _this.tableData = _this.clusterTableData.slice(0, _this.pageSize);
    } else {
      _this.tableData = _this.clusterTableData.slice(
        0,
        _this.clusterTableData.length
      );
    }

    let count = 0;

    _this.clusterTableData.forEach((tableItem) => {
      tableItem.id = count;
      count++;
    });

    _this.refreshList();
  },
};
</script>

<style scoped lang="scss">
.cluster {
  margin-top: 10px;
  .cluster-title {
    height: 44px;
    background: #3e9df9;
    line-height: 44px;
    padding: 0 10px;
    .cluster-title__text {
      font-size: 14px;
      color: #fff;
    }

    .cluster-title__operate {
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

  .cluster-footer {
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

  .deleteHostText {
    margin-left: 10px;
    .del-host {
      color: red;
      font-weight: 600;
      font-size: 18px;
    }
  }
}
</style>

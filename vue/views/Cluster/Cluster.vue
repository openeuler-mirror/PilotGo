<template>
  <div class="cluster">
    <div class="dept">
      <ky-tree :getData="getChildNode" ref="tree" @nodeClick="handleSelectDept"></ky-tree>
    </div>
    <div class="info">
      <ky-table
        class="cluster-table"
        ref="table"
        :getData="getClusters"
        :searchData="searchData"
      >
        <template v-slot:table_search>
          <div>{{ departName }}</div>
        </template>
        <template v-slot:table_action>
          <el-button  @click="handleAddIp"> 注册 </el-button>
          <el-button  @click="handleAddBatch"> 创建批次 </el-button>
          <el-popconfirm title="确定删除所选项目吗？" @confirm="handleDelete">
            <el-button  slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </el-button>
          </el-popconfirm>
        </template>
        <template v-slot:table>
          <el-table-column label="IP" width="90">
            <template slot-scope="scope">
              <a @click.stop="handleSelectIp(scope.row.ip)">
                {{ scope.row.ip }}
              </a>
            </template>
          </el-table-column>
          <el-table-column prop="system_cpu" label="cpu" width="90"> 
          </el-table-column>
          <el-table-column label="状态" width="150">
            <template slot-scope="scope">
              <span v-if="scope.row.system_status == 0">异常</span>
              <span v-else>正常</span>
            </template>
          </el-table-column>
           <el-table-column prop="system_info" label="系统信息" width="150"> 
          </el-table-column>
          <el-table-column label="防火墙配置" width="120">
            <template slot-scope="scope">
              <el-button
                size="mini"
                @click="handleFireWall(scope.row.ip)">
                <em class="el-icon-edit-outline"></em>
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template slot-scope="scope">
              <el-button size="mini" type="primary" plain 
                @click="handleUpdateIp(scope.row.ip)"> 
                编辑 </el-button>
            </template>
          </el-table-column>
        </template>
      </ky-table>
    </div>

    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="display" 
      width="560px"
    >
     <add-form v-if="type === 'create'" @click="handleClose"></add-form>
     <update-form v-if="type === 'update'" :ip='ip' @click="handleClose"></update-form>   
     <batch-form v-if="type === 'batch'" :departInfo='departInfo' :machineIds='machineIds' @click="handleClose"></batch-form>   
     <device-detail v-if="type === 'disk'" :ip='ip'></device-detail>
    </el-dialog>
  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import kyTree from "@/components/KyTree";
import AddForm from "./form/addForm";
import UpdateForm from "./form/updateForm";
import BatchForm from "./form/batchForm";
import DeviceDetail from "./detail/index";
import { getClusters, deleteIp, getChildNode } from "@/request/cluster";
export default {
  name: "Cluster",
  components: {
    AddForm,
    UpdateForm,
    BatchForm,
    DeviceDetail,
    kyTable,
    kyTree,
  },
  data() {  
    return {
      title: '',
      type: '',
      ip: '',
      departName: '',
      departInfo: {},
      machineIds: [],
      display: false,
      disabled: false,
      searchData: {
        DepartId: 1
      },
    };
  },
  mounted() {
    this.departName = '机器列表';
  },
  methods: {
    getClusters,
    getChildNode,
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
    },
    handleAddIp() {
      this.display = true;
      this.title = "注册IP";
      this.type = "create"; 
    },
    handleUpdateIp(ip) {
      this.display = true;
      this.title = "编辑IP";
      this.type = "update"; 
      this.ip = ip;
    },
    handleAddBatch() {
      this.machineIds = [];
      this.display = true;
      this.title = "创建批次";
      this.type = "batch"; 
      let selects = this.$refs.table.selectRow.rows;
      selects.forEach(item => {
        this.machineIds.push(item.machineuuid);
      })
    },
    handleDelete() {
      let ids = this.$refs.table.selectRow.rows[0];
      deleteIp({ uuid: ids }).then((res) => {
        if (res.data.code === 200) {
          this.$refs.table.handleSearch();
          this.$message.success("删除成功");
        }
      });
    },
    handleSelectDept(data) {
      if(data) {
        this.departName = data.label + '机器列表';
        this.searchData.DepartId = data.id;
        this.departInfo = data;
        this.$refs.table.handleSearch();
      }
    },
    handleSelectIp(ip) {
      this.display = true;
      this.title = "磁盘使用";
      this.type = "disk"; 
      this.ip = ip;
      // this.$store.commit("SET_SELECTIP", row.ip);
    },
    handleFireWall(ip) {
      this.$router.push({
        name: 'Firewall',
        query: { ip: ip }
      })
    }
  },
};
</script>

<style scoped lang="scss">
.cluster {
  margin-top: 10px;
  .dept {
    width: 36%;
    display: inline-block;
  }
  .info {
    width: 60%;
    float: right;
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

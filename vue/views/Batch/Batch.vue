<template>
  <div style="height:100%">
  <router-view v-if="$route.name == 'BatchDetail'"></router-view>
  <div v-if="$route.name == 'Batch'" class="panel">
    <ky-table
        class="cluster-table"
        ref="table"
        :getData="getBatches"
      >
        <template v-slot:table_search>
          <div>批次列表</div>
        </template>
        <template v-slot:table_action>
          <el-popconfirm 
          title="确定删除此批次？"
          cancel-button-type="default"
          confirm-button-type="danger"
          @confirm="handleDelete">
          <auth-button name="batch_delete" slot="reference" :disabled="$refs.table && $refs.table.selectRow.rows.length == 0"> 删除 </auth-button>
        </el-popconfirm>
        </template>
        <template v-slot:table>
          <el-table-column label="批次名称">
            <template slot-scope="scope">
              <router-link :to="$route.path + scope.row.ID">
                {{ scope.row.name }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="manager" label="创建者"> 
          </el-table-column>
          <el-table-column prop="DepartName" label="部门"> 
          </el-table-column>
          <el-table-column prop="CreatedAt" label="创建时间" sortable>
            <template slot-scope="scope">
              <span>{{scope.row.CreatedAt | dateFormat}}</span>
            </template>
          </el-table-column>
           <el-table-column prop="description" label="备注"> 
          </el-table-column>
          <el-table-column prop="operation" label="操作" width="200">
            <template slot-scope="scope">
              <auth-button name="batch_edit" size="mini" type="primary" plain @click="handleEdit(scope.row)">
                编辑
              </auth-button>
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
        <update-form :row="rowData" v-if="type === 'update'" @click="handleClose"></update-form>
      </el-dialog>
      </div>
  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import UpdateForm from "./form/updateForm.vue"
import { getBatches, delBatches } from "@/request/batch";
export default {
  name: "Batch",
  components: {
    kyTable,
    AuthButton,
    UpdateForm,
  },
  data() {
    return {
      display: false,
      title: "",
      type: "",
      rowData: {},
    }
  },
  methods: {
    getBatches,
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
    handleEdit(row) {
      this.display = true;
      this.title = "编辑批次";
      this.type = "update";
      this.rowData = row;
    },
    handleDelete() {
      let delDatas = [];
      delDatas = this.$refs.table.selectRow.ids.map(item => {
        return item;
      });
      delBatches({BatchID: delDatas}).then(res => {
        if(res.data.code === 200) {
          this.$message.success(res.data.msg);
        } else {
          this.$message.error(res.data.msg);
        }
        this.refresh();
      })
    },
    
  },
  filters: {
    dateFormat: function(value) {
      let date = new Date(value);
      let y = date.getFullYear();
      let MM = date.getMonth() + 1;
      MM = MM < 10 ? "0" + MM : MM;
      let d = date.getDate();
      d = d < 10 ? "0" + d : d;
      let h = date.getHours();
      h = h < 10 ? "0" + h : h;
      let m = date.getMinutes();
      m = m < 10 ? "0" + m : m;
      let s = date.getSeconds();
      s = s < 10 ? "0" + s : s;
      return y + "-" + MM + "-" + d + " " + h + ":" + m;
    }
  },
}
</script>

<style scoped>
.panel {
  height: 100%;
}
</style>

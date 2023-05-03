<template>
  <div class="history">
    <ky-table
      ref="table"
      :getData="lastFileList"
      :searchData="searchData"
      v-if="!showCompare"
    >
      <template v-slot:table_search>
        <span>历史版本</span>
      </template>
      <template v-slot:table_action>
        <auth-button name="user_del" :disabled="$refs.table && $refs.table.selectRow.rows.length != 2" @click="handleCompare"> 对比 </auth-button>
      </template>
      <template v-slot:table>
        <el-table-column  prop="name" label="名称">
        </el-table-column>
        <el-table-column  prop="user" label="修改人">
        </el-table-column>
        <el-table-column  prop="userDept" label="修改人部门">
        </el-table-column>
        <el-table-column prop="UpdatedAt" label="更新时间">
          <template slot-scope="scope">
            <span>{{scope.row.UpdatedAt | dateFormat}}</span>
          </template>
        </el-table-column>
        <el-table-column  prop="description" label="描述">
        </el-table-column>
        <el-table-column label="操作" fixed="right">
          <template slot-scope="scope">
            <el-popconfirm 
              title="确定要回滚到此版本?"
              cancel-button-type="default"
              confirm-button-type="danger"
              @confirm="handleRollBack(scope.row)">
            <auth-button name="user_del" type="primary" plain size="mini" slot="reference"> 回滚 </auth-button>
          </el-popconfirm>
          </template>
        </el-table-column>
      </template>
    </ky-table>
    <div class="top" v-if="showCompare">
      <span style="color: rgb(242, 150, 38)">提示：请选择左右对比文件进行对比</span>
      <el-button plain type="primary"  @click="handleExit" icon="el-icon-back">返回</el-button>
    </div>
    <compare-form :id="searchData.id" :leftFile="leftFile" :rightFile="rightFile" v-if="showCompare" @click="handleClose"></compare-form>

  </div>
</template>

<script>
import CompareForm from "../form/compareForm.vue";
import kyTable from "@/components/KyTable";
import AuthButton from "@/components/AuthButton";
import { lastFileList,fileRollback } from "@/request/config"
export default {
  components: {
    kyTable,
    CompareForm,
    AuthButton,
  },
  props: {
    row: {
      type: Object,
      default: function() {
        return {}
      }
    }
  },
  data() {
    return {
      showCompare: false,
      leftFile: '',
      rightFile: '',
      searchData: {
        id: this.row.id,
      },
    }
  },
  methods: {
    lastFileList,
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
    }, 
    handleCompare() {
      this.leftFile = this.$refs.table.selectRow.rows[0];
      this.rightFile = this.$refs.table.selectRow.rows[1];
      this.showCompare = true;
    },
    handleExit() {
      this.leftFile = '';
      this.rightFile = '';
      this.showCompare = false;
    },
    handleRollBack(row) {
      let params = {
        id: row && row.id,
        filePId: row && row.filePId,
        user: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      fileRollback(params).then(res => {
        if(res.data.code === 200) {
          this.$emit("click");
          this.$message.success(res.data.msg);
        } else {
          this.$message.error(res.data.msg)
        }
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

<style scoped lang="scss">
.history {
  width: 100%;
  height: 70vh;
  .top {
    width: 100%;
    height: 4%;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>

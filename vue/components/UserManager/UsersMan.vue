<template>
  <div>
    <div style="display: flex;justify-content: space-between" class="search-form">
      <div>
        <el-input placeholder="请输入用户名或邮箱名进行搜索..." prefix-icon="el-icon-search"
                  clearable
                  @clear="initEmps"
                  style="width: 350px;margin-right: 10px" v-model="keyword"
                  @keydown.enter.native="initEmps" :disabled="showAdvanceSearchView"></el-input>
        <el-button icon="el-icon-search" type="primary" @click="initEmps" :disabled="showAdvanceSearchView">
          搜索
        </el-button>
      </div>
      <div>
        <el-button type="primary" icon="el-icon-plus" @click="showAddEmpView">
          添加用户
        </el-button>
        <el-button type="success" @click="exportData" icon="el-icon-download">
          导出用户数据
        </el-button>
        <el-upload
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="onSuccess"
          :on-error="onError"
          :disabled="importDataDisabled"
          style="display: inline-flex;margin-right: 8px"
          action="/employee/basic/import">
          <el-button :disabled="importDataDisabled" type="success" :icon="importDataBtnIcon">
            {{importDataBtnText}}批量添加用户
          </el-button>
        </el-upload>
      </div>
    </div>
    <el-table
      :data="tableData"
      style="width: 100%"
      :row-class-name="tableRowClassName">
      <el-table-column
        prop="id"
        label="编号"
        width="80">
      </el-table-column>
      <el-table-column
        prop="username"
        label="用户名"
        width="180">
      </el-table-column>
      <el-table-column
        prop="name"
        label="手机号"
        width="180">
      </el-table-column>
      <el-table-column
        prop="email"
        label="邮箱">
      </el-table-column>
      <el-table-column
        prop="enable"
        label="是否启用"
        width="80">
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: "UsersMan"
}
</script>

<style scoped>
.search-form{
  margin-bottom: 12px;
}
.el-table .warning-row {
  background: oldlace;
}

.el-table .success-row {
  background: #f0f9eb;
}
</style>

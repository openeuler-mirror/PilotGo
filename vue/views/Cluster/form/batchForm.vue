<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="批次名称:" prop="batchName">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.batchName"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="已选机器:" prop="mechines">
          <el-table
            :header-cell-style="hStyle"
            :cell-style="bstyle"
            :data="machineArr"
            height="180"
            border
            style="width: 100%; padding:0px">
            <el-table-column
              prop="ip"
              label="ip">
            </el-table-column>
            <el-table-column
              prop="uuid"
              label="id">
            </el-table-column>
            <el-table-column
              prop="departid"
              label="部门">
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.description"
            autocomplete="off"
          ></el-input>
        </el-form-item>
      </el-form>

      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleContinue">继续选择</el-button>
        <el-button type="primary" @click="handleConfirm">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  createBatch  } from "@/request/batch";
export default {
  props: {
    departInfo: {
      type: Object,
      default: function() {
        return {
          id: 1
        }
      }
    },
    machines: {
      type: Array,
    }
  },
  data() {
    return {
      machineArr: [],
      hStyle: {
        background:'#F3F4F7',
        color:'#555',
        textAlign:'center',
        padding:'0'
      },
      bstyle: {
         textAlign:'center',
      },
      form: {
        batchName: "",
        description: ""
      },
      rules: {
        batchName: [{ 
          required: true, 
          message: "请填写批次名称", 
          trigger: "blur" 
        }]
      },
    }
  },
  mounted() {
    let keys = {};
    this.machines.forEach((item) => keys[item.ip]=item);
    for(let key in keys) {
      this.machineArr.push(keys[key])
    }
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click",{isBatch:false});
    },
    handleContinue() {
      this.$emit("click",{isBatch:true});
    },
    handleConfirm() {
      let machineuuids = [];
      let deptids = [];
      let deptNames = [];
      this.machineArr.forEach(item => {
        machineuuids.push(item.uuid);
        deptids.push(item.departid+'');
        deptNames.push(item.departname);
      })
      let data = {
        'Name': this.form.batchName, 
        'Descrip': this.form.description, 
        'Manager': this.$store.getters.userName, 
        "DepartID": [...new Set(deptids)],
        "DepartName": deptNames,
        "Machine": machineuuids || [],
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          createBatch(data)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click",{isBatch:false});
                this.$refs.form.resetFields();
                this.$message.success(res.data.msg);
              } else {
                this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("创建失败，请检查输入内容");
            });
        }
      });
    },
  },
};
</script>
<style scoped>
.el-table td.el-table__cell, 
.el-table th.el-table__cell.is-leaf {
  border-bottom: 0px;
  text-align: center;
}
</style>
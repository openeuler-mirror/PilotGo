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
        <el-button type="primary" @click="handleConfirm()">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  createBatch  } from "@/request/batch";
export default {
  props: {
    departInfo: {
      type: Object,
      default: {}
    },
    machineIds: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      idArray: [],
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
    this.idArray.push(this.departInfo.id+'');
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleConfirm() {
      let data = {
        'Name': this.form.batchName, 
        'Description': this.form.description, 
        'Manager': this.$store.getters.userName, 
        "DepartID": this.idArray,
        "Machine": this.machineIds || [],
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          createBatch(data)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click");
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
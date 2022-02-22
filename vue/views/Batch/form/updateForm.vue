<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="批次名称" prop="name">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.name"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="备注" prop="description">
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
        <el-button type="primary" @click="handleSubmitForm">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  updateBatch  } from "@/request/batch";
export default {
  props: {
    row: {
      type: Object,
      default: {
         function(){
           return {}
        }
      }
    } 
  },
  data() {
    return {
      form: {
        name: "",
        description: "",
      },
      rules: {
        name: [{ 
          required: true, 
          message: "请输入名称",
          trigger: "blur"
        }],
      },
      disabled: true,
    }
  },
  mounted() {
    this.form.name = this.row.name;
    this.form.description = this.row.description;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSubmitForm() {
      let params = {
        BatchID: this.row.ID + '',
        BatchName: this.form.name,
        Descrip: this.form.description,
      }
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateBatch(params)
            .then((res) => {
              if (res.data.code === 200) {
                this.$emit("click",'success');
                 this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.error);
              }
            })
        }
      });
    },
  },
};
</script>
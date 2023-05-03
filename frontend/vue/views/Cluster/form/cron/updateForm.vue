<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="任务名称:" prop="taskname">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.taskname"
            autocomplete="off"
          ></el-input>
        </el-form-item>

        <el-form-item label="脚本命令:" prop="cmd">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.cmd"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item style="margin-top: -10px; margin-bottom:0px;">
          <span style="color: rgb(241, 139, 14); font-size: 14px;">cron从左到右(用空格隔开):秒 分 小时 月份中的日期 月份 星期中的日期 年份</span>
          <cron v-if="showCronBox" v-model="form.spec"></cron>
        </el-form-item>
        <el-form-item label="Cron:" prop="spec">
          <el-input
            @focus="showCronBox = true"
            class="ipInput"
            controls-position="right"
            v-model="form.spec"
            autocomplete="off"
          ><el-button slot="append" @click="showCronBox = false" title="确定">确定</el-button></el-input>
        </el-form-item>
        <el-form-item label="任务状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="true">启用</el-radio>
            <el-radio :label="false">暂停</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述:">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.description"
            autocomplete="off"
          ></el-input>
        </el-form-item>
      </el-form>

      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleSubmitForm()">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  updateCron  } from "@/request/cluster";
import cron  from '@/components/VueCron';
export default {
   components: { 
    cron,
  },
  props: {
    row: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      showCronBox: false,
      form: {
        taskname: "",
        spec: "",
        cmd: "",
        status: true,
        description: "",
      },
      rules: {
        taskname: [{ required: true, trigger: "blur", message: "请输入任务名称" }],
        spec: [{ required: true, trigger: "blur",message: "请输入执行脚本" }],
        cmd: [{ required: true, trigger: "blur", message: "请输入cron表达式" }],
      },
    }
  },
  mounted() {
    this.form.taskname = this.row.taskname;
    this.form.spec = this.row.spec;
    this.form.cmd = this.row.cmd;
    this.form.status = this.row.status;
    this.form.description = this.row.description;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSubmitForm() {
      let _this = this; 
      _this.form.id = this.row.ID;
      _this.form.uuid = this.$route.params.detail;
      console.log(_this.form)
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateCron({..._this.form})
            .then((res) => {
              if (res.data.code === 200) {
                _this.$emit("click");
                _this.$message.success(res.data.msg);
                _this.$refs.form.resetFields();
              } else {
                _this.$emit("click");
                _this.$message.error(res.data.msg);
              }
            })
            .catch((res) => {
              this.$message.error("添加失败，请检查输入内容");
            });
        }
      });
    },
  },
};
</script>
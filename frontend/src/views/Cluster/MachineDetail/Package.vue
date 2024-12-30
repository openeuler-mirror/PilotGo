<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="content">
    <div class="repo">
      <span>repo源：</span>
      <el-table
        :data="allRepos"
        height="260px"
        :header-cell-style="{ border: '1px solid #dcdfe6', color: '#000' }"
      >
        <el-table-column
          align="center"
          prop="File"
          label="文件"
          width="400px"
        ></el-table-column>
        <el-table-column
          align="center"
          prop="Enabled"
          label="enabled"
          width="100px"
        ></el-table-column>
        <el-table-column
          align="center"
          prop="URL"
          label="repo地址"
        ></el-table-column>
      </el-table>
    </div>
    <div class="packages">
      <span>请选择软件包：</span>
      <el-autocomplete
        style="width: 30%; margin-right: 10px"
        class="inline-input"
        v-model="packageName"
        @select="onPackageSelected"
        :fetch-suggestions="querySuggestions"
        placeholder="请输入内容"
      ></el-autocomplete>
      <auth-button
        auth="button/rpm_install"
        type="primary"
        name="rpm_install"
        @click="onInstallPackage"
        >安装</auth-button
      >
      <auth-button
        auth="button/rpm_uninstall"
        type="danger"
        name="rpm_uninstall"
        @click="onUnInstallPackage"
        >卸载</auth-button
      >
    </div>
    <div class="info">
      <div class="detail" v-if="display">
        <el-divider content-position="left">软件包详情</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="软件包名">{{
            packageInfo.Name
          }}</el-descriptions-item>
          <el-descriptions-item label="Version">{{
            packageInfo.Version
          }}</el-descriptions-item>
          <el-descriptions-item label="Release">{{
            packageInfo.Release
          }}</el-descriptions-item>
          <el-descriptions-item label="Architecture">{{
            packageInfo.Architecture
          }}</el-descriptions-item>
          <el-descriptions-item label="说明">{{
            packageInfo.Summary
          }}</el-descriptions-item>
        </el-descriptions>
      </div>
      <div class="result" v-else>
        <el-divider content-position="left">执行结果</el-divider>
        <el-descriptions :column="2" border v-loading="loading">
          <el-descriptions-item label="软件包名">{{
            packageName
          }}</el-descriptions-item>
          <el-descriptions-item label="执行动作">{{
            action
          }}</el-descriptions-item>
          <el-descriptions-item label="结果">
            <el-tag
              effect="plain"
              round
              :type="result === '成功' ? 'success' : 'danger'"
              >{{ result }}</el-tag
            >
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";

import AuthButton from "@/components/AuthButton.vue";

import {
  getRepos,
  getInstalledPackages,
  getPackageDetail,
  installPackage,
  removePackage,
} from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute();

// 机器UUID
const machineID = ref(route.params.uuid);

const allRepos = ref<any>([]);
const allPackages = ref<any>([]);

const display = ref(true);
const packageName = ref("");
const packageInfo = ref<any>({});

onMounted(() => {
  getRepos({ uuid: machineID.value })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        allRepos.value = resp.data;

        allRepos.value = [];
        let data = resp.data;
        for (let i = 0; i < data.length; i++) {
          let url = "";
          if (data[i].BaseURL !== "") {
            url = data[i].BaseURL;
          } else if (data[i].MirrorList !== "") {
            url = data[i].MirrorList;
          } else if (data[i].MetaLink !== "") {
            url = data[i].MetaLink;
          }
          allRepos.value.push({
            File: data[i].File,
            ID: data[i].Name,
            URL: url,
            Enabled: data[i].Enabled ? "是" : "否",
          });
        }
      } else {
        ElMessage.error("failed to get machine repo info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get machine repo info:" + err.msg);
    });

  updateInstalledPackage();
});

function updateInstalledPackage() {
  getInstalledPackages({ uuid: machineID.value })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        allPackages.value = resp.data.rpm_all;
      } else {
        ElMessage.error(
          "failed to get machine installed packages info: " + resp.msg
        );
      }
    })
    .catch((err: any) => {
      ElMessage.error(
        "failed to get machine installed packages info:" + err.msg
      );
    });
}

function querySuggestions(query: string, callback: Function) {
  let result: any[] = [];
  allPackages.value.forEach((name: string) => {
    if (name.indexOf(query) === 0) {
      result.push({
        value: name,
      });
    }
  });
  callback(result);
}

function onPackageSelected() {
  getPackageDetail({
    uuid: machineID.value,
    rpm: packageName.value,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        packageInfo.value = resp.data.rpm_info;

        display.value = true;
      } else {
        ElMessage.error(
          "failed to get machine package detail info: " + resp.msg
        );
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get machine package detail info:" + err.msg);
    });
}

const action = ref("");
const result = ref("");
const loading = ref(false);
const handleBaseInfo = (type: string) => {
  loading.value = true;
  action.value = type;
  display.value = false;
  result.value = "";
  let params = {
    uuid: [machineID.value],
    rpm: packageName.value,
  };
  return params;
};

const handleResult = (resp: any, actione_type: string) => {
  display.value = false;
  result.value = "";
  if (resp.code === RespCodeOK) {
    loading.value = false;
    result.value = "成功";
    updateInstalledPackage();
  } else {
    loading.value = false;
    result.value = "失败"
    ElMessage.error(
      "failed to " + actione_type + " machine package detail info: " + resp.msg
    );
  }
};

const rpmRegexTest = () => {
  let regex = /^[A-Za-z0-9+-._]+$/;
  if (!regex.test(packageName.value)) {
    ElMessage.error("请输入合法的软件包名称");
    return false;
  }
  return true;
};
function onInstallPackage() {
  if (!rpmRegexTest()) return;
  installPackage(handleBaseInfo("软件包安装"))
    .then((resp: any) => handleResult(resp, "get"))
    .catch((err: any) => {
      ElMessage.error("failed to get machine package detail info:" + err.msg);
    });
}

const onUnInstallPackage = () => {
  if (!rpmRegexTest()) return;
  removePackage(handleBaseInfo("软件包卸载"))
    .then((resp: any) => handleResult(resp, "remove"))
    .catch((err: any) => {
      ElMessage.error(
        "failed to remove machine package detail info:" + err.msg
      );
    });
};
</script>

<style lang="scss" scoped>
.content {
  width: 80%;
  display: flex;
  flex-direction: column;
  justify-content: space-evenly;

  .repo {
    height: 300px;
  }

  .packages {
    height: 40px;
  }
}
</style>

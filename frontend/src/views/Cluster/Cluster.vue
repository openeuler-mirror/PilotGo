<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <div class="department">
      <PGTree :editable="true" @onNodeClicked="onDepartmentClicked">
        <template v-slot:header>
          <p>部门</p>
        </template>
      </PGTree>
    </div>
    <div class="cluster">
      <PGTable
        :data="machines"
        title="机器列表"
        :showSelect="showSelect"
        :total="total"
        v-model:page="page"
        v-model:selectedData="selectedMachines"
      >
        <template v-slot:action>
          <div class="search">
            <el-input
              class="search_input"
              v-model.trim="searchInput"
              placeholder="请输入关键字进行搜索..."
              @change="onSearchHost"
            />&nbsp;
            <el-button @click="onSearchHost">搜索</el-button>
            <el-divider direction="vertical" style="height: 2.5em" />

            <el-dropdown>
              <el-button>
                操作
                <el-icon>
                  <ArrowDown />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item>
                    <auth-button auth="button/dept_change" link :show="true" @click="showChangeDepartDialog = true">
                      变更部门
                    </auth-button>
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <auth-button auth="button/machine_delete" link :show="true" @click="handleDeleteMachine">
                      删除
                    </auth-button>
                  </el-dropdown-item>

                  <!-- 插件扩展点 -->
                  <el-dropdown-item
                    v-if="hasPermisson('button/monitor_operate')"
                    v-for="item in pluginBtns.filter((item:Extention)=>item.permission.split('.')[1] === 'prometheus')"
                  >
                    <el-button link @click="handlePluginAPI(item.url)">{{ item.name }}</el-button>
                  </el-dropdown-item>
                   <el-dropdown-item
                    v-if="hasPermisson('button/atune_operate')"
                    v-for="item in pluginBtns.filter((item:Extention)=>item.permission.split('.')[1] === 'atune')"
                  >
                    <el-button link @click="handlePluginAPI(item.url)">{{ item.name }}</el-button>
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="hasPermisson('button/logs_operate')"
                    v-for="item in pluginBtns.filter((item:Extention)=>item.permission.split('.')[1] === 'logs')"
                  >
                    <el-button link @click="handlePluginAPI(item.url)">{{ item.name }}</el-button>
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="hasPermisson('button/topology_operate')"
                    v-for="item in pluginBtns.filter((item:Extention)=>item.permission.split('.')[1] === 'topology')"
                  >
                    <el-button link @click="handlePluginAPI(item.url)">{{ item.name }}</el-button>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
        <template v-slot:content>
          <el-table-column align="center" label="ip">
            <template #default="data">
              <el-button link type="primary" title="查看机器详情" @click="machineDetail(data.row)">
                {{ data.row.ip }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="departname" label="部门"> </el-table-column>
          <el-table-column align="center" prop="cpu" label="cpu"> </el-table-column>
          <el-table-column align="center" label="状态">
            <template #default="scope">
              <state-dot :runstatus="scope.row.runstatus" :maintstatus="scope.row.maintstatus"></state-dot>
            </template>
          </el-table-column>
          <el-table-column align="center" label="标签">
            <template #default="scope">
              <el-tag v-show="item.data.length > 0" v-for="item in scope.row.tags">{{ item.data }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="systeminfo" label="系统"> </el-table-column>
          <el-table-column align="center" label="操作">
            <template #default="{ row }">
              <el-button size="small" link type="primary" @click="handleTerminal(row.ip)">进入终端</el-button>
            </template>
          </el-table-column>
        </template>
      </PGTable>
    </div>

    <el-dialog title="主机部门变更" v-model="showChangeDepartDialog" destroy-on-close>
      <change-depart
        :machines="selectedMachines"
        @depart-updated="updateDepartmentMachines({ DepartId: departmentID })"
        @close="showChangeDepartDialog = false"
      />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive, watch, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import axios from "axios";
import AuthButton from "@/components/AuthButton.vue";
import PGTable from "@/components/PGTable.vue";
import PGTree from "@/components/PGTree.vue";
import StateDot from "@/components/StateDot.vue";
import ChangeDepart from "./components/ChangeDepart.vue";
import { directTo } from "@/router/index";
import { getPagedDepartMachines, getMachineTags, deleteMachine } from "@/request/cluster";
import { RespCodeOK, type RespInterface } from "@/request/request";
import { usePluginStore } from "@/stores/plugin";
import { useTerminalStore } from "@/stores/terminal";
import { updatePlugins } from "@/views/Plugin/plugin";
import { hasPermisson } from "@/module/permission";
import type { Extention } from "@/types/plugin";
// 部门树
const departmentID = ref(1);

// 机器列表
const showSelect = ref(true);
const machines = ref<any>([]);
const total = ref(0);
const page = ref({ pageSize: 10, currentPage: 1 });
const searchInput = ref("");

const showChangeDepartDialog = ref(false);
let pluginBtns = ref([] as any);

onMounted(() => {
  updatePlugins();
  updateDepartmentMachines({ DepartId: departmentID.value });
  pluginBtns.value = usePluginStore().extention;
});

const extentions = computed(() => {
  return [...usePluginStore().extention];
});
watch(
  extentions,
  (newV, oldV) => {
    if (newV) {
      pluginBtns.value = newV;
    }
  },
  { immediate: true, deep: true }
);

function updateDepartmentMachines(params: any) {
  getPagedDepartMachines({
    page: page.value.currentPage,
    size: page.value.pageSize,
    ...params,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total;
        machines.value = resp.data;

        // 获取机器节点的tags标签
        let uuids = [];
        for (let i in resp.data) {
          uuids.push(resp.data[i].uuid);
        }
        // let result = resp
        getMachineTags({ uuids: uuids }).then((resp: any) => {
          if (resp.code != 200) {
            ElMessage.error("failed to get machine tags: " + resp.msg);
          }

          for (let n in resp.data) {
            for (let i in machines.value) {
              if (resp.data[n].machineuuid === machines.value[i].uuid) {
                if (!("tags" in machines.value[i])) {
                  machines.value[i].tags = [resp.data[n]];
                } else {
                  machines.value[i].tags.push(resp.data[n]);
                }
              }
            }
          }
        });
      } else {
        ElMessage.error("failed to get machines overview info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get machines overview info:" + err.msg);
    });
}

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      searchInput.value ? onSearchHost() : updateDepartmentMachines({ DepartId: departmentID.value });
    }
  },
  { deep: true }
);
function machineDetail(info: any) {
  directTo("/cluster/machine/" + info.uuid);
}

function onDepartmentClicked(depart: any) {
  departmentID.value = depart.id;
  updateDepartmentMachines({ DepartId: depart.id });
}

// 发送插件的请求
const handlePluginAPI = (url: string) => {
  let uuidArr = Array<string>();
  selectedMachines.value.forEach((item: any) => {
    uuidArr.push(item.uuid);
  });
  axios.post(window.location.origin + url, { uuids: uuidArr }).then((response: any) => {
    if (response.data.data.code === 200) {
      setTimeout(() => {
        updateDepartmentMachines({ DepartId: departmentID.value });
        ElMessage.success(response.data.data.msg);
      }, 2000);
    }
  });
};

/**
 * 模糊搜索
 * @params search:string 模糊关键字
 * 其中状态一栏对应关系：
 * 在线-online
 * | 离线-offline
 * | normal-正常使用
 * | maintenance-维护中
 */
const onSearchHost = () => {
  let searchKey: string = "";
  let stateDict = [
    {
      label: "在线",
      value: "online",
    },
    {
      label: "离线",
      value: "offline",
    },
    {
      label: "正常使用",
      value: "normal",
    },
    {
      label: "维护中",
      value: "maintenance",
    },
  ];

  if (searchInput.value) {
    let filterStates = stateDict.filter((item) => item.label.match(searchInput.value));
    searchKey = filterStates.length > 0 ? filterStates[0].value : searchInput.value;
  }
  updateDepartmentMachines({ search: searchKey });
};

/*
 * 删除机器
 * @params deluuid[] 机器的uuid列表
 *  */
const handleDeleteMachine = () => {
  ElMessageBox.confirm("Are you sure yoou want to delete these machines. Continue?", "Warning", {
    confirmButtonText: "OK",
    cancelButtonText: "Cancel",
    type: "warning",
  })
    .then(() => {
      let uuidArr = Array<string>();
      selectedMachines.value.forEach((item: any) => {
        uuidArr.push(item.uuid);
      });
      deleteMachine({ deluuid: uuidArr }).then((res: RespInterface) => {
        if (res.code === RespCodeOK) {
          page.value.currentPage = 1;
          page.value.pageSize = 10;
          updateDepartmentMachines({ DepartId: 1 });
          ElMessage.success(res.msg);
        } else {
          ElMessage.success(res.msg);
        }
      });
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "Delete canceled",
      });
    });
};

const selectedMachines = ref([]);

// 跳转终端
const handleTerminal = (ip: string) => {
  useTerminalStore().setTerminalIp(ip, "host");
  directTo("/terminal");
};
</script>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  display: flex;

  .department {
    width: 20%;
    height: 100%;
    margin-right: 5px;
  }

  .cluster {
    width: 80%;
    height: 100%;
    .search {
      height: 100%;
      display: flex;
      flex-direction: row;
      align-items: center;
      &_input {
        width: 300px;
      }
    }
  }
}
</style>

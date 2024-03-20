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
      <PGTable :data="machines" title="机器列表" :showSelect="showSelect" :total="total" :onPageChanged="onPageChanged"
        v-model:selectedData="selectedMachines">
        <template v-slot:action>
          <el-dropdown>
            <el-button>
              操作 <el-icon>
                <ArrowDown />
              </el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <auth-button auth="button/dept_change" :show="true" @click="showChangeDepartDialog = true">
                    变更部门
                  </auth-button>
                </el-dropdown-item>
                <el-dropdown-item>
                  <auth-button auth="button/machine_delete" :show="true">
                    删除
                  </auth-button>
                </el-dropdown-item>
                <el-dropdown-item v-for="item in extButtons">
                  <auth-button :auth="item.permission" :show="true" @click="handlePluginAPI(item.url)">
                    {{ item.name }}
                  </auth-button>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-slot:content>
          <el-table-column align="center" label="ip">
            <template #default="data">
              <el-button link type="primary" title="查看机器详情" @click="machineDetail(data.row)">
                {{ data.row.ip }}
              </el-button>
            </template>
          </el-table-column>
          <el-table-column align="center" prop="departname" label="部门">
          </el-table-column>
          <el-table-column align="center" prop="cpu" label="cpu">
          </el-table-column>
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
          <el-table-column align="center" prop="systeminfo" label="系统">
          </el-table-column>
        </template>
      </PGTable>
    </div>

    <el-dialog title="主机部门变更" v-model="showChangeDepartDialog" destroy-on-close>
      <change-depart :machines="selectedMachines" @depart-updated="updateDepartmentMachines(departmentID)"
        @close="showChangeDepartDialog = false" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive, watch } from "vue";
import { ElMessage } from 'element-plus';
import axios from "axios";
import AuthButton from "@/components/AuthButton.vue";
import PGTable from "@/components/PGTable.vue";
import PGTree from "@/components/PGTree.vue";
import StateDot from "@/components/StateDot.vue";
import ChangeDepart from "./components/ChangeDepart.vue";
import { directTo } from "@/router/index"
import { getPagedDepartMachines, getMachineTags } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";
import { usePluginStore } from "@/stores/plugin";
// 部门树
const departmentID = ref(1)

// 机器列表
const showSelect = ref(true)
const machines = ref<any>([])
const total = ref(0)

const showChangeDepartDialog = ref(false)
let extButtons = reactive([] as any)

onMounted(() => {
  updateDepartmentMachines(departmentID.value)
  extButtons = usePluginStore().extention;
})

watch(usePluginStore().extention, (newV, oldV) => {
  if (newV) {
    extButtons = newV;
  }
}, { immediate: true })

function updateDepartmentMachines(departID: number, page: number = 1, size: number = 10) {
  getPagedDepartMachines({
    page: page,
    size: size,
    DepartId: departID,
  }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      total.value = resp.total
      machines.value = resp.data

      // 获取机器节点的tags标签
      let uuids = []
      for (let i in resp.data) {
        uuids.push(resp.data[i].uuid)
      }
      // let result = resp
      getMachineTags({ "uuids": uuids }).then((resp: any) => {
        if (resp.code != 200) {
          ElMessage.error("failed to get machine tags: " + resp.msg)
        }

        for (let n in resp.data) {
          for (let i in machines.value) {
            if (resp.data[n].machineuuid === machines.value[i].uuid) {
              if (!("tags" in machines.value[i])) {
                machines.value[i].tags = [resp.data[n]]
              } else {
                machines.value[i].tags.push(resp.data[n])
              }
            }
          }
        }
      })
    } else {
      ElMessage.error("failed to get machines overview info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get machines overview info:" + err.msg)
  })
}

function onPageChanged(currentPage: number, currentSize: number) {
  updateDepartmentMachines(departmentID.value, currentPage, currentSize)
}

function machineDetail(info: any) {
  directTo("/cluster/machine/" + info.uuid)
}

function onDepartmentClicked(depart: any) {
  departmentID.value = depart.id
  updateDepartmentMachines(depart.id)
}

// 发送插件的请求
const handlePluginAPI = (url: string) => {
  let uuidArr = Array<string>();
  selectedMachines.value.forEach((item: any) => {
    uuidArr.push(item.uuid);
  })
  axios.post(window.location.origin + url, { uuids: uuidArr }).then((response: any) => {
    if (response.data.data.code === 200) {
      setTimeout(() => {
        updateDepartmentMachines(departmentID.value)
        ElMessage.success(response.data.data.msg)
      }, 2000)
    }
  })
}

const selectedMachines = ref([])


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
  }
}
</style>
<template>
  <div class="cluster">
    <div class="dept">
      <el-tree
       class="treeitems"
       :data="data"
       node-key="id"
       :props="defaultProps"
       :load="loadNode"
       lazy
       :default-expanded-keys="[0]"
       @node-click="handleNodeClick"
       draggable
       ref="tree"
     >
     <span class="custom-tree-node" slot-scope="{ node, data }">
       <span>{{ node.label }}</span>
       <span>
         <em @click="() => append(node,data)" class="el-icon-plus"></em>
         <em v-if="data.id !== 0" @click="() => remove(node,data)" class="el-icon-delete"></em>
         <em v-if="data.id !== 0" @click="() => rename(node,data)" class="el-icon-edit"></em>
       </span>
     </span>
    </el-tree>

    </div>
    <div class="info">
      <div class="cluster-title">
        <span class="cluster-title__text">主机列表</span>
        <div class="cluster-title__operate">
          <el-button size="mini" @click="handleAddIp"
            >添加IP</el-button>
            <el-button size="mini" @click="handleUpdateIp"
            >编辑IP</el-button>
          <el-popconfirm
            title="确定删除所选项目吗？"
            @confirm="handleDeleteItems"
          >
            <el-button size="mini" slot="reference">删除</el-button>
          </el-popconfirm>
        </div>
      </div>
      <ky-table
        class="cluster-table"
        ref="table"
        :getData="getClusters"
      >
        <template v-slot:table>
          <el-table-column label="IP" width="90">
            <template slot-scope="scope">
              <router-link
                :to="$route.path + scope.row.id + `?ip=${scope.row.ip}`"
                @click.native="handleSelectIp(scope.row)"
              >
                {{ scope.row.ip }}
              </router-link>
            </template>
          </el-table-column>
          <el-table-column prop="system_cpu" label="cpu" width="90"> 
          </el-table-column>
          <el-table-column label="状态" width="150">
            <template slot-scope="scope">
              <span v-if="scope.row.system_status == 0">异常</span>
              <span v-else>正常</span>
            </template>
          </el-table-column>
           <el-table-column prop="system_info" label="系统信息" width="150"> 
          </el-table-column>
          <el-table-column label="防火墙配置" width="150">
            <template slot-scope="scope">
              <el-button
                size="mini"
                @click="handleFireWall(scope.row.ip)">
                <i class="el-icon-edit-outline"></i>
              </el-button>
            </template>
          </el-table-column>
          <!-- <el-table-column label="类型" width="150">
            <template slot-scope="scope">
              <span v-if="scope.row.machine_type == 0">虚拟机</span>
              <span v-else>物理机</span>
            </template>
          </el-table-column>
          <el-table-column prop="system_version" label="系统版本" width="150">
          </el-table-column>
          <el-table-column prop="arch" label="架构" width="150"> </el-table-column>
          <el-table-column prop="installation_time" label="安装时间">
          </el-table-column> -->
        </template>
      </ky-table>
    </div>

    <el-dialog 
      :title="title"
      :before-close="handleClose" 
      :visible.sync="addIPDialogVisible" 
      width="560px"
    >
     <add-form v-if="type === 'create'" @click="handleClose"></add-form>
     <update-form v-if="type === 'update'" @click="handleClose"></update-form>   
    </el-dialog>
  </div>
</template>

<script>
import kyTable from "@/components/KyTable";
import AddForm from "./form/addForm";
import UpdateForm from "./form/updateForm";
import { getClusters, deleteIp } from "@/request/cluster";
export default {
  name: "Cluster",
  components: {
    AddForm,
    UpdateForm,
    kyTable,
  },
  data() {  
    return {
      title: '',
      type: '',
      dirIdSuff: 0,
      editDirId: 0,
      addIPDialogVisible: false,
      delHostDialogVisible: false,
      multipleSelection: [],
      deleteIPList: [],
      cockpitUrl: "",
      filterText: '',
      data: [{
        id:0,
        label: '中国',     
      }],
      children: [{
        id:1,
        label: '北京',
        children: [{
          id:11,
          label: '通州'
        }]
      },
      {  
        id:2,
        label: '上海',
        leaf: true,
      },
      {  
        id:3,
        label: '山西',
        children:[{
          id: 13,
          label: '太原'
        },{
          id: 14,
          label: '阳泉'
        }]
      },{
        id:4,
        label: '黑龙江',
        children: [{
          id:12,
          label: '哈尔滨'
        }]
      }],
      defaultProps: {
        children: 'children',
        label: 'label',
        isLeaf: 'leaf'
        },
    };
  },
  mounted() {
    // 加载部门节点的接口数据

  },
  methods: {
    getClusters,
    handleAddIp() {
      this.addIPDialogVisible = true;
      this.title = "添加IP";
      this.type = "create"; 
    },
    handleUpdateIp() {
      this.addIPDialogVisible = true;
      this.title = "编辑IP";
      this.type = "update"; 
    },
    handleClose() {
      this.addIPDialogVisible = false;
      this.title = "";
      this.type = "";
    }, 
    append(node,data) {
      this.$prompt('输入节点名字', '新建节点', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({ value }) => {
         this.$message({
             type: 'success',
             message: '新建成功'
           });
         /* http().then((data)=>{
           this.$message({
             type: 'success',
             message: '新建成功'
           }); 
           this.partialRefresh(node)
         })
         //请求失败
         .catch(()=>{
           this.$message({
             type: 'info',
             message: '新建失败'
           }); 
         }) */
       }).catch(() => {}); 
   },
   
   refreshNode(){
     // 封装请求整个树结构数据
   },
   //懒加载
   loadNode(node, resolve){
     if (node.level === 0) {
       return resolve(this.data);
     }
     else if(node.level === 1){
       return resolve(this.children)
     }
     else{
       return resolve([])
     }
   },
   rename(node,data) {
      this.$prompt('输入节点名字', '编辑节点', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({ value }) => {
         this.$message({
             type: 'success',
             message: '修改成功'
           });
         /* http().then((data)=>{
           this.$message({
             type: 'success',
             message: '修改成功'
           }); 
           this.partialRefresh(node)
         })
         //请求失败
         .catch(()=>{
           this.$message({
             type: 'info',
             message: '修改失败'
           }); 
         }) */
       }).catch(() => {}); 
   },
   remove(node,data) {
     this.$confirm('确定删除该节点？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          this.$message({
            type: 'success',
            message: '删除成功!'
          });
        }).catch(() => {});
   },
   //拖拽==>拖拽时判定目标节点能否被放置
   allowDrop(draggingNode, dropNode, type){
     //参数：被拖拽节点，要拖拽到的位置
     //因为根目录是我本地写死的，不能有同级，所以我设置凡是拖拽到的level===1都存放到根节点的下面；
     if(dropNode.level===1){
       return type == 'inner';
     }
     else {
       return true;
     }
   },
   //拖拽==>判断节点能否被拖拽
   allowDrag(draggingNode){
    //第一级节点不允许拖拽
     return draggingNode.level !== 1;
   },
    handleNodeClick(node,data) {
      // 获取当前分支与上级分支的数据
      console.log(node,data);
    },
    // 删除主机界面
    handleDeleteItems() {
      let _this = this;
      let ids = [];
      for (let i of _this.multipleSelection) {
        ids.push(i.id);
      }
      deleteIp({ id: ids }).then((res) => {
        if (res.data.status === "success") {
          this.$refs.table.handleSearch();
          this.$message.success("删除成功");
        }
      });
    },
    handleSelectIp(row) {
      this.$store.commit("SET_SELECTIP", row.ip);
    },
    handleFireWall(ip) {
      this.$router.push({
        name: 'FireWall',
        query: { ip: ip }
      })
    }
  },
};
</script>

<style scoped lang="scss">
.cluster {
  margin-top: 10px;
  .dept {
    width: 36%;
    display: inline-block;
    .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 14px;
    padding-right: 8px;
  }
  }
  .info {
    width: 60%;
    float: right;
  }
  .cluster-title {
    height: 44px;
    background: #3e9df9;
    line-height: 44px;
    padding: 0 10px;
    .cluster-title__text {
      font-size: 14px;
      color: #fff;
    }

    .cluster-title__operate {
      float: right;
      .el-button {
        color: #fff;
        background: transparent;
        border: none;
      }
      .el-button:hover,
      .el-button:active {
        background: #fff;
        color: #3e9df9;
      }
    }
  }

  .deleteHostText {
    margin-left: 10px;
    .del-host {
      color: red;
      font-weight: 600;
      font-size: 18px;
    }
  }
}
</style>

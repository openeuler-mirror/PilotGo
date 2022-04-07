<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-02-22 16:43:19
  LastEditTime: 2022-04-07 09:43:02
 -->
<template>
  <div class="ky-tree">
    <div class="content">
      <el-tree
       class="treeitems"
       empty-text="暂无数据"
       node-key="id"
       :props="defaultProps"
       :load="loadNode"
       :highlight-current="true"
       lazy
       :default-expanded-keys="[1]"
       :expand-on-click-node="false"
       :render-after-expand="isRender"
       @node-click="handleNodeClick"
       ref="multipleTree"
      >
      <span class="custom-tree-node" slot-scope="{ node, data }">
        <span>{{ node.label }}</span>
        <span v-if="showEdit">
          <em @click.stop="() => append(node,data)" class="el-icon-plus"></em>
          <em v-if="data.id !== 1" @click.stop="() => remove(node,data)" class="el-icon-delete"></em>
          <em @click.stop="() => rename(node,data)" class="el-icon-edit"></em>
        </span>
      </span>
      </el-tree>
    </div>
  </div>
</template>
<script>
import { addDepart, deleteDepart, updateDepart } from "@/request/cluster";
export default {
  props: {
    getData: {
      type: Function,
    },
    showEdit: {
      type: Boolean,
      default(){
        return true
      }
    },
  },
  data() {
    return {
      isRender: false,
      defaultProps: {
        children: 'children',
        label: 'label',
        isLeaf: 'isLeaf',
        id: 'id'
      },
    };
  },
  methods: {
    loadNode(node, resolve){
     if (node.level === 0) {
       this.loadFirstNode(resolve);
     }
     else if(node.level >= 1){
       this.loadChildrenNode(node,resolve);
     }
     else{
       return resolve([])
     }
   },
    // 加载根
    async loadFirstNode(resolve) {
      let res = await this.getData({'DepartID': this.$store.getters.UserDepartId});
      if(res.data.code === 200) {
        let data = res.data.data;
        let rootNode = [{
          id: data.id,
          label: data.label,
          pid: data.pid
        }]
        return resolve(rootNode);
      } else {
        this.$message({
          type: 'error',
          message: res.data.msg
        });
      }
    },
    // 加载子树
    async loadChildrenNode(node,resolve) {
      let deptId = node.key;
      let res = await this.getData({'DepartID': deptId});
      if(res.data.code === 200) {
        if(node.childNodes.length == 0) {
          node.isLeaf = true;
        }
        let children = res.data.data.children == null ? [] : res.data.data.children
        return resolve(children);
      } else {
        this.$message({
          type: 'error',
          message: res.data.msg
        });
      }
    },
    // 追加节点
    append(node,data) {
      this.$prompt('输入节点名字', '新建节点', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({value}) => {
         addDepart({'PID':data.id,'ParentDepart':data.label,'Depart':value}).then((data)=>{
           this.$message({
             type: 'success',
             message: '新建成功'
           }); 
          node.loaded = false;
          node.expand();
         })
         .catch(()=>{
           this.$message({
             type: 'info',
             message: '新建失败'
           }); 
         })
       }).catch(() => {}); 
    },
    // 修改节点名
    rename(node,data) {
      this.$prompt('输入节点名字', '编辑节点', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({ value }) => {
         updateDepart({'DepartID': data.id, 'DepartName': value}).then((data)=>{
           this.$message({
             type: 'success',
             message: '修改成功'
           }); 
          node.parent.loaded = false;
          node.parent.expand();
         })
         .catch(()=>{
           this.$message({
             type: 'info',
             message: '修改失败'
           }); 
         })
       }).catch(() => {}); 
    },
    // 删除节点
    remove(node,data) {
     this.$confirm('确定删除该节点？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          deleteDepart({'DepartID':data.id}).then((data)=>{
           this.$message({
             type: 'success',
             message: '删除成功'
           }); 
          node.parent.loaded = false;
          node.parent.expand();
         })
         .catch(()=>{
           this.$message({
             type: 'info',
             message: '删除失败'
           }); 
         })
       }).catch(() => {}); 
    },
    //拖拽==>拖拽时判定目标节点能否被放置  draggable属性最后做
    allowDrop(draggingNode, dropNode, type){
     //参数：被拖拽节点，要拖拽到的位置
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
      this.$emit("nodeClick",node);
    },
  },
};
</script>

<style rel="stylesheet/scss" lang="scss">
.ky-tree {
  .custom-tree-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 16px;
    padding-right: 8px;
    em {
      color: rgb(11, 35, 117)
    }
  }
  .el-icon-caret-right:before {
    content: "\27a7";
    font-size: 18px;
    font-weight: bold;
    color: rgb(11, 35, 117);
}
}
</style>
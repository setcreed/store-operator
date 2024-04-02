<template>
    <div>
        <div style="width: 95%;margin:10px auto">
            <el-table
                    :data="dblist"
                    border
                    style="width: 100%">
              <el-table-column label="名称" width="150">
                <template slot-scope="scope">
                  <p>  {{ scope.row.metadata.name }}  </p>
                </template>
              </el-table-column>
                <el-table-column label="命名空间" width="150">
                  <template slot-scope="scope">
                    <p>  {{ scope.row.metadata.namespace }}  </p>
                  </template>
                </el-table-column>
              <el-table-column label="副本数" width="150">
                <template slot-scope="scope">
                  <p>  {{ scope.row.spec.replicas }}  </p>
                </template>
              </el-table-column>
              <el-table-column label="DSN" width="150">
                <template slot-scope="scope">
                  <p>  {{ scope.row.spec.dsn }}  </p>
                </template>
              </el-table-column>
              <el-table-column label="状态" width="150">
                <template slot-scope="scope">
                  <p>  {{ scope.row.status.ready }}  </p>
                </template>
              </el-table-column>

              <el-table-column label="操作" width="150">
                <template slot-scope="scope">
                  <p>   <el-link  @click="()=>showEvent(scope.row)" >
                    查看事件<i class="el-icon-info el-icon--right"></i></el-link>
                  </p>
                  <p>   <el-link  @click="()=>rmConfig(scope.row)" >
                    删除<i class="el-icon-s-custom el-icon--right"></i></el-link>
                  </p>
                </template>
              </el-table-column>
            </el-table>
        </div>

    </div>

</template>
<script>
   module.exports ={
       data(){
           return {
               dblist: [],
               lock:false
           }
       },
       mounted(){
         this.loadData()
       },
       methods: {
         showEvent(row){
           const {namespace,name}=row.metadata
           axios.get("/events/"+namespace+"/"+name).then(rsp=>{
             const h = this.$createElement;
             let events=[]
             for(var i=0;i<rsp.data.length;i++){
               let e=rsp.data[i]
               events.push(h("p",{style:"margin-bottom:10px"},
                   "type:"+e.type+" reason:"
                   +e.reason+" message:"+e.message+" lastTimestamp"+e.lastTimestamp))
             }
             this.$notify({
               title: '事件',
               message: h('div', { style: 'color: teal'},
                   events)
             });
           })
         },

         loadDBList(){ //加载列表
            axios.get("/configs").then(rsp=>{
              this.dblist=rsp.data
              this.lock=false
              this.loadData()
            })
         },
         loadData(){
           if(this.lock) return ;
           setTimeout( ()=>{
             this.lock=true
             this.loadDBList()
           },1000)
         },
         rmConfig(row){
           this.$confirm('是否删除配置?', '提示', {
             confirmButtonText: '确定',
             cancelButtonText: '取消',
             type: 'warning',
             center: true
           }).then(()=>{
              axios.delete("/configs/"+row.metadata.namespace+"/"+row.metadata.name)
           })
         }
       }
   }
</script>
<style>
.cell{text-align: center}
  .sdt{margin: 10px auto;width:90%;border-radius: 5px;display: block;float:left;margin-left: 50px}
  .sdt dt{width:100%;display: block;color:#3A7B43;font-size:16px;font-weight: bold;margin-bottom: 10px}
  .sdt dd{float:left;margin:0 auto;text-indent:1em}
  .sdt .row{width:100%;}
  .sdt dd .search{width:50%;}
  .sdt dd a{color: #3a8ee6}
  .sdt dd a:hover{background: #eb5975;color:#fff}
  .sdt dd .select{background: #eb5975;color:#fff}
  .sdt dd .numtext{width:100px}
  .sdt dd span{margin: 0 auto}
    a{cursor:pointer}
    .el-pagination{margin:0 auto;float: left}
</style>
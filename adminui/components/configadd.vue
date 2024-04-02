<template>
    <div>
        <div style="width: 95%;margin:10px auto">
          <div slot="header" class="clearfix">
            <div  >
              <el-form    >
                <el-form-item label="填写配置内容" style="width: 100%">
                  <textarea ref="myeditor"></textarea>
                </el-form-item>
                <el-form-item style="text-align: center">
                  <el-button type="primary" @click="addConfig">提交</el-button>
                </el-form-item>
              </el-form>

            </div>
          </div>
        </div>

    </div>

</template>
<script>
   module.exports ={
       data(){
           return {
              config: ''
           }
       },
       mounted(){
         this.initEditor()
       },
       methods: {
         addConfig(){
           axios.post("/configs",this.config).then(rsp=>{
                alert("ok") //这里可以做一些跳转等操作
           }).catch(e=>{
              console.log(e)
           })
         },
         initEditor(){
           //初始化 yaml 编辑器
           this.editor = CodeMirror.fromTextArea(this.$refs.myeditor, {
             mode: 'yaml', // 语法model
             theme: 'monokai', // 编辑器主题
             tabSize: 2, // 缩进格式
           })
           this.editor.on("change",  (editor, changes)=>{
             this.config=editor.getValue()
           });
         },
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
<template>

<div>
<iframe  :src="urlPath" frameborder="0" width="100%" height="900"></iframe>
</div>
</template>


<script>
export default{
    data(){
        return{
            urlPath:"",
        }
    },
    created(){
        this.urlPath=this.$route.query.urlPath,
        console.log(this.urlPath)
    }
}
</script>




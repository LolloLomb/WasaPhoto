<template>
    <nav class="banner fixed-top" style="text-align: center">
        <a>
            <img src="../assets/Wasa_Logo.png" class="logo" @click="refresh">
        </a>
    </nav>
</template>

<script>
export default{
  data(){
    return{
    }
  },
  methods:{
    refresh(){
      this.$router.replace('/home')
    },
  },
}
</script>

<style>

.logo{
    max-height: 60px;
    padding-right: 10px;
    cursor:pointer;
}

.banner {
  justify-content: space-around; /* Distribuzione uniforme degli elementi */
  align-content: center;
  width:100px;
  height: 100px;
  background-color: white;
  border: 5px solid #abd3da;
  border-bottom-right-radius: 50px;
}

</style>
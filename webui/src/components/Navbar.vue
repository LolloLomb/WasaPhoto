<script>
export default{
  data(){
    return{
    }
  },
  methods:{
    logout(){
      localStorage.clear()
      this.$emit('logoutNavbar',false)
    },
  },
}
</script>



<template>
  <nav class="navbar fixed-bottom mx-auto mb-1">
      <a href="#home">
        <img src="../assets/home.svg" :class="{'iconNormalselect': $route.path === '/home', 'iconNormal': $route.path !== '/home'}" alt="">
      </a>
      <a href="#search">
        <img src="../assets/search.svg" :class="{'iconNormalselect': $route.path === '/search', 'iconNormal': $route.path !== '/search'}" alt="">
      </a>
      <a href="#upload">
        <img src="../assets/upload.svg" :class="{'iconNormalselect': $route.path === '/upload', 'iconNormal': $route.path !== '/upload'}" alt="">
      </a>
      <a href="#profile">
        <img src="../assets/profile.svg" :class="{'iconNormalselect': $route.path === '/profile', 'iconNormal': $route.path !== '/profile'}" alt="">
      </a>
      <a href="#">
        <img src="../assets/logout.svg" class="iconOut" alt="" @click="logout">
      </a>
  </nav>
</template>


<style scoped>
/* Stili opzionali per il footer */

.iconNormal, .iconOut {
  width: 60px;
  height: 60px;
  filter: drop-shadow(3px 5px 2px rgb(0 0 0 / 0.2));
}

.iconOut:hover {
  filter: invert(39%) sepia(50%) saturate(6189%) hue-rotate(342deg) brightness(96%) contrast(109%)
}

.iconNormalselect {
  width: 60px;
  height: 60px;
  filter: invert(76%) sepia(61%) saturate(3180%) hue-rotate(143deg) brightness(92%) contrast(81%) drop-shadow(3px 5px 2px rgb(0 0 0 / 0.2));
}

.navbar {
  justify-content: space-around; /* Distribuzione uniforme degli elementi */
  width: 40%;
}

.iconNormal:hover, .iconOut:hover,
.iconNormalselect:hover, .iconOutselect:hover {
  cursor: pointer;
  display: inline-block;
  animation: bounce 1.5s;
  -webkit-animation: bounce 0.5s forwards;
  -moz-animation: bounce 0.5s forwards;
  -o-animation: bounce 0.5s forwards;
}

@-webkit-keyframes bounce {
    0% {
      -webkit-transform: translateY(0);
    } 
    100% {
      -webkit-transform: translateY(-10px);
    }
}
@-moz-keyframes bounce {
    0%{
       -moz-transform: translateY(0);
    }
    100% {
       -moz-transform: translateY(-10px);
    }
}
@keyframes bounce {
    0% {
        transform: translateY(0);
    }
    100% {
        transform: translateY(-10px);
    }
}

</style>

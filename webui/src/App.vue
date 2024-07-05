<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
			logged: false,
		}
	},
	methods:{
		logout(newValue){
			this.logged = newValue
			this.$router.replace("/login")
		},
		updateLogged(newLogged){
			this.logged = newLogged
		},
		updateView(newRoute){
			this.$router.replace(newRoute)
		},
	},
	
	mounted(){
		if (localStorage.getItem('remember') == 'false') {
			this.logged = false
			this.$router.replace('/login')
		}
		
		if (localStorage.getItem('remember') == 'true') {
			this.logged = true
		}
	},
}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main >
					<Navbar v-if="logged" 
					@logoutNavbar="logout" 
					@requestUpdateView="updateView"/>

					<RouterView 
					@updatedLoggedChild="updateLogged" 
					@requestUpdateView="updateView"/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
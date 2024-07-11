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
		refresh(){
			location.reload()
		},
		updateLogged(newLogged){
			this.logged = newLogged
		},
		updateView(newRoute){
			this.$router.replace(newRoute)
		},
		leaving(){
			localStorage.clear()
			this.logged = false
			this.$router.replace('login')
		}
	},
	
	mounted(){
		if (localStorage.getItem('remember') == 'false') {
			this.logged = true
		}
	},

	ready(){
		window.onbeforeunload = this.leaving
	}
}
</script>

<template>
	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main >
					<Navbar v-if="logged" 
					@logoutNavbar="logout" />
					<Banner v-if="logged" 
					@refresh="refresh"/>
					<RouterView 
					@updatedLoggedChild="updateLogged"/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
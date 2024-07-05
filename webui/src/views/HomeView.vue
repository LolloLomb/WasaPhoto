<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async streamLoader() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/user/" + localStorage.getItem('token') + "/stream")
				if (response.data != null){
					this.photos = response.data
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		if (!localStorage.getItem('token')) {
			this.$router.replace('/login')
		}
		this.streamLoader()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3">
			<h1 class="h2">Homepage</h1>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

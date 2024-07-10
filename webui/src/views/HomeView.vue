<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			photos: [],
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
		<div class="row mx-auto" style="margin-top: 100px;">
            <div class="col-md-4 d-flex justify-content-center" style="margin-bottom: 100px;" v-for="(photo, index) in photos" :key="index">
                <Photo 
                    :owner="photo.username_owner" 
                    :photo_id="photo.ID" 
                    :comments="photo.comments" 
                    :likes="photo.like_username"
					:upload_date="photo.upload_date"
					:isOwner="false"/>
            </div>
        </div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

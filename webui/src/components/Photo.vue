<script>
export default {
    data(){
        return {
            photoPath: null,
            liked: false,
            AllComments: [],
            AllLikes: [],
        }
    },

    props: ['owner','likes','comments',"upload_date","photo_id","isOwner"], 

    methods: {
        async loadPhoto() {
            try {
                let response = await this.$axios.get("/photo/" + this.photo_id);
                const photoData = `data:image/png;base64,${response.data.content}`;
                this.photoPath = photoData;
            } catch (error) {
                console.error("Error loading photo:", error);
            }
        }   
    },

    async mounted(){
		this.loadPhoto()
		if (this.likes != null){
			this.AllLikes = this.likes
		}

		if (this.likes != null){
			this.liked = this.allLikes.some(obj => obj.username === localStorage.getItem('username'))
		}
		if (this.comments != null){
			this.AllComments = this.comments
		}
	},
}

</script>

<template>
    <div>
        <img :src="photoPath" alt="Photo" v-if="photoPath"/>
    </div>
</template>

<style>

</style>
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

    props: ['owner','likes','comments',"photo_id","isOwner"], 

    methods: {
        async loadPhoto() {
            try {
                let response = await this.$axios.get("/photo/" + this.photo_id);
                const photoData = `data:image/png;base64,${response.data.content}`;
                this.photoPath = photoData;
            } catch (error) {
                console.error("Error loading photo:", error);
            }
        },
        async toggleLike(){
            if (!this.liked){
                try{
                    let response = await this.$axios.post("/photo/" + this.photo_id + "/likes", {username:localStorage.getItem('username')})
                    this.liked = !this.liked
                }
                catch (error) {
                    console.error(error)
                }
            }
            else {
                try{
                    let response = await this.$axios.delete("/photo/" + this.photo_id + "/likes/" + localStorage.getItem("token"))
                    this.liked = !this.liked
                }
                catch(error) {
                    console.error(error)
                }
            }
        }   
    },

    async mounted(){
		this.loadPhoto()
		if (this.likes != null){
			this.AllLikes = this.likes
		}

		if (this.likes != null){
			this.liked = this.AllLikes.some(obj => obj === localStorage.getItem('username'))
		}
		if (this.comments != null){
			this.AllComments = this.comments
		}
	},
}

</script>

<template>
    <div class="post-container">
        <div class="photo-container">
            <img :src="photoPath" alt="Photo" v-if="photoPath"/>
        </div>
        <div>
           <!-- Bottone per aprire i commenti e per mettere like -->
           <div class="row mx-auto" style="text-align: center; width: 80%; margin-top: 10px;">
            <div class="col" style="text-align: center;">
                <button class="commentsButton" @click="openCommentsModal">Comments
                </button>
            </div>
            <div class="col my-auto" v-if="!isOwner">
                <svg @click="toggleLike" xmlns="http://www.w3.org/2000/svg" width="40" height="40" fill="#3d93ad" class="bi-heart" viewBox="0 0 16 16" v-if="!liked">
                    <path d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143q.09.083.176.171a3 3 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15"/>
                </svg>
                <svg @click="toggleLike" xmlns="http://www.w3.org/2000/svg" width="40" height="40" fill="#fd5252" class="bi-heart-fill" viewBox="0 0 16 16" v-if="liked">
                    <path fill-rule="evenodd" d="M8 1.314C12.438-3.248 23.534 4.735 8 15-7.534 4.736 3.562-3.248 8 1.314"/>
                </svg>
            </div>
        </div>
        </div>
    </div>
</template>

<style scoped>

.bi-heart:hover, .bi-heart-fill:hover {
    cursor:pointer;
}

.commentsButton{
    margin: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 12px 30px;
    font-family: -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
    border-radius: 6px;
    border: none;
    color: #fff;
    background: linear-gradient(180deg, #3d93ad 0%, #2a6679 100%);
    background-origin: border-box;
    box-shadow: 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2);
    user-select: none;
    -webkit-user-select: none;
    touch-action: manipulation;
    display: block;
}

.commentsButton:hover {
  box-shadow: inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2), 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), 0px 0px 0px 3.5px rgba(58, 108, 217, 0.5);
  outline: 0;
}

.post-container {
    border:5px solid #abd3da;
    padding:8px;
    border-radius: 25px;
}

.photo-container {
    width: 420px; /* Sostituisci con la dimensione desiderata del quadrato */
    height: 420px; /* Sostituisci con la dimensione desiderata del quadrato */
    display: flex;
    padding: 2px;
    justify-content: center;
    align-items: center;
    overflow: hidden; /* Nascondi eventuali parti dell'immagine che fuoriescono */
    border: 2px solid rgb(197, 243, 255);
    border-radius: 25px;
}

.photo-container img {
    width: 100%;
    height: 100%;
    object-fit: contain; /* Adatta l'immagine nel contenitore mantenendo le proporzioni */
}
</style>

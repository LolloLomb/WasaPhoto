<script>
export default {
    data(){
        return {
            photoPath: null,
            UploadDate: null,
            liked: false,
            AllComments: [],
            AllLikes: [],
            isCommentsMode: false,
            newComment: "",
            currentCommentId: 0,
        }
    },

    props: ['owner','likes','comments',"photo_id","isOwner", "upload_date"], 

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
        },

        isCommentOwner(comment){
            return comment.username_owner == localStorage.getItem('username') || this.isOwner
        },  

        async deleteComment(comment){
            try{
                let response = await this.$axios.delete("/photo/" + this.photo_id + "/comment/" + comment.ID)
                this.AllComments = this.AllComments.filter(obj => obj.ID !== comment.ID)

            }
            catch(error){
                console.log(error)
            }
        },

        commentsModalIn(){
            this.isCommentsMode = true
        },

        commentModalOut() {
            this.isCommentsMode = false
        },

        async submitComment() {
            try {
                let response = await this.$axios.post("/photo/" + this.photo_id + "/comment", {
                    comment_content: this.newComment,
                    username_owner: localStorage.getItem('username'),
                });
                this.currentCommentId = response.data.success
                this.AllComments.push({username_owner: localStorage.getItem('username'), comment_content: this.newComment, ID: this.currentCommentId});
                this.newComment = ""; // Reset textarea
                this.currentCommentId = 0;
            } catch (error) {
                console.error("Error submitting comment:", error);
            }
        },

        async removePhoto() {
            try {
                let response = await this.$axios.delete("/photo/" + this.photo_id)
                this.$emit('removePhoto', this.photo_id)
            }
            catch(error) {
                console.log(error)
            }
        },
    },

    async mounted(){
		this.loadPhoto()
		if (this.likes != null){
			this.AllLikes = this.likes
		}

        if(this.upload_date != null){
            this.UploadDate = this.upload_date
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

    <div v-if="isCommentsMode" class="popup-overlay" @click.self="commentModalOut">
        <div class="popup-content" style="width: 40%;">
            <span class="close-button" @click="commentModalOut">&times;</span>
                <div style="word-wrap: break-word;">
                    <p v-for="(comment, index) in AllComments" :key="index" style="margin-top: 0.2rem; margin-left: 10px; font-size: 16px;">
                        <svg @click="deleteComment(comment)" v-if="isCommentOwner(comment)" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="red" class="bi-trash" viewBox="0 0 16 16" style="margin-right: 8px;">
                            <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
                            <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
                        </svg>
                        {{ comment.username_owner }} : {{ comment.comment_content}}
                    </p>
                </div>
            <!-- Form to add new comment -->
            <form @submit.prevent="submitComment" style="text-align: center;">
                <input type="text" class="inputComment" v-model="newComment" placeholder="Add a comment..." rows="1" cols="50" style="width: 100%;" spellcheck="false" minlength="1"/>
                <button class="submit" :disabled="newComment.trim().length < 1">Submit</button>
            </form>
        </div>
    </div>
    <div class="post-container">

        <div v-if="!isOwner">
            <div>
                {{  }}
            </div>
            <div>
                {{ UploadDate }}
            </div>
        </div>

        <div class="photo-container">
            <img :src="photoPath" alt="Photo" v-if="photoPath"/>
        </div>
        <div>
           <div class="row mx-auto" style="text-align: center; width: 80%; margin-top: 10px;">
            <div class="col" style="text-align: center;">
                <button class="commentsButton" @click="commentsModalIn">Comments
                </button>
            </div>
            <div class="col" style="text-align: center;" v-if="isOwner">
                <button class="removeButton" @click="removePhoto">Remove
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

.inputComment {
    border: 2px solid #3d93ad;
    border-radius: 15px;
    padding-left: 10px;
    }

.inputComment:focus::placeholder{
    color:transparent
}

.bi-trash:hover {
    cursor:pointer;
}

.submit:disabled{
    background: linear-gradient(180deg, #a9acad 0%, #dfdfdf 100%);
    border: 2px solid white;
}

.submit{
    margin-top:20px;
    margin-bottom: 20px;
    border: 2px solid #2a6679;
    border-radius: 15px;
    padding: 8px 17px;
    font-size: 15px;;
    background: linear-gradient(180deg, #3d93ad 0%, #2a6679 100%);
    color:white;
}

.close-button {
  position: absolute;
  top: 8px;
  right: 20px;
  font-size: 30px;
  cursor: pointer;
}

.popup-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.popup-content {
  background: white;
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 60px;
  border-radius: 20px;
  border: 10px solid #abd3da;
  position: relative;
}


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

.removeButton{
    margin: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 12px 30px;
    font-family: -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
    border-radius: 6px;
    border: none;
    color: #fff;
    background: linear-gradient(180deg, #db1919 0%, #7c3333 100%);
    background-origin: border-box;
    box-shadow: 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2);
    user-select: none;
    -webkit-user-select: none;
    touch-action: manipulation;
    display: block;
}

.removeButton:hover {
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

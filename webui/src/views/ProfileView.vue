<script>
export default {
	data: function() {
		return {
            username: '',
            followersCount: 0,
            followingCount: 0,
            posts_amount: 0,
            followers: [],
            following: [],
            currentBanned: false,
            token: null,
            banStatus: false,
            userExists: false,
            followStatus: false,
            isChangingName: false,
            newUsarname: null,
            alreadyTaken: false,
		}
	},

    watch: {
        currentPath(newName, oldName){
            if (newName !== oldName){
                this.loadInfo()
            }
        }
    },

	methods: {
        async loadInfo() {
            if (this.$route.params.username === undefined){
                return
            }
            try{
                this.username = this.$route.path.split("/")[2]
                let id = await this.$axios.get("/get_id?username=" + this.username)
                this.userExists = true
				let response = await this.$axios.get("/user/" + id.data.success);

                this.banStatus = false
                this.currentBanned = false
                
                this.nickname = response.data.username
				this.followersCount = response.data.followers != null ? response.data.followers.length : 0
				this.followingCount = response.data.following != null ? response.data.following.length : 0
				this.posts_amount = response.data.postsAmount != 0 ? response.data.posts : 0
                
			}catch(e){
				this.currentBanned = true
			}
        },
        toggleBan(){
            this.banStatus = !this.banStatus
            if (this.banStatus == true) {
                this.followStatus = false
            }
        },
        toggleFollow(){
            this.followStatus = !this.followStatus
        },
        changeNameModalIn() {
            this.isChangingName = true
        },
        changeNameModalOut() {
            this.isChangingName = false
        },
        toLowerCase(event) {
			this.newUsarname = this.newUsarname.toLowerCase();
		},
        async changeUsername(){
			try{
				let resp = await this.$axios.put("/user/" + localStorage.getItem("token") + "/username",{
					newUsername: this.newUsarname,
				})

                this.username = this.newUsarname
                localStorage.setItem('username', this.username)
                this.newUsarname=""     
                this.isChangingName=false
                this.alreadyTaken=false           

			} 
            catch(e){
                this.alreadyTaken = true;
			}
		},
	},

	async mounted() {
        if (!localStorage.getItem('token')) {
            this.$router.replace('/login')
        }
        await this.loadInfo()
	},

    computed: {
        currentPath() {
            return this.$route.params.username
        },
        isOwner() {
            return this.$route.params.username === localStorage.getItem('username')
        }
    }
}
</script>

<template>
    <div class="popup-overlay" v-if="isChangingName" @click.self="changeNameModalOut">
      <div class="popup-content">
        <span class="close-button" @click="changeNameModalOut">&times;</span>
        <form @submit.prevent="submitNewUsername">
          <div>
            <label for="username">New Username:</label>
            <input type="text" class="form-control" v-model="newUsarname" @input="toLowerCase"
								maxlength="16" minlength="3" placeholder="Insert new username" />
          </div>
          <button class="confirmButton mt-2" @click="changeUsername">Confirm</button>
          <div class="mx-auto" v-if="this.alreadyTaken">
            <p style="color:red; margin-top: 10px;">Nome utente già in uso!</p>
        </div>
        </form>
      </div>
    </div>

	<div class = "info" v-if="userExists && !currentBanned">
        <div class="row" style="width: 90%;">
            <div class="col" style="margin-top: 15px;">
                <div style="margin-left: 350px;">
                <svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" class="bi bi-person-square" viewBox="0 0 16 16" style="width: 200px;">
                    <path d="M11 6a3 3 0 1 1-6 0 3 3 0 0 1 6 0"/>
                    <path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2zm12 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1v-1c0-1-1-4-6-4s-6 3-6 4v1a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1z"/>
                </svg>
                <h3 class="name_under_propic"> {{ username }} </h3>
            </div>
            </div>
            <div class="col my-auto" style="padding-right: 300px;">
            <div class="row">
            <div class="col my-auto">
                <p class="textUp" style="text-align: center">Posts</p>
                <p class="textDown" style="text-align: center;">{{ posts_amount }}</p>
            </div>
            <div class="col my-auto">
                <p class="textUp" style="text-align: center;">Followers</p>
                <p class="textDown" style="text-align: center;">{{ followersCount }}</p>
            </div>
            <div class="col my-auto">
                <p class="textUp" style="text-align: center;">Following</p>
                <p class="textDown" style="text-align: center;">{{ followingCount }}</p>
            </div>
        </div>
        <div class="row mx-auto" style="text-align: center; width: 50%;" v-if="!isOwner">
            <div class="col">
                <button class="followButton" @click="toggleFollow">
                    {{ followStatus ? 'Unfollow' : 'Follow' }}
                </button>
            </div>
            <div class="col">
                <button class="banButton" @click="toggleBan">
                    {{ banStatus ? 'Unban' : 'Ban' }}
                </button>
            </div>
        </div>
        <div v-else class="row mx-auto" style="text-align: center; width: 50%;">
            <div class="col">
                <button class="setButton" @click="changeNameModalIn">
                    Change username
                </button>
            </div>
        </div>
    </div>
        </div>
    </div>
    <div>
        <!--Qua ci va la logica di griglia per gestire le foto-->

    </div>
</template>

<style>

.setButton {
    display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 14px;
  font-family: -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
  border-radius: 6px;
  border: none;
  color: #fff;
  background: linear-gradient(180deg, #d1ad48 0%, #d1a837 100%);
   background-origin: border-box;
  box-shadow: 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2);
  user-select: none;
  -webkit-user-select: none;
  touch-action: manipulation;
  margin: 0 auto;
  display: block;
}

.banButton {
    display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 14px;
  font-family: -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
  border-radius: 6px;
  border: none;
  color: #fff;
  background: linear-gradient(180deg, #fd5252 0%, #812e2e 100%);
   background-origin: border-box;
  box-shadow: 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2);
  user-select: none;
  -webkit-user-select: none;
  touch-action: manipulation;
  margin: 0 auto;
  display: block;
}

.followButton, .confirmButton {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 6px 14px;
  font-family: -apple-system, BlinkMacSystemFont, 'Roboto', sans-serif;
  border-radius: 6px;
  border: none;
  color: #fff;
  background: linear-gradient(180deg, #3dad42 0%, #31812e 100%);
   background-origin: border-box;
  box-shadow: 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2);
  user-select: none;
  -webkit-user-select: none;
  touch-action: manipulation;
  margin: 0 auto;
  display: block;
}

.followButton:hover, .banButton:hover, .setButton:hover {
  box-shadow: inset 0px 0.8px 0px -0.25px rgba(255, 255, 255, 0.2), 0px 0.5px 1.5px rgba(54, 122, 246, 0.25), 0px 0px 0px 3.5px rgba(58, 108, 217, 0.5);
  outline: 0;
}

.textDown {
    font-size: 24px;
}

.textUp {
    font-size: 26px;
}

.name_under_propic {
    text-align: center;
    padding-top: 10px;
}

.bi-person-square{
    text-align: center;
    margin: auto;
    display: block;
    color: #4295a3;
}

.info {
    padding-top: 100px;
    animation: fadeIn 1s ease-in-out;
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
  padding: 20px;
  border-radius: 8px;
  position: relative;
  width: 300px;
  text-align: center;
}

.close-button {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 24px;
  cursor: pointer;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

</style>

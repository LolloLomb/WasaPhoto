<script>
export default {
    data: function (){
        return {
            search_username: "",
            search_results: []
        }
    },

    emits: ['updatedLoggedChild'],

    methods: {
        async search() {
            this.search_username = this.search_username.toLowerCase();
            try {
                let response = await this.$axios.get("/user?username=" + this.search_username)
                this.search_results = response.data
            }
            catch(e) {
                this.search_results = []
            }
        },
        goTo(user) {
            this.$router.replace("/profile/" + user)
        },
    }
}
</script>

<template>
    <div class="search-bar-container">
        <input type="text" class="input-search" placeholder="Search user" v-model="search_username" @input="search"
        minlength='3' maxlength='16' />
    </div>
    <div class="results-container" v-if="search_username.length > 0 && search_results.length > 0">
        <div class="result" v-for="user in search_results" :key="user" @click="goTo(user); search_results=[]">
            {{ user }}
        </div>
    </div>
</template>

<style scoped>

.results-container {
    position: fixed;
    top: 430px;
    left: 50%;
    transform: translateX(-50%);
    width: 100%;
    max-width: 600px;
    background-color: white;
    border: 1px solid #ddd;
    border-radius: 25px;
    border-color: #1ebbd7;
    text-align: center;
    letter-spacing: 2px;
    box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
    overflow-y: auto;
    max-height: 400px;
    z-index: 1000;
}

.result {
    padding: 10px;
    border-bottom: 1px solid #eee;
    cursor: pointer;
    z-index: 100px;
}

.result:hover {
    background-color:rgb(30,87,215,0.1);
}

.search-bar-container {
    position: fixed;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 100%;
    padding: 370px;
    display: flex;
    justify-content: center;
}

.search-bar {
    width: 50%;
    max-width: 600px;
    padding: 10px;
    border-radius: 25px;
}

.input-search{
    height: 50px;
    width: 500px;
    padding: 10px;
    font-size: 18px;
    letter-spacing: 1px;
    outline: none;
    border-radius: 25px;
    transition: all .5s ease-in-out;
    background-color: white;
    border-color: #1ebbd7;
    border-width: 4px;
    text-align: center;
    color:black;
}
.input-search::placeholder{
    color:black;
    font-size: 18px;
    letter-spacing: 2px;
    font-weight: 500;
    text-align: center;
    opacity: 0.3;
}

.input-search:hover{
    width: 900px;
    border-radius: 25px;
    background-color: white;
    transition: all 1200ms;
    border-color: #1ebbd7;
}

.input-search:focus{
    width: 900px;
}

.input-search:focus::placeholder {
    color: transparent;
}

</style>
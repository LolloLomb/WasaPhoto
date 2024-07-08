<script>
export default {
    data(){
        return {
            photoPath: null,
            liked: false,
            Comments: [],
            Likes: [],
        }
    },
    props: ['owner','likes','comments',"upload_date","photo_id","isOwner"],

    methods: {
        loadPhoto(){
            //this.photoPath = 
        }
    }
}

</script>

<template>

</template>

<style>

</style>
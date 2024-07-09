  
<script>
    export default {
    data(){
        return {
            message: "",
            showMessage: false,
            timeout: null,
            showMessageBad: false,
        }
    },

    methods: {
        triggerFileInput() {
            this.$refs.fileInput.click();
        },
        async uploadFile(){
            let fileInput = document.getElementById('fileUploader')
            const file = fileInput.files[0];
            const reader = new FileReader();
            const allowedExtensions = ["image/png", "image/jpeg", "image/jpg"];
            
            if (!allowedExtensions.includes(file.type)) {
                this.message = "File type not allowed. Please upload a .png, .jpeg, or .jpg file.";
                this.showMessageBad = true;
                setTimeout(() => {
                this.showMessageBad = false;
                }, 3000);
                return;
            }

            reader.onload = async () => {
                if (typeof reader.result === 'string') {
                    const base64String = reader.result.split(',')[1]; 
                    const payload = {
                        username_owner: localStorage.getItem('username'),
                        content: base64String
                    };

                    let response = await this.$axios.post("/photo", payload, {
                        headers: {
                            'Content-Type': 'application/json'
                        },
                    });
                    this.message = "Photo uploaded successfully!"
                    this.showMessage = true

                    if (this.timeout) {
                        clearTimeout(this.timeout)  
                    }

                    this.timeout = setTimeout(() => {
                        this.showMessage = false;
                        }, 5000);

                } else {
                    console.error("Tipo di dati non gestito:", typeof reader.result);
                    this.message = "An error occured during upload"
                    this.showMessageBad = true

                    if (this.timeout) {
                        clearTimeout(this.timeout)  
                    }

                    this.timeout = setTimeout(() => {
                        this.showMessageBad = false;
                        }, 5000);
                }
            };

            // Scegli il metodo di lettura appropriato
            reader.readAsDataURL(file);
        },

    },
}

</script>

<template>
    <div id="upload-container">
        <div id="mydiv">
            <input type="file" ref="fileInput" @change="uploadFile" id="fileUploader" accept=".png, .jpeg, .jpg"/>
            <button @click="triggerFileInput">Choose File</button>
        </div>
        <div class="upload-message" v-if="showMessage">
            {{ message }}
        </div>
        <div class="unable-message" v-if="showMessageBad">
            {{ message }}
        </div>
    </div>
</template>

<style scoped>
    .upload-message {
        position: fixed;
        top: 20px;
        left: 50%;
        transform: translateX(-50%);
        background-color: rgb(99, 196, 99);
        color: white;
        padding: 10px 20px;
        border-radius: 5px;
        text-align: center;
        opacity: 1;
    }

    .unable-message {
        position: fixed;
        top: 20px;
        left: 50%;
        transform: translateX(-50%);
        background-color: rgb(230, 18, 18);
        color: white;
        padding: 10px 20px;
        border-radius: 5px;
        text-align: center;
        opacity: 1;
    }

    #upload-container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
    }

    #mydiv {
        display: flex;
        width: 600px;
        height: 300px;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        border: dashed 5px #1ebbd7;
        border-radius: 50px;
        padding: 20px;
        text-align: center;
    }

    input[type="file"] {
        display: none;
        margin: 100px auto 0;
    }

    button {
        margin: auto 0;
        padding: 20px 40px;
        border: 2px solid #1ebbd7;
        background-color: white;
        color: #1ebbd7;
        font-size: 16px;
        cursor: pointer;
        border-radius: 20px;
        transition: background-color 0.3s, color 0.3s;
}

    button:hover {
        background-color: #1ebbd7;
        color: white;
    }

</style>
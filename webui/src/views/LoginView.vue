<template>
    <div class="login-container">
      <h2>Login a WasaPhoto</h2>
      <form @submit.prevent="handleSubmit">
        <input type="text" v-model="username" placeholder="Nome utente" required />
        <input type="password" v-model="password" placeholder="Password" required />
        <input type="submit" value="Login" />
      </form>
    </div>
</template>
  
<script>
  export default {
    data() {
      return {
        username: '',
        password: ''
      };
    },
    methods: {
      handleSubmit() {
        // Aggiungi qui la logica per gestire il login
        console.log('Username:', this.username);
        console.log('Password:', this.password);
      }
    }
  };
</script>
  
<style scoped>
  body {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background-color: #f0f0f0;
    margin: 0;
    font-family: Arial, sans-serif;
  }
  .login-container {
    background-color: white;
    border-radius: 15px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    padding: 30px;
    max-width: 400px;
    width: 100%;
    text-align: center;
  }
  .login-container h2 {
    margin-bottom: 20px;
  }
  .login-container form {
    display: flex;
    flex-direction: column;
  }
  .login-container input[type="text"],
  .login-container input[type="password"] {
    padding: 10px;
    margin-bottom: 15px;
    border: 1px solid #ccc;
    border-radius: 5px;
  }
  .login-container input[type="submit"] {
    padding: 10px;
    border: none;
    border-radius: 5px;
    background-color: #333;
    color: white;
    cursor: pointer;
  }
  .login-container input[type="submit"]:hover {
    background-color: #555;
  }
</style>
  
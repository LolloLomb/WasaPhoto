import {createApp, reactive} from 'vue'

import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import Navbar from './components/Navbar.vue'
import Banner from './components/Banner.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Photo from './components/Photo.vue'
// importare tutti i componenti nuovi


import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("Navbar", Navbar);
app.component("Banner", Banner)
app.component("LoadingSpinner", LoadingSpinner);
app.component('Photo', Photo)
app.use(router)
app.mount('#app')

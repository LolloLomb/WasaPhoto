import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import SearchView from '../views/SearchView.vue'
import ProfileView from '../views/ProfileView.vue'
import UploadView from '../views/UploadView.vue'
/* nb importa tutti i comp che crei */

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      	path: '/',
	    redirect: '/login'
    },
    {
      	path: '/login',
      	component: LoginView
    },
	{
		path: '/home',
		component: HomeView
	  },
	{
		path: '/upload',
		component: UploadView
	},
	{
		path: '/search',
		component: SearchView
	},
	{
		path: '/profile/:username',
		component: ProfileView
	},
	{
		path: "/:catchAll(.*)",
		component: HomeView
	}
  ]
})

export default router

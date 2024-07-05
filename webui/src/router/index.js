import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import SearchView from '../views/SearchView.vue'
/* Assicurati di importare correttamente i componenti delle view */

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
		path: '/search',
		component: SearchView
	},
	{
		path: "/:catchAll(.*)",
		component: HomeView
	}
  ]
})

export default router

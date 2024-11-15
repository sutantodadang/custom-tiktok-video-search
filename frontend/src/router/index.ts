import { createRouter, createWebHistory } from 'vue-router'
// import AppView from '../App.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/video/:id',
      name: 'video',
      props: true,

      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/VideoDetailView.vue'),
    },
  ],
})



export default router

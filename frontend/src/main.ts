import './index.css'
// import VueVideoPlayer from '@videojs-player/vue'
import 'video.js/dist/video-js.css'
// import Vue from 'vue'
import VueCookies from 'vue-cookies'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

const app = createApp(App)


app.use(VueCookies, {
    domain: 'www.tiktok.com',
    path: '/',
    secure: true,
    sameSite: 'None',
    expires: 'Session',
})

app.use(createPinia())
app.use(router)

app.mount('#app')

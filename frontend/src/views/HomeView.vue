<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8">
    <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-center mb-4 sm:mb-6 lg:mb-8">Custom TikTok Video Search
    </h1>

    <!-- Search Bar - Responsive width -->
    <div class="max-w-md mx-auto">
      <label for="default-search" class="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white">Search</label>
      <div class="relative">
        <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
          <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true" xmlns="http://www.w3.org/2000/svg"
            fill="none" viewBox="0 0 20 20">
            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z" />
          </svg>
        </div>
        <input v-model="keyword" type="search" id="default-search"
          class="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Search video tiktok you like..." required />
        <button @click="handleSearch" type="button"
          class="text-white absolute end-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Search</button>
      </div>
    </div>

    <!-- Skeleton Loading - Responsive grid -->
    <div v-if="videoStore.loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
      <div v-for="n in 6" :key="n" class="bg-gray-100 p-3 sm:p-4 rounded animate-pulse">
        <div class="flex flex-col">
          <!-- Responsive height for skeleton -->
          <div class="bg-gray-200 h-[400px] sm:h-[450px] lg:h-[500px] rounded-t-xl"></div>
          <div class="p-3 sm:p-4">
            <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div class="h-4 bg-gray-200 rounded w-1/2"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Video Grid - Responsive -->
    <div v-if="videoStore.hasVideos" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
      ref="videoContainer">
      <div v-for="(video, index) in videoStore.videos" :key="index"
        class="p-3 sm:p-4 rounded transition-transform duration-500 ease-in-out hover:scale-[1.02] focus:scale-[1.02]">
        <a @click="navigateToVideo(video.Item.video.id)"
          class="flex flex-col bg-white border shadow-sm rounded-xl overflow-hidden hover:shadow-lg dark:bg-neutral-900 dark:border-neutral-700 dark:shadow-neutral-700/70">
          <!-- Video Container with responsive height -->
          <div class="relative rounded-t-xl overflow-hidden">
            <iframe class="w-full h-[400px] sm:h-[450px] lg:h-[500px] object-cover"
              :src="`https://www.tiktok.com/player/v1/${video.Item.video.id}?&music_info=1&description=1&volume_control=1`"
              title="TikTok video player"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowfullscreen>
            </iframe>
          </div>

          <!-- Video Info Section -->
          <div class="p-3 sm:p-4">
            <h3 class="text-base sm:text-lg font-bold text-gray-800 dark:text-white mb-2">
              <p class="line-clamp-2">{{ video.Item.desc }}</p>
            </h3>
            <a :href="`https://www.tiktok.com/@${video.Item.author.uniqueId}`" target="_blank"
              class="text-sm sm:text-base text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
              @{{ video.Item.author.uniqueId }}
            </a>
          </div>
        </a>
      </div>
    </div>

    <!-- Loading More Indicator - Centered -->
    <div v-if="loadingMore" class="flex justify-center mt-6">
      <div
        class="animate-spin inline-block size-6 border-[3px] border-current border-t-transparent text-blue-600 rounded-full dark:text-blue-500"
        role="status" aria-label="loading">
        <span class="sr-only">Loading more...</span>
      </div>
    </div>

    <!-- No Results Message -->
    <div v-if="!videoStore.loading && !videoStore.hasVideos"
      class="text-center text-gray-600 dark:text-gray-400 mt-6 text-base sm:text-lg">
      No videos found.
    </div>

    <!-- Scroll to Top Button -->
    <button v-if="showScrollTop" @click="scrollToTop"
      class="fixed bottom-6 right-6 p-3 bg-blue-600 text-white rounded-full shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-opacity duration-300"
      aria-label="Scroll to top">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
      </svg>
    </button>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, onUnmounted, computed, defineComponent } from 'vue'
import { useVideoStore } from '../stores/videos'
import { useRouter } from 'vue-router'
import type { Ref } from 'vue'

export default defineComponent({
  name: 'VideoList',
  setup() {
    const videoStore = useVideoStore()
    const router = useRouter()

    // Refs
    const keyword: Ref<string> = ref('')
    const page: Ref<number> = ref(0)
    const loadingMore: Ref<boolean> = ref(false)
    const videoContainer: Ref<HTMLElement | null> = ref(null)
    const showScrollTop: Ref<boolean> = ref(false)

    // Throttle function
    const throttle = (fn: Function, delay: number) => {
      let lastCall = 0
      return (...args: any[]) => {
        const now = Date.now()
        if (now - lastCall >= delay) {
          fn(...args)
          lastCall = now
        }
      }
    }

    // Check scroll position
    const checkScrollPosition = () => {
      showScrollTop.value = window.scrollY > 500
    }

    // Infinite scroll handler
    const handleScroll = async () => {
      if (!videoContainer.value) return

      const bottomOfWindow = Math.ceil(window.innerHeight + window.scrollY) >=
        document.documentElement.scrollHeight - 1000

      if (bottomOfWindow && !loadingMore.value && videoStore.hasVideos && videoStore.hasMore) {
        loadingMore.value = true
        page.value += 10
        await videoStore.loadMoreVideos(keyword.value, page.value)
        loadingMore.value = false
      }

      checkScrollPosition()
    }

    // Scroll to top function
    const scrollToTop = () => {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    }

    // Search handler
    const handleSearch = async () => {
      if (!keyword.value.trim()) return

      page.value = 0
      await videoStore.searchVideos(keyword.value, page.value)
    }

    // Navigation handler
    const navigateToVideo = (videoId: string) => {
      videoStore.setCurrentVideo(videoId)
      router.push(`/video/${videoId}`)
    }

    // Throttled scroll handler
    const throttledScrollHandler = throttle(handleScroll, 200)

    // Lifecycle hooks
    onMounted(() => {
      window.addEventListener('scroll', throttledScrollHandler)
    })

    onUnmounted(() => {
      window.removeEventListener('scroll', throttledScrollHandler)
    })

    return {
      videoStore,
      keyword,
      page,
      loadingMore,
      videoContainer,
      showScrollTop,
      scrollToTop,
      handleSearch,
      navigateToVideo
    }
  }
})
</script>

<style scoped>
/* Optional: Smooth transitions for dark mode */
@media (prefers-color-scheme: dark) {
  .dark\:bg-neutral-900 {
    transition: background-color 0.3s ease;
  }
}

/* Hide scrollbar for cleaner look */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
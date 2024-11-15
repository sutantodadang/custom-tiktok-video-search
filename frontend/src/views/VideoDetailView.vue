<template>
  <div class="container mx-auto p-4 sm:p-6 lg:p-8">
    <!-- Back button - responsive padding and width -->
    <button @click="$router.push('/')"
      class="w-full sm:w-auto mb-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
      Back to Search
    </button>

    <!-- Video container - responsive width and padding -->
    <div v-if="video" class="max-w-3xl mx-auto bg-white rounded-xl shadow-lg overflow-hidden">
      <!-- Video player container with responsive aspect ratio -->
      <div class="relative pt-[177.77%] sm:pt-[150%] md:pt-[120%] lg:pt-[100%]">
        <iframe class="absolute top-0 left-0 w-full h-full"
          :src="`https://www.tiktok.com/player/v1/${video.Item.video.id}?&music_info=1&description=1&volume_control=1`"
          allowfullscreen scrolling="no"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share">
        </iframe>
      </div>

      <!-- Content section - responsive padding and typography -->
      <div class="p-4 sm:p-6">
        <!-- Title - responsive font size -->
        <h1 class="text-xl sm:text-2xl font-bold mb-4 break-words">{{ video.Item.desc }}</h1>

        <div class="mt-4 space-y-2">
          <!-- Author link - responsive styling -->
          <a :href="`https://www.tiktok.com/@${video.Item.author.uniqueId}`" target="_blank"
            class="text-blue-500 hover:text-blue-600 text-sm sm:text-base inline-block">
            @{{ video.Item.author.uniqueId }}
          </a>

          <!-- Stats container - responsive layout -->
          <div class="flex flex-wrap gap-3 sm:gap-5">
            <!-- Like count -->
            <p class="text-gray-600 flex items-center text-sm sm:text-base">
              <svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-800 dark:text-white mr-1" aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                <path fill-rule="evenodd"
                  d="M15.03 9.684h3.965c.322 0 .64.08.925.232.286.153.532.374.717.645a2.109 2.109 0 0 1 .242 1.883l-2.36 7.201c-.288.814-.48 1.355-1.884 1.355-2.072 0-4.276-.677-6.157-1.256-.472-.145-.924-.284-1.348-.404h-.115V9.478a25.485 25.485 0 0 0 4.238-5.514 1.8 1.8 0 0 1 .901-.83 1.74 1.74 0 0 1 1.21-.048c.396.13.736.397.96.757.225.36.32.788.269 1.211l-1.562 4.63ZM4.177 10H7v8a2 2 0 1 1-4 0v-6.823C3 10.527 3.527 10 4.176 10Z"
                  clip-rule="evenodd" />
              </svg>
              <span class="ml-1">{{ formatNumber(video.Item.stats.diggCount) }}</span>
            </p>

            <!-- Comment count -->
            <p class="text-gray-600 flex items-center text-sm sm:text-base">
              <svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-800 dark:text-white mr-1" aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                <path fill-rule="evenodd"
                  d="M3 6a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-6.616l-2.88 2.592C8.537 20.461 7 19.776 7 18.477V17H5a2 2 0 0 1-2-2V6Zm4 2a1 1 0 0 0 0 2h5a1 1 0 1 0 0-2H7Zm8 0a1 1 0 1 0 0 2h2a1 1 0 1 0 0-2h-2Zm-8 3a1 1 0 1 0 0 2h2a1 1 0 1 0 0-2H7Zm5 0a1 1 0 1 0 0 2h5a1 1 0 1 0 0-2h-5Z"
                  clip-rule="evenodd" />
              </svg>
              <span class="ml-1">{{ formatNumber(video.Item.stats.commentCount) }}</span>
            </p>

            <!-- Share count -->
            <p class="text-gray-600 flex items-center text-sm sm:text-base">
              <svg class="w-5 h-5 sm:w-6 sm:h-6 text-gray-800 dark:text-white mr-1" aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                <path fill-rule="evenodd"
                  d="M14.516 6.743c-.41-.368-.443-1-.077-1.41a.99.99 0 0 1 1.405-.078l5.487 4.948.007.006A2.047 2.047 0 0 1 22 11.721a2.06 2.06 0 0 1-.662 1.51l-5.584 5.09a.99.99 0 0 1-1.404-.07 1.003 1.003 0 0 1 .068-1.412l5.578-5.082a.05.05 0 0 0 .015-.036.051.051 0 0 0-.015-.036l-5.48-4.942Zm-6.543 9.199v-.42a4.168 4.168 0 0 0-2.715 2.415c-.154.382-.44.695-.806.88a1.683 1.683 0 0 1-2.167-.571 1.705 1.705 0 0 1-.279-1.092V15.88c0-3.77 2.526-7.039 5.967-7.573V7.57a1.957 1.957 0 0 1 .993-1.838 1.931 1.931 0 0 1 2.153.184l5.08 4.248a.646.646 0 0 1 .012.011l.011.01a2.098 2.098 0 0 1 .703 1.57 2.108 2.108 0 0 1-.726 1.59l-5.08 4.25a1.933 1.933 0 0 1-2.929-.614 1.957 1.957 0 0 1-.217-1.04Z"
                  clip-rule="evenodd" />
              </svg>
              <span class="ml-1">{{ formatNumber(video.Item.stats.shareCount) }}</span>
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Error message - responsive text size -->
    <div v-if="check" class="text-center text-sm sm:text-base text-gray-600 mt-4">
      Video not found
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useVideoStore } from '../stores/videos'

export default defineComponent({
  name: 'VideoDetail',
  props: {
    id: {
      type: String,
      required: true
    }
  },

  setup(props) {
    const router = useRouter()
    const videoStore = useVideoStore()

    const video = computed(() => {
      const data = videoStore.getVideoById(props.id)
      return data
    })

    const check = computed(() => {
      return videoStore.hasVideos
    })

    // Format numbers for better readability
    const formatNumber = (num: number): string => {
      if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M'
      }
      if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K'
      }
      return num.toString()
    }

    if (!video.value) {
      router.push('/')
    }

    return {
      video,
      check,
      videoStore,
      formatNumber
    }
  }
})
</script>
<template>
    <div class="video-container">
        <!-- Video Player -->
        <iframe ref="videoFrame"
            :src="`https://www.tiktok.com/player/v1/${videoId}?&music_info=1&description=1&volume_control=1`"
            class="w-full aspect-video rounded-lg shadow-lg" allow="encrypted-media;" @load="setupPlayer"></iframe>

        <!-- Volume Controls -->
        <div class="flex items-center space-x-2 mt-4 p-2 bg-gray-50 rounded-lg">
            <!-- Mute/Unmute Button -->
            <button @click="toggleMute" type="button"
                class="text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-4 focus:ring-gray-200 rounded-lg p-2.5 inline-flex items-center justify-center"
                :class="{ 'bg-gray-100': isMuted }">
                <!-- Volume High Icon -->
                <svg v-if="!isMuted" class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                    viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 6h3l7-5v22l-7-5H4a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1Zm14 6a3 3 0 0 1 0 6m2-9a6 6 0 0 1 0 12m2-15c5 1 8.6 5.6 8 11-.6 5.3-5.2 9.4-10.5 9l.5-2c4 .4 7.6-2.9 8-6.9.4-4-2.4-7.7-6.3-8.1l.3-3Z" />
                </svg>
                <!-- Volume Muted Icon -->
                <svg v-else class="w-5 h-5" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="none"
                    viewBox="0 0 24 24">
                    <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M4 6h3l7-5v22l-7-5H4a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1Zm14.7 1.7L14 12.4m0-4.8 4.7 4.8" />
                </svg>
                <span class="sr-only">Volume control</span>
            </button>

            <!-- Volume Slider -->
            <div class="relative flex-grow">
                <input type="range" min="0" max="100" :value="volume" @input="updateVolume"
                    class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700">
            </div>

            <!-- Volume Percentage -->
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300 min-w-[3rem] text-center">
                {{ volume }}%
            </span>
        </div>

        <!-- Current State (optional, for debugging) -->
        <div v-if="debug" class="mt-2 p-2 bg-gray-50 rounded text-sm text-gray-500">
            <p>Player State: {{ playerState }}</p>
            <p>Current Time: {{ currentTime }}s</p>
        </div>
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'

export default {
    name: 'TikTokPlayer',

    props: {
        videoId: {
            type: String,
            required: true
        },
        debug: {
            type: Boolean,
            default: false
        }
    },

    setup() {
        const videoFrame = ref(null)
        const player = ref(null)
        const volume = ref(50)
        const isMuted = ref(false)
        const playerState = ref('init')
        const currentTime = ref(0)
        let messageHandler = null

        // Handle messages from the iframe
        const handleMessage = (event) => {
            const { type, value } = event.data

            switch (type) {
                case 'onPlayerReady':
                    player.value = {
                        postMessage: (msg) => {
                            videoFrame.value.contentWindow.postMessage(msg, '*')
                        }
                    }
                    setVolume(volume.value)
                    break

                case 'onVolumeChange':
                    volume.value = value
                    isMuted.value = value === 0
                    break

                case 'onStateChange':
                    const states = ['init', 'ended', 'playing', 'paused', 'buffering']
                    playerState.value = states[value + 1] || 'unknown'
                    break

                case 'onCurrentTime':
                    currentTime.value = Math.round(value.currentTime)
                    break
            }
        }

        // Set volume through postMessage
        const setVolume = (value) => {
            if (player.value) {
                player.value.postMessage({
                    method: 'setVolume',
                    value: value
                })
            }
        }

        // Update volume from slider
        const updateVolume = (event) => {
            const newVolume = parseInt(event.target.value)
            volume.value = newVolume
            setVolume(newVolume)
            isMuted.value = newVolume === 0
        }

        // Toggle mute state
        const toggleMute = () => {
            const newVolume = isMuted.value ? 50 : 0
            volume.value = newVolume
            setVolume(newVolume)
            isMuted.value = !isMuted.value
        }

        // Setup player and message listener
        const setupPlayer = () => {
            messageHandler = (event) => handleMessage(event)
            window.addEventListener('message', messageHandler)
        }

        onUnmounted(() => {
            if (messageHandler) {
                window.removeEventListener('message', messageHandler)
            }
        })

        return {
            videoFrame,
            volume,
            isMuted,
            playerState,
            currentTime,
            setupPlayer,
            updateVolume,
            toggleMute
        }
    }
}
</script>

<style scoped>
/* Custom range slider styling */
input[type="range"]::-webkit-slider-thumb {
    @apply w-4 h-4 bg-blue-600 rounded-full shadow appearance-none cursor-pointer;
    -webkit-appearance: none;
    margin-top: -4px;
}

input[type="range"]::-moz-range-thumb {
    @apply w-4 h-4 bg-blue-600 rounded-full shadow cursor-pointer border-0;
}

input[type="range"]::-webkit-slider-runnable-track {
    @apply h-2 bg-gray-200 rounded-lg dark:bg-gray-700;
}

input[type="range"]::-moz-range-track {
    @apply h-2 bg-gray-200 rounded-lg dark:bg-gray-700;
}

/* Focus styles */
input[type="range"]:focus {
    @apply outline-none;
}

input[type="range"]:focus::-webkit-slider-thumb {
    @apply ring-4 ring-blue-300;
}

input[type="range"]:focus::-moz-range-thumb {
    @apply ring-4 ring-blue-300;
}
</style>
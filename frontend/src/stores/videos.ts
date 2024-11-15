import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Video {
    Item: {
        video: {
            id: string
        }
        desc: string
        author: {
            uniqueId: string
        }
        stats: {
            diggCount: number
            commentCount: number
            shareCount: number
        }
    }
}

export const useVideoStore = defineStore('videos', () => {
    const videos = ref<Video[]>([])
    const loading = ref(false)
    const currentVideo = ref<string | null>(null)
    const hasMore = ref(true)

    const searchVideos = async (keyword: string, offset: number) => {
        try {
            loading.value = true
            // Reset the videos array when starting a new search
            videos.value = []

            keyword = keyword.split(" ").join("-")

            const response = await fetch(`http://localhost:5555/api/videos?keyword=${keyword}&offset=${offset}`)
            const data = await response.json()

            console.log(data)

            videos.value = data.Data
            hasMore.value = data.hasMore // Assuming your API returns this flag
        } catch (error) {
            console.error('Error searching videos:', error)
        } finally {
            loading.value = false
        }
    }

    const loadMoreVideos = async (keyword: string, page: number) => {
        if (!hasMore.value) return

        try {
            const response = await fetch(`/api/search?keyword=${keyword}&offset=${page}`)
            const data = await response.json()

            videos.value = [...videos.value, ...data.videos]
            hasMore.value = data.hasMore
        } catch (error) {
            console.error('Error loading more videos:', error)
        }
    }

    const getVideoById = (id: string) => {

        return videos.value.find((x) => x.Item.video.id === id)

    }

    const setCurrentVideo = (videoId: string) => {
        currentVideo.value = videoId
    }

    const hasVideos = computed(() => videos.value.length > 0)

    return {
        videos,
        loading,
        currentVideo,
        hasMore,
        hasVideos,
        searchVideos,
        loadMoreVideos,
        setCurrentVideo,
        getVideoById
    }
})
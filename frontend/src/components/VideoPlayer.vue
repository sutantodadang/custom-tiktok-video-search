<!-- src/components/VideoPlayer.vue -->
<template>
    <div>
        <videojs-player class="video-js vjs-default-skin" :options="playerOptions" @ready="onPlayerReady"
            @play="onPlay" />
    </div>
</template>

<script>
import { VideoJsPlayer } from '@videojs-player/vue';

export default {
    name: 'VideoPlayer',
    components: {
        VideoJsPlayer
    },
    props: {
        src: {
            type: String,
            required: true
        }
    },
    computed: {
        playerOptions() {
            return {
                controls: true,
                autoplay: false,
                preload: 'auto',
                sources: [
                    {
                        src: this.src,
                        type: 'video/mp4'
                    }
                ]
            };
        }
    },
    methods: {
        onPlayerReady(player) {
            console.log('Video.js player is ready!', player);
        },
        onPlay() {
            console.log('Video is playing!');
        }
    }
};
</script>

<style>
.video-js {
    width: 100%;
    height: 400px;
}
</style>
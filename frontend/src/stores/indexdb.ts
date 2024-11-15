
import { openDB } from 'idb';

export interface VideoData {
    Item: {
        video: {
            id: string;
        };
        desc: string;
        author: {
            uniqueId: string;
        };
    };
}

const dbPromise = openDB('video-store', 1, {
    upgrade(db) {
        db.createObjectStore('videos');
    },
});

export async function saveVideo(id: string, videoData: VideoData): Promise<void> {
    const db = await dbPromise;
    await db.put('videos', videoData, id);
}

export async function getVideo(id: string): Promise<VideoData | undefined> {
    const db = await dbPromise;
    return await db.get('videos', id);
}
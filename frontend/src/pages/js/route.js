import {createRouter, createWebHistory} from 'vue-router';
import mergeVideo from '../mergeVideo.vue';
import videoClip from '../videoClip.vue';
import videoBrowse from '../videoBrowse.vue';

const routes = [
    {path: '/', component: videoClip},
    {path: '/mergeVideo', component: mergeVideo},
    {path: '/videoClip', component: videoClip},
    {path: '/videoBrowse', component: videoBrowse},
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;

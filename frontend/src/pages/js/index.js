import { createApp } from 'vue'
import App from '../index.vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import routes from './route.js';

let app = createApp(App)
app.use(ElementPlus)
app.use(routes);
app.mount('#app')

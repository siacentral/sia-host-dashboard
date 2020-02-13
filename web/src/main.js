import Vue from 'vue';
import App from './App.vue';
import { library } from '@fortawesome/fontawesome-svg-core';
import { faFile, faFileExport, faUnlock, faLock, faEllipsisV, faChevronLeft, faChevronRight, faEye, faCogs, faPlus, faTimes, faRedo, faWifi } from '@fortawesome/free-solid-svg-icons';
import { faUsb, faGithub } from '@fortawesome/free-brands-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

import store from '@/store';

library.add(faFile, faFileExport, faUnlock, faLock, faEllipsisV, faChevronLeft, faChevronRight, faWifi, faEye, faUsb, faGithub, faCogs, faPlus, faTimes, faRedo);

Vue.component('icon', FontAwesomeIcon);

Vue.config.productionTip = false;

new Vue({
	store,
	render: h => h(App)
}).$mount('#app');

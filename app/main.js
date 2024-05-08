import Vue from 'vue';
import VueRouter from 'vue-router';
/*import Vuex from 'vuex';
import {sync} from 'vuex-router-sync';
import {routerHistory, writeHistory} from 'vue-router-back-button';*/


import VueResource from 'vue-resource';
import VueLazyload from 'vue-lazyload';

//import VueSimpleMarkdown from 'vue-simple-markdown';
var VueSimpleMarkdown = require('vue-simple-markdown');

window.moment = require('moment');
window.utils = {
	title(name) {
		if(name) {
			document.title = name + ' - Openverse\u00b2'
		} else {
			document.title = 'Openverse\u00b2'
		}
	},
	properNowStamp() {
		// moment hates timezones so we're using hacks to make it proper
		return moment(window.moment().utc().format('YYYY-MM-DDTHH:mm:ss+04:00'));
	},
	properNow() {
		// moment hates timezones so we're using hacks to make it proper
		return moment(window.moment().local().format('YYYY-MM-DDTHH:mm:ss+00:00'));
	},
	doTime(time) {
		// moment hates timezones so we're using hacks to make it proper
		var time = moment(time);
		if(moment() - time._d < 156084000) {
			return moment(moment(time).utc().format('YYYY-MM-DDTHH:mm:ss-04:00')).fromNow();
		} else {
			return time.utc().format('MM/DD/YYYY hh:mm A');
		}
	},
	doAvatar(avatar) {
		if(!avatar) {
			return '/static/img/anonymous-mii.png';
		}
		if(avatar.indexOf('http') != 0) {
			if(avatar.length == 16 && parseInt(avatar, 16)) {
				return 'https://cdn-mii.accounts.nintendo.com/2.0.0/mii_images/' + avatar + '/00000000000000000000000000000000.png?type=face&width=96';
			}
			return 'https://mii-secure.cdn.nintendo.net/' + avatar + '_normal_face.png';
		}
		return avatar;
	},
	doAvatarFeeling(avatar, feelingId) {
		if(!avatar) {
			return '/static/img/anonymous-mii.png';
		}
		if(avatar.indexOf('http') != 0) {
			// If the avatar is a 16-length hex string, or a cdn-mii Mii, then do this instead
			if(avatar.length == 16 && parseInt(avatar, 16)) {
				var feeling = 'normal';
				switch(feelingId) {
					case 1:
						feeling = 'smile_open_mouth';
						break;
					case 2:
						feeling = 'like_wink_left';
						break;
					case 3:
						feeling = 'surprise_open_mouth';
						break;
					case 4:
						feeling = 'frustrated';
						break;
					case 5:
						feeling = 'sorrow';
						break;
				}
				return 'https://cdn-mii.accounts.nintendo.com/2.0.0/mii_images/' + avatar + '/00000000000000000000000000000000.png?type=face&expression=' + feeling + '&width=96';
			}
			var feeling = '_normal';
			switch(feelingId) {
				case 1:
					feeling = '_happy';
					break;
				case 2:
					feeling = '_like';
					break;
				case 3:
					feeling = '_surprised';
					break;
				case 4:
					feeling = '_frustrated';
					break;
				case 5:
					feeling = '_puzzled';
					break;
			}
			return 'https://mii-secure.cdn.nintendo.net/' + avatar + feeling + '_face.png';
		}
		return avatar;
	},
	doRank(rank) {
		switch(rank) {
			case 1:
				return 'donator';
				break;
			case 2:
				return 'tester';
				break;
			case 3:
				return 'moderator';
				break;
			case 4:
				return 'administrator';
				break;
			case 5:
				return 'developer';
				break;
			default:
				return '';
		}
	},
	doRankText(rank) {
		switch(rank) {
			case 1:
				return 'Donator';
				break;
			case 2:
				return 'Tester';
				break;
			case 3:
				return 'Moderator';
				break;
			case 4:
				return 'Administrator';
				break;
			case 5:
				return 'Developer';
				break;
			default:
				return '';
		}
	},
	yeahButton(feeling, hasYeah) {
		if(hasYeah) {
			return 'Unyeah';
		}
		switch(feeling) {
			case 2:
				return 'Yeah\u2665';
				break;
			case 3:
				return 'Yeah!?';
				break;
			case 4:
				return 'Yeah...';
				break;
			case 5:
				return 'Yeah...';
				break;
			default:
				return 'Yeah!';
		}
	},
	// CSS doesn't let us have multi-line ellipsising so we need this
	textTrun(text, max) {
		if(!text) {
			return text;
		}
		if(text.length > max) {
			return text.substring(0, max) + '...';
		} else {
			return text;
		}
	},
	isAudio(url) {
		return url.indexOf('.m4a') >= 0 || url.indexOf('.bfstm') >= 0;
	}
}

// layout stuff
import Layout from './components/layout.vue';
import NotFound from './components/not_found.vue';
// auth stuff
import Login from './components/login.vue';
		// TODO: make logout better
import Logout from './components/logout.vue';
import Signup from './components/signup.vue';
// community stuff
import CommunitiesList from './components/communities_list.vue';
import CommunityPosts from './components/community_posts.vue';
import PostView from './components/post_view.vue';
import CommentView from './components/comment_view.vue';
import NotificationsView from './components/notifications_view.vue';
import ProfileSettings from './components/profile_settings.vue';
import UserView from './components/user_view.vue';
import UserSidebar from './components/user_sidebar.vue';
Vue.component('user-sidebar', UserSidebar);

Vue.use(VueRouter);
Vue.use(VueResource);
//Vue.use(Vuex);
/*Vue.use(VueLazyload, {
	preLoad: 1.0,
	attempt: 1
});*/
VueSimpleMarkdown.VueSimpleMarkdown.props.heading = {type: Boolean, default: false};
VueSimpleMarkdown.VueSimpleMarkdown.props.horizontalLine = {type: Boolean, default: false};
VueSimpleMarkdown.VueSimpleMarkdown.props.image = {type: Boolean, default: false};
VueSimpleMarkdown.VueSimpleMarkdown.props.lists = {type: Boolean, default: false};
Vue.component('markdown', VueSimpleMarkdown.VueSimpleMarkdown);

Vue.http.interceptors.push((request, next) => {
	if(request.method == 'POST') {
		if(request.body == undefined) {
			request.body = {};
		}
		request.body._csrf = window.csrf;
	}
	next();
});
Vue.http.options.emulateJSON = true;

const router = new VueRouter({
	mode: 'history',
	scrollBehavior(to, from, savedPosition) {
		if(savedPosition) {
			return savedPosition;
  	} else {
    	return {x: 0, y: 0};
  	}
	},
	routes: [
		{path: '/', component: CommunitiesList, name: 'communities_list'},
		{path: '/communities/:id', component: CommunityPosts, name: 'community_posts'},
		{path: '/titles/*/:id', component: {
			created() {this.$router.push('/communities/' + this.$route.params.id)}, template: ''},
			name: 'title_view_redirect'},
		{path: '/posts/:id', component: PostView, name: 'post_view'},
		{path: '/replies/:id', component: CommentView, name: 'comment_view'},
		{path: '/news/my_news', component: NotificationsView, name: 'notifications_view'},
		{path: '/settings/profile', component: ProfileSettings, name: 'profile_settings'},
		{path: '/users/:id', component: UserView, name: 'user_view'},
		{path: '/login', component: Login, name: 'login'},
		{path: '/logout', component: Logout, name: 'logout'},
		{path: '/signup', component: Signup, name: 'signup'},
		{path: '*', component: NotFound, name: 'not_found'}
	]
});

/*const store = new Vuex.Store({
  state: {
    count: 0
  },
  mutations: {
    increment(state) {
      state.count++;
    }
  }
});
const unsync = sync(store, router);
*/

const app = new Vue({
  el: '#app',
  router,
  render: h => h(Layout),
	//destroyed: unsync
});

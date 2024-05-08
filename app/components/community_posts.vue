<template>
	<div id="community-top">
		<div class="no-content" v-show="loading">
			<p>Now loading...</p>
		</div>
		<div class="no-content" v-if="!loading && error">
			<p>{{ error }}</p>
		</div>
		<div id="sidebar" v-if="!loading && !error">
			<section class="sidebar-container" id="sidebar-community">
				<span id="sidebar-cover" v-show="community.banner"><router-link :to="'/communities/' + community.id"><img v-bind:src="community.banner"></router-link></span>
				<header id="sidebar-community-body">
					<span id="sidebar-community-img"><span class="icon-container"><router-link :to="'/communities/' + community.id"><img class="icon" v-bind:src="community.icon"></router-link></span></span>
					<h1 class="community-name"><router-link :to="'/communities/' + community.id">{{ community.name }}</router-link></h1>
				</header>
				<div class="community-description" v-show="community.description">
					<p class="text">{{ community.description }}</p>
				</div>
			</section>
		</div>
		<div class="main-column post-list-outline" v-if="!loading && !error">
			<!--<div class="tab-container">
				<div class="tab2">
					<a class="selected">All Posts</a><a>Popular Posts</a>
				</div>
			</div>-->
			<form id="post-form" v-if="loggedIn && (userLevel || 0) >= (community.permission || 0)" v-bind:class="{folded: formFolded}" v-on:submit.prevent="sendForm()">
					<div class="feeling-selector">
						<label class="symbol feeling-button feeling-button-normal" v-bind:class="{checked: feelingId == 0}"><input checked type="radio" value="0" v-model="feelingId"><span class="symbol-label">normal</span></label>
						<label class="symbol feeling-button feeling-button-happy" v-bind:class="{checked: feelingId == 1}"><input type="radio" value="1" v-model="feelingId"><span class="symbol-label">happy</span></label>
						<label class="symbol feeling-button feeling-button-like" v-bind:class="{checked: feelingId == 2}"><input type="radio" value="2" v-model="feelingId"><span class="symbol-label">like</span></label>
						<label class="symbol feeling-button feeling-button-surprised" v-bind:class="{checked: feelingId == 3}"><input type="radio" value="3" v-model="feelingId"><span class="symbol-label">surprised</span></label>
						<label class="symbol feeling-button feeling-button-frustrated" v-bind:class="{checked: feelingId == 4}"><input type="radio" value="4" v-model="feelingId"><span class="symbol-label">frustrated</span></label>
						<label class="symbol feeling-button feeling-button-puzzled" v-bind:class="{checked: feelingId == 5}"><input type="radio" value="5" v-model="feelingId"><span class="symbol-label">puzzled</span></label>
					</div>
				<div class="textarea-with-menu active-text">
					<!--<menu class="textarea-menu">
						<ul>
							<li><label class="textarea-menu-text"><input name="_post_type" type="radio" value="body"></label></li>
							<li><label class="textarea-menu-memo checked" data-modal-open="#memo-drawboard-page"><input name="_post_type" type="radio" value="painting"></label></li>
						</ul>
					</menu>-->
					<div class="textarea-container">
						<textarea class="textarea-text textarea" maxlength="2000" placeholder="Share your thoughts in a post to this community." v-model="body" v-on:click="unfoldForm()" v-on:input="formInput()" v-on:paste="pasteImage($event)"></textarea>
					</div>
				</div>
			<!--
				<input class="textarea-line url-form" maxlength="1024" placeholder="URL" type="text" v-model="url">
			-->
				<label class="file-button-container" v-on:drop.prevent="pasteImage($event)" v-on:dragover.prevent="">
					<span class="input-label">Image <span v-show="!image && !imageDimensions">powered by catgirl.host</span><span v-show="image && !imageDimensions">...</span><span v-show="imageDimensions">{{ imageDimensions }}</span></span>
					<span class="button file-upload-button">Upload</span>
					<input accept="image/*" class="file-button none" type="file" v-on:change="imageUpload($event)">
					<div class="screenshot-container still-image" v-show="image">
						<img v-bind:src="image">
					</div>
				</label>
				<div class="post-form-footer-options">
					<label class="spoiler-button symbol" v-bind:class="{checked: isSpoiler}"><input id="is_spoiler" type="checkbox" value="1" v-model="isSpoiler"> Spoilers</label>
				</div>
				<div class="form-buttons">
					<input class="black-button post-button" type="submit" value="Send" v-bind:class="{disabled: formDisabled}" v-bind:disabled="formDisabled" v-show="!formSending">
					<input class="black-button post-button disabled" disabled type="submit" value="Please wait..." v-show="formSending">
				</div>
			</form>
			<div class="body-content" id="community-post-list">
				<div class="no-content" v-show="posts.length <= 0">
					<p>This community doesn't have any posts yet.</p>
				</div>
				<transition-group class="list post-list" v-show="posts.length > 0" name="fade" tag="ul">
					<li class="post trigger" v-bind:class="{hidden: post.spoiler && post.user_id != userId && !post.spoilerRevealed}" v-bind:key="post.id" v-for="post in posts" v-on:click="goToPost($event, post)">
						<router-link class="icon-container" v-bind:class="utils.doRank(post.user.rank)" :to="'/users/' + post.user.name"><img class="icon" v-bind:src="utils.doAvatarFeeling(post.user.avatar, post.feeling)"></router-link>
						<p class="user-name"><router-link :to="'/users/' + post.user.name">{{ post.user.nick }}</router-link></p>
						<p class="timestamp-container"><router-link class="timestamp" :to="'/posts/' + post.id">{{ post.timestamp }}</router-link>
						<span class="spoiler-status spoiler" v-show="post.spoiler">· Spoilers</span>
						<span class="spoiler-status spoiler" v-show="post.edited && post.edited.Time != post.date.Time && (new Date(post.edited.Time) - new Date(post.date.Time) > 300000)">· Edited</span></p>
						<div class="body post-content">
							<p class="post-content-text" v-if="post.content"><markdown :source="utils.textTrun(post.content, 350)"></markdown></p>
							<router-link class="screenshot-container still-image" :to="'/posts/' + post.id" v-show="post.screenshot"><img v-bind:src="post.screenshot"></router-link>
							<audio v-if="post.url && utils.isAudio(post.url)" v-bind:src="post.url" preload="none" controls></audio>
							<div class="hidden-content" v-show="post.spoiler && post.user_id != userId && !post.spoilerRevealed">
								<p>This post contains spoilers.</p>
								<button type="button" class="hidden-content-button" v-on:click="revealSpoiler(post)">View Post</button>
							</div>
							<div class="post-meta">
								<button type="button" class="symbol submit empathy-button" v-bind:class="{disabled: !post.can_yeah || post.yeahSending, 'empathy-added': post.has_yeah}" v-bind:disabled="!post.can_yeah || post.yeahSending" v-on:click="sendYeah(post)"><span class="empathy-button-text">{{ utils.yeahButton(post.feeling, post.has_yeah) }}</span></button>
								<div class="empathy symbol">
									<span class="empathy-count">{{ post.yeah_count || 0 }}</span>
								</div>
								<div class="reply symbol">
									<span class="reply-count">{{ post.comment_count || 0 }}</span>
								</div>
							</div>
							<div class="recent-reply-content" v-if="post.recent_comment" v-on:click="goToRecentComment($event, post.recent_comment)">
								<div class="recent-reply-read-more-container" v-show="post.comment_count > 1">View {{ post.comment_count > 1 ? 'More Comments' : 'Other Comment' }} ({{ post.comment_count }})</div>
								<div class="recent-reply trigger">
									<router-link class="icon-container" v-bind:class="utils.doRank(post.recent_comment.user.rank)" :to="'/users/' + post.recent_comment.user.name"><img class="icon" v-bind:src="utils.doAvatarFeeling(post.recent_comment.user.avatar, post.recent_comment.feeling)"></router-link>
									<p class="user-name"><router-link :to="'/users/' + post.recent_comment.user.name">{{ post.recent_comment.user.nick }}</router-link></p>
									<p class="timestamp-container"><router-link class="timestamp" :to="'/replies/' + post.recent_comment.id">{{ post.recent_comment.timestamp }}</router-link>
									<span class="spoiler-status spoiler" v-show="post.recent_comment.edited && post.recent_comment.edited.Time != post.recent_comment.date.Time && (new Date(post.recent_comment.edited.Time) - new Date(post.recent_comment.date.Time) > 300000)">· Edited</span></p>
									<div class="body">
										<div class="post-content">
											<p class="recent-reply-content-text" v-if="post.recent_comment.content"><markdown :source="utils.textTrun(post.recent_comment.content, 100)"></markdown></p>
											<router-link class="screenshot-container still-image" :to="'/replies/' + post.recent_comment.id" v-show="post.recent_comment.screenshot"><img v-bind:src="post.recent_comment.screenshot"></router-link>
										</div>
									</div>
								</div>
							</div>
						</div>
					</li>
					<!--<div class="post-list-loading" v-show="scrollLoading" :key="postCount">
						<img src="/static/img/loading-image-blue.gif" alt="Now loading...">
					</div>-->
				</transition-group>
			</div>
		</div>
		<div v-if="dialog">
			<div class="mask"></div>
			<div class="dialog active-dialog">
				<div class="dialog-inner">
					<div class="window">
						<h1 class="window-title">{{ dialog.title }}</h1>
						<div class="window-body">
							<p class="window-body-content">{{ dialog.text }}</p>
							<div class="form-buttons">
								<button class="black-button" type="button" v-on:click="closeDialog()">OK</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
//import saveState from 'vue-save-state';
export default {
	//mixins: [saveState],
	computed: {
		loggedIn() {
			return window.user.id > 0;
		},
		userId() {
			return window.user.id;
		},
		userLevel() {
			if(window.user) {
				return window.user.rank;
			}
		},
		utils() {
			return window.utils;
		}
	},
	created() {
		if(!this.community) {
			var __this = this;
			this.$http.get('/app/communities/' + this.$route.params.id).then(response => {
		    this.$set(this, 'loading', false);
			  if(response.body.community) {
					response.body.posts.forEach(post => {
						// Time functions
						function updateTime() {
							__this.$set(post, 'timestamp', __this.utils.doTime(post.date.Time));
						}
						updateTime();
						if(__this.utils.properNow() - moment(moment(post.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 336084000) {
							setInterval(updateTime, 1000);
						}
						if(post.recent_comment) {
							function updateTime() {
								__this.$set(post.recent_comment, 'timestamp', __this.utils.doTime(post.recent_comment.date.Time));
							}
							updateTime();
							if(__this.utils.properNow() - moment(moment(post.recent_comment.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 336084000) {
								setInterval(updateTime, 1000);
							}
						}
					});
			  	this.$set(this, 'community', response.body.community);
					utils.title(response.body.community.name);
		      this.$set(this, 'posts', response.body.posts);
					this.postCount = response.body.posts.length;
					this.postCountTotal = response.body.posts.length;
					this.timeLoaded = window.moment().format();
					window.addEventListener('scroll', this.onScrollEvent);
			  } else {
					this.error = 'An invalid response was recieved from the server.';
			  }
			}, response => {
				this.$set(this, 'loading', false);
				if(response.body.error) {
		      this.error = response.body.error;
		    } else {
					if(response.status) {
						this.error = 'Error: ' + response.status + ' ' + response.statusText;
					} else {
		      	this.error = 'The server seems to be down right now, try again in a moment.';
					}
		    }
		  });
		}
	},
	mounted() {
		if(!this.loggedIn) {
			return;
		}
		var _this = this;
		this.messageStream = new WebSocket('ws' + ((location.protocol == 'https:') ? 's' : '') + '://' + location.host + '/app/communities/' + this.$route.params.id + '/stream');
		this.messageStream.onmessage = function(event) {
			var type = event.data.slice(0, event.data.indexOf(':'));
			var data = JSON.parse(event.data.substring(event.data.indexOf(':') + 1));
			switch(type) {
				case 'post':
					if(data.user_id != _this.userId) {
						_this.addPost(data);
					}
					break;
				case 'edit':
					if(data.user_id != _this.userId) {
						var post = _this.posts.filter(function(post) {
							return post.id == data.topic;
						})[0];
						_this.$set(post, 'content', data.post.content);
						_this.$set(post, 'feeling', data.post.feeling);
						_this.$set(post, 'spoiler', data.post.spoiler);
						_this.$set(post, 'screenshot', data.post.screenshot);
						_this.$set(post, 'edited', {Time: _this.utils.properNowStamp()});
					}
					break;
				case 'yeah':
					if(data.user_id != _this.userId) {
						var post = _this.posts.filter(function(post) {
							return post.id == data.topic;
						})[0];
						_this.$set(post, 'yeah_count', (post.yeah_count || 0) + 1);
					}
					break;
				case 'unyeah':
					if(data.user_id != _this.userId) {
						var post = _this.posts.filter(function(post) {
							return post.id == data.topic;
						})[0];
						_this.$set(post, 'yeah_count', (post.yeah_count || 0) - 1);
					}
					break;
				case 'comment':
					if(data.user_id != _this.userId) {
						var post = _this.posts.filter(function(post) {
							return post.id == data.topic;
						})[0];
						if(!post) {
							return;
						}
						_this.$set(post, 'comment_count', (post.comment_count || 0) + 1);

						if(data.comment) {
							data.comment.date.Time = _this.utils.properNowStamp();
							// Update time here too
							var __this = _this;
							function updateTime() {
								__this.$set(data.comment, 'timestamp', __this.utils.doTime(data.comment.date.Time));
							}
							updateTime();
							setInterval(updateTime, 1000);

							_this.$set(post, 'recent_comment', data.comment);
						}
					}
					break;
			}
		}
	},
	destroyed() {
		window.removeEventListener('scroll', this.onScrollEvent);
		if(this.messageStream) {
			this.messageStream.close();
		}
	},
	data() {
		return {
			loading: true,
			error: '',
			timeLoaded: '',
			community: null,
			posts: [],
			postCount: 0,
			postCountTotal: 0,
			feelingId: 0,
			body: '',
			isSpoiler: false,
			image: '',
			imageDimensions: '',
			formFolded: true,
			formDisabled: true,
			formSending: false,
			scrollLoading: false,
			dialog: null,
			messageStream: null
		}
	},
	methods: {
		/*getSaveStateConfig() {
			return {
				cacheKey: 'CommunityPosts'
			}
		},*/
		goToPost(event, post) {
			if(!event.target.closest("a, button, .recent-reply, audio")) {
				if(post.spoiler && post.user_id != this.userId && !post.spoilerRevealed) {
					return false;
				}
				this.$router.push('/posts/' + post.id);
			}
		},
		goToRecentComment(event, comment) {
			if(!event.target.closest("a")) {
				this.$router.push('/replies/' + comment.id);
			}
		},
		unfoldForm() {
			this.$set(this, 'formFolded', false);
		},
		formInput() {
			if(this.image) {
				this.$set(this, 'formDisabled', false);
				return;
			}
			if(/^[\s\u00A0\u3000]*$/.test(this.body)) {
				this.$set(this, 'formDisabled', true);
				return;
			}
			this.$set(this, 'formDisabled', false);
		},
		addPost(post) {
			if(this.loggedIn && post.user_id != this.userId) {
        post.can_yeah = true;
      }
			// moment hates timezones so we're using hacks to make it proper
			post.date.Time = this.utils.properNowStamp();
			// Update time here too
			var _this = this;
			function updateTime() {
				_this.$set(post, 'timestamp', _this.utils.doTime(post.date.Time));
			}
			updateTime();
			setInterval(updateTime, 1000);

			if(this.posts.length > 49) {
				this.posts.pop();
			}
			this.posts.unshift(post);
		},
		sendForm() {
			if(this.formDisabled) {
				return false;
			}
			this.$set(this, 'formSending', true);
			this.$http.post('/app/communities/' + this.community.id + '/post', {
				body: this.body,
				feeling_id: this.feelingId,
				is_spoiler: this.isSpoiler,
				screenshot: this.image,
			}).then(response => {
				this.$set(this, 'formSending', false);
				this.addPost(response.body);
				this.resetForm();
			}, response => {
				this.$set(this, 'formSending', false);
				var error;
				if(response.body.error) {
		      error = response.body.error;
		    } else {
		      error = response.status + ' ' + response.statusText;
		    }
				this.$set(this, 'dialog', {
					title: 'Error',
					text: error
				});
			});
		},
		imageUpload(event) {
			if(navigator.userAgent.indexOf('Nintendo') > 0) {
				alert('todo: file uploading is broken on Wii U! make file uploading use FormData() somehow?');
			}
			var reader = new FileReader();
			reader.readAsDataURL(event.target.files[0]);
			var _this = this;
			reader.onload = function() {
				_this.$set(_this, 'image', reader.result);
				_this.$set(_this, 'formDisabled', false);
				var imageDimension = new Image();
				imageDimension.src = reader.result;
				imageDimension.onerror = function() {
					_this.$set(_this, 'image', '');
					_this.$set(_this, 'imageDimensions', 'Invalid image.');
					_this.formInput();
				}
				imageDimension.onload = function() {
					_this.$set(_this, 'imageDimensions', imageDimension.width + ' x ' + imageDimension.height);
				}
			}
		},
		pasteImage(event) {
			var dataSource;
			if(event.clipboardData) {
				dataSource = event.clipboardData;
			} else if(event.dataTransfer) {
				dataSource = event.dataTransfer;
			} else {
				return;
			}
			var _this = this;
			Array.prototype.forEach.call(dataSource.items, function(item) {
				if(item.kind == 'file' && item.type.indexOf('image') >= 0) {
					// found a file
					var reader = new FileReader();
					reader.readAsDataURL(item.getAsFile());
					reader.onload = function() {
						_this.$set(_this, 'image', reader.result);
						_this.$set(_this, 'formDisabled', false);
						var imageDimension = new Image();
						imageDimension.src = reader.result;
						imageDimension.onerror = function() {
							_this.$set(_this, 'image', '');
							_this.$set(_this, 'imageDimensions', 'Invalid image.');
							_this.formInput();
						}
						imageDimension.onload = function() {
							_this.$set(_this, 'imageDimensions', imageDimension.width + ' x ' + imageDimension.height);
						}
					}
				}
			});
		},
		resetForm() {
			this.$set(this, 'feelingId', 0);
			this.$set(this, 'body', '');
			this.$set(this, 'isSpoiler', false);
			this.$set(this, 'formDisabled', true);
			document.querySelector('input[type=file]').value = '';
			this.$set(this, 'image', '');
			this.$set(this, 'imageDimensions', '');
		},
		sendYeah(post) {
			this.$set(post, 'yeahSending', true);
			if(post.has_yeah) {
				// Unyeah
				this.$http.post('/app/posts/' + post.id + '/unyeah')
				.then(response => {
					this.$set(post, 'yeahSending', false);
					this.$set(post, 'has_yeah', false);
					post.yeah_count -= 1;
				}, response => {
					this.$set(post, 'yeahSending', false);
					var error;
					if(response.body.error) {
			      error = response.body.error;
			    } else {
			      error = response.status + ' ' + response.statusText;
			    }
					this.$set(this, 'dialog', {
						title: 'Error',
						text: error
					});
				});
			} else {
				// Yeah
				this.$http.post('/app/posts/' + post.id + '/yeah')
				.then(response => {
					this.$set(post, 'yeahSending', false);
					this.$set(post, 'has_yeah', true);
					if(!post.yeah_count) {
						post.yeah_count = 0;
					}
					post.yeah_count += 1;
				}, response => {
					this.$set(post, 'yeahSending', false);
					var error;
					if(response.body.error) {
			      error = response.body.error;
			    } else {
			      error = response.status + ' ' + response.statusText;
			    }
					this.$set(this, 'dialog', {
						title: 'Error',
						text: error
					});
				});
			}
		},
		revealSpoiler(post) {
			this.$set(post, 'spoilerRevealed', true);
		},
		onScrollEvent() {
			// when called, infinite scroll stuff will work
			if(!this.scrollLoading && this.postCount > 49) {
				var atBottom = ((document.documentElement.scrollHeight - document.documentElement.scrollTop) - window.innerHeight) < 200;
				if(atBottom) {
					this.$set(this, 'scrollLoading', true);
					this.$http.get('/app/communities/' + this.$route.params.id + '?time=' + this.timeLoaded + '&offset=' + this.postCountTotal).then(response => {
						this.$set(this, 'scrollLoading', false);
						var oldPostCount = this.postCountTotal;
						response.body.posts.forEach(post => {
							// Time functions
							var _this = this;
							function updateTime() {
								_this.$set(post, 'timestamp', _this.utils.doTime(post.date.Time));
							}
							updateTime();
							if(this.utils.properNow() - moment(moment(post.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 336084000) {
								setInterval(updateTime, 1000);
							}
							if(post.recent_comment) {
								function updateTime() {
									_this.$set(post.recent_comment, 'timestamp', _this.utils.doTime(post.recent_comment.date.Time));
								}
								updateTime();
								if(this.utils.properNow() - moment(moment(post.recent_comment.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 336084000) {
									setInterval(updateTime, 1000);
								}
							}
						});
						this.$set(this, 'community', response.body.community);
						this.$set(this, 'posts', this.posts.concat(response.body.posts));
						this.postCount = response.body.posts.length;
						this.postCountTotal = response.body.posts.length + oldPostCount;
					});
				}
			}
		},
		closeDialog() {
			this.$set(this, 'dialog', null);
		}
	}
}
</script>

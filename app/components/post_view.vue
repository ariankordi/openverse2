<template>
  <div>
    <div class="no-content" v-show="loading">
      <p>Now loading...</p>
    </div>
    <div class="no-content" v-if="!loading && error">
      <p>{{ error }}</p>
    </div>
    <div class="main-column" v-if="!loading && !error">
    	<div class="post-list-outline">
    		<section class="post post-subtype-default" id="post-content">
    			<header class="community-container">
    				<h1 class="community-container-heading"><router-link :to="'/communities/' + post.community_id"><img class="community-icon" v-bind:src="post.community.icon">{{ post.community.name }}</router-link></h1>
    			</header>
    			<div class="edit-buttons-content">
    				<button class="symbol button edit-button remove" type="button" v-show="!post.deleted" v-if="post.user_id == user.id || user.rank > 3" v-on:click="deletePostDialog()"><span class="symbol-label">Delete</span></button>
            <button class="symbol button edit-button edit" type="button" v-if="post.user_id == user.id" v-on:click="postEditToggle()"><span class="symbol-label">Edit Post</span></button>
    			</div>
    			<div class="user-content">
    				<router-link class="icon-container" v-bind:class="utils.doRank(post.user.rank)" :to="'/users/' + post.user.name"><img class="icon" v-bind:src="utils.doAvatarFeeling(post.user.avatar, post.feeling)"></router-link>
    				<div class="user-name-content">
              <p class="user-organization" v-if="post.user.rank > 0">{{ utils.doRankText(post.user.rank) }}</p>
    					<p class="user-name"><router-link :to="'/users/' + post.user.name">{{ post.user.nick }}</router-link><span class="user-id">{{ post.user.name }}</span></p>
    					<p class="timestamp-container">
                <span class="timestamp">{{ post.timestamp }}</span>
                <span class="spoiler-status spoiler" v-show="post.spoiler">· Spoilers</span>
                <transition name="fade">
                  <span class="spoiler-status spoiler" v-show="post.edited && post.edited.Time != post.date.Time && (new Date(post.edited.Time) - new Date(post.date.Time) > 300000)">· Edited ({{ post.editstamp }})</span>
                </transition><transition name="fade">
                <span class="spoiler-status spoiler" v-show="post.deleted">· <span class="red">Deleted</span></span>
                </transition></p>
    				</div>
    			</div>
    			<div class="body" v-show="!formEditing">
    				<p class="post-content-text" v-if="post.content"><markdown :source="post.content"></markdown></p>
            <div class="screenshot-container still-image" v-show="post.screenshot"><img v-bind:src="post.screenshot"></div>
            <audio v-if="post.url && utils.isAudio(post.url)" v-bind:src="post.url" preload="none" controls></audio>
    				<div class="post-meta">
              <button type="button" class="symbol submit empathy-button" v-bind:class="{disabled: !post.can_yeah || post.yeahSending, 'empathy-added': post.has_yeah}" v-bind:disabled="!post.can_yeah || post.yeahSending" v-on:click="sendYeah(post)"><span class="empathy-button-text">{{ utils.yeahButton(post.feeling, post.has_yeah) }}</span></button>
              <div class="empathy symbol">
    						<span class="symbol-label">Yeahs</span><span class="empathy-count">{{ post.yeah_count || 0 }}</span>
    					</div>
    					<div class="reply symbol">
    						<span class="symbol-label">Comments</span><span class="reply-count">{{ post.comment_count || 0 }}</span>
    					</div>
    				</div>
    			</div>
          <form id="edit-form" v-show="formEditing" v-on:submit.prevent="sendForm()">
            <div class="feeling-selector">
    					<label class="symbol feeling-button feeling-button-normal" v-bind:class="{checked: feelingId == 0}"><input checked type="radio" value="0" v-model="feelingId"><span class="symbol-label">normal</span></label>
    					<label class="symbol feeling-button feeling-button-happy" v-bind:class="{checked: feelingId == 1}"><input type="radio" value="1" v-model="feelingId"><span class="symbol-label">happy</span></label>
    					<label class="symbol feeling-button feeling-button-like" v-bind:class="{checked: feelingId == 2}"><input type="radio" value="2" v-model="feelingId"><span class="symbol-label">like</span></label>
    					<label class="symbol feeling-button feeling-button-surprised" v-bind:class="{checked: feelingId == 3}"><input type="radio" value="3" v-model="feelingId"><span class="symbol-label">surprised</span></label>
    					<label class="symbol feeling-button feeling-button-frustrated" v-bind:class="{checked: feelingId == 4}"><input type="radio" value="4" v-model="feelingId"><span class="symbol-label">frustrated</span></label>
    					<label class="symbol feeling-button feeling-button-puzzled" v-bind:class="{checked: feelingId == 5}"><input type="radio" value="5" v-model="feelingId"><span class="symbol-label">puzzled</span></label>
    				</div>
    				<div class="textarea-with-menu active-text">
    					<div class="textarea-container">
    						<textarea class="textarea-text textarea" maxlength="2000" placeholder="Share your thoughts in a post to this community." v-model="body" v-on:input="formInput()" v-on:paste="pasteImage($event)"></textarea>
    					</div>
    				</div>
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
    		</section>
    		<div id="empathy-content" class="post-permalink" v-bind:class="{none: yeahs.length <= 0 && !post.has_yeah}">
    			<router-link class="post-permalink-feeling-icon visitor" v-bind:class="utils.doRank(user.rank)" v-show="post.has_yeah" :to="'/users/' + user.name"><img class="user-icon" v-bind:src="utils.doAvatarFeeling(user.avatar, post.feeling)"></router-link>
          <router-link class="post-permalink-feeling-icon visitor" v-for="yeah in yeahs" :key="yeah.user_id" v-if="user.id != yeah.user_id" v-bind:class="utils.doRank(yeah.user.rank)" :to="'/users/' + yeah.user.name"><img class="user-icon" v-bind:src="utils.doAvatarFeeling(yeah.user.avatar, post.feeling)" v-bind:alt="yeah.user.nick"></router-link>
    		</div>
    		<div id="reply-content">
    			<h2 class="reply-label">Comments</h2>
    			<div class="no-reply-content" v-show="comments.length <= 0">
    				<div>
    					<p>This post has no comments.</p>
    				</div>
    			</div>
          <button class="more-button active oldest-replies-button" v-show="post.comment_count > 20 && comments.length <= 20 && !commentsLoading" v-on:click="getMoreComments(true)"><span class="symbol">Show all comments ({{ post.comment_count }})</span></button>
          <button class="more-button active older-replies-button" v-bind:class="{'oldest-replies-button': post.comment_count < 40}" v-show="post.comment_count > 40 && comments.length < post.comment_count && !commentsLoading" v-on:click="getMoreComments()"><span class="symbol">Show more comments</span></button>
          <button class="more-button all-replies-button disabled" disabled v-show="post.comment_count > 20 && comments.length < post.comment_count && commentsLoading"><span class="symbol">Now loading...</span></button>
    			<transition-group class="list reply-list test-reply-list" v-show="comments.length > 0" name="fade" tag="ul">
            <li class="post trigger" v-bind:class="{my: comment.user_id == post.user_id, other: comment.user_id != post.user_id, hidden: comment.spoiler && comment.user_id != user.id && !comment.spoilerRevealed}" v-for="comment in comments" v-bind:key="comment.id" v-on:click="goToComment($event, comment)">
            	<router-link class="icon-container" v-bind:class="utils.doRank(comment.user.rank)" :to="'/users/' + comment.user.name"><img class="icon" v-bind:src="utils.doAvatarFeeling(comment.user.avatar, comment.feeling)"></router-link>
            	<div class="body">
            		<div class="header">
            			<p class="user-name"><router-link :to="'/users/' + comment.user.name">{{ comment.user.nick }}</router-link></p>
            			<p class="timestamp-container">
                    <router-link class="timestamp" :to="'/replies/' + comment.id">{{ comment.timestamp }}</router-link>
                    <span class="spoiler-status spoiler" v-show="comment.spoiler">· Spoilers</span>
                    <span class="spoiler-status spoiler" v-show="comment.edited && comment.edited.Time != comment.date.Time && (new Date(comment.edited.Time) - new Date(comment.date.Time) > 300000)">· Edited ({{ comment.editstamp }})</span></p>
            		</div>
            		<p class="reply-content-text" v-if="comment.content"><markdown :source="comment.content"></markdown></p>
                <router-link class="screenshot-container still-image" :to="'/replies/' + comment.id" v-show="comment.screenshot"><img v-bind:src="comment.screenshot"></router-link>
  							<div class="hidden-content" v-show="comment.spoiler && comment.user_id != user.id && !comment.spoilerRevealed">
  								<p>This comment contains spoilers.</p>
  								<button type="button" class="hidden-content-button" v-on:click="revealCommentSpoiler(comment)">View Comment</button>
  							</div>
            		<div class="reply-meta">
            		  <button type="button" class="symbol submit empathy-button" v-bind:class="{disabled: !comment.can_yeah || comment.yeahSending, 'empathy-added': comment.has_yeah}" v-bind:disabled="!comment.can_yeah || comment.yeahSending" v-on:click="sendCommentYeah(comment)"><span class="empathy-button-text">{{ utils.yeahButton(comment.feeling, comment.has_yeah) }}</span></button>
            			<div class="empathy symbol">
            				<span class="symbol-label">Yeahs</span><span class="empathy-count">{{ comment.yeah_count || 0 }}</span>
            			</div>
            		</div>
            	</div>
            </li>
          </transition-group>
    		</div>
    		<h2 class="reply-label">Add a Comment</h2>
        <form id="reply-form" v-if="loggedIn" v-on:submit.prevent="commentSendForm()">
        	<div class="feeling-selector">
            <label class="symbol feeling-button feeling-button-normal" v-bind:class="{checked: commentFeelingId == 0}"><input checked type="radio" value="0" v-model="commentFeelingId"><span class="symbol-label">normal</span></label>
            <label class="symbol feeling-button feeling-button-happy" v-bind:class="{checked: commentFeelingId == 1}"><input type="radio" value="1" v-model="commentFeelingId"><span class="symbol-label">happy</span></label>
            <label class="symbol feeling-button feeling-button-like" v-bind:class="{checked: commentFeelingId == 2}"><input type="radio" value="2" v-model="commentFeelingId"><span class="symbol-label">like</span></label>
            <label class="symbol feeling-button feeling-button-surprised" v-bind:class="{checked: commentFeelingId == 3}"><input type="radio" value="3" v-model="commentFeelingId"><span class="symbol-label">surprised</span></label>
            <label class="symbol feeling-button feeling-button-frustrated" v-bind:class="{checked: commentFeelingId == 4}"><input type="radio" value="4" v-model="commentFeelingId"><span class="symbol-label">frustrated</span></label>
            <label class="symbol feeling-button feeling-button-puzzled" v-bind:class="{checked: commentFeelingId == 5}"><input type="radio" value="5" v-model="commentFeelingId"><span class="symbol-label">puzzled</span></label>
          </div>
        	<div class="textarea-with-menu active-text">
        		<!--<menu class="textarea-menu">
        			<ul>
        				<li><label class="textarea-menu-text"><input name="_post_type" type="radio" value="body"></label></li>
        				<li><label class="textarea-menu-memo"><input name="_post_type" type="radio" value="painting"></label></li>
        			</ul>
        		</menu>-->
            <div class="textarea-container">
  						<textarea class="textarea-text textarea" maxlength="2000" placeholder="Add a comment here." v-model="commentBody" v-on:input="commentFormInput()" v-on:paste="commentPasteImage($event)"></textarea>
  					</div>
        	</div>
          <label class="file-button-container" v-on:drop.prevent="commentPasteImage($event)" v-on:dragover.prevent="">
            <span class="input-label">Image <span v-show="!commentImage && !commentImageDimensions">powered by catgirl.host</span><span v-show="commentImage && !commentImageDimensions">...</span><span v-show="commentImageDimensions">{{ commentImageDimensions }}</span></span>
            <span class="button file-upload-button">Upload</span>
            <input accept="image/*" class="file-button none" type="file" v-on:change="commentImageUpload($event)">
            <div class="screenshot-container still-image" v-show="commentImage">
              <img v-bind:src="commentImage">
            </div>
          </label>
          <div class="post-form-footer-options">
            <label class="spoiler-button symbol" v-bind:class="{checked: commentIsSpoiler}"><input id="is_spoiler" type="checkbox" value="1" v-model="commentIsSpoiler"> Spoilers</label>
          </div>
          <div class="form-buttons">
  					<input class="black-button post-button" type="submit" value="Send" v-bind:class="{disabled: commentFormDisabled}" v-bind:disabled="commentFormDisabled" v-show="!commentFormSending">
  					<input class="black-button post-button disabled" disabled type="submit" value="Please wait..." v-show="commentFormSending">
  				</div>
        </form>
        <div class="cannot-reply" v-show="!loggedIn">
          <p>You cannot comment on this post.</p>
        </div>
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
              <div class="form-buttons" v-if="dialog.postDelete">
                <button class="gray-button" type="button" v-bind:disabled="deleteSending" v-on:click="closeDialog()">No</button>
                <button class="black-button" type="button" v-on:click="deletePost()" v-show="!deleteSending">Yes</button>
                <button class="black-button disabled" disabled type="button" v-show="deleteSending">Please wait...</button>
              </div>
              <div class="form-buttons" v-if="!dialog.postDelete">
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
export default {
  computed: {
    user() {
      return window.user;
    },
    loggedIn() {
			return window.user.id > 0;
		},
    utils() {
      return window.utils;
    }
  },
  created() {
    //if(!this.post) {
      var __this = this;
      this.$http.get('/app/posts/' + this.$route.params.id).then(response => {
        this.$set(this, 'loading', false);
        if(response.body.post) {
          // Time functions
          function updateTime() {
            __this.$set(response.body.post, 'timestamp', __this.utils.doTime(response.body.post.date.Time));
          }
          updateTime();
          if(__this.utils.properNow() - moment(moment(response.body.post.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
            setInterval(updateTime, 1000);
          }
          if(response.body.post.edited && response.body.post.edited.Time != response.body.post.date.Time) {
            // If post edited, then put "editstamp" there
            function updateEditTime() {
              __this.$set(response.body.post, 'editstamp', __this.utils.doTime(response.body.post.edited.Time));
            }
            updateEditTime();
            if(__this.utils.properNow() - moment(moment(response.body.post.edited.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
              setInterval(updateEditTime, 1000);
            }
          }

          this.$set(this, 'post', response.body.post);
          if(response.body.post.user_id == user.id) {
            utils.title('Your Post');
          } else {
            utils.title(response.body.post.user.nick + '\'s Post');
          }
  	      this.$set(this, 'yeahs', response.body.yeahs);
          response.body.comments.forEach(comment => {
  					// Time functions
  					function updateTime() {
  						__this.$set(comment, 'timestamp', __this.utils.doTime(comment.date.Time));
  					}
  					updateTime();
  					if(__this.utils.properNow() - moment(moment(comment.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
  						setInterval(updateTime, 1000);
  					}
            if(comment.edited && comment.edited.Time != comment.date.Time) {
              // If comment edited, then put "editstamp" there
              function updateEditTime() {
                __this.$set(comment, 'editstamp', __this.utils.doTime(comment.edited.Time));
              }
              updateEditTime();
              if(__this.utils.properNow() - moment(moment(comment.edited.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
                setInterval(updateEditTime, 1000);
              }
            }
  				});
          this.$set(this, 'comments', response.body.comments);
          this.postEdited = this.post;
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
    //}
  },
  mounted() {
    if(!this.user) {
			return;
		}
    var _this = this;
    this.messageStream = new WebSocket('ws' + ((location.protocol == 'https:') ? 's' : '') + '://' + location.host + '/app/posts/' + this.$route.params.id + '/stream');
		this.messageStream.onmessage = function(event) {
			var type = event.data.slice(0, event.data.indexOf(':'));
			var data = JSON.parse(event.data.substring(event.data.indexOf(':') + 1));
      switch(type) {
        case 'comment':
          if(data.user_id != _this.user.id) {
            _this.addComment(data);
          }
          break;
        case 'edit':
          _this.$set(_this.post, 'content', data.content);
          _this.$set(_this.post, 'feeling', data.feeling);
          _this.$set(_this.post, 'spoiler', data.spoiler);
          _this.$set(_this.post, 'screenshot', data.screenshot);
          _this.$set(_this.post, 'edited', {Time: _this.utils.properNowStamp()});
          break;
        case 'yeah':
          if(data.id != _this.user.id) {
            if(!_this.yeahs) {
              _this.yeahs = [];
            }
            _this.yeahs.unshift({user: data, user_id: data.id});
            _this.$set(_this.post, 'yeah_count', (_this.post.yeah_count || 0) + 1);
          }
          break;
        case 'unyeah':
          if(data != _this.user.id) {
            console.log(data)
            var yeah = _this.yeahs.filter(function(yeah) {
              return yeah.user_id == data;
            })[0];
            var index = _this.yeahs.indexOf(yeah);
            if(index > -1) {
              _this.yeahs.splice(index, 1);
            }
            _this.$set(_this.post, 'yeah_count', (_this.post.yeah_count || 0) - 1);
          }
          break;
        case 'commentyeah':
          if(data.user_id != _this.user.id) {
            var comment = _this.comments.filter(function(comment) {
              return comment.id == data.topic;
            })[0];
            _this.$set(comment, 'yeah_count', (comment.yeah_count || 0) + 1);
          }
          break;
        case 'commentunyeah':
          if(data.user_id != _this.user.id) {
            var comment = _this.comments.filter(function(comment) {
              return comment.id == data.topic;
            })[0];
            _this.$set(comment, 'yeah_count', (comment.yeah_count || 0) - 1);
          }
          break;
          case 'commentedit':
            if(data.user_id != _this.user.id) {
              var comment = _this.comments.filter(function(comment) {
                return comment.id == data.topic;
              })[0];
  						_this.$set(comment, 'content', data.comment.content);
  						_this.$set(comment, 'feeling', data.comment.feeling);
  						_this.$set(comment, 'spoiler', data.comment.spoiler);
  						_this.$set(comment, 'screenshot', data.comment.screenshot);
  						_this.$set(comment, 'edited', {Time: _this.utils.properNowStamp()});
            }
            break;
      }
    }
  },
  destroyed() {
    if(this.messageStream) {
      this.messageStream.close();
    }
  },
  data() {
    return {
      loading: true,
      error: '',
      post: null,
      yeahs: [],
      comments: [],
      commentFeelingId: 0,
			commentBody: '',
			commentIsSpoiler: false,
      commentImage: '',
      commentImageDimensions: '',
			commentFormFolded: true,
			commentFormDisabled: true,
			commentFormSending: false,
      feelingId: 0,
			body: '',
			isSpoiler: false,
			image: '',
			imageDimensions: '',
      formEditing: false,
      formSending: false,
      formDisabled: false,
      deleteSending: false,
      dialog: null,
      commentsLoading: false
    }
  },
  methods: {
    /*getSaveStateConfig() {
			return {
				cacheKey: 'PostView'
			}
		},*/
    goToComment(event, comment) {
      if(!event.target.closest("a, button")) {
        if(comment.spoiler && comment.user_id != this.user.id && !comment.spoilerRevealed) {
          return false;
        }
        this.$router.push('/replies/' + comment.id);
      }
    },
    sendYeah(post) {
      this.$set(post, 'yeahSending', true);
      if(post.has_yeah) {
        // Unyeah
        this.$http.post('/app/posts/' + post.id + '/unyeah')
        .then(response => {
          this.$set(post, 'yeahSending', false);
          this.$set(post, 'has_yeah', false);
          this.$set(post, 'yeah_count', (post.yeah_count || 0) - 1);
          if(post.yeah_count <= 0 && this.yeahs.length > 0) {
            this.yeahs = [];
          }
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
          this.$set(post, 'yeah_count', (post.yeah_count || 0) + 1);
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
    sendCommentYeah(comment) {
      this.$set(comment, 'yeahSending', true);
      if(comment.has_yeah) {
        // Unyeah
        this.$http.post('/app/replies/' + comment.id + '/unyeah')
        .then(response => {
          this.$set(comment, 'yeahSending', false);
          this.$set(comment, 'has_yeah', false);
          comment.yeah_count -= 1;
        }, response => {
          this.$set(comment, 'yeahSending', false);
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
        this.$http.post('/app/replies/' + comment.id + '/yeah')
        .then(response => {
          this.$set(comment, 'yeahSending', false);
          this.$set(comment, 'has_yeah', true);
          if(!comment.yeah_count) {
            comment.yeah_count = 0;
          }
          comment.yeah_count += 1;
        }, response => {
          this.$set(comment, 'yeahSending', false);
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
    revealCommentSpoiler(comment) {
			this.$set(comment, 'spoilerRevealed', true);
		},
    commentFormInput() {
      if(this.commentImage) {
				this.$set(this, 'commentFormDisabled', false);
				return;
			}
			if(/^[\s\u00A0\u3000]*$/.test(this.commentBody)) {
				this.$set(this, 'commentFormDisabled', true);
				return;
			}
			this.$set(this, 'commentFormDisabled', false);
		},
    commentSendForm() {
      if(this.formDisabled) {
				return false;
			}
      this.$set(this, 'commentFormSending', true);
			this.$http.post('/app/posts/' + this.post.id + '/reply', {
				body: this.commentBody,
				feeling_id: this.commentFeelingId,
				is_spoiler: this.commentIsSpoiler,
        screenshot: this.commentImage
			}).then(response => {
				this.$set(this, 'commentFormSending', false);
				this.addComment(response.body);
				this.commentResetForm();
			}, response => {
				this.$set(this, 'commentFormSending', false);
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
    commentImageUpload(event) {
			if(navigator.userAgent.indexOf('Nintendo') > 0) {
				alert('todo: file uploading is broken on Wii U! make file uploading use FormData() somehow?');
			}
			var reader = new FileReader();
			reader.readAsDataURL(event.target.files[0]);
			var _this = this;
			reader.onload = function() {
				_this.$set(_this, 'commentImage', reader.result);
        _this.$set(_this, 'commentFormDisabled', false);
				var imageDimension = new Image();
				imageDimension.src = reader.result;
        imageDimension.onerror = function() {
					_this.$set(_this, 'commentImage', '');
					_this.$set(_this, 'commentImageDimensions', 'Invalid image.');
          _this.commentFormInput();
				}
				imageDimension.onload = function() {
					_this.$set(_this, 'commentImageDimensions', imageDimension.width + ' x ' + imageDimension.height);
				}
			}
		},
    commentPasteImage(event) {
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
            _this.$set(_this, 'commentImage', reader.result);
            _this.$set(_this, 'commentFormDisabled', false);
            var imageDimension = new Image();
            imageDimension.src = reader.result;
            imageDimension.onerror = function() {
    					_this.$set(_this, 'commentImage', '');
    					_this.$set(_this, 'commentImageDimensions', 'Invalid image.');
    					_this.commentFormInput();
    				}
            imageDimension.onload = function() {
              _this.$set(_this, 'commentImageDimensions', imageDimension.width + ' x ' + imageDimension.height);
            }
          }
        }
      });
    },
    commentResetForm() {
			this.$set(this, 'commentFeelingId', 0);
			this.$set(this, 'commentBody', '');
			this.$set(this, 'commentIsSpoiler', false);
			this.$set(this, 'commentFormDisabled', true);
      document.querySelector('#reply-form input[type=file]').value = '';
			this.$set(this, 'commentImage', '');
			this.$set(this, 'commentImageDimensions', '');
		},
    addComment(comment) {
      if(this.user.id > 0 && comment.user_id != this.user.id) {
        comment.can_yeah = true;
      }
			// moment hates timezones so we're using hacks to make it proper
			comment.date.Time = this.utils.properNowStamp();
			// Update time here too
			var _this = this;
			function updateTime() {
				_this.$set(comment, 'timestamp', _this.utils.doTime(comment.date.Time));
			}
			updateTime();
			setInterval(updateTime, 1000);

      if(this.comments.length > 19) {
				this.comments.shift();
			}
      this.$set(this.post, 'comment_count', (this.post.comment_count || 0) + 1);
			this.comments.push(comment);
		},
    deletePost() {
      this.$set(this, 'deleteSending', true);
      this.$http.post('/app/posts/' + this.post.id + '/delete')
      .then(response => {
        this.$set(this, 'deleteSending', false);
        this.$set(this, 'dialog', null);
        this.$set(this.post, 'deleted', true);
      }, response => {
        this.$set(this, 'deleteSending', false);
        this.$set(this, 'dialog', null);
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
    deletePostDialog() {
      this.$set(this, 'dialog', {
        title: 'Delete Post',
        text: 'Do you really want to delete this post?',
        postDelete: true
      });
    },
    postEditToggle() {
      if(!this.feelingId) {
        this.feelingId = this.post.feeling || 0;
      }
      if(!this.body) {
        this.body = this.post.content;
      }
      if(!this.isSpoiler) {
        this.isSpoiler = this.post.spoiler;
      }
      this.$set(this, 'formEditing', !this.formEditing);
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
    sendForm() {
      if(this.formDisabled) {
        return false;
      }
      // if feeling id, body, and spoiler matches post, and there's no image being sent
      if(this.feelingId == (this.post.feeling || 0) && this.body == this.post.content && Boolean(this.isSpoiler) == Boolean(this.post.spoiler) && (!this.image || this.image == this.post.screenshot)) {
        this.postEditToggle();
        return;
      }
      this.$set(this, 'formSending', true);
      this.$http.post('/app/posts/' + this.post.id + '/edit', {
        body: this.body,
        feeling_id: this.feelingId,
        is_spoiler: this.isSpoiler,
        screenshot: this.image,
      }).then(response => {
        this.$set(this, 'formSending', false);
        // after post edit, change fields
        this.$set(this.post, 'feeling', +this.feelingId);
        this.$set(this.post, 'content', this.body);
        this.$set(this.post, 'spoiler', this.isSpoiler);
        if(this.image) {
          this.$set(this.post, 'screenshot', {
            String: this.image
          });
        }
        this.$set(this.post, 'edited', {
          Time: this.utils.properNowStamp().format(),
          Valid: true
        });
        var _this = this;
        function updateEditTime() {
          _this.$set(_this.post, 'editstamp', _this.utils.doTime(_this.post.edited.Time));
        }
        updateEditTime();
        if(this.utils.properNow() - moment(moment(this.post.edited.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
          setInterval(updateEditTime, 1000);
        }
        document.querySelector('#edit-form input[type=file]').value = '';
        this.$set(this, 'image', '');
        this.$set(this, 'imageDimensions', '');
        this.postEditToggle();
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
    closeDialog() {
			this.$set(this, 'dialog', null);
		},
    getMoreComments(all) {
      this.$set(this, 'commentsLoading', true);
      var url;
      if(all) {
        url = '/app/posts/' + this.post.id + '/replies?limit=' + (this.comments.length > 20 ? this.comment_count - this.comments.length - 20 : this.post.comment_count - 20);
      } else if(((this.post.comment_count - this.comments.length) - 20) < 20) {
        url = '/app/posts/' + this.post.id + '/replies?limit=' + (this.post.comment_count - this.comments.length);
      } else {
        url = '/app/posts/' + this.post.id + '/replies?offset=' + ((this.post.comment_count - this.comments.length) - 20);
      }
      this.$http.get(url).then(response => {
        this.$set(this, 'commentsLoading', false);
        var _this = this;
        response.body.forEach(comment => {
          function updateTime() {
            _this.$set(comment, 'timestamp', _this.utils.doTime(comment.date.Time));
          }
          updateTime();
          if(_this.utils.properNow() - moment(moment(comment.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
            setInterval(updateTime, 1000);
          }
        });
        this.$set(this, 'comments', response.body.concat(this.comments));
      }, response => {
        this.$set(this, 'commentsLoading', false);
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
  }
}
</script>

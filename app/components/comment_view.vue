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
      	<router-link class="post-permalink-button info-ticker" :to="'/posts/' + post.id">
          <span class="icon-container" v-bind:class="utils.doRank(post.user.rank)">
            <img class="icon" v-bind:src="utils.doAvatarFeeling(post.user.avatar, post.feeling)">
          </span>
          <span>View <span class="post-user-description">{{ post.user.nick }}'s post ({{ utils.textTrun(post.content, 35) }})</span> for this comment.</span>
        </router-link>
      </div>
      <div class="post-list-outline more">
      	<div class="post reply-permalink-post" id="post-content">
      		<div>
      			<p class="community-container">
              <router-link :to="'/communities/' + post.community.id">
                <img class="community-icon" v-bind:src="post.community.icon">{{ post.community.name }}</router-link>
            </p>
            <div class="edit-buttons-content">
      				<button class="symbol button edit-button remove" type="button" v-show="!comment.deleted" v-if="comment.user_id == user.id || user.rank > 3" v-on:click="deleteCommentDialog()"><span class="symbol-label">Delete</span></button>
              <button class="symbol button edit-button edit" type="button" v-if="comment.user_id == user.id" v-on:click="commentEditToggle()"><span class="symbol-label">Edit Comment</span></button>
      			</div>
      			<div class="user-content">
      				<router-link class="icon-container" v-bind:class="utils.doRank(comment.user.rank)" :to="'/users/' + comment.user.name"><img class="icon" v-bind:src="utils.doAvatarFeeling(comment.user.avatar, comment.feeling)"></router-link>
              <div class="user-name-content">
      					<p class="user-name"><router-link :to="'/users/' + comment.user.name">{{ comment.user.nick }}</router-link><span class="user-id">{{ comment.user.name }}</span></p>
      					<p class="timestamp-container"><span class="timestamp">{{ comment.timestamp }}</span>
                <span class="spoiler-status spoiler" v-show="comment.spoiler">· Spoilers</span>
                <transition name="fade">
                  <span class="spoiler-status spoiler" v-show="comment.edited && comment.edited.Time != comment.date.Time && (new Date(comment.edited.Time) - new Date(comment.date.Time) > 300000)">· Edited ({{ comment.editstamp }})</span>
                </transition><transition name="fade">
                  <span class="spoiler-status spoiler" v-show="comment.deleted">· <span class="red">Deleted</span></span>
                </transition></p>
      				</div>
      			</div>
            <div class="body" v-show="!formEditing">
      				<p class="reply-content-text" v-if="comment.content"><markdown :source="comment.content"></markdown></p>
      				<div class="screenshot-container still-image" v-show="comment.screenshot"><img v-bind:src="comment.screenshot"></div>
      				<div class="post-meta">
                <button type="button" class="symbol submit empathy-button" v-bind:class="{disabled: !comment.can_yeah || comment.yeahSending, 'empathy-added': comment.has_yeah}" v-bind:disabled="!comment.can_yeah || comment.yeahSending" v-on:click="sendYeah(comment)"><span class="empathy-button-text">{{ utils.yeahButton(comment.feeling, comment.has_yeah) }}</span></button>
                <div class="empathy symbol">
                  <span class="symbol-label">Yeahs</span><span class="empathy-count">{{ comment.yeah_count || 0 }}</span>
                </div>
      				</div>
              <div id="empathy-content" class="post-permalink" v-bind:class="{none: yeahs.length <= 0 && !comment.has_yeah}">
          			<router-link class="post-permalink-feeling-icon visitor" v-bind:class="utils.doRank(user.rank)" v-show="comment.has_yeah" :to="'/users/' + user.name"><img class="user-icon" v-bind:src="utils.doAvatarFeeling(user.avatar, comment.feeling)"></router-link>
                <router-link class="post-permalink-feeling-icon visitor" v-for="yeah in yeahs" :key="yeah.user_id" v-if="user.id != yeah.user_id" v-bind:class="utils.doRank(yeah.user.rank)" :to="'/users/' + yeah.user.name"><img class="user-icon" v-bind:src="utils.doAvatarFeeling(yeah.user.avatar, comment.feeling)" v-bind:alt="yeah.user.nick"></router-link>
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
      		</div>
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
              <div class="form-buttons" v-if="dialog.commentDelete">
                <button class="gray-button" type="button" v-bind:disabled="deleteSending" v-on:click="closeDialog()">No</button>
                <button class="black-button" type="button" v-on:click="deleteComment()" v-show="!deleteSending">Yes</button>
                <button class="black-button disabled" disabled type="button" v-show="deleteSending">Please wait...</button>
              </div>
              <div class="form-buttons" v-if="!dialog.commentDelete">
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
    //if(!this.comment) {
      var __this = this;
      this.$http.get('/app/replies/' + this.$route.params.id).then(response => {
        this.$set(this, 'loading', false);
        if(response.body.post) {
          // Time functions
          function updateTime() {
            __this.$set(response.body.comment, 'timestamp', __this.utils.doTime(response.body.comment.date.Time));
          }
          updateTime();
          if(__this.utils.properNow() - moment(moment(response.body.comment.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
            setInterval(updateTime, 1000);
          }
          if(response.body.comment.edited && response.body.comment.edited.Time != response.body.comment.date.Time) {
            // If comment edited, then put "editstamp" there
            function updateEditTime() {
              __this.$set(response.body.comment, 'editstamp', __this.utils.doTime(response.body.comment.edited.Time));
            }
            updateEditTime();
            if(__this.utils.properNow() - moment(moment(response.body.comment.edited.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
              setInterval(updateEditTime, 1000);
            }
          }

          this.$set(this, 'comment', response.body.comment);
          this.$set(this, 'post', response.body.post);
          if(response.body.post.user_id == user.id && response.body.comment.user_id == user.id) {
            utils.title('Your Comment on Your Post');
          } else if(response.body.comment.user_id == user.id) {
            utils.title('Your Comment on ' + response.body.post.user.nick + '\'s Post');
          } else if(response.body.post.user_id == user.id) {
            utils.title(response.body.comment.user.nick + '\'s Comment on Your Post');
          } else {
            utils.title(response.body.comment.user.nick + '\'s Comment on ' + response.body.post.user.nick + '\'s Post');
          }
  	      this.$set(this, 'yeahs', response.body.yeahs);
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
  data() {
    return {
      loading: true,
      error: '',
      comment: null,
      post: {},
      yeahs: [],
      feelingId: 0,
			body: '',
			isSpoiler: false,
			image: '',
			imageDimensions: '',
      formEditing: false,
      formSending: false,
      formDisabled: false,
      deleteSending: false,
      dialog: null
    }
  },
  methods: {
    /*getSaveStateConfig() {
			return {
				cacheKey: 'CommentView'
			}
		},*/
    sendYeah(comment) {
      this.$set(comment, 'yeahSending', true);
      if(comment.has_yeah) {
        // Unyeah
        this.$http.post('/app/replies/' + comment.id + '/unyeah')
        .then(response => {
          this.$set(comment, 'yeahSending', false);
          this.$set(comment, 'has_yeah', false);
          this.$set(comment, 'yeah_count', (comment.yeah_count || 0) - 1);
          if(comment.yeah_count <= 0 && this.yeahs.length > 0) {
            this.yeahs = [];
          }
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
          this.$set(comment, 'yeah_count', (comment.yeah_count || 0) + 1);
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
    deleteComment() {
      this.$set(this, 'deleteSending', true);
      this.$http.post('/app/replies/' + this.comment.id + '/delete')
      .then(response => {
        this.$set(this, 'deleteSending', false);
        this.$set(this, 'dialog', null);
        this.$set(this.comment, 'deleted', true);
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
    deleteCommentDialog() {
      this.$set(this, 'dialog', {
        title: 'Delete Comment',
        text: 'Do you really want to delete this comment?',
        commentDelete: true
      });
    },
    commentEditToggle() {
      if(!this.feelingId) {
        this.feelingId = this.comment.feeling || 0;
      }
      if(!this.body) {
        this.body = this.comment.content;
      }
      if(!this.isSpoiler) {
        this.isSpoiler = this.comment.spoiler;
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
      if(this.feelingId == (this.comment.feeling || 0) && this.body == this.comment.content && Boolean(this.isSpoiler) == Boolean(this.comment.spoiler) && (!this.image || this.image == this.post.screenshot)) {
        this.commentEditToggle();
        return;
      }
      this.$set(this, 'formSending', true);
      this.$http.post('/app/replies/' + this.comment.id + '/edit', {
        body: this.body,
        feeling_id: this.feelingId,
        is_spoiler: this.isSpoiler,
        screenshot: this.image,
      }).then(response => {
        this.$set(this, 'formSending', false);
        // after comment edit, change fields
        this.$set(this.comment, 'feeling', this.feelingId);
        this.$set(this.comment, 'content', this.body);
        this.$set(this.comment, 'spoiler', this.isSpoiler);
        if(this.image) {
          this.$set(this.comment, 'screenshot', {
            String: this.image
          });
        }
        this.$set(this.comment, 'edited', {
          Time: this.utils.properNowStamp().format(),
          Valid: true
        });
        var _this = this;
        function updateEditTime() {
          _this.$set(_this.comment, 'editstamp', _this.utils.doTime(_this.comment.edited.Time));
        }
        updateEditTime();
        if(this.utils.properNow() - moment(moment(this.comment.edited.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
          setInterval(updateEditTime, 1000);
        }
        document.querySelector('input[type=file]').value = '';
  			this.$set(this, 'image', '');
  			this.$set(this, 'imageDimensions', '');
        this.commentEditToggle();
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
		}
  }
}
</script>

<template>
  <div>
    <user-sidebar :user="user" :general="true"></user-sidebar>
    <div class="no-content" v-show="loading">
      <p>Now loading...</p>
    </div>
    <div class="no-content" v-if="!loading && error">
      <p>{{ error }}</p>
    </div>
    <div class="main-column" v-if="!loading && !error">
    	<div class="post-list-outline">
    		<h2 class="label">Notifications</h2>
    		<!--<div class="tab-container" id="notification-tab-container">
    			<div class="tab2">
    				<a class="tab-icon-my-news selected" href="/notifications"><span class="symbol nf"></span> <span>Updates</span></a> <a class="tab-icon-my-news" href="/notifications/friend_requests"><span class="symbol fr"></span> <span>Friend Requests</span></a>
    			</div>
    		</div>-->
    		<div class="list news-list">
          <div class="no-content" v-show="notifications.length <= 0">
  					<p>No notifications.</p>
  				</div>
    			<div class="news-list-content trigger" v-bind:class="{notify: !notification.read}" v-on:click="goToNotification($event, notification)" v-for="notification in notifications">
    				<router-link class="icon-container" v-bind:class="utils.doRank(notification.users[0].rank)" :to="'/users/' + notification.users[0].name"><img class="icon" v-bind:src="utils.doAvatarFeeling(notification.users[0].avatar, notification.feeling)"></router-link>
    				<div class="body">
              <div>
                <div v-if="notification.type < 4">
                  <div v-if="notification.users.length > 5">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>, and {{ notification.users.length - 4 }} others
                  </div>
                  <div v-else-if="notification.users.length > 4">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>, and 1 other person
                  </div>
                  <div v-else-if="notification.users.length > 3">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, and <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>
                  </div>
                  <div v-else-if="notification.users.length > 2">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link> and <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>
                  </div>
                  <div v-else-if="notification.users.length > 1">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link> and <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>
                  </div>
                  <div v-else-if="notification.users.length > 0">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>
                  </div>
                </div>

                <span v-if="notification.type == 0">gave <router-link class="link" :to="'/posts/' + notification.topic">your post ({{ utils.textTrun(notification.post.content, 35) }})</router-link> a Yeah.</span>
                <span v-else-if="notification.type == 1">gave <router-link class="link" :to="'/replies/' + notification.topic">your comment ({{ utils.textTrun(notification.comment.content, 35) }})</router-link> a Yeah.</span>
                <span v-else-if="notification.type == 2">commented on <router-link class="link" :to="'/posts/' + notification.topic">your post ({{ utils.textTrun(notification.post.content, 35) }})</router-link></span>
                <span v-else-if="notification.type == 3">commented on <router-link class="link" :to="'/posts/' + notification.topic">{{ notification.users[0].nick }}'s post ({{ utils.textTrun(notification.post.content, 35) }})</router-link></span>
                <span v-else-if="notification.type == 4">Followed by
                  <div v-if="notification.users.length > 5">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>, and {{ notification.users.length - 4 }} others.
                  </div>
                  <div v-else-if="notification.users.length > 4">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>, and 1 other person.
                  </div>
                  <div v-else-if="notification.users.length > 3">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>, and <router-link class="nick-name" :to="'/users/' + notification.users[3].name">{{ notification.users[3].nick }}</router-link>.
                  </div>
                  <div v-else-if="notification.users.length > 2">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>, <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link> and <router-link class="nick-name" :to="'/users/' + notification.users[2].name">{{ notification.users[2].nick }}</router-link>.
                  </div>
                  <div v-else-if="notification.users.length > 1">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link> and <router-link class="nick-name" :to="'/users/' + notification.users[1].name">{{ notification.users[1].nick }}</router-link>.
                  </div>
                  <div v-else-if="notification.users.length > 0">
                    <router-link class="nick-name" :to="'/users/' + notification.users[0].name">{{ notification.users[0].nick }}</router-link>.
                  </div>
                </span>
              </div>
              <span class="timestamp">{{ notification.timestamp }}</span>
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
    utils() {
      return window.utils;
    }
  },
  created() {
    //if(!this.notifications) {
      this.$http.get('/app/news/my_news').then(response => {
        this.$set(this, 'loading', false);
        var _this = this;
        response.body.forEach(notification => {
          // Time functions
          function updateTime() {
            _this.$set(notification, 'timestamp', _this.utils.doTime(notification.date.Time));
          }
          updateTime();
          if(_this.utils.properNow() - moment(moment(notification.date.Time).utc().format('YYYY-MM-DDTHH:mm:ss')) < 156084000) {
            setInterval(updateTime, 1000);
          }
          // get feeling ID from either post, comment or none
          if(notification.post) {
            notification.feeling = notification.post.feeling;
          } else if(notification.comment) {
            notification.feeling = notification.comment.feeling;
          } else {
            notification.feeling = 0;
          }
          if(!notification.type) {
            notification.type = 0;
          }
        });
        this.$set(this, 'notifications', response.body);
        utils.title('Notifications');
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
      notifications: null
    }
  },
  methods: {
    goToNotification(event, notification) {
      if(!event.target.closest("a")) {
        switch(notification.type) {
          case 0:
            this.$router.push('/posts/' + notification.topic);
            break;
          case 1:
            this.$router.push('/replies/' + notification.topic);
            break;
          case 2:
            this.$router.push('/posts/' + notification.topic);
            break;
          case 3:
            this.$router.push('/posts/' + notification.topic);
            break;
          case 4:
            this.$router.push('/users/' + notification.users[0].name);
            break;
        }
      }
    }
  }
}
</script>

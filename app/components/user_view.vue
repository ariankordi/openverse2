<template>
	<div id="profile-top">
    <div class="no-content" v-show="loading">
      <p>Now loading...</p>
    </div>
    <div class="no-content" v-if="!loading && error">
      <p>{{ error }}</p>
    </div>
		<user-sidebar :user="user" :general="false" v-if="!loading && !error"></user-sidebar>
    <div class="main-column" v-if="!loading && !error">
		</div>
	</div>
</template>
<script>
export default {
	computed: {
    utils() {
      return window.utils;
    }
  },
	created() {
    //if(!this.user) {
      this.$http.get('/app/users/' + this.$route.params.id).then(response => {
        this.$set(this, 'loading', false);
        this.$set(this, 'user', response.body);
        utils.title(this.user.nick + '\'s Profile');
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
			user: null
		}
	}
}
</script>

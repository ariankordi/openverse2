<template>
  <div>
    <div class="no-content" v-show="loading">
      <p>Now loading...</p>
    </div>
    <div class="no-content" v-if="!loading && error">
      <p>{{ error }}</p>
    </div>
    <user-sidebar :user="user" :general="false" v-if="!loading && !error"></user-sidebar>
    <div class="main-column" v-if="!loading && !error">
    	<div class="post-list-outline">
    		<h2 class="label">Profile Settings</h2>
        <form class="setting-form" v-on:submit.prevent="sendForm()">
        	<ul class="settings-list">
            <li class="setting-username">
        			<p class="settings-label">Username</p>
        			<div class="center center-input">
        				<input type="text" placeholder="Username" maxlength="64" v-model="name">
        			</div>
        			<p class="note">You can edit your User ID here. This will break previous profile URLs, and I'll be restricting this to certain users once the full version is out.</p>
        		</li>
        		<li class="setting-nickname">
        			<p class="settings-label">Nickname</p>
        			<div class="center center-input">
        				<input type="text" placeholder="Nickname" maxlength="64" v-model="nick">
        			</div>
        			<p class="note">The Mii name of your account. Max length of 32 characters.</p>
        		</li>
        		<li class="setting-avatar">
              <p class="settings-label">Avatar</p>
              <div class="center center-input">
                <input maxlength="255" type="text" placeholder="Avatar" v-model="avatar">
              </div>
        			<p class="note">This should be either a URL link to your avatar or the letters and numbers of a Mii URL (like "3oyie3zp1pdev"). or a <a href="https://www.pf2m.com/tools/mii/">Nintendo Account Mii ID</a></p>
        		</li>
          </ul>
        	<div class="form-buttons">
        		<input class="black-button apply-button" type="submit" v-show="!formSending" value="Save Settings">
            <input class="black-button apply-button disabled" disabled type="submit" v-show="formSending" value="Please wait...">
        	</div>
        </form>
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
export default {
  computed: {
    utils() {
      return window.utils;
    }
  },
  created() {
    //if(!this.profile) {
      this.$http.get('/app/settings/profile').then(response => {
        this.$set(this, 'loading', false);
        this.$set(this, 'user', Object.assign(response.body, window.user));
        this.$set(this, 'name', this.user.name);
        this.$set(this, 'nick', this.user.nick);
        this.$set(this, 'avatar', this.user.avatar);
        this.$set(this, 'profile', response.body);
        utils.title('Profile Settings');
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
      profile: null,
      user: null,
      name: '',
      nick: '',
      avatar: '',
      formSending: false,
      dialog: null
    }
  },
  methods: {
    sendForm() {
      var error;
      if(this.name == '') {
        error = 'The username field cannot be empty.';
        this.$set(this, 'dialog', {
          title: 'Error',
          text: error
        });
        return;
      }
      if(this.nick == '') {
        error = 'The nickname field cannot be empty.';
        this.$set(this, 'dialog', {
          title: 'Error',
          text: error
        });
        return;
      }
      if(this.avatar == '') {
        error = 'The avatar field cannot be empty.';
        this.$set(this, 'dialog', {
          title: 'Error',
          text: error
        });
        return;
      }
      if(!/^[^/]{2,32}$/.test(this.name)) {
        error = 'The username you entered contains an invalid character (forward slash), or is too long or short.';
        this.$set(this, 'dialog', {
          title: 'Error',
          text: error
        });
        return;
      }
      if(this.nick.length > 64) {
        error = 'The nickname you entered is too long.';
        this.$set(this, 'dialog', {
          title: 'Error',
          text: error
        });
        return;
      }
      this.$set(this, 'formSending', true);
      this.$http.post('/app/settings/profile', {
        name: this.name,
        nick: this.nick,
        avatar: this.avatar
      }).then(response => {
        this.$set(this, 'formSending', false);
        this.$set(this.user, 'name', this.name);
        this.$set(this.user, 'nick', this.nick);
        this.$set(this.user, 'avatar', this.avatar);
        this.$set(this.$parent.user, 'name', this.name);
        this.$set(this.$parent.user, 'nick', this.nick);
        this.$set(this.$parent.user, 'avatar', this.avatar);
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

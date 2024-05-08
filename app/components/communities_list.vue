<template>
  <div>
    <div class="welcome">
      <span class="y" style="font-family:serif;text-transform:uppercase;position:absolute;text-shadow:4px 5px rgba(0,0,0,.2);letter-spacing:10px;font-size:12px;text-orientation:inherit;top:60px">welcome to</span> <span class="yy" style="text-transform:uppercase;position:absolute;text-shadow:3.5px 5px rgba(0,0,0,.2);letter-spacing:3px;font-size:25px;top:45px;left:35vw;font-family:sans-serif;vertical-align:middle;font-weight:700;transform:scale(1.2,1)"><span class="ye" style="position:absolute;font-size:50px;transform:scale(1.4,1);font-family:Arial">O</span><span class="yb" style=position:absolute;top:10px;left:47px;font-family:Arial>penverse<span class="yx" style=vertical-align:super;font-family:Arial>&sup2;</span></span></span><div v-show="loggedIn"><router-link to="/logout"><img src="https://kek.gg/i/4jMC_p.png" style="max-height:200px"><span class="sgi-big-spaced-shadowed">log out...</span></router-link></div><div v-show="!loggedIn"><img src="https://kek.gg/i/35pxQQ.png" style="max-width:200px"><router-link to="/login"><span class="sgi-big-spaced-shadowed">log in...</span></router-link></div>
    </div>
    <!-- eeeee -->
    <div style="margin-top:55px"></div>
    <div class="post-list-outline" v-show="loading">
      <div class="no-content"><p>Now loading...</p></div>
    </div>
    <div class="post-list-outline" v-show="!loading && error">
      <div class="no-content"><p>{{ error }}</p></div>
    </div>
    <div v-if="!loading && !error">
    	<h3 class="community-title">General Communities</h3>
    	<ul class="list community-list community-card-list device-new-community-list">
            <div class="post-list-outline" v-if="communities.general.length == 0">
              <div class="no-content"><p>No communities of this type have been created yet.</p></div>
            </div>
            <li class="trigger" v-for="community in communities.general" v-on:click="goToCommunity(community.id)">
    			<div class="community-list-body">
    				<span class="icon-container"><img class="icon" v-bind:src="community.icon"></span>
    				<div class="body">
    					<router-link class="title" :to="'/communities/' + community.id">{{ community.name }}</router-link><span class="platform-tag" v-if="community.platform < 3"><img v-bind:src="'/static/img/platform-tag-' + getPlatformTag(community.platform) + '.png'"></span><span class="text">General Community</span>
    				</div>
    			</div>
    		</li>
       </ul>
       <h3 class="community-title">Game Communities</h3>
         <ul class="list community-list community-card-list device-new-community-list">
           <div class="post-list-outline" v-if="communities.game.length == 0">
             <div class="no-content"><p>No communities of this type have been created yet.</p></div>
           </div>
            <li class="trigger" v-for="community in communities.game" v-on:click="goToCommunity(community.id)">
    			<div class="community-list-body">
    				<span class="icon-container"><img class="icon" v-bind:src="community.icon"></span>
    				<div class="body">
    					<router-link class="title" :to="'/communities/' + community.id">{{ community.name }}</router-link><span class="platform-tag" v-if="community.platform < 3"><img v-bind:src="'/static/img/platform-tag-' + getPlatformTag(community.platform) + '.png'"></span><span class="text">{{ getPlatformText(community.platform) || 'General Community' }}</span>
    				</div>
    			</div>
    		</li>
       </ul>
       <h3 class="community-title">Special Communities</h3>
         <ul class="list community-list community-card-list device-new-community-list">
            <div class="post-list-outline" v-if="communities.special.length == 0">
               <div class="no-content"><p>No communities of this type have been created yet.</p></div>
            </div>
            <li class="trigger" v-for="community in communities.special" v-on:click="goToCommunity(community.id)">
    			<div class="community-list-body">
    				<span class="icon-container"><img class="icon" v-bind:src="community.icon"></span>
    				<div class="body">
    					<router-link class="title" :to="'/communities/' + community.id">{{ community.name }}</router-link><span class="platform-tag" v-if="community.platform < 3"><img v-bind:src="'/static/img/platform-tag-' + getPlatformTag(community.platform) + '.png'"></span><span class="text">Special Community</span>
    				</div>
    			</div>
    		</li>
       </ul>
    </div>
  </div>
</template>
<script>
export default {
  computed: {
    loggedIn() {
      return window.user.id > 0;
    }
  },
  data() {
    return {
      loading: true,
      error: '',
      communities: {general: [], game: [], special: []}
    }
  },
  created() {
    //if(!this.community) {
      this.$http.get('/app/communities').then(response => {
        this.$set(this, 'loading', false);
        this.$set(this, 'communities', response.body);
        utils.title('All Communities');
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
  methods: {
    /*getSaveStateConfig() {
			return {
				cacheKey: 'CommunityList'
			}
		},*/
      getPlatformTag(platform) {
            switch(platform) {
                 case 0:
                    return '3ds';
                    break;
                 case 1:
                    return 'wiiu';
                    break;
                 case 2:
                    return 'wiiu-3ds';
                    break;
            }
       },
       getPlatformText(platform) {
            switch(platform) {
                 case 0:
                    return '3DS Games';
                    break;
                 case 1:
                    return 'Wii U Games';
                    break;
                 case 2:
                    return 'Wii U Games\u30FB3DS Games';
                    break;
            }
       },
       goToCommunity(id) {
            this.$router.push('/communities/' + id);
       }
  }
}
</script>

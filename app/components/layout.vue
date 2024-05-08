<template>
  <div id="wrapper" v-bind:class="{guest: !loggedIn}" v-on:click="dismissMenu($event)">
    <div id="sub-body">
      <menu id="global-menu">
        <ul>
          <li id="global-menu-logo">
            <h1><router-link to="/"><img alt="Openverse" src="/static/img/menu-logo.png"></router-link></h1>
          </li>
          <li id="global-menu-list" v-if="loggedIn">
            <ul>
              <li id="global-menu-mymenu" v-bind:class="{selected: selected == 'mymenu'}">
                <router-link :to="'/users/' + user.name"><span class="icon-container" v-bind:class="utils.doRank(user.rank)"><img alt="User Page" v-bind:src="utils.doAvatar(user.avatar)"></span><span>User Page</span></router-link>
              </li>
              <li id="global-menu-feed" v-bind:class="{selected: selected == 'feed'}">
                <router-link class="symbol" to="/activity"><span>Activity Feed</span></router-link>
              </li>
              <li id="global-menu-community" v-bind:class="{selected: selected == 'community'}">
                <router-link class="symbol" to="/"><span>Communities</span></router-link>
              </li>
              <li id="global-menu-news" v-bind:class="{selected: selected == 'news'}">
                <router-link class="symbol" to="/news/my_news">
                  <p class="badge" v-show="unread">{{ unread }}</p>
                </router-link>
              </li>
              <li id="global-menu-my-menu">
                <button v-on:click="toggleMenu()" class="symbol js-open-global-my-menu open-global-my-menu" id="my-menu-btn"></button>
                <menu v-bind:class="{invisible: !menuOpen, none: !menuOpen}" id="global-my-menu">
                  <ul>
                    <li>
                      <router-link class="symbol my-menu-profile-setting" to="/settings/profile"><span>Profile Settings</span></router-link>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-miiverse-setting" to="/settings/account"><span>Account Settings</span></router-link>
                    </li>
                    <li>
                      <a class="symbol my-menu-info" href="javascript:alert('Not%20really.')"><span>Openverse Announcements</span></a>
                    </li>
                    <li>
                      <a class="symbol my-menu-info" href="javascript:alert('Not%20really.')"><span>Openverse Changelog</span></a>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-guide" to="/guide/"><span>Openverse Code of Conduct</span></router-link>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-guide" to="/guide/legal"><span>Legal Stuff</span></router-link>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-guide" to="/guide/faq"><span>Frequently Asked Questions (FAQ)</span></router-link>
                    </li>
                    <li class="dark-toggle">
                      <a class="symbol my-menu-info dark-toggle" href="#" v-on:click.prevent="toggleDark()"><span class="dark-toggle">Toggle Dark Mode</span></a>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-miiverse-setting" to="/admin/"><span>Admin Tools</span></router-link>
                    </li>
                    <li>
                      <router-link class="symbol my-menu-logout" to="/logout"><span>Log Out</span></router-link>
                    </li>
                  </ul>
                </menu>
              </li>
            </ul>
          </li>
          <li id="global-menu-login" v-if="!loggedIn">
            <router-link class="login" to="/login"><input alt="Sign in" src="/static/img/sign-in.png" type="image"></router-link>
          </li>
        </ul>
      </menu>
    </div>
    <div id="main-body" v-bind:class="{guest: !loggedIn}">
      <router-view :key="$route.fullPath"></router-view>
    </div>
    <div id="footer">
      <div id="footer-inner">
        <div class="link-container">
          <p><a href="javascript:alert('There%20is%20no%20Github%20yet.%20Please%20understand.');">Github</a></p>
          <p><router-link to="/guide/contact">Contact Us</router-link></p>
          <p><a href="https://www.paypal.me/PF2M">Donate</a></p>
          <p id="copyright">Openverse is not-for-profit and is not associated with Miiverse, Nintendo, or Hatena.<br>
          Please support the products of these companies, as they deserve your money!</p>
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
  data() {
    return {
      loggedIn: Boolean(user.id),
      menuOpen: false,
      selected: '',
      unread: 0,
      messageStream: null
    }
  },
  watch: {
    $route(to, from) {
      this.setSelected(to);
    }
  },
  mounted() {
    this.setSelected(this.$route);
    if(this.user.id > 0) {
      this.unread = window.unread;
      this.messageStream = new WebSocket('ws' + ((location.protocol == 'https:') ? 's' : '') + '://' + location.host + '/app/news/my_news/stream');
      var _this = this;
  		this.messageStream.onmessage = function(event) {
        var data = JSON.parse(event.data);
        if(data == 0) {
          _this.$set(_this, 'unread', 0);
        } else {
          var hasUnread = true;
          if(_this.$route.params.id) {
            switch(data.type) {
              case 0:
                if(_this.$route.params.id == data.topic) {
                  hasUnread = false;
                }
                break;
              case 1:
                if(_this.$route.params.id == data.topic || _this.$route.params.id == data.other_topic) {
                  hasUnread = false;
                }
                break;
              case 2:
                if(_this.$route.params.id == data.topic) {
                  hasUnread = false;
                }
                break;
              case 3:
                if(_this.$route.params.id == data.topic) {
                  hasUnread = false;
                }
                break;
            }
          }
          if(hasUnread) {
            _this.$set(_this, 'unread', (_this.unread || 0) + 1);
          }
        }
      }
    }
  },
  destroyed() {
    if(this.messageStream) {
      this.messageStream.close();
    }
  },
  methods: {
    setSelected(route) {
      var select = '';
      switch(route.name) {
        case 'communities_list':
          select = 'community';
          break;
        case 'community_posts':
          select = 'community';
          break;
        case 'title_view_redirect':
          select = 'community';
          break;
        case 'notifications_view':
          select = 'news';
          this.unread = 0;
          break;
        case 'user_view':
          if(this.user.id > 0 && route.params.id == this.user.name) {
            select = 'mymenu';
          }
          break;
      }
      this.$set(this, 'selected', select);
    },
    toggleMenu() {
      this.$set(this, 'menuOpen', !this.menuOpen);
    },
    dismissMenu(event) {
      if(!this.menuOpen) {
        return;
      }
      if(event.target.classList.contains('dark-toggle')) {
        return;
      }
      if(event.target.id != 'my-menu-btn') {
        this.$set(this, 'menuOpen', false);
      }
    },
    closeDialog() {
      this.$set(window, 'dialog', null);
    },
    toggleDark() {
      // Get cookie
      var hasDarkCookie = document.cookie.indexOf('openverse_dark=') >= 0;
      if(hasDarkCookie) {
        // Disable dark mode
        document.getElementById('dark-mode').disabled = true;
        document.cookie = 'openverse_dark=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT';
      } else {
        // Enable dark mode
        document.getElementById('dark-mode').disabled = false;
        document.cookie = 'openverse_dark=1; path=/';
      }
    }
  }
}
</script>

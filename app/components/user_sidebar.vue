<template>
	<div v-bind:class="{'general-sidebar': general, 'user-sidebar': !general}" id="sidebar">
		<div class="sidebar-container">
			<router-link id="sidebar-cover" class="sidebar-cover-image" :to="'/posts/' + user.favorite_post_id" v-if="user.favorite_post"><img v-bind:src="user.favorite_post.screenshot"></router-link>
			<div id="sidebar-profile-body" v-bind:class="{'with-profile-post-image': user.favorite_post}">
				<div class="icon-container" v-bind:class="utils.doRank(user.rank)">
					<router-link :to="'/users/' + user.name"><img class="icon" v-bind:src="utils.doAvatar(user.avatar)"></router-link>
				</div>
				<p class="user-organization" v-if="utils.doRankText(user.rank)">{{ utils.doRankText(user.rank) }}</p>
				<router-link class="nick-name" :to="'/users/' + user.name">{{ user.nick }}</router-link>
				<p class="id-name">{{ user.name }}</p>
			</div>
			<div id="edit-profile-settings" v-if="!general && meUserId == user.id">
				<router-link to="/settings/profile" class="button symbol">Profile Settings</router-link>
			</div>
			<div class="user-action-content" v-if="!general && meUserId != user.id">
				<div class="toggle-button">
					<button class="follow-button button symbol" type="button">foo</button>
				</div>
			</div>
			<ul id="sidebar-profile-status">
				<li>
					<a><span class="number">0</span>bar</a>
				</li>
				<li>
					<a><span class="number">0</span>baz</a>
				</li>
				<li>
					<a><span class="number">0</span>qux</a>
				</li>
			</ul>
		</div>
		<div class="sidebar-container sidebar-setting" v-if="!general">
			<div class="sidebar-post-menu">
				<a class="sidebar-menu-post with-count symbol">
					<span>All Posts</span><span class="post-count"><span class="test-post-count">0</span></span>
				</a>
				<a class="sidebar-menu-empathies with-count symbol">
					<span>Yeahs</span><span class="post-count"><span class="test-empathy-count">0</span></span>
				</a>
			</div>
		</div>
		<div class="sidebar-container sidebar-profile" v-if="!general">
			<div class="profile-comment" v-show="user.profile_comment">
				<p>{{ user.profile_comment }}</p>
			</div>
			<div class="user-data">
				<div class="data-content">
					<h4><span>Region</span></h4>
					<div class="note">
						<span v-show="user.country">{{ user.country }}</span>
						<span v-show="!user.country">Not Set</span>
					</div>
				</div>
				<div class="data-content">
					<h4><span>NNID</span></h4>
					<div class="note">
						<span v-show="user.nnid">{{ user.nnid }}</span>
						<span v-show="!user.nnid">Not Set</span>
					</div>
				</div>
				<div class="data-content game-skill">
					<h4><span>Game Experience</span></h4>
					<div class="note">
						<span class="expert" v-show="user.skill == 3">Expert</span>
						<span class="intermediate" v-show="user.skill == 2">Intermediate</span>
						<span class="beginner" v-show="user.skill == 1">Beginner</span>
						<span v-show="!user.skill">Not Set</span>
					</div>
				</div>
				<div class="data-content">
					<h4><span>Member ID</span></h4>
					<div class="note">
						<span>#{{ user.id }}</span>
					</div>
				</div>
				<div class="data-content">
					<h4><span>Website</span></h4>
					<div class="note">
						<span v-show="user.website"><a v-bind:href="user.website">{{ user.website }}</a></span>
						<span v-show="!user.website">Not Set</span>
					</div>
				</div>
			</div>
		</div>
		<div style="text-align:center" v-if="general"><span>hello, {{ user.nick.charAt(0).toUpperCase() + user.nick.slice(1) }}!<br>implement sidebars when user pages are implemented</span></div>
		<div style="text-align:center" v-else><span>THIS IS A USER PAGE. Welcome to {{ user.nick.charAt(0).toUpperCase() + user.nick.slice(1) }}!'s user page!</span></div>
	</div>
</template>
<script>
export default {
	name: 'user-sidebar',
	computed: {
		meUserId() {
			return window.user.id;
		},
		utils() {
			return window.utils;
		}
	},
	props: {
		user: Object,
		general: {
			type: Boolean,
			default: true
		}
	}
}
</script>

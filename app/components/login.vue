<template>
<div class="main-column center">
	<div class="post-list-outline auth-page">
		<form method="post" v-on:submit.prevent="submitForm()">
			<img src="/static/img/menu-logo.png">
			<h1>Log In</h1>
			<p>Sign in with an Openverse account to create and give Yeahs to posts and comments, edit your profile, and connect with other Openverse users.</p>
			<h3 class="label"><label>Username <input class="auth-input" v-bind:class="{bad: errorFields.includes('name')}" maxlength="32" placeholder="Username" type="text" v-model="name"></label></h3>
			<h3 class="label"><label>Password <input class="auth-input" v-bind:class="{bad: errorFields.includes('pass')}" placeholder="Password" type="password" v-model="pass"></label></h3>
			<p class="error" v-show="error">{{ error }}</p>
			<button class="big-button" type="submit" v-bind:disabled="formSending" v-bind:class="{disabled: formSending}">Sign In</button>
			<footer>
				<p>Don't have an account? <router-link to="/signup">Click here to make one.</router-link></p>
				<p>Forgot your password? <router-link to="/forgot">Reset it here.</router-link></p>
			</footer>
		</form>
	</div>
</div>
</template>
<script>
export default {
	data() {
		return {
			name: '',
			pass: '',
			error: '',
			errorFields: [],
			formSending: false
		}
	},
	created() {
		window.utils.title('Log In');
	},
	methods: {
		submitForm() {
			this.errorFields = [];
			if(this.name == '') {
				this.errorFields.push('name');
				this.error = 'The username field cannot be empty.';
				return;
			}
			if(this.pass == '') {
				this.errorFields.push('pass');
				this.error = 'The password field cannot be empty.';
				return;
			}
			this.$set(this, 'formSending', true);
			this.error = '';
			this.$http.post('/app/login', {
				name: this.name,
				pass: this.pass
			}).then(response => {
				console.log(response.body)
				location.href = '/';
			}, response => {
				this.$set(this, 'formSending', false);
				if(response.body.error) {
					if(this.error == 'Invalid password.') {
						this.errorFields.push('pass');
					}
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
	}
}
</script>

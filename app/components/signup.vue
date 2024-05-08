<template>
<div class="main-column center">
	<div class="post-list-outline auth-page">
		<form method="post" v-on:submit.prevent="submitForm()" v-on:paste="pasteImage($event)">
			<img src="/static/img/menu-logo.png">
			<h1>Sign Up</h1>
			<p>Create an Openverse account to make posts and comments to various communities, give Yeahs to other users' content, and interact with other members of the Openverse community.</p>
			<h3 class="label"><span class="red">*</span><label>Username <span class="note">(This is what you will use to sign in.)</span> <input class="auth-input" v-bind:class="{bad: errorFields.includes('name')}" maxlength="32" placeholder="Username" type="text" v-model="name"></label></h3>
			<h3 class="label"><span class="red">*</span><label>Nickname <span class="note">(This is the name displayed beside your username.)</span> <input class="auth-input" v-bind:class="{bad: errorFields.includes('nick')}" maxlength="64" placeholder="Nickname" type="text" v-model="nick"></label></h3>
			<h3 class="label"><label>Nintendo Network ID <input class="auth-input" v-bind:class="{bad: errorFields.includes('nnid')}" maxlength="16" placeholder="NNID" type="text" v-model="nnid"></label></h3>
			<h3 class="label"><label>Email Address <input class="auth-input" v-bind:class="{bad: errorFields.includes('email')}" placeholder="Email" type="email" v-model="email"></label></h3>
			<h3 class="label">Custom Avatar <span class="note">(This will be scaled down to 128x128.)</span></h3>
			<label class="file-button-container" v-on:drop.prevent="pasteImage($event)" v-on:dragover.prevent="">
				<span class="button file-upload-button">Upload</span>
				<input accept="image/*" class="file-button none" type="file" v-on:change="imageUpload($event)">
				<div class="screenshot-container still-image" v-show="image">
					<img v-bind:src="image">
				</div>
			</label>
			<h3 class="label"><span class="red">*</span><label>Password <input class="auth-input" v-bind:class="{bad: errorFields.includes('pass')}" placeholder="Password" type="password" v-model="pass"></label></h3>
			<h3 class="label"><span class="red">*</span><label>Confirm Password <input class="auth-input" v-bind:class="{bad: errorFields.includes('passAgain')}" placeholder="Confirm Password" type="password" v-model="passAgain"></label></h3>

			<p class="error" v-show="error">{{ error }}</p>

			<button class="big-button" type="submit" v-bind:disabled="formSending" v-bind:class="{disabled: formSending}">Create Account</button>
			<footer>
				<p>All fields with a red asterisk (<span class="red">*</span>) are required.</p>
				<p>Nintendo Network IDs are optionally required if you want to get your Mii, and can be hidden from the public.</p>
				<p>You don't have to enter an email address, but it's required to reset your password.</p>
				<p>If you don't submit an avatar, a Gravatar will be grabbed from your email address, or your Mii will be grabbed from your NNID.</p>
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
			nick: '',
			nnid: '',
			email: '',
			pass: '',
			passAgain: '',
			image: '',
			imageDimensions: '',
			error: '',
			errorFields: [],
			formSending: false
		}
	},
	created() {
		window.utils.title('Sign Up');
	},
	methods: {
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
						imageDimension.onload = function() {
							_this.$set(_this, 'imageDimensions', imageDimension.width + ' x ' + imageDimension.height);
						}
					}
				}
			});
		},
		submitForm() {
			this.errorFields = [];
			if(this.name == '') {
				this.errorFields.push('name');
				this.error = 'The username field cannot be empty.';
				return;
			}
			if(this.nick == '' && this.nnid == '') {
				this.errorFields.push('nick');
				this.error = 'The nickname field cannot be empty.';
				return;
			}
			if(this.pass == '') {
				this.errorFields.push('pass');
				this.error = 'The password field cannot be empty.';
				return;
			}
			if(this.pass != this.passAgain) {
				this.errorFields.push('pass');
				this.errorFields.push('passAgain');
				this.error = 'The passwords you entered do not match.';
				return;
			}
			if(!/^[^/]{2,32}$/.test(this.name)) {
				this.errorFields.push('name');
				this.error = 'The username you entered contains an invalid character (forward slash), or is too long or short.';
				return;
			}
			if(this.nick.length > 64) {
				this.errorFields.push('nick');
				this.error = 'The nickname you entered is too long.';
				return;
			}
			if(this.nnid && !/^[A-Za-z0-9-._]{6,16}$/.test(this.nnid)) {
				this.errorFields.push('nnid');
				this.error = 'The Nintendo Network ID you entered is invalid.';
				return;
			}
			if(this.email && !/^[A-Za-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$/.test(this.email)) {
				this.errorFields.push('email');
				this.error = 'The email address you entered is invalid.';
				return;
			}
			this.error = '';
			this.$set(this, 'formSending', true);
			this.$http.post('/app/signup', {
				name: this.name,
				nick: this.nick,
				nnid: this.nnid,
				email: this.email,
				avatar: this.image,
				pass: this.pass,
				pass_again: this.passAgain
			}).then(response => {
				location.href = '/';
			}, response => {
				this.$set(this, 'formSending', false);
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
	}
}
</script>

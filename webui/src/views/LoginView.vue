<script>
export default {
	data: function () {
		return {
			errormsg: null,
			username: "",
			remember: false,
		}
	},
	methods: {
		async login() {
			this.errormsg = null;
			try {
				let response = await this.$axios.post("/session", {username: this.username.trim()});
				let token = response.data.success.split('ID: ')[1];

				localStorage.setItem('username', this.username);
				localStorage.setItem('token', token);
				localStorage.setItem('remember', this.remember);
				this.$router.replace("/home")
				this.$emit('updatedLoggedChild', true)

			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		toLowerCase(event) {
			this.username = this.username.toLowerCase();
		},
	},

	mounted() {
		if (localStorage.getItem('remember') == "true") {
			this.$router.replace("/home")
		}

		else {
			localStorage.clear()
      		this.$emit('logoutNavbar',false) 
		}

	}
}
</script>

<template>
	<div class="row my-auto" style="height: 100vh;">

		<div class="row">
			<div class="col">
				<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			</div>
		</div>

		<div class="row my-auto h-100 w-100">

			<form @submit.prevent="login" class="d-flex flex-column align-items-center justify-content-center p-0">

				<div class="row mx-auto my-auto">
					<div class="col my-auto">
						<img src="../assets/WasaPhoto_Logo.png" class="img-fluid" alt="Responsive image">
					</div>
					<div class="col my-auto">
						<div class="row mx-auto p-1" style="width: 300px;">
							<input type="text" class="form-control" v-model="username" @input="toLowerCase"
								maxlength="16" minlength="3" placeholder="Username" />
						</div>
						<div class="row mx-auto p-1" style="width: 300px;">
							<button class="btn btn-primary"
								:disabled="username == null || username.length > 16 || username.length < 3 || username.trim().length < 3 || username.split(' ').length - 1 > 0">
								Register/Login
							</button>
						</div>
						<div class="row mx-auto">
							<table cellspacing="0" cellpadding="0">
								<tr>
									<td>
										<input type="checkbox" v-model="remember"/> Remember me
									</td>
								</tr>
							</table>
						</div>
					</div>
				</div>
			</form>
		</div>
	</div>
</template>

<style>

* {
	font-family: Verdana, Geneva, Tahoma, sans-serif;
}

.login {
	height: 100vh;
}

[v-cloak] {
	display: none
}

td {
  text-align: center;
  vertical-align: middle;
}

table {
  height: 20px;

}

</style>
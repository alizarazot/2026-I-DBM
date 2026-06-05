<script lang="ts">
	const handleSignIn = async (e: SubmitEvent) => {
		e.preventDefault();

		const target = e.target as HTMLFormElement;

		const formData = new FormData(target);
		const data = Object.fromEntries(formData);

		try {
			const res = await fetch('/api/auth', {
				method: target.method,
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			});

			console.log(await res.json());
		} catch (err) {
			console.error(err);
		}
	};
</script>

<form method="POST" onsubmit={handleSignIn}>
	<input name="email" type="email" placeholder="Email" />
	<input name="password" type="password" placeholder="Password" />
	<button type="submit">Sign In</button>
</form>

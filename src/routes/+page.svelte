<script>
import { goto } from "$app/navigation";
import { userState } from "$lib/user.svelte";

$effect(() => {
	if (userState.userType == null) {
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto("/auth");
	} else {
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`/${userState.userType}`);
	}
});
</script>

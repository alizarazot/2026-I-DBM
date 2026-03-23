<script>
import { goto } from "$app/navigation";
import { userState } from "$lib/user.svelte.ts";

$effect((_) => {
	if (userState.userType == null) {
		goto("/auth");
	} else {
		goto(`/${userState.userType}`);
	}
});
</script>

<script lang="ts">
import { enhance } from "$app/forms";
import type { PageServerData } from "./$types";

let { data } = $props<{ data: PageServerData }>();
</script>

Hi, {data.user.name}


<form method="post" action="/auth?/signOut" use:enhance>
  <button> Sign Out </button>
</form>

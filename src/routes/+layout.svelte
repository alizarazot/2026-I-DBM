<script lang="ts">
	import './layout.css';

	let { children } = $props();
</script>

<div class="app"><main class="h-dvh">{@render children()}</main></div>

<script lang="ts">
	let { id, title, description, status, dueDateTime } = $props();
</script>

<p>{id}{title}{description}{status}{dueDateTime}</p>

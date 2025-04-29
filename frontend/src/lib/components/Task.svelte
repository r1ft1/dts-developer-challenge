<script lang="ts">
	let { id, title, description, status, dueDateTime } = $props();
</script>

<div>{id}</div>
<div>{title}</div>
<div>{description}</div>
<div>{status}</div>
<div>{dueDateTime}</div>

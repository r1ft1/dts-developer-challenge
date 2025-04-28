<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import Task from "$lib/components/Task.svelte";
	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();

	console.log(data);

	let titleInput = $state("");
	let descInput = $state("");
	let statusInput = $state("");
	let dateInput = $state("");

	const createTaskButtonHandler = async () => {
		const form = new FormData();

		// the go api wants form data in this format
		/*
		title := r.FormValue("title")
		description := r.FormValue("description")
		status := r.FormValue("status")
		dueDateTime := r.FormValue("due_date_time")
		*/
		form.append("title", titleInput);
		form.append("description", descInput);
		form.append("status", statusInput);
		form.append("due_date_time", dateInput);

		const createResponse = await fetch("http://localhost:8080/create", {
			method: "POST",
			body: form,
		});

		const data = await createResponse.json();
		console.log(data);

		if (createResponse.ok) {
			console.log("FE: Task created successfully");
		} else {
			console.error("FE: Error creating task:", data);
		}

		//Refreshes the page which rerequests the index route, which reads all tasks from the DB
		invalidateAll();
	};
</script>

<h1>HMCTS Caseworker Task Management</h1>

<div class="taskForm">
	<input type="text" placeholder="Title" bind:value={titleInput} />
	<input type="text" placeholder="Description" bind:value={descInput} />
	<input type="text" placeholder="Status" bind:value={statusInput} />
	<input
		type="datetime-local"
		placeholder="Due Date Time"
		bind:value={dateInput}
	/>
	<button onclick={createTaskButtonHandler}>Create Task</button>
</div>

{#each data.tasks as task}
	<Task
		id={task.id}
		title={task.title}
		description={task.description}
		status={task.status}
		dueDateTime={task.due_date_time}
	/>
{/each}

<style>
	.taskForm {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 300px;
	}
</style>

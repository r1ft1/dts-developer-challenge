<script lang="ts">
	import { invalidateAll } from "$app/navigation";
	import Task from "$lib/components/Task.svelte";
	import type { PageProps } from "./$types";
	import { PUBLIC_API_URL } from "$env/static/public";

	let { data }: PageProps = $props();
	const API_URL = PUBLIC_API_URL;

	console.log(data);

	let titleInput = $state("");
	let descInput = $state("");
	let statusInput = $state("");
	let dateInput = $state("");

	const createTaskButtonHandler = async () => {
		// validation of inputs to check if title, status and date are not empty (description optional)
		if (!titleInput || !statusInput || !dateInput) {
			console.error("FE: Title, Status and Due Date Time are required");
			alert("Title, Status and Due Date Time are required");
			return;
		}
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

		const createResponse = await fetch(`${API_URL}/create`, {
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

	const deleteTaskButtonHandler = async (id: number) => {
		const deleteResponse = await fetch(`${API_URL}/delete?id=${id}`, {
			method: "DELETE",
		});

		const data = await deleteResponse.json();
		console.log(data);

		if (deleteResponse.ok) {
			console.log("FE: Task deleted successfully");
		} else {
			console.error("FE: Error deleting task:", data);
		}

		invalidateAll();
	};

	let editing = $state(0);
	let editStatusInput = $state("");

	const editTaskButtonHandler = (id: number) => {
		editing = id;
	};

	const updateTaskButtonHandler = async (id: number) => {
		console.log("update clicked");
		const updateResponse = await fetch(
			`${API_URL}/update?id=${id}&status=${editStatusInput}`,
			{
				method: "PUT",
			},
		);
		const data = await updateResponse.json();
		console.log(data);

		if (updateResponse.ok) {
			console.log("FE: Task updated successfully");
		} else {
			console.error("FE: Error updated task:", data);
		}

		editing = 0;

		invalidateAll();
	};
</script>

<h1>HMCTS Caseworker Task Management</h1>

<div class="taskForm">
	<input type="text" placeholder="Title" bind:value={titleInput} />
	<input
		type="text"
		placeholder="Description (OPTIONAL)"
		bind:value={descInput}
	/>
	<input type="text" placeholder="Status" bind:value={statusInput} />
	<input
		type="datetime-local"
		placeholder="Due Date Time"
		bind:value={dateInput}
	/>

	<button onclick={createTaskButtonHandler}>Create Task</button>
</div>

{#each data.tasks as task (task.id)}
	{#if editing != task.id}
		<div class="task">
			<Task
				id={task.id}
				title={task.title}
				description={task.description}
				status={task.status}
				dueDateTime={task.due_date_time}
			/>
			<button onclick={() => editTaskButtonHandler(task.id)}>
				Edit Task
			</button>
			<button onclick={() => deleteTaskButtonHandler(task.id)}>
				Delete Task
			</button>
		</div>
	{:else}
		<div class="task">
			<div class="column">
				{task.id}
			</div>
			<div class="column">
				{task.title}
			</div>
			<div class="column">
				{task.description}
			</div>
			<div class="column">
				<input
					type="text"
					placeholder={task.status}
					bind:value={editStatusInput}
				/>
			</div>
			<div class="column">
				{task.due_date_time}
			</div>
			<button onclick={() => updateTaskButtonHandler(task.id)}>
				Update Task
			</button>
			<button onclick={() => (editing = 0)}>Cancel</button>
		</div>
	{/if}
{/each}

<style>
	.taskForm {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		width: 300px;
	}
	.task {
		display: flex;
		flex-direction: row;
		gap: 1rem;
		width: 10px;
	}
</style>

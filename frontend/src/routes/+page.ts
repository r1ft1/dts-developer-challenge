import type { PageLoad } from './$types';
import { PUBLIC_API_URL } from "$env/static/public";

type APIResponse = {
	message: string;
	tasks: Task[];
}

type Task = {
	id: number;
	title: string;
	description: string;
	status: string;
	due_date_time: string;
};

let apiUrl = PUBLIC_API_URL;

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`${apiUrl}`);
	const data: APIResponse = await res.json();

	console.log(data);

	return data;

};

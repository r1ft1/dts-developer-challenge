import type { PageLoad } from './$types';
import { dev } from "$app/environment";


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

let apiUrl = "http://localhost:8080/";

if (!dev) {
	apiUrl = "https://backend-production-93b4.up.railway.app";
}
export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`${apiUrl}`);
	const data: APIResponse = await res.json();

	console.log(data);

	return data;

};

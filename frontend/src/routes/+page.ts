import type { PageLoad } from './$types';


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

const apiUrl = process.env.API_URL || "http://localhost:8080/";
export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`${apiUrl}`);
	const data: APIResponse = await res.json();

	console.log(data);

	return data;

};

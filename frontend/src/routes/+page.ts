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

export const load: PageLoad = async ({ fetch }) => {
	const res = await fetch(`http://localhost:8080/`);
	const data: APIResponse = await res.json();

	console.log(data);

	return data;

};

import { ref } from "vue";
import axios from "axios";

export function fetcher(url: string) {
	const data = ref(null);
	const err = ref<unknown | null>(null);
	const isLoading = ref(true);

	const fetchData = async () => {
		isLoading.value = true;

		try {
			const res = await axios.get(url);
			data.value = res.data;
		} catch (e) {
			err.value = e;
		} finally {
			isLoading.value = false;
		}
	};

	fetchData();

	return { data, err, isLoading, fetchData };
}

export function poster(
	url: string,
	body: { username: string; password: string },
) {
	const data = ref(null);
	const err = ref<unknown | null>(null);
	const isLoading = ref(true);

	const fetchData = async () => {
		isLoading.value = true;

		try {
			const res = await axios.post(url, body);
			data.value = res.data;
		} catch (e) {
			err.value = e;
		} finally {
			isLoading.value = false;
		}
	};

	fetchData();

	return { data, err, isLoading, fetchData };
}

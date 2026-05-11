import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/api";

/**
 * Store quản lý phòng ban và chức vụ
 */
export const useDepartmentStore = defineStore("departments", () => {
	// --- State ---
	const departments = ref([]);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Tải danh sách phòng ban kèm chức vụ
	 */
	async function fetchDepartments() {
		loading.value = true;
		try {
			const res = await api.get("/departments");
			departments.value = res.data || [];
		} catch (error) {
			console.error("[DepartmentStore] fetchDepartments failed:", error);
			departments.value = [];
		} finally {
			loading.value = false;
		}
	}

	/**
	 * Quản lý Phòng ban
	 */
	async function createDepartment(data) {
		return await api.post("/departments", data);
	}

	async function updateDepartment(id, data) {
		return await api.put(`/departments/${id}`, data);
	}

	async function deleteDepartment(id) {
		return await api.delete(`/departments/${id}`);
	}

	/**
	 * Quản lý Chức vụ
	 */
	async function createPosition(data) {
		return await api.post("/positions", data);
	}

	async function updatePosition(id, data) {
		return await api.put(`/positions/${id}`, data);
	}

	async function deletePosition(id) {
		return await api.delete(`/positions/${id}`);
	}

	return {
		departments,
		loading,
		fetchDepartments,
		createDepartment,
		updateDepartment,
		deleteDepartment,
		createPosition,
		updatePosition,
		deletePosition,
	};
});

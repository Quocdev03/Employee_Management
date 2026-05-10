import { defineStore } from "pinia";
import api from "../api";
import { ref } from "vue";

export const useDepartmentStore = defineStore("departments", () => {
	// --- State ---
	const departments = ref([]);

	// --- Actions ---

	/**
	 * Lấy danh sách toàn bộ phòng ban từ hệ thống
	 */
	async function fetchDepartments() {
		try {
			const res = await api.get("/departments");

			// Chuẩn hóa dữ liệu trả về thành mảng
			if (Array.isArray(res.data)) {
				departments.value = res.data;
				return;
			}

			if (Array.isArray(res)) {
				departments.value = res;
				return;
			}

			departments.value = [];
		} catch (error) {
			console.error("Lỗi khi lấy danh sách phòng ban:", error);
			departments.value = [];
		}
	}

	/**
	 * Tạo mới một phòng ban
	 */
	async function createDepartment(data) {
		const res = await api.post("/departments", data);
		return res.data;
	}

	/**
	 * Cập nhật thông tin phòng ban
	 */
	async function updateDepartment(id, data) {
		const res = await api.put(`/departments/${id}`, data);
		return res.data;
	}

	/**
	 * Xoá một phòng ban khỏi hệ thống
	 */
	async function deleteDepartment(id) {
		const res = await api.delete(`/departments/${id}`);
		return res.data;
	}

	/**
	 * Tạo mới một chức vụ trong phòng ban
	 */
	async function createPosition(data) {
		const res = await api.post("/positions", data);
		return res.data;
	}

	/**
	 * Cập nhật tên chức vụ
	 */
	async function updatePosition(id, data) {
		const res = await api.put(`/positions/${id}`, data);
		return res.data;
	}

	/**
	 * Xoá một chức vụ
	 */
	async function deletePosition(id) {
		const res = await api.delete(`/positions/${id}`);
		return res.data;
	}

	return {
		departments,
		fetchDepartments,
		createDepartment,
		updateDepartment,
		deleteDepartment,
		createPosition,
		updatePosition,
		deletePosition,
	};
});

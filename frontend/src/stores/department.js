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

	return {
		departments,
		fetchDepartments,
	};
});


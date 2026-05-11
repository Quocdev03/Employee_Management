import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/api";

/**
 * Store quản lý dữ liệu nhân viên
 */
export const useEmployeeStore = defineStore("employee", () => {
	// --- State ---
	const employees = ref([]);
	const total = ref(0);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Tải danh sách nhân viên theo bộ lọc
	 * @param {Object} params
	 */
	async function fetchEmployees(params = {}) {
		loading.value = true;
		try {
			const res = await api.get("/employees", { params });
			
			// Cấu trúc: { success: true, data: [], total: ... }
			employees.value = res.data || [];
			total.value = res.total || employees.value.length || 0;
		} catch (error) {
			console.error("[EmployeeStore] fetchEmployees failed:", error);
			employees.value = [];
			total.value = 0;
		} finally {
			loading.value = false;
		}
	}

	/**
	 * Tạo nhân viên mới
	 * @param {Object} data
	 */
	async function createEmployee(data) {
		return await api.post("/employees", data);
	}

	/**
	 * Cập nhật thông tin nhân viên
	 * @param {number|string} id
	 * @param {Object} data
	 */
	async function updateEmployee(id, data) {
		return await api.put(`/employees/${id}`, data);
	}

	/**
	 * Xóa nhân viên khỏi hệ thống
	 * @param {number|string} id
	 */
	async function deleteEmployee(id) {
		return await api.delete(`/employees/${id}`);
	}

	return {
		employees,
		total,
		loading,
		fetchEmployees,
		createEmployee,
		updateEmployee,
		deleteEmployee,
	};
});

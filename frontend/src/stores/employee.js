import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/api";

export const useEmployeeStore = defineStore("employee", () => {
	// --- State ---
	const employees = ref([]);
	const total = ref(0);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Lấy danh sách nhân viên với các tham số lọc và phân trang
	 * @param {Object} params - Tham số lọc (page, limit, search, department_id, ...)
	 */
	async function fetchEmployees(params = {}) {
		loading.value = true;

		try {
			const res = await api.get("/employees", { params });

			// Chuẩn hóa dữ liệu trả về thành mảng
			let data = [];
			if (Array.isArray(res.data)) {
				data = res.data;
			} else if (Array.isArray(res)) {
				data = res;
			}

			employees.value = data;
			total.value = res.total || data.length || 0;
		} catch (error) {
			console.error("Lỗi khi lấy danh sách nhân viên:", error);
			employees.value = [];
			total.value = 0;
		} finally {
			loading.value = false;
		}
	}

	/**
	 * Tạo mới thông tin nhân viên
	 * @param {Object} employeeData - Dữ liệu nhân viên cần tạo
	 */
	async function createEmployee(employeeData) {
		return await api.post("/employees", employeeData);
	}

	/**
	 * Cập nhật thông tin nhân viên hiện có
	 * @param {number|string} id - ID của nhân viên
	 * @param {Object} employeeData - Dữ liệu cập nhật
	 */
	async function updateEmployee(id, employeeData) {
		return await api.put(`/employees/${id}`, employeeData);
	}

	/**
	 * Xóa thông tin nhân viên khỏi hệ thống
	 * @param {number|string} id - ID của nhân viên cần xóa
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


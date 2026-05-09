import api from "@/api";
import { defineStore } from "pinia";
import { ref } from "vue";

export const useUserStore = defineStore("user", () => {
	// --- State ---
	const users = ref([]);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Lấy danh sách toàn bộ người dùng từ hệ thống
	 */
	async function fetchUsers() {
		loading.value = true;

		try {
			const { success, data } = await api.get("/users");

			const isValidResponse = success && Array.isArray(data);
			if (!isValidResponse) {
				users.value = [];
				return;
			}

			users.value = data;
		} catch (error) {
			console.error("Lỗi khi lấy danh sách người dùng:", error);
			users.value = [];
		} finally {
			loading.value = false;
		}
	}

	/**
	 * Tạo tài khoản người dùng mới
	 * @param {Object} userData - Thông tin người dùng cần tạo
	 */
	async function createUser(userData) {
		return await api.post("/users", userData);
	}

	/**
	 * Cập nhật thông tin tài khoản người dùng
	 * @param {number|string} id - ID của người dùng
	 * @param {Object} userData - Dữ liệu cập nhật
	 */
	async function updateUser(id, userData) {
		return await api.put(`/users/${id}`, userData);
	}

	/**
	 * Xóa tài khoản người dùng khỏi hệ thống
	 * @param {number|string} id - ID của người dùng cần xóa
	 */
	async function deleteUser(id) {
		return await api.delete(`/users/${id}`);
	}

	return {
		users,
		loading,
		fetchUsers,
		createUser,
		updateUser,
		deleteUser,
	};
});


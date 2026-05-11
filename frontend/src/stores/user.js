import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/api";

/**
 * Store quản lý tài khoản người dùng
 */
export const useUserStore = defineStore("user", () => {
	// --- State ---
	const users = ref([]);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Tải danh sách người dùng hệ thống
	 */
	async function fetchUsers() {
		loading.value = true;
		try {
			const res = await api.get("/users");
			users.value = res.data || [];
		} catch (error) {
			console.error("[UserStore] fetchUsers failed:", error);
			users.value = [];
		} finally {
			loading.value = false;
		}
	}

	/**
	 * Tạo tài khoản mới
	 * @param {Object} data
	 */
	async function createUser(data) {
		return await api.post("/users", data);
	}

	/**
	 * Cập nhật thông tin tài khoản
	 * @param {number|string} id
	 * @param {Object} data
	 */
	async function updateUser(id, data) {
		return await api.put(`/users/${id}`, data);
	}

	/**
	 * Xóa tài khoản khỏi hệ thống
	 * @param {number|string} id
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

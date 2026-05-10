import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "../api";

export const useAuthStore = defineStore("auth", () => {
	// --- State ---
	const token = ref(localStorage.getItem("token"));
	const user = ref(null);

	// Khởi tạo user từ storage
	try {
		const savedUser = localStorage.getItem("user");
		user.value = savedUser ? JSON.parse(savedUser) : null;
	} catch (error) {
		user.value = null;
	}

	// --- Getters ---
	const isLoggedIn = computed(() => !!token.value);

	// --- Actions ---

	/**
	 * Xử lý đăng nhập và lưu trữ thông tin xác thực
	 * @param {string} email - Email đăng nhập
	 * @param {string} password - Mật khẩu
	 */
	async function login(email, password) {
		const response = await api.post("/auth/login", { email, password });
		
		const data = response.data;
		token.value = data.token;
		user.value = data.user;

		localStorage.setItem("token", token.value);
		localStorage.setItem("user", JSON.stringify(user.value));
	}

	/**
	 * Đăng xuất và xóa sạch thông tin xác thực
	 */
	function logout() {
		token.value = null;
		user.value = null;

		localStorage.removeItem("token");
		localStorage.removeItem("user");
		sessionStorage.removeItem("token");
		sessionStorage.removeItem("user");
	}

	/**
	 * Xác minh token hiện tại với Server
	 * @returns {Promise<boolean>} - Trả về true nếu token hợp lệ
	 */
	async function verifyToken() {
		if (!token.value) return false;
		
		try {
			const response = await api.get("/auth/me");
			user.value = response.data;
			
			localStorage.setItem("user", JSON.stringify(user.value));
			return true;
		} catch (error) {
			logout();
			return false;
		}
	}

	return {
		token,
		user,
		isLoggedIn,
		login,
		logout,
		verifyToken,
	};
});

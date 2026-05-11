import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/api";
import { useMeStore } from "./me";

/**
 * Store quản lý xác thực (Token và Account cơ bản)
 */
export const useAuthStore = defineStore("auth", () => {
	const meStore = useMeStore();

	// --- State ---
	const token = ref(localStorage.getItem("token"));
	const user = ref(null);

	// Khởi tạo user cơ bản từ localStorage
	try {
		const savedUser = localStorage.getItem("user");
		user.value = savedUser ? JSON.parse(savedUser) : null;
	} catch (error) {
		user.value = null;
	}

	// --- Getters ---
	const isLoggedIn = computed(() => !!token.value);
	const userRole = computed(() => user.value?.role || "");

	// --- Actions ---

	/**
	 * Đăng nhập và lưu token
	 */
	async function login(email, password) {
		const res = await api.post("/auth/login", { email, password });
		const { token: tokenData, user: userData } = res.data;

		token.value = tokenData;
		user.value = userData;

		localStorage.setItem("token", tokenData);
		localStorage.setItem("user", JSON.stringify(userData));

		// Sau khi đăng nhập, nạp hồ sơ chi tiết vào meStore
		await meStore.fetchProfile();
	}

	/**
	 * Đăng xuất: Xóa sạch dấu vết
	 */
	function logout() {
		token.value = null;
		user.value = null;
		meStore.clearProfile();

		localStorage.removeItem("token");
		localStorage.removeItem("user");

		// Xóa các key cũ nếu còn tồn tại
		localStorage.removeItem("employee");
		sessionStorage.clear();
	}

	/**
	 * Xác minh Token và làm mới dữ liệu
	 */
	async function verifyToken() {
		if (!token.value) return false;

		try {
			const data = await meStore.fetchProfile();
			// API /me trả về: { user: {...}, employee: {...} }
			user.value = data.user;
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
		userRole,
		login,
		logout,
		verifyToken,
	};
});

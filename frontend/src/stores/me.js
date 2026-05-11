import { defineStore } from "pinia";
import { ref, computed } from "vue";
import api from "@/api";

/**
 * Store quản lý thông tin cá nhân (Me)
 * Cấu trúc API thực tế bạn vừa gửi: 
 * { 
 *   success: true, 
 *   data: { 
 *     user: { id, name, email, avatar, role }, 
 *     employee: { ... } 
 *   } 
 * }
 */
export const useMeStore = defineStore("me", () => {
	// --- State ---
	const profile = ref(null); 
	const loading = ref(false);

	// --- Getters ---
	const isLoggedIn = computed(() => !!profile.value);
	
	/**
	 * Thông tin tài khoản (Lấy từ object con 'user')
	 */
	const user = computed(() => profile.value?.user || {});
	
	/**
	 * Thông tin nhân viên chi tiết (Lấy từ object con 'employee')
	 */
	const employee = computed(() => profile.value?.employee || {});

	// --- Actions ---

	async function fetchProfile() {
		if (loading.value) return;
		loading.value = true;
		try {
			const res = await api.get("/auth/me");
			// res = { success, data: { user, employee } }
			profile.value = res.data;
			return res.data;
		} catch (error) {
			console.error("[MeStore] fetchProfile failed:", error);
			profile.value = null;
			throw error;
		} finally {
			loading.value = false;
		}
	}

	function clearProfile() {
		profile.value = null;
	}

	return {
		profile,
		user,
		employee,
		loading,
		isLoggedIn,
		fetchProfile,
		clearProfile,
	};
});

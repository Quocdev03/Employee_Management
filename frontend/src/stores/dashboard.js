import { defineStore } from "pinia";
import { ref } from "vue";
import api from "../api";

export const useDashboardStore = defineStore("dashboard", () => {
	// --- State ---
	const dashboardData = ref(null);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Lấy dữ liệu tổng quan và thống kê cho Dashboard
	 */
	async function fetchDashboardData() {
		loading.value = true;
		try {
			const res = await api.get("/dashboard");
			
			// Ưu tiên lấy dữ liệu từ trường data, nếu không thì lấy toàn bộ res
			dashboardData.value = res.data || res;
		} catch (error) {
			console.error("Lỗi khi lấy dữ liệu Dashboard:", error);
			dashboardData.value = null;
		} finally {
			loading.value = false;
		}
	}

	return {
		dashboardData,
		loading,
		fetchDashboardData,
	};
});


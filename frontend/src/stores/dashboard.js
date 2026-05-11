import { defineStore } from "pinia";
import { ref } from "vue";
import api from "@/api";

/**
 * Store quản lý dữ liệu thống kê bảng điều khiển
 */
export const useDashboardStore = defineStore("dashboard", () => {
	// --- State ---
	const dashboardData = ref(null);
	const loading = ref(false);

	// --- Actions ---

	/**
	 * Tải dữ liệu tổng quan cho dashboard
	 */
	async function fetchDashboardData() {
		loading.value = true;
		try {
			const res = await api.get("/dashboard");
			dashboardData.value = res.data;
		} catch (error) {
			console.error("[DashboardStore] fetchDashboardData failed:", error);
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

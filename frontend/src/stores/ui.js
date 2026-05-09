import { defineStore } from "pinia";
import { ref } from "vue";

export const useUIStore = defineStore("ui", () => {
	// --- State ---
	const mobileOpen = ref(false);
	const isCollapsed = ref(false);

	// --- Actions ---

	/**
	 * Bật/tắt thanh Sidebar trên giao diện di động
	 */
	function toggleMobileSidebar() {
		mobileOpen.value = !mobileOpen.value;
	}

	/**
	 * Đóng thanh Sidebar trên di động
	 */
	function closeMobileSidebar() {
		mobileOpen.value = false;
	}

	/**
	 * Mở thanh Sidebar trên di động
	 */
	function openMobileSidebar() {
		mobileOpen.value = true;
	}

	/**
	 * Bật/tắt trạng thái thu gọn của thanh Sidebar chính
	 */
	function toggleCollapse() {
		isCollapsed.value = !isCollapsed.value;
	}

	return {
		mobileOpen,
		isCollapsed,
		toggleMobileSidebar,
		closeMobileSidebar,
		openMobileSidebar,
		toggleCollapse,
	};
});


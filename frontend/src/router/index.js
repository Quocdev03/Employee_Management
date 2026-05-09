import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";

// --- Danh sách các Routes ---
const routes = [
	{
		path: "/login",
		name: "login",
		component: () => import("../views/LoginView.vue"),
		meta: { public: true },
	},
	{
		path: "/",
		component: () => import("../layouts/MainLayout.vue"),
		redirect: "/dashboard",
		children: [
			{
				path: "dashboard",
				name: "dashboard",
				component: () => import("../views/DashboardView.vue"),
			},
			{
				path: "employees",
				name: "employees",
				component: () => import("@/views/EmployeeListView.vue"),
			},
			{
				path: "profile",
				name: "profile",
				component: () => import("@/views/ProfileView.vue"),
			},
			{
				path: "users",
				name: "users",
				component: () => import("@/views/UserManagementView.vue"),
				meta: { adminOnly: true },
			},
			{
				path: "departments",
				name: "departments",
				component: () => import("@/views/DepartmentView.vue"),
			},
		],
	},
];

// --- Cấu hình Router ---
const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes,
});

// Biến đánh dấu đã xác thực token lần đầu hay chưa
let isTokenVerified = false;

// --- Navigation Guards (Bảo mật tuyến đường) ---

router.beforeEach(async to => {
	const auth = useAuthStore();

	// 1. Xác thực Token lần đầu khi vào các route được bảo vệ
	const needsTokenVerification =
		!to.meta.public && auth.token && !isTokenVerified;
	if (needsTokenVerification) {
		const isValid = await auth.verifyToken();
		isTokenVerified = true;

		if (!isValid) {
			return { name: "login" };
		}
	}

	// 2. Kiểm tra trạng thái đăng nhập
	const isProtectedRoute = !to.meta.public;
	if (isProtectedRoute && !auth.isLoggedIn) {
		return { name: "login" };
	}

	// 3. Chống quay lại trang Login khi đã đăng nhập
	if (to.name === "login" && auth.isLoggedIn) {
		return { name: "dashboard" };
	}

	// 4. Kiểm tra quyền Admin cho các route đặc thù
	const isAdminRequired = to.meta.adminOnly;
	if (isAdminRequired && auth.user?.role !== "admin") {
		return { name: "dashboard" };
	}
});

export default router;

import axios from "axios";

// --- Cấu hình Instance ---
const api = axios.create({
	baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api",
	timeout: 10000,
});

// --- Interceptors (Yêu cầu) ---

/**
 * Tự động gắn Token vào Header trước khi gửi yêu cầu nếu có
 */
api.interceptors.request.use((config) => {
	const token = localStorage.getItem("token");
	if (token) {
		config.headers.Authorization = `Bearer ${token}`;
	}
	return config;
});

// --- Interceptors (Phản hồi) ---

/**
 * Xử lý dữ liệu trả về và các lỗi hệ thống (như 401 Unauthorized)
 */
api.interceptors.response.use(
	(response) => response.data,
	(error) => {
		const status = error.response?.status;

		// 1. Xử lý lỗi 401: Token hết hạn hoặc không hợp lệ -> Đăng xuất tự động
		if (status === 401) {
			const storageKeys = ["token", "user"];
			
			storageKeys.forEach(key => {
				localStorage.removeItem(key);
				sessionStorage.removeItem(key);
			});

			// Điều hướng về trang đăng nhập nếu không phải đang ở đó
			const isAtLoginPage = window.location.pathname === "/login";
			if (!isAtLoginPage) {
				window.location.href = "/login";
			}
		}

		// 2. Trả về đối tượng lỗi chuẩn hóa cho các service/store xử lý
		const errorMessage = 
			error.response?.data?.message || 
			error.response?.data?.error || 
			"Đã có lỗi xảy ra, vui lòng thử lại sau";

		return Promise.reject({
			status,
			message: errorMessage,
		});
	},
);

export default api;


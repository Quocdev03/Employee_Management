<template>
	<!-- ===== Container Chính ===== -->
	<main class="login-container">
		<!-- ===== Thẻ Đăng Nhập ===== -->
		<section class="login-card">
			<!-- ===== Phần Tiêu Đề ===== -->
			<header class="logo-section">
				<div class="logo-icon">
					<img :src="usersIcon" alt="logo" class="logo-img" />
				</div>
				<h2>Welcome Back</h2>
				<p>Đăng nhập vào Hệ thống Quản lý Nhân sự</p>
			</header>

			<!-- ===== Biểu Mẫu Đăng Nhập ===== -->
			<form @submit.prevent="handleLogin" class="login-form">
				<!-- Nhập Email -->
				<div class="input-group">
					<label for="email">Địa chỉ Email</label>
					<div class="input-wrapper">
						<img :src="mailIcon" class="input-icon" />
						<input
							type="email"
							id="email"
							v-model="credentials.email"
							placeholder="admin@company.com"
							required
						/>
					</div>
				</div>

				<!-- Nhập Mật Khẩu -->
				<div class="input-group">
					<label for="password">Mật khẩu</label>
					<div class="input-wrapper">
						<img :src="lockIcon" class="input-icon" />
						<input
							:type="isPasswordVisible ? 'text' : 'password'"
							id="password"
							v-model="credentials.password"
							placeholder="••••••••"
							required
						/>
						<button
							type="button"
							class="toggle-password"
							@click="togglePasswordVisibility"
						>
							<img
								:src="!isPasswordVisible ? eyeIcon : eyeOffIcon"
								alt="toggle"
							/>
						</button>
					</div>
				</div>

				<!-- Hành Động Phụ -->
				<div class="form-actions">
					<a href="#" class="forgot-password">Quên mật khẩu?</a>
				</div>

				<!-- Nút Gửi -->
				<button type="submit" class="submit-btn" :disabled="isLoading">
					<span v-if="!isLoading">Đăng Nhập</span>
					<span v-else class="loader"></span>
				</button>
			</form>
		</section>
	</main>
</template>
<script setup>
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";
import { useAuthStore } from "../stores/auth";

// --- Tài nguyên và Biểu tượng ---
import usersIcon from "@/assets/icons/users.svg";
import mailIcon from "@/assets/icons/mail.svg";
import lockIcon from "@/assets/icons/lock.svg";
import eyeIcon from "@/assets/icons/eye.svg";
import eyeOffIcon from "@/assets/icons/eye-off.svg";

// --- Store và Router ---
const authStore = useAuthStore();
const router = useRouter();
const toast = useToast();

// --- Trạng thái Reactive ---
const isLoading = ref(false);
const isPasswordVisible = ref(false);
const credentials = reactive({
	email: "",
	password: "",
});

// --- Tiện ích bổ trợ ---

/**
 * Tạo độ trễ giả lập để tăng trải nghiệm người dùng
 */
const wait = ms => new Promise(resolve => setTimeout(resolve, ms));

// --- Xử lý Logic nghiệp vụ ---

/**
 * Xử lý đăng nhập hệ thống
 */
async function handleLogin() {
	if (isLoading.value) return;

	isLoading.value = true;

	try {
		// Tạo hiệu ứng chờ mượt mà cho UI
		await wait(800);

		await authStore.login(credentials.email, credentials.password);

		toast.success("Đăng nhập thành công!");
		router.push("/dashboard");
	} catch (error) {
		const message =
			error.message || "Đăng nhập thất bại. Vui lòng kiểm tra lại!";
		toast.error(message);
	} finally {
		isLoading.value = false;
	}
}

/**
 * Chuyển đổi trạng thái hiển thị mật khẩu
 */
function togglePasswordVisibility() {
	isPasswordVisible.value = !isPasswordVisible.value;
}
</script>
<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

/* ===== Base Styles & Container ===== */
.login-container {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f1f5f9;
	font-family: "Inter", sans-serif;
	padding: 1rem;
}

/* ===== Login Card ===== */
.login-card {
	width: 100%;
	max-width: 450px;
	background: rgba(255, 255, 255, 0.98);
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	border: 1px solid #e2e8f0;
	border-radius: 16px;
	padding: 2.5rem;
	box-shadow:
		0 10px 25px -5px rgba(0, 0, 0, 0.03),
		0 8px 10px -6px rgba(0, 0, 0, 0.01);
}

/* ===== Header Section ===== */
.logo-section {
	text-align: center;
	margin-bottom: 2rem;
}

.logo-icon {
	margin-bottom: 1rem;
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 60px;
	height: 60px;
	background-color: #eff6ff;
	color: #3b82f6;
	border-radius: 16px;
}

.logo-img {
	width: 32px;
	height: 32px;
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%); /* #3b82f6 */
}

.logo-section h2 {
	color: #1e293b;
	font-size: 1.5rem;
	font-weight: 700;
	margin-bottom: 0.5rem;
	letter-spacing: -0.01em;
}

.logo-section p {
	color: #64748b;
	font-size: 0.95rem;
}

/* ===== Form Structure ===== */
.login-form {
	display: flex;
	flex-direction: column;
	gap: 1.25rem;
}

.input-group {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.input-group label {
	font-size: 0.875rem;
	font-weight: 600;
	color: #475569;
}

/* ===== Input Fields ===== */
.input-wrapper {
	position: relative;
	display: flex;
	align-items: center;
}

.input-icon {
	position: absolute;
	left: 1rem;
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg)
		brightness(88%) contrast(89%); /* #94a3b8 */
	transition: all 0.2s ease;
}

.input-wrapper input {
	width: 100%;
	padding: 0.75rem 1rem 0.75rem 2.75rem;
	border: 1px solid #cbd5e1;
	border-radius: 8px;
	font-size: 0.95rem;
	color: #1e293b;
	background: #ffffff;
	transition: all 0.2s ease;
	outline: none;
}

.input-wrapper input::placeholder {
	color: #94a3b8;
}

.input-wrapper input:focus {
	border-color: #3b82f6;
	box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-wrapper input:focus ~ .input-icon {
	filter: invert(44%) sepia(87%) saturate(2258%) hue-rotate(200deg)
		brightness(101%) contrast(92%); /* #3b82f6 */
}

/* ===== Password Toggle ===== */
.toggle-password {
	position: absolute;
	right: 0.75rem;
	background: none;
	border: none;
	color: #94a3b8;
	cursor: pointer;
	padding: 0.25rem;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: color 0.2s ease;
}

.toggle-password:hover {
	color: #64748b;
}

.toggle-password img {
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg)
		brightness(88%) contrast(89%);
}

/* ===== Form Actions ===== */
.form-actions {
	display: flex;
	align-items: center;
	justify-content: flex-end;
}

.forgot-password {
	font-size: 0.875rem;
	color: #3b82f6;
	font-weight: 500;
	text-decoration: none;
	transition: color 0.2s ease;
}

.forgot-password:hover {
	color: #2563eb;
	text-decoration: underline;
}

/* ===== Submit Button ===== */
.submit-btn {
	background-color: #3b82f6;
	color: white;
	border: none;
	border-radius: 8px;
	padding: 0.75rem 1rem;
	font-size: 1rem;
	font-weight: 600;
	cursor: pointer;
	transition: background-color 0.2s ease;
	display: flex;
	justify-content: center;
	align-items: center;
	min-height: 48px;
	margin-top: 0.5rem;
}

.submit-btn:hover {
	background-color: #2563eb;
	transition: all 0.2s linear;
}

.submit-btn:active {
	background-color: #1d4ed8;
}

.submit-btn:disabled {
	opacity: 0.7;
	cursor: not-allowed;
}

/* ===== Loading Spinner ===== */
.loader {
	width: 20px;
	height: 20px;
	border: 2px solid rgba(255, 255, 255, 0.4);
	border-radius: 50%;
	border-top-color: white;
	animation: spin 0.8s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

/* ===== Responsive Design ===== */
@media (max-width: 480px) {
	.login-card {
		padding: 2rem 1.5rem;
		border: none;
		border-radius: 0;
		box-shadow: none;
		background: transparent;
		backdrop-filter: none;
	}
}
</style>

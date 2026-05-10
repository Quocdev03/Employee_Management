<template>
	<div class="department-view">
		<!-- ===== Tiêu đề trang ===== -->
		<header class="page-header">
			<div class="header-content">
				<h1 class="page-title">Quản lý phòng ban</h1>
				<p class="page-subtitle">
					Quản lý cơ cấu tổ chức và các phòng ban trong hệ thống
				</p>
			</div>
			<button
				v-if="isAdmin"
				class="btn btn--primary"
				@click="openForm(null)"
			>
				<img :src="plusIcon" alt="add" class="btn__icon" />
				Thêm phòng ban
			</button>
		</header>

		<!-- ===== Nội dung chính ===== -->
		<main class="content-card">
			<!-- Thanh công cụ -->
			<div class="toolbar">
				<div class="search-box">
					<img :src="searchIcon" class="search-box__icon" alt="search" />
					<input
						v-model="searchKeyword"
						class="form-control search-box__input"
						placeholder="Tìm tên phòng ban..."
					/>
				</div>
			</div>

			<!-- Bảng dữ liệu -->
			<BaseTable>
				<template #header>
					<th>Tên phòng ban</th>
					<th>Chức vụ</th>
					<th>Ngày tạo</th>
					<th v-if="isAdmin" class="text-right">Thao tác</th>
				</template>
				<template #body>
					<tr v-for="dept in filteredDepartments" :key="dept.id">
						<td>
							<div class="dept-info">
								<div class="dept-icon-box">
									<img :src="deptIcon" alt="dept" />
								</div>
								<div class="dept-name">{{ dept.name }}</div>
							</div>
						</td>
						<td>
							<div class="position-tags">
								<span 
									v-for="pos in dept.positions" 
									:key="pos.id"
									class="position-tag"
								>
									{{ pos.name }}
								</span>
								<span v-if="!dept.positions?.length" class="text-muted text-xs italic">
									Chưa có chức vụ
								</span>
							</div>
						</td>
						<td class="text-muted">
							{{ formatVietnameseDate(dept.created_at) }}
						</td>
						<td v-if="isAdmin" class="text-right">
							<div class="action-group">
								<button
									class="btn-icon btn-icon--edit"
									title="Chỉnh sửa"
									@click="openForm(dept)"
								>
									<img :src="editIcon" alt="edit" />
								</button>
								<button
									class="btn-icon btn-icon--delete"
									title="Xoá"
									@click="prepareDelete(dept)"
								>
									<img :src="deleteIcon" alt="delete" />
								</button>
							</div>
						</td>
					</tr>
				</template>
			</BaseTable>

			<!-- Trạng thái trống -->
			<div v-if="filteredDepartments.length === 0" class="empty-state">
				<div class="empty-state__icon">🏢</div>
				<p class="empty-state__text">Không tìm thấy phòng ban nào.</p>
			</div>
		</main>

		<!-- ===== Modal: Thêm/Sửa phòng ban ===== -->
		<BaseModal
			:visible="isFormVisible"
			:title="editingDept ? 'Chỉnh sửa phòng ban' : 'Thêm phòng ban mới'"
			subtitle="Thông tin phòng ban trong sơ đồ tổ chức"
			is-form
			@submit="handleFormSubmit"
			@close="isFormVisible = false"
		>
			<div class="form-group">
				<label class="form-label">Tên phòng ban <span class="required">*</span></label>
				<div class="input-group">
					<img :src="deptIcon" class="input-group__icon" alt="dept" />
					<input
						v-model="form.name"
						class="form-control input-group__input"
						placeholder="Ví dụ: Phòng Kỹ thuật"
						ref="nameInput"
					/>
				</div>
			</div>

			<!-- ===== Section: Quản lý chức vụ (Chỉ hiện khi Edit) ===== -->
			<div v-if="editingDept" class="positions-section">
				<div class="section-divider"></div>
				<label class="form-label">Danh sách chức vụ</label>
				
				<!-- Danh sách chức vụ hiện có -->
				<div class="positions-list">
					<div 
						v-for="pos in editingDept.positions" 
						:key="pos.id" 
						class="position-item"
					>
						<template v-if="editingPosId !== pos.id">
							<span class="position-item__name">{{ pos.name }}</span>
							<div class="position-item__actions">
								<button 
									type="button" 
									class="btn-icon btn-icon--sm" 
									@click="startEditPosition(pos)"
								>
									<img :src="editIcon" alt="edit" />
								</button>
								<button 
									type="button" 
									class="btn-icon btn-icon--sm btn-icon--delete" 
									@click="handleDeletePosition(pos.id)"
								>
									<img :src="deleteIcon" alt="delete" />
								</button>
							</div>
						</template>
						<template v-else>
							<input 
								v-model="editingPosName" 
								class="form-control form-control--sm" 
								@keyup.enter="handleUpdatePosition"
								@keyup.esc="cancelEditPosition"
								ref="editPosInput"
							/>
							<div class="position-item__actions">
								<button type="button" class="btn-icon btn-icon--sm" @click="handleUpdatePosition">✅</button>
								<button type="button" class="btn-icon btn-icon--sm" @click="cancelEditPosition">❌</button>
							</div>
						</template>
					</div>
				</div>

				<!-- Ô thêm chức vụ mới -->
				<div class="add-position-box">
					<input 
						v-model="newPosName" 
						class="form-control form-control--sm" 
						placeholder="Thêm chức vụ mới..." 
						@keyup.enter="handleAddPosition"
					/>
					<button 
						type="button" 
						class="btn btn--primary btn--sm" 
						@click="handleAddPosition"
						:disabled="!newPosName.trim()"
					>
						Thêm
					</button>
				</div>
			</div>

			<template #footer>
				<button type="button" class="btn btn--secondary" @click="isFormVisible = false">
					Huỷ bỏ
				</button>
				<button type="submit" class="btn btn--primary" :disabled="isSubmitting">
					{{ isSubmitting ? "Đang lưu..." : (editingDept ? "Lưu thay đổi" : "Tạo phòng ban") }}
				</button>
			</template>
		</BaseModal>

		<!-- ===== Modal: Xác nhận xoá ===== -->
		<ConfirmModal
			:visible="isDeleteModalVisible"
			title="Xác nhận xoá"
			confirm-text="Đồng ý xoá"
			@cancel="isDeleteModalVisible = false"
			@confirm="executeDelete"
		>
			Bạn có chắc chắn muốn xoá phòng ban <strong>{{ deptToDelete?.name }}</strong>? 
			Hành động này không thể hoàn tác.
		</ConfirmModal>
	</div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from "vue";
import { useDepartmentStore } from "@/stores/department";
import { useAuthStore } from "@/stores/auth";
import { useToast } from "vue-toastification";

import BaseTable from "@/components/BaseTable.vue";
import BaseModal from "@/components/BaseModal.vue";
import ConfirmModal from "@/components/ConfirmModal.vue";

// Icons
import plusIcon from "@/assets/icons/plus.svg";
import searchIcon from "@/assets/icons/search.svg";
import editIcon from "@/assets/icons/edit.svg";
import deleteIcon from "@/assets/icons/delete.svg";
import deptIcon from "@/assets/icons/department.svg";

const deptStore = useDepartmentStore();
const authStore = useAuthStore();
const toast = useToast();

const searchKeyword = ref("");
const isFormVisible = ref(false);
const isDeleteModalVisible = ref(false);
const isSubmitting = ref(false);
const editingDept = ref(null);
const deptToDelete = ref(null);

// State cho quản lý chức vụ
const newPosName = ref("");
const editingPosId = ref(null);
const editingPosName = ref("");
const editPosInput = ref(null);

const form = reactive({
	name: ""
});

const isAdmin = computed(() => authStore.user?.role === "admin");

const filteredDepartments = computed(() => {
	const query = searchKeyword.value.toLowerCase().trim();
	if (!query) return deptStore.departments;
	return deptStore.departments.filter(d => d.name.toLowerCase().includes(query));
});

function openForm(dept = null) {
	editingDept.value = dept;
	form.name = dept ? dept.name : "";
	
	// Reset state chức vụ
	newPosName.value = "";
	editingPosId.value = null;
	editingPosName.value = "";
	
	isFormVisible.value = true;
}

// --- Logic Quản lý Chức vụ ---

async function handleAddPosition() {
	if (!newPosName.value.trim()) return;
	
	try {
		await deptStore.createPosition({
			name: newPosName.value.trim(),
			department_id: editingDept.value.id
		});
		newPosName.value = "";
		toast.success("Đã thêm chức vụ");
		
		// Refresh data để cập nhật list chức vụ trong modal
		await deptStore.fetchDepartments();
		// Cập nhật lại editingDept để UI trong modal refresh
		editingDept.value = deptStore.departments.find(d => d.id === editingDept.value.id);
	} catch (error) {
		toast.error(error.response?.data?.error || "Không thể thêm chức vụ");
	}
}

function startEditPosition(pos) {
	editingPosId.value = pos.id;
	editingPosName.value = pos.name;
}

function cancelEditPosition() {
	editingPosId.value = null;
	editingPosName.value = "";
}

async function handleUpdatePosition() {
	if (!editingPosName.value.trim()) return;
	
	try {
		await deptStore.updatePosition(editingPosId.value, {
			name: editingPosName.value.trim()
		});
		editingPosId.value = null;
		toast.success("Đã cập nhật chức vụ");
		
		await deptStore.fetchDepartments();
		editingDept.value = deptStore.departments.find(d => d.id === editingDept.value.id);
	} catch (error) {
		toast.error(error.response?.data?.error || "Không thể cập nhật chức vụ");
	}
}

async function handleDeletePosition(posId) {
	if (!confirm("Bạn có chắc muốn xoá chức vụ này?")) return;
	
	try {
		await deptStore.deletePosition(posId);
		toast.success("Đã xoá chức vụ");
		
		await deptStore.fetchDepartments();
		editingDept.value = deptStore.departments.find(d => d.id === editingDept.value.id);
	} catch (error) {
		toast.error(error.response?.data?.error || "Không thể xoá chức vụ");
	}
}

async function handleFormSubmit() {
	if (!form.name.trim()) {
		toast.warning("Vui lòng nhập tên phòng ban");
		return;
	}

	isSubmitting.value = true;
	try {
		if (editingDept.value) {
			await deptStore.updateDepartment(editingDept.value.id, { name: form.name });
			toast.success("Cập nhật phòng ban thành công");
		} else {
			await deptStore.createDepartment({ name: form.name });
			toast.success("Thêm phòng ban mới thành công");
		}
		await deptStore.fetchDepartments();
		isFormVisible.value = false;
	} catch (error) {
		toast.error(error.response?.data?.error || "Thao tác thất bại");
	} finally {
		isSubmitting.value = false;
	}
}

function prepareDelete(dept) {
	deptToDelete.value = dept;
	isDeleteModalVisible.value = true;
}

async function executeDelete() {
	if (!deptToDelete.value) return;

	try {
		await deptStore.deleteDepartment(deptToDelete.value.id);
		toast.success(`Đã xoá phòng ban "${deptToDelete.value.name}"`);
		await deptStore.fetchDepartments();
	} catch (error) {
		toast.error(error.response?.data?.error || "Không thể xoá phòng ban này");
	} finally {
		isDeleteModalVisible.value = false;
		deptToDelete.value = null;
	}
}

function formatVietnameseDate(dateStr) {
	if (!dateStr) return "—";
	return new Date(dateStr).toLocaleDateString("vi-VN");
}

onMounted(() => {
	deptStore.fetchDepartments();
});
</script>

<style scoped>
@import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap");

.department-view {
	font-family: "Inter", sans-serif;
	color: var(--text-main);
}

.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 2rem;
	flex-wrap: wrap;
	gap: 1rem;
}

.page-title {
	font-size: 1.75rem;
	font-weight: 700;
	letter-spacing: -0.02em;
	margin: 0 0 4px 0;
}

.page-subtitle {
	color: var(--text-muted);
	font-size: 0.95rem;
	margin: 0;
}

.content-card {
	background: var(--bg-card);
	border-radius: var(--radius-lg);
	border: 1px solid var(--border-color);
	box-shadow: var(--shadow-sm);
	overflow: hidden;
}

.toolbar {
	padding: 1.5rem;
	border-bottom: 1px solid var(--bg-main);
}

.search-box {
	position: relative;
	max-width: 400px;
}

.search-box__icon {
	position: absolute;
	left: 1rem;
	top: 50%;
	transform: translateY(-50%);
	width: 18px;
	height: 18px;
	filter: invert(72%) sepia(10%) saturate(415%) hue-rotate(182deg) brightness(88%) contrast(89%);
}

.search-box__input {
	padding-left: 2.75rem;
}

.dept-info {
	display: flex;
	align-items: center;
	gap: 12px;
}

.dept-icon-box {
	width: 36px;
	height: 36px;
	background: #f1f5f9;
	border-radius: var(--radius-md);
	display: flex;
	align-items: center;
	justify-content: center;
}

.dept-icon-box img {
	width: 20px;
	height: 20px;
	filter: invert(44%) sepia(12%) saturate(545%) hue-rotate(183deg) brightness(95%) contrast(87%);
}

.dept-name {
	font-weight: 600;
	color: var(--text-main);
}

.position-tags {
	display: flex;
	flex-wrap: wrap;
	gap: 4px;
}

.position-tag {
	font-size: 0.75rem;
	background: #f8fafc;
	color: #64748b;
	padding: 2px 8px;
	border-radius: 4px;
	border: 1px solid #e2e8f0;
}

.action-group {
	display: flex;
	justify-content: flex-end;
	gap: 8px;
}

.empty-state {
	padding: 4rem;
	text-align: center;
	color: #94a3b8;
}

.empty-state__icon {
	font-size: 3rem;
	margin-bottom: 1rem;
}

@media (max-width: 768px) {
	.page-header {
		flex-direction: column;
		align-items: flex-start;
	}
	.btn--primary {
		width: 100%;
		justify-content: center;
	}
	.search-box {
		max-width: none;
	}
}

/* ===== Positions Management Styles ===== */
.positions-section {
	margin-top: 1.5rem;
}

.section-divider {
	height: 1px;
	background: var(--bg-main);
	margin-bottom: 1.5rem;
}

.positions-list {
	display: flex;
	flex-direction: column;
	gap: 8px;
	margin-bottom: 1rem;
	max-height: 200px;
	overflow-y: auto;
	padding-right: 4px;
}

.position-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 8px 12px;
	background: #f8fafc;
	border-radius: var(--radius-md);
	border: 1px solid #e2e8f0;
}

.position-item__name {
	font-size: 0.875rem;
	font-weight: 500;
	color: #334155;
}

.position-item__actions {
	display: flex;
	gap: 4px;
}

.add-position-box {
	display: flex;
	gap: 8px;
	margin-top: 12px;
}

.btn-icon--sm {
	width: 28px;
	height: 28px;
}

.btn-icon--sm img {
	width: 14px;
	height: 14px;
}

.form-control--sm {
	padding: 0.4rem 0.75rem;
	font-size: 0.875rem;
}
</style>

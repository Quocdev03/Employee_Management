<template>
	<div class="table-wrapper">
		<table class="data-table">
			<thead>
				<tr>
					<slot name="header"></slot>
				</tr>
			</thead>
			<tbody>
				<slot name="body"></slot>
			</tbody>
		</table>
	</div>
</template>

<style scoped>
.table-wrapper {
	width: 100%;
	overflow-x: auto;
}

.data-table {
	width: 100%;
	border-collapse: collapse;
	font-size: 0.95rem;
}

:deep(.data-table th),
:deep(.data-table td) {
	padding: 1rem 1.5rem;
	text-align: left;
	border-bottom: 1px solid var(--border-color, #e2e8f0);
}

:deep(.data-table th) {
	background: var(--bg-main, #f8fafc);
	font-weight: 600;
	color: #475569;
	text-transform: uppercase;
	font-size: 0.75rem;
	letter-spacing: 0.05em;
	white-space: nowrap;
}

:deep(.data-table tbody tr) {
	transition: background-color 0.2s;
}

:deep(.data-table tbody tr:hover) {
	background: #f8fafc;
}
</style>

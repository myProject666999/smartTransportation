<template>
  <div class="carousels-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>轮播图管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加轮播图
          </el-button>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="启用" :value="1" />
              <el-option label="禁用" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="title" label="标题" min-width="150" />
          <el-table-column prop="image_url" label="图片" width="120">
            <template #default="scope">
              <el-image
                v-if="scope.row.image_url"
                :src="scope.row.image_url"
                :preview-src-list="[scope.row.image_url]"
                fit="cover"
                style="width: 100px; height: 60px;"
              >
                <template #placeholder>
                  <div class="image-placeholder">暂无图片</div>
                </template>
              </el-image>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="link_url" label="跳转链接" min-width="200" show-overflow-tooltip>
            <template #default="scope">
              <span v-if="scope.row.link_url" class="link-url">{{ scope.row.link_url }}</span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="sort_order" label="排序" width="80">
            <template #default="scope">
              <el-tag type="primary" size="small">{{ scope.row.sort_order || 0 }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                {{ scope.row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleMove(scope.row, 'up')" :disabled="scope.row.sort_order <= 1">
                <el-icon><ArrowUp /></el-icon>上移
              </el-button>
              <el-button type="primary" link size="small" @click="handleMove(scope.row, 'down')">
                <el-icon><ArrowDown /></el-icon>下移
              </el-button>
              <el-button type="success" link size="small" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button type="danger" link size="small" @click="handleDelete(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadData"
          @current-change="loadData"
          style="margin-top: 20px; justify-content: flex-end"
        />
      </el-card>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入轮播图标题" />
        </el-form-item>
        <el-form-item label="图片URL" prop="image_url">
          <el-input v-model="form.image_url" placeholder="请输入图片URL" />
          <div class="upload-tip" v-if="form.image_url">
            <el-image :src="form.image_url" fit="contain" style="width: 100%; height: 150px; margin-top: 10px;">
              <template #error>
                <div class="image-error">图片预览失败</div>
              </template>
            </el-image>
          </div>
        </el-form-item>
        <el-form-item label="跳转链接" prop="link_url">
          <el-input v-model="form.link_url" placeholder="请输入跳转链接（可选）" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" style="width: 100%" />
          <span class="tip-text">数字越小越靠前</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getCarousels, createCarousel, updateCarousel, deleteCarousel } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)

const searchForm = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  title: '',
  image_url: '',
  link_url: '',
  sort_order: 0,
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑轮播图' : '添加轮播图')

const rules = {
  title: [
    { required: true, message: '请输入轮播图标题', trigger: 'blur' }
  ],
  image_url: [
    { required: true, message: '请输入图片URL', trigger: 'blur' }
  ]
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getCarousels(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载轮播图列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.title = ''
  form.image_url = ''
  form.link_url = ''
  form.sort_order = tableData.value.length + 1
  form.status = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.title = row.title
  form.image_url = row.image_url || ''
  form.link_url = row.link_url || ''
  form.sort_order = row.sort_order || 0
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除轮播图"${row.title}"吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteCarousel(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (e) {
      console.error('删除失败', e)
    }
  }).catch(() => {})
}

const handleMove = (row, direction) => {
  const newOrder = direction === 'up' ? (row.sort_order || 0) - 1 : (row.sort_order || 0) + 1
  updateCarousel(row.id, { sort_order: newOrder }).then(res => {
    if (res.code === 200) {
      ElMessage.success(direction === 'up' ? '上移成功' : '下移成功')
      loadData()
    }
  })
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = {
      title: form.title,
      image_url: form.image_url,
      link_url: form.link_url,
      sort_order: form.sort_order,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateCarousel(form.id, data)
    } else {
      res = await createCarousel(data)
    }

    if (res.code === 200) {
      ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
      dialogVisible.value = false
      loadData()
    }
  } catch (e) {
    console.error('提交失败', e)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.carousels-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content h3 {
  margin: 0;
  color: #303133;
}

.search-form {
  margin-bottom: 20px;
}

.text-muted {
  color: #909399;
  font-size: 14px;
}

.link-url {
  color: #409EFF;
  font-size: 13px;
}

.image-placeholder {
  width: 100%;
  height: 60px;
  background: #F5F7FA;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 12px;
  border-radius: 4px;
}

.image-error {
  width: 100%;
  height: 150px;
  background: #F5F7FA;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 14px;
  border-radius: 4px;
}

.tip-text {
  margin-left: 10px;
  color: #909399;
  font-size: 12px;
}

.upload-tip {
  margin-top: 10px;
}
</style>

<template>
  <div class="tickets-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>票务管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加车票
          </el-button>
        </div>
      </el-card>
    </div>

    <div class="page-content">
      <el-card>
        <el-form :inline="true" :model="searchForm" class="search-form">
          <el-form-item label="关键词">
            <el-input
              v-model="searchForm.keyword"
              placeholder="车票名称/路线"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="在售" :value="1" />
              <el-option label="下架" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="name" label="车票名称" min-width="150" />
          <el-table-column prop="route_name" label="路线" min-width="150" />
          <el-table-column prop="start_station" label="始发站" width="100" />
          <el-table-column prop="end_station" label="终点站" width="100" />
          <el-table-column prop="price" label="票价" width="100">
            <template #default="scope">
              <span style="color: #F56C6C; font-weight: bold;">¥{{ scope.row.price }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="stock" label="库存" width="80">
            <template #default="scope">
              <el-tag :type="scope.row.stock > 100 ? 'success' : scope.row.stock > 0 ? 'warning' : 'danger'" size="small">
                {{ scope.row.stock }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                {{ scope.row.status === 1 ? '在售' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">查看</el-button>
              <el-button type="primary" link size="small" @click="handleEdit(scope.row)">编辑</el-button>
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
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="车票名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入车票名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路线" prop="route_name">
              <el-input v-model="form.route_name" placeholder="请输入路线名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="始发站" prop="start_station">
              <el-input v-model="form.start_station" placeholder="请输入始发站" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="终点站" prop="end_station">
              <el-input v-model="form.end_station" placeholder="请输入终点站" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="票价" prop="price">
              <el-input-number v-model="form.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="库存" prop="stock">
              <el-input-number v-model="form.stock" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入车票描述"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">在售</el-radio>
            <el-radio :label="0">下架</el-radio>
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
import { getTickets, createTicket, updateTicket, deleteTicket } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const form = reactive({
  id: null,
  name: '',
  route_name: '',
  start_station: '',
  end_station: '',
  price: 0,
  stock: 0,
  description: '',
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑车票' : '添加车票')

const rules = {
  name: [
    { required: true, message: '请输入车票名称', trigger: 'blur' }
  ],
  route_name: [
    { required: true, message: '请输入路线名称', trigger: 'blur' }
  ],
  start_station: [
    { required: true, message: '请输入始发站', trigger: 'blur' }
  ],
  end_station: [
    { required: true, message: '请输入终点站', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入票价', trigger: 'blur' }
  ]
}

const handleView = (row) => {
  ElMessage.info(`车票名称: ${row.name}\n路线: ${row.route_name}\n票价: ¥${row.price}\n库存: ${row.stock}`)
}

const loadData = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }
    if (searchForm.status !== '') {
      params.status = searchForm.status
    }

    const res = await getTickets(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载车票列表失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetSearch = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  pagination.page = 1
  loadData()
}

const handleAdd = () => {
  isEdit.value = false
  formRef.value?.resetFields()
  form.id = null
  form.name = ''
  form.route_name = ''
  form.start_station = ''
  form.end_station = ''
  form.price = 0
  form.stock = 0
  form.description = ''
  form.status = 1
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.route_name = row.route_name
  form.start_station = row.start_station
  form.end_station = row.end_station
  form.price = row.price
  form.stock = row.stock
  form.description = row.description
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除车票"${row.name}"吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteTicket(row.id)
      if (res.code === 200) {
        ElMessage.success('删除成功')
        loadData()
      }
    } catch (e) {
      console.error('删除失败', e)
    }
  }).catch(() => {})
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data = {
      name: form.name,
      route_name: form.route_name,
      start_station: form.start_station,
      end_station: form.end_station,
      price: form.price,
      stock: form.stock,
      description: form.description,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateTicket(form.id, data)
    } else {
      res = await createTicket(data)
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
.tickets-page {
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
</style>

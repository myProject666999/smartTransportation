<template>
  <div class="routes-page">
    <div class="page-header">
      <el-card>
        <div class="header-content">
          <h3>路线管理</h3>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon> 添加路线
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
              placeholder="路线名称/站点"
              clearable
              style="width: 180px"
              @keyup.enter="loadData"
            />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 100px">
              <el-option label="运营中" :value="1" />
              <el-option label="暂停" :value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="loadData"><el-icon><Search /></el-icon> 搜索</el-button>
            <el-button @click="handleResetSearch">重置</el-button>
          </el-form-item>
        </el-form>

        <el-table :data="tableData" style="width: 100%" v-loading="loading" stripe>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="route_name" label="路线名称" min-width="180">
            <template #default="scope">
              <div class="route-name">
                <span class="route-label">线路</span>
                <span>{{ scope.row.route_name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="start_station" label="起点站" width="120" />
          <el-table-column prop="end_station" label="终点站" width="120" />
          <el-table-column prop="stations" label="途经站点" min-width="200" show-overflow-tooltip>
            <template #default="scope">
              <span v-if="scope.row.stations">
                {{ scope.row.stations }}
              </span>
              <span v-else class="text-muted">-</span>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="票价" width="100">
            <template #default="scope">
              <span class="price-text">¥{{ scope.row.price || 0 }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="duration" label="预计时长" width="100">
            <template #default="scope">
              <span>{{ scope.row.duration || 0 }}分钟</span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="scope.row.status === 1 ? 'success' : 'info'">
                {{ scope.row.status === 1 ? '运营中' : '暂停' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="scope">
              <el-button type="primary" link size="small" @click="handleView(scope.row)">详情</el-button>
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
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路线名称" prop="route_name">
              <el-input v-model="form.route_name" placeholder="如：1号线、K1路" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :value="1">运营中</el-radio>
                <el-radio :value="0">暂停</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="起点站" prop="start_station">
              <el-input v-model="form.start_station" placeholder="请输入起点站" />
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
            <el-form-item label="票价(元)" prop="price">
              <el-input-number v-model="form.price" :min="0" :precision="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="预计时长(分钟)" prop="duration">
              <el-input-number v-model="form.duration" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="发车间隔(分钟)" prop="interval">
              <el-input-number v-model="form.interval" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="车辆类型" prop="vehicle_type">
              <el-select v-model="form.vehicle_type" placeholder="请选择" style="width: 100%">
                <el-option label="公交" value="公交" />
                <el-option label="地铁" value="地铁" />
                <el-option label="大巴" value="大巴" />
                <el-option label="轻轨" value="轻轨" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="途经站点" prop="stations">
          <el-input
            v-model="form.stations"
            type="textarea"
            :rows="3"
            placeholder="请输入途经站点，用逗号分隔，如：站1,站2,站3"
          />
        </el-form-item>
        <el-form-item label="首末班车" prop="schedule">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-time-picker
                v-model="form.first_time"
                placeholder="首班车时间"
                format="HH:mm"
                value-format="HH:mm"
                style="width: 100%"
              />
            </el-col>
            <el-col :span="12">
              <el-time-picker
                v-model="form.last_time"
                placeholder="末班车时间"
                format="HH:mm"
                value-format="HH:mm"
                style="width: 100%"
              />
            </el-col>
          </el-row>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="路线详情" width="600px">
      <el-descriptions :column="2" border v-if="currentItem">
        <el-descriptions-item label="路线名称">{{ currentItem.route_name }}</el-descriptions-item>
        <el-descriptions-item label="车辆类型">{{ currentItem.vehicle_type || '-' }}</el-descriptions-item>
        <el-descriptions-item label="起点站">{{ currentItem.start_station }}</el-descriptions-item>
        <el-descriptions-item label="终点站">{{ currentItem.end_station }}</el-descriptions-item>
        <el-descriptions-item label="票价">¥{{ currentItem.price || 0 }}</el-descriptions-item>
        <el-descriptions-item label="预计时长">{{ currentItem.duration || 0 }}分钟</el-descriptions-item>
        <el-descriptions-item label="发车间隔">{{ currentItem.interval || 0 }}分钟</el-descriptions-item>
        <el-descriptions-item label="运营状态">
          <el-tag :type="currentItem.status === 1 ? 'success' : 'info'">
            {{ currentItem.status === 1 ? '运营中' : '暂停' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="首班车">{{ currentItem.first_time || '-' }}</el-descriptions-item>
        <el-descriptions-item label="末班车">{{ currentItem.last_time || '-' }}</el-descriptions-item>
        <el-descriptions-item label="途经站点" :span="2">
          <div v-if="currentItem.stations" class="stations-list">
            <el-tag v-for="(station, index) in currentItem.stations.split(',')" :key="index" style="margin: 4px;">
              {{ station }}
            </el-tag>
          </div>
          <span v-else class="text-muted">-</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { getRoutes, createRoute, updateRoute, deleteRoute } from '@/api/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const tableData = ref([])
const formRef = ref(null)
const currentItem = ref({})

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
  route_name: '',
  start_station: '',
  end_station: '',
  stations: '',
  price: 0,
  duration: 0,
  interval: 0,
  vehicle_type: '',
  first_time: '',
  last_time: '',
  status: 1
})

const dialogTitle = computed(() => isEdit.value ? '编辑路线' : '添加路线')

const rules = {
  route_name: [
    { required: true, message: '请输入路线名称', trigger: 'blur' }
  ],
  start_station: [
    { required: true, message: '请输入起点站', trigger: 'blur' }
  ],
  end_station: [
    { required: true, message: '请输入终点站', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入票价', trigger: 'blur' }
  ]
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

    const res = await getRoutes(params)
    if (res.code === 200) {
      tableData.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    console.error('加载路线列表失败', e)
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
  form.route_name = ''
  form.start_station = ''
  form.end_station = ''
  form.stations = ''
  form.price = 0
  form.duration = 0
  form.interval = 10
  form.vehicle_type = '公交'
  form.first_time = '06:00'
  form.last_time = '22:00'
  form.status = 1
  dialogVisible.value = true
}

const handleView = (row) => {
  currentItem.value = row
  detailDialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.route_name = row.route_name
  form.start_station = row.start_station
  form.end_station = row.end_station
  form.stations = row.stations || ''
  form.price = row.price
  form.duration = row.duration || 0
  form.interval = row.interval || 0
  form.vehicle_type = row.vehicle_type || '公交'
  form.first_time = row.first_time || ''
  form.last_time = row.last_time || ''
  form.status = row.status
  dialogVisible.value = true
}

const handleDelete = (row) => {
  ElMessageBox.confirm(`确认要删除路线"${row.route_name}"吗？`, '确认删除', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await deleteRoute(row.id)
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
      route_name: form.route_name,
      start_station: form.start_station,
      end_station: form.end_station,
      stations: form.stations,
      price: form.price,
      duration: form.duration,
      interval: form.interval,
      vehicle_type: form.vehicle_type,
      first_time: form.first_time,
      last_time: form.last_time,
      status: form.status
    }

    let res
    if (isEdit.value) {
      res = await updateRoute(form.id, data)
    } else {
      res = await createRoute(data)
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
.routes-page {
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

.route-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.route-label {
  background: #409EFF;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.price-text {
  color: #F56C6C;
  font-weight: bold;
  font-size: 16px;
}

.text-muted {
  color: #909399;
}

.stations-list {
  line-height: 1.8;
}
</style>

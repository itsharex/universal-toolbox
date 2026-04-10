<template>
  <ToolPage title="定时任务" description="CRON 定时任务管理与监控">
    <!-- 顶部操作栏 -->
    <div class="flex gap-2 mb-4">
      <button @click="showAddModal = true" class="btn btn-primary">
        <Plus :size="14"/>添加任务
      </button>
      <button @click="loadJobs" class="btn btn-secondary" :disabled="loading">
        <RefreshCw :size="14" :class="loading ? 'loading-spin' : ''"/>刷新
      </button>
    </div>

    <!-- 任务列表 -->
    <div v-if="jobs.length > 0" class="space-y-3">
      <div v-for="job in jobs" :key="job.id" class="card p-4">
        <!-- 任务头部：名称 + 状态 -->
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <span class="status-dot" :class="job.enabled ? (job.lastError ? 'red' : 'green') : 'gray'" />
            <span class="font-semibold text-sm">{{ job.name }}</span>
          </div>
          <div class="flex items-center gap-1.5">
            <button @click="toggleJob(job)" class="btn btn-secondary py-0.5 px-2 text-xs" :title="job.enabled ? '禁用' : '启用'">
              <component :is="job.enabled ? Pause : Play" :size="12"/>
              {{ job.enabled ? '禁用' : '启用' }}
            </button>
            <button @click="runNow(job)" class="btn btn-secondary py-0.5 px-2 text-xs" title="立即执行">
              <Play :size="12"/>执行
            </button>
            <button @click="confirmDelete(job)" class="btn btn-danger py-0.5 px-2 text-xs" title="删除">
              <Trash2 :size="12"/>删除
            </button>
          </div>
        </div>

        <!-- cron 表达式 + 描述 -->
        <div class="mb-3">
          <code class="text-xs px-2 py-0.5 rounded" style="background:var(--bg-hover)">{{ job.expression }}</code>
          <span v-if="job.description" class="text-xs opacity-50 ml-2">{{ job.description }}</span>
        </div>

        <!-- 执行信息 -->
        <div class="flex gap-4 text-xs opacity-60">
          <span>上次执行: {{ job.lastRun || '从未执行' }}</span>
          <span>下次执行: {{ job.nextRun || '-' }}</span>
          <span>执行次数: {{ job.runCount ?? 0 }}</span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading" class="empty-state">
      <Timer :size="40" class="opacity-20 mb-3" />
      <div class="text-sm opacity-40">暂无定时任务</div>
      <div class="text-xs opacity-25 mt-1">点击"添加任务"创建第一个定时任务</div>
    </div>

    <!-- 添加任务弹窗 -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
      <div class="modal-box card w-[480px]">
        <div class="font-semibold mb-4">添加定时任务</div>

        <div class="space-y-3">
          <!-- 任务名称 -->
          <div>
            <label class="text-xs opacity-50 mb-1 block">任务名称 <span class="text-red-400">*</span></label>
            <input v-model="newJob.name" class="input-field w-full" placeholder="例如：每日备份" />
          </div>

          <!-- cron 表达式 -->
          <div>
            <label class="text-xs opacity-50 mb-1 block">CRON 表达式 <span class="text-red-400">*</span></label>
            <input v-model="newJob.expression" class="input-field w-full font-mono" placeholder="例如：*/5 * * * *" @input="onExprInput" />
            <div v-if="exprPreview.description" class="text-xs mt-1.5 px-2 py-1 rounded" style="background:var(--bg-hover)">
              <span class="opacity-50">解析：</span>{{ exprPreview.description }}
            </div>
            <div v-if="exprPreview.nextRuns && exprPreview.nextRuns.length > 0" class="text-xs mt-1 opacity-40">
              下次执行：{{ exprPreview.nextRuns[0] }}
            </div>
            <div v-if="exprError" class="text-xs mt-1 text-red-400">{{ exprError }}</div>
          </div>

          <!-- 执行命令 -->
          <div>
            <label class="text-xs opacity-50 mb-1 block">执行命令 <span class="text-red-400">*</span></label>
            <input v-model="newJob.command" class="input-field w-full font-mono" placeholder="例如：/usr/bin/backup.sh" />
          </div>

          <!-- 描述 -->
          <div>
            <label class="text-xs opacity-50 mb-1 block">描述</label>
            <input v-model="newJob.description" class="input-field w-full" placeholder="任务描述（可选）" />
          </div>
        </div>

        <div class="flex gap-2 justify-end mt-5">
          <button @click="showAddModal = false" class="btn btn-secondary">取消</button>
          <button @click="addJob" class="btn btn-primary" :disabled="adding">保存</button>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div v-if="deleteConfirm" class="modal-overlay" @click.self="deleteConfirm = null">
      <div class="modal-box card w-[400px]">
        <div class="font-semibold mb-3 text-red-400">确认删除任务</div>
        <div class="text-sm mb-4">
          确定要删除任务 <span class="font-mono text-primary-400">{{ deleteConfirm.name }}</span> 吗？
          <div class="mt-2 text-xs opacity-50">此操作不可撤销。</div>
        </div>
        <div class="flex gap-2 justify-end">
          <button @click="deleteConfirm = null" class="btn btn-secondary">取消</button>
          <button @click="doDelete" class="btn btn-danger">确认删除</button>
        </div>
      </div>
    </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, reactive, onMounted } from 'vue'
import { Plus, RefreshCw, Play, Pause, Trash2, Timer } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetCronJobs, AddCronJob, RemoveCronJob, ToggleCronJob, RunCronJobNow, ParseCronExpr } from '../../../wailsjs/go/sysinfo/SysInfo'

const appStore = useAppStore()
const jobs = ref<any[]>([])
const loading = ref(false)
const adding = ref(false)
const showAddModal = ref(false)
const deleteConfirm = ref<any>(null)

// 新任务表单
const newJob = reactive({
  name: '',
  expression: '',
  command: '',
  description: '',
})

// cron 表达式解析预览
const exprPreview = reactive({
  description: '',
  nextRuns: [] as string[],
})
const exprError = ref('')

// 防抖定时器
let debounceTimer: ReturnType<typeof setTimeout> | null = null

async function loadJobs() {
  loading.value = true
  try {
    jobs.value = (await GetCronJobs() as any[]) || []
  } catch (e) {
    appStore.showToast('error', '加载任务列表失败: ' + String(e))
  } finally {
    loading.value = false
  }
}

async function addJob() {
  if (!newJob.name.trim() || !newJob.expression.trim() || !newJob.command.trim()) {
    appStore.showToast('error', '请填写任务名称、CRON 表达式和执行命令')
    return
  }
  adding.value = true
  try {
    await AddCronJob({
      name: newJob.name.trim(),
      expression: newJob.expression.trim(),
      command: newJob.command.trim(),
      description: newJob.description.trim(),
    })
    appStore.showToast('success', '任务添加成功')
    showAddModal.value = false
    // 重置表单
    newJob.name = ''
    newJob.expression = ''
    newJob.command = ''
    newJob.description = ''
    exprPreview.description = ''
    exprPreview.nextRuns = []
    exprError.value = ''
    loadJobs()
  } catch (e) {
    appStore.showToast('error', '添加任务失败: ' + String(e))
  } finally {
    adding.value = false
  }
}

function confirmDelete(job: any) {
  deleteConfirm.value = job
}

async function doDelete() {
  if (!deleteConfirm.value) return
  const id = deleteConfirm.value.id
  const name = deleteConfirm.value.name
  deleteConfirm.value = null
  try {
    await RemoveCronJob(id)
    appStore.showToast('success', `任务 "${name}" 已删除`)
    loadJobs()
  } catch (e) {
    appStore.showToast('error', '删除任务失败: ' + String(e))
  }
}

async function toggleJob(job: any) {
  try {
    await ToggleCronJob(job.id)
    appStore.showToast('success', `任务 "${job.name}" 已${job.enabled ? '禁用' : '启用'}`)
    loadJobs()
  } catch (e) {
    appStore.showToast('error', '操作失败: ' + String(e))
  }
}

async function runNow(job: any) {
  try {
    const result = await RunCronJobNow(job.id)
    appStore.showToast('success', `任务 "${job.name}" 执行成功: ${result}`)
  } catch (e) {
    appStore.showToast('error', `任务执行失败: ${e}`)
  }
}

function onExprInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(async () => {
    if (!newJob.expression.trim()) {
      exprPreview.description = ''
      exprPreview.nextRuns = []
      exprError.value = ''
      return
    }
    try {
      const result = await ParseCronExpr(newJob.expression.trim()) as any
      exprPreview.description = result.description || ''
      exprPreview.nextRuns = result.nextRuns || []
      exprError.value = ''
    } catch (e) {
      exprPreview.description = ''
      exprPreview.nextRuns = []
      exprError.value = '表达式无效: ' + String(e)
    }
  }, 300)
}

onMounted(loadJobs)
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}

.status-dot.green {
  background: #4ade80;
  box-shadow: 0 0 6px rgba(74, 222, 128, 0.5);
}

.status-dot.gray {
  background: #6b7280;
}

.status-dot.red {
  background: #f87171;
  box-shadow: 0 0 6px rgba(248, 113, 113, 0.5);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
}
</style>

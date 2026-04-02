<template>
  <!-- 备忘录工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <StickyNote :size="20" class="text-primary-400" />
        备忘录
      </div>
      <div class="page-desc">快速记录想法与待办事项 · 本地存储 · 随时查看</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-4">
      <button @click="openAdd" class="btn btn-primary">
        <Plus :size="14" />
        新建备忘
      </button>
      <div class="flex-1" />
      <input
        v-model="searchQuery"
        class="input-field w-48"
        placeholder="搜索备忘..."
      />
      <select v-model="sortBy" class="input-field w-28">
        <option value="updated">最近更新</option>
        <option value="created">创建时间</option>
        <option value="pinned">置顶优先</option>
      </select>
    </div>

    <!-- 备忘列表 -->
    <div class="flex-1 overflow-auto">
      <div v-if="filteredNotes.length === 0" class="flex-1 flex items-center justify-center opacity-30 pt-16">
        <div class="text-center">
          <StickyNote :size="48" class="mx-auto mb-2" />
          <div class="text-sm">{{ searchQuery ? '未找到匹配的备忘' : '暂无备忘，点击新建' }}</div>
        </div>
      </div>
      <div v-else class="grid grid-cols-3 gap-3 content-start">
        <div
          v-for="note in filteredNotes"
          :key="note.id"
          class="note-card group"
          :style="{ borderLeftColor: note.color, borderLeftWidth: '4px' }"
        >
          <!-- 标题和操作 -->
          <div class="flex items-start justify-between mb-2">
            <div class="font-medium text-sm flex items-center gap-1.5 truncate flex-1">
              <Pin v-if="note.pinned" :size="12" class="text-yellow-400 shrink-0" />
              <span class="truncate">{{ note.title || '无标题' }}</span>
            </div>
            <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
              <button @click="editNote(note)" class="btn btn-secondary py-0.5 px-1.5" title="编辑">
                <Pencil :size="11" />
              </button>
              <button @click="togglePin(note)" class="btn btn-secondary py-0.5 px-1.5" :title="note.pinned ? '取消置顶' : '置顶'">
                <Pin :size="11" :class="note.pinned ? 'text-yellow-400' : ''" />
              </button>
              <button @click="deleteNote(note.id)" class="btn btn-danger py-0.5 px-1.5" title="删除">
                <Trash2 :size="11" />
              </button>
            </div>
          </div>

          <!-- 内容 -->
          <div class="text-xs opacity-70 leading-relaxed whitespace-pre-wrap line-clamp-4 mb-2">
            {{ note.content || '无内容' }}
          </div>

          <!-- 底部信息 -->
          <div class="flex items-center justify-between text-xs opacity-30">
            <span>{{ formatDate(note.updatedAt) }}</span>
            <div class="flex gap-1">
              <span
                class="w-3 h-3 rounded-full cursor-pointer"
                :style="{ background: note.color }"
                :title="'颜色: ' + note.color"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showEdit" class="modal-overlay" @click.self="showEdit = false">
      <div class="modal-box card w-[480px]">
        <div class="font-semibold mb-3 flex items-center gap-2">
          <Pencil :size="16" />
          {{ editingNote.id ? '编辑备忘' : '新建备忘' }}
        </div>

        <!-- 标题 -->
        <input
          v-model="editingNote.title"
          class="input-field mb-3"
          placeholder="标题（可选）"
          @keyup.ctrl.enter="saveNote"
        />

        <!-- 内容 -->
        <textarea
          v-model="editingNote.content"
          class="textarea-field mb-3"
          rows="8"
          placeholder="输入备忘内容..."
          @keyup.ctrl.enter="saveNote"
        />

        <!-- 颜色选择 -->
        <div class="mb-3">
          <div class="text-xs opacity-60 mb-2">选择颜色</div>
          <div class="flex gap-2">
            <button
              v-for="c in colors"
              :key="c"
              @click="editingNote.color = c"
              :class="[
                'w-6 h-6 rounded-full transition-all',
                editingNote.color === c ? 'ring-2 ring-white ring-offset-2 ring-offset-[var(--bg-card)]' : ''
              ]"
              :style="{ background: c }"
            />
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="flex gap-2 justify-end">
          <button @click="showEdit = false" class="btn btn-secondary">取消</button>
          <button @click="saveNote" :disabled="!editingNote.content.trim()" class="btn btn-primary">
            <Check :size="14" />
            保存
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { StickyNote, Plus, Pin, Pencil, Trash2, Check } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetNotes, SaveNote, PinNote, DeleteNote } from '../../../wailsjs/go/daily/DailyTools'

const appStore = useAppStore()

// 备忘接口
interface Note {
  id: number
  title: string
  content: string
  color: string
  pinned: boolean
  createdAt: string
  updatedAt: string
}

// 状态
const notes = ref<Note[]>([])
const searchQuery = ref('')
const sortBy = ref('updated')
const showEdit = ref(false)
const editingNote = ref<Partial<Note>>({
  title: '',
  content: '',
  color: '#6366f1'
})

// 预设颜色
const colors = [
  '#6366f1', '#8b5cf6', '#ec4899', '#ef4444',
  '#f59e0b', '#10b981', '#06b6d4', '#64748b'
]

// 过滤和排序
const filteredNotes = computed(() => {
  let result = notes.value

  // 搜索
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(n =>
      n.title.toLowerCase().includes(query) ||
      n.content.toLowerCase().includes(query)
    )
  }

  // 排序
  result = [...result].sort((a, b) => {
    if (sortBy.value === 'pinned') {
      if (a.pinned !== b.pinned) return a.pinned ? -1 : 1
    }
    const dateA = new Date(sortBy.value === 'created' ? a.createdAt : a.updatedAt).getTime()
    const dateB = new Date(sortBy.value === 'created' ? b.createdAt : b.updatedAt).getTime()
    return dateB - dateA
  })

  return result
})

// 加载备忘
async function loadNotes() {
  try {
    notes.value = await GetNotes() as Note[]
  } catch {
    notes.value = []
  }
}

// 新增
function openAdd() {
  editingNote.value = {
    id: 0,
    title: '',
    content: '',
    color: '#6366f1'
  }
  showEdit.value = true
}

// 编辑
function editNote(note: Note) {
  editingNote.value = { ...note }
  showEdit.value = true
}

// 保存
async function saveNote() {
  if (!editingNote.value.content?.trim()) return

  try {
    await SaveNote(
      editingNote.value.title || '',
      editingNote.value.content,
      editingNote.value.color || '#6366f1'
    )
    showEdit.value = false
    loadNotes()
    appStore.showToast('success', '备忘已保存')
  } catch {
    appStore.showToast('error', '保存失败')
  }
}

// 切换置顶
async function togglePin(note: Note) {
  try {
    await PinNote(note.id)
    loadNotes()
  } catch {
    appStore.showToast('error', '操作失败')
  }
}

// 删除
async function deleteNote(id: number) {
  if (!confirm('确定删除此备忘？')) return
  try {
    await DeleteNote(id)
    loadNotes()
    appStore.showToast('success', '已删除')
  } catch {
    appStore.showToast('error', '删除失败')
  }
}

// 格式化日期
function formatDate(dateStr: string) {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天 ' + date.toLocaleTimeString().slice(0, 5)
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString()
}

onMounted(loadNotes)
</script>

<style scoped>
.note-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 12px;
  transition: all 0.2s;
}

.note-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  backdrop-filter: blur(2px);
}

.line-clamp-4 {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>

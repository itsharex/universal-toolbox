// 路由配置文件
// 定义所有页面路由，支持懒加载以提升首屏性能

import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// 路由表（全部懒加载）
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/devtools/json',
  },

  // ============ 开发工具 ============
  {
    path: '/devtools',
    redirect: '/devtools/json',
  },
  {
    path: '/devtools/json',
    name: 'JsonTool',
    component: () => import('@/views/devtools/JsonTool.vue'),
    meta: { title: 'JSON 工具', category: 'devtools', icon: 'Braces' },
  },
  {
    path: '/devtools/jsondiff',
    name: 'JsonDiffTool',
    component: () => import('@/views/devtools/JsonDiffTool.vue'),
    meta: { title: 'JSON 对比', category: 'devtools', icon: 'GitCompare' },
  },
  {
    path: '/devtools/xml',
    name: 'XmlTool',
    component: () => import('@/views/devtools/XmlTool.vue'),
    meta: { title: 'XML 工具', category: 'devtools', icon: 'Code2' },
  },
  {
    path: '/devtools/base64',
    name: 'Base64Tool',
    component: () => import('@/views/devtools/Base64Tool.vue'),
    meta: { title: 'Base64 编解码', category: 'devtools', icon: 'Binary' },
  },
  {
    path: '/devtools/url',
    name: 'UrlTool',
    component: () => import('@/views/devtools/UrlTool.vue'),
    meta: { title: 'URL 编解码', category: 'devtools', icon: 'Link' },
  },
  {
    path: '/devtools/hash',
    name: 'HashTool',
    component: () => import('@/views/devtools/HashTool.vue'),
    meta: { title: '哈希计算', category: 'devtools', icon: 'Hash' },
  },
  {
    path: '/devtools/text',
    name: 'TextTool',
    component: () => import('@/views/devtools/TextTool.vue'),
    meta: { title: '文本工具', category: 'devtools', icon: 'FileText' },
  },
  {
    path: '/devtools/qrcode',
    name: 'QRCodeTool',
    component: () => import('@/views/devtools/QRCodeTool.vue'),
    meta: { title: '二维码工具', category: 'devtools', icon: 'QrCode' },
  },
  {
    path: '/devtools/uuid',
    name: 'UuidTool',
    component: () => import('@/views/devtools/UuidTool.vue'),
    meta: { title: 'UUID 生成', category: 'devtools', icon: 'Fingerprint' },
  },
  {
    path: '/devtools/timestamp',
    name: 'TimestampTool',
    component: () => import('@/views/devtools/TimestampTool.vue'),
    meta: { title: '时间戳转换', category: 'devtools', icon: 'Clock' },
  },
  {
    path: '/devtools/regex',
    name: 'RegexTool',
    component: () => import('@/views/devtools/RegexTool.vue'),
    meta: { title: '正则测试', category: 'devtools', icon: 'Regex' },
  },
  {
    path: '/devtools/snippets',
    name: 'SnippetsTool',
    component: () => import('@/views/devtools/SnippetsTool.vue'),
    meta: { title: '代码片段', category: 'devtools', icon: 'BookMarked' },
  },
  {
    path: '/devtools/crypto',
    name: 'CryptoTool',
    component: () => import('@/views/devtools/CryptoTool.vue'),
    meta: { title: '加密解密', category: 'devtools', icon: 'Shield' },
  },
  {
    path: '/devtools/radix',
    name: 'RadixTool',
    component: () => import('@/views/devtools/RadixTool.vue'),
    meta: { title: '进制转换', category: 'devtools', icon: 'Binary' },
  },
  {
    path: '/devtools/dummydata',
    name: 'DummyDataTool',
    component: () => import('@/views/devtools/DummyDataTool.vue'),
    meta: { title: '占位文本生成', category: 'devtools', icon: 'Database' },
  },
  {
    path: '/devtools/apidoc',
    name: 'APIDocTool',
    component: () => import('@/views/devtools/APIDocTool.vue'),
    meta: { title: '接口文档生成', category: 'devtools', icon: 'FileCode' },
  },
  {
    path: '/devtools/obfuscator',
    name: 'CodeObfuscator',
    component: () => import('@/views/devtools/CodeObfuscator.vue'),
    meta: { title: '代码混淆', category: 'devtools', icon: 'Lock' },
  },
  {
    path: '/devtools/formatter',
    name: 'CodeFormatter',
    component: () => import('@/views/devtools/CodeFormatter.vue'),
    meta: { title: '代码格式化', category: 'devtools', icon: 'AlignLeft' },
  },

  // ============ 系统工具 ============
  {
    path: '/sysinfo/process',
    name: 'ProcessTool',
    component: () => import('@/views/sysinfo/ProcessTool.vue'),
    meta: { title: '进程管理', category: 'sysinfo', icon: 'Cpu' },
  },
  {
    path: '/sysinfo/ports',
    name: 'PortsTool',
    component: () => import('@/views/sysinfo/PortsTool.vue'),
    meta: { title: '端口查看', category: 'sysinfo', icon: 'Network' },
  },
  {
    path: '/sysinfo/sysinfo',
    name: 'SystemInfo',
    component: () => import('@/views/sysinfo/SystemInfo.vue'),
    meta: { title: '系统信息', category: 'sysinfo', icon: 'Monitor' },
  },
  {
    path: '/sysinfo/rename',
    name: 'BatchRename',
    component: () => import('@/views/sysinfo/BatchRename.vue'),
    meta: { title: '批量重命名', category: 'sysinfo', icon: 'FilePen' },
  },
  {
    path: '/sysinfo/colorpicker',
    name: 'ColorPicker',
    component: () => import('@/views/sysinfo/ColorPicker.vue'),
    meta: { title: '取色器', category: 'sysinfo', icon: 'Pipette' },
  },
  {
    path: '/sysinfo/imagetool',
    name: 'ImageTool',
    component: () => import('@/views/sysinfo/ImageTool.vue'),
    meta: { title: '图片工具', category: 'sysinfo', icon: 'Image' },
  },
  {
    path: '/sysinfo/filebatch',
    name: 'FileBatchTool',
    component: () => import('@/views/sysinfo/FileBatchTool.vue'),
    meta: { title: '文件批量处理', category: 'sysinfo', icon: 'FolderSearch' },
  },
  {
    path: '/sysinfo/diskcleaner',
    name: 'DiskCleaner',
    component: () => import('@/views/sysinfo/DiskCleaner.vue'),
    meta: { title: '磁盘清理', category: 'sysinfo', icon: 'Trash2' },
  },
  {
    path: '/sysinfo/clipboard',
    name: 'ClipboardManager',
    component: () => import('@/views/sysinfo/ClipboardManager.vue'),
    meta: { title: '剪贴板管理', category: 'sysinfo', icon: 'Clipboard' },
  },
  {
    path: '/sysinfo/cron',
    name: 'CronTool',
    component: () => import('@/views/sysinfo/CronTool.vue'),
    meta: { title: '定时任务', category: 'sysinfo', icon: 'Timer' },
  },

  // ============ 日常工具 ============
  {
    path: '/daily/calculator',
    name: 'Calculator',
    component: () => import('@/views/daily/Calculator.vue'),
    meta: { title: '计算器', category: 'daily', icon: 'Calculator' },
  },
  {
    path: '/daily/converter',
    name: 'Converter',
    component: () => import('@/views/daily/Converter.vue'),
    meta: { title: '单位换算', category: 'daily', icon: 'ArrowLeftRight' },
  },
  {
    path: '/daily/notes',
    name: 'Notes',
    component: () => import('@/views/daily/Notes.vue'),
    meta: { title: '备忘录', category: 'daily', icon: 'StickyNote' },
  },
  {
    path: '/daily/password',
    name: 'PasswordGenerator',
    component: () => import('@/views/daily/PasswordGenerator.vue'),
    meta: { title: '密码生成', category: 'daily', icon: 'KeyRound' },
  },
  {
    path: '/daily/qrbatch',
    name: 'QRBatchTool',
    component: () => import('@/views/daily/QRBatchTool.vue'),
    meta: { title: '二维码批量处理', category: 'daily', icon: 'QrCode' },
  },
  {
    path: '/daily/docconverter',
    name: 'DocConverter',
    component: () => import('@/views/daily/DocConverter.vue'),
    meta: { title: '文档转换', category: 'daily', icon: 'FileInput' },
  },

  // ============ 网络工具 ============
  {
    path: '/network/ping',
    name: 'PingTool',
    component: () => import('@/views/network/PingTool.vue'),
    meta: { title: 'Ping 测试', category: 'network', icon: 'Signal' },
  },
  {
    path: '/network/scan',
    name: 'ScanTool',
    component: () => import('@/views/network/ScanTool.vue'),
    meta: { title: '内网扫描', category: 'network', icon: 'Radar' },
  },
  {
    path: '/network/http',
    name: 'HttpTool',
    component: () => import('@/views/network/HttpTool.vue'),
    meta: { title: 'HTTP 测试', category: 'network', icon: 'Globe' },
  },
  {
    path: '/network/hosts',
    name: 'HostsTool',
    component: () => import('@/views/network/HostsTool.vue'),
    meta: { title: 'Hosts 编辑', category: 'network', icon: 'Server' },
  },
  {
    path: '/network/dns',
    name: 'DnsTool',
    component: () => import('@/views/network/DnsTool.vue'),
    meta: { title: 'DNS 查询', category: 'network', icon: 'Globe2' },
  },
  {
    path: '/network/whois',
    name: 'WhoisTool',
    component: () => import('@/views/network/WhoisTool.vue'),
    meta: { title: 'WHOIS 查询', category: 'network', icon: 'SearchCode' },
  },
  {
    path: '/network/speedtest',
    name: 'SpeedTest',
    component: () => import('@/views/network/SpeedTest.vue'),
    meta: { title: '网络测速', category: 'network', icon: 'Gauge' },
  },
  {
    path: '/network/portforward',
    name: 'PortForward',
    component: () => import('@/views/network/PortForward.vue'),
    meta: { title: '端口转发', category: 'network', icon: 'ArrowLeftRight' },
  },

  // ============ 设置 ============
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/settings/Settings.vue'),
    meta: { title: '设置', category: 'settings', icon: 'Settings' },
  },
  {
    path: '/settings/theme',
    name: 'ThemeEditor',
    component: () => import('@/views/settings/ThemeEditor.vue'),
    meta: { title: '主题编辑器', category: 'settings', icon: 'Palette' },
  },
  {
    path: '/settings/shortcuts',
    name: 'ShortcutManager',
    component: () => import('@/views/settings/ShortcutManager.vue'),
    meta: { title: '快捷键管理', category: 'settings', icon: 'Keyboard' },
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router

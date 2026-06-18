<template>
  <div class="app" @keydown.esc="onEsc" tabindex="0">
    <div class="update-bar" v-if="updateInfo.has_update" @click="doUpdate">
      发现新版本 {{ updateInfo.latest }} — 当前 {{ updateInfo.current }}，点击下载
    </div>
    <TopBar :system="data.system" :cpu="data.cpu" :gpu="data.gpu" :memory="data.memory" :total-power="data.total_power" :network="data.network" />
    <div class="main-grid">
      <CpuBlock :cpu="data.cpu" />
      <GpuBlock :gpu="data.gpu" />
    </div>
    <div class="bottom-grid">
      <MemoryBlock :memory="data.memory" />
      <StorageBlock :storage="data.storage" :disk-i-o="data.disk_io" />
      <FansBlock :fans="data.fans" />
    </div>
    <div class="overlay-btns">
      <div class="overlay-btn" @click="toggleFs" :title="fullscreen ? '退出全屏' : '全屏'">
        <svg v-if="fullscreen" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M8 3v3a2 2 0 01-2 2H3m18 0h-3a2 2 0 01-2-2V3m0 18v-3a2 2 0 012-2h3M3 16h3a2 2 0 012 2v3"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M8 3H5a2 2 0 00-2 2v3m18 0V5a2 2 0 00-2-2h-3m0 18h3a2 2 0 002-2v-3M3 16v3a2 2 0 002 2h3"/>
        </svg>
      </div>
      <div class="overlay-btn" @click="doShutdown" title="关闭">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6L6 18M6 6l12 12"/>
        </svg>
      </div>
    </div>
  </div>
</template>

<script>
import { reactive, ref, onMounted, onUnmounted } from 'vue'
import TopBar from './components/TopBar.vue'
import CpuBlock from './components/CpuBlock.vue'
import GpuBlock from './components/GpuBlock.vue'
import MemoryBlock from './components/MemoryBlock.vue'
import StorageBlock from './components/StorageBlock.vue'
import FansBlock from './components/FansBlock.vue'

export default {
  components: { TopBar, CpuBlock, GpuBlock, MemoryBlock, StorageBlock, FansBlock },
  setup() {
    const data = reactive({
      system: { time: '--:--:--', date: '----/--/--', uptime: '--', process_count: 0, thread_count: 0, ip_address: '', display_info: '', fps: 0 },
      cpu: { model: '', max_freq: '', cores_threads: '', package_temp: 0, usage: 0, clock_speed: 0, vcore: 0, fan_speed: 0, power: 0, per_core_usage: [], mb_temp: 0, vrm_temp: 0 },
      gpu: { model: '', vram_spec: '', vram_type: '', pcie_version: '', temp: 0, usage: 0, clock: 0, mem_used: 0, mem_total: 0, fan_speed: 0, power: 0 },
      memory: { used_gb: 0, total_gb: 0, type: '', frequency: '', channel: '', brand: '' },
      storage: [],
      fans: [],
      total_power: 0,
      network: { download_mbps: 0, upload_mbps: 0, connection_count: 0 },
      disk_io: { read_mbps: 0, write_mbps: 0 }
    })

    const fullscreen = ref(true)
    const updateInfo = reactive({ has_update: false, latest: '', current: '', download_url: '' })
    let eventCleanup = null
    let clockTimer = null
    let fpsTimer = null
    let fpsFrames = 0

    function mergeSystem(sys) {
      if (sys.uptime !== undefined) data.system.uptime = sys.uptime
      if (sys.process_count !== undefined) data.system.process_count = sys.process_count
      if (sys.ip_address !== undefined) data.system.ip_address = sys.ip_address
      if (sys.display_info !== undefined) data.system.display_info = sys.display_info
    }

    function startFpsTracker() {
      fpsFrames = 0
      const tick = () => {
        fpsFrames++
        requestAnimationFrame(tick)
      }
      requestAnimationFrame(tick)
      fpsTimer = setInterval(() => {
        data.system.fps = fpsFrames
        fpsFrames = 0
      }, 1000)
    }

    function updateClock() {
      const now = new Date()
      const h = String(now.getHours()).padStart(2, '0')
      const m = String(now.getMinutes()).padStart(2, '0')
      const s = String(now.getSeconds()).padStart(2, '0')
      data.system.time = `${h}:${m}:${s}`
      const weekdays = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
      const y = now.getFullYear()
      const mo = String(now.getMonth() + 1).padStart(2, '0')
      const d = String(now.getDate()).padStart(2, '0')
      data.system.date = `${y}/${mo}/${d} ${weekdays[now.getDay()]}`
    }

    onMounted(async () => {
      updateClock()
      clockTimer = setInterval(updateClock, 1000)
      if (window.runtime) startFpsTracker()

      try {
        const d = await window.go.main.App.GetDashboard()
        if (d) {
          if (d.system) { mergeSystem(d.system); delete d.system }
          Object.assign(data, d)
        }
        fullscreen.value = await window.go.main.App.GetIsFullscreen()
      } catch (e) { console.error('initial dashboard error:', e) }

      if (window.runtime) {
        eventCleanup = window.runtime.EventsOn('dashboard-update', (raw) => {
          try {
            const d = typeof raw === 'string' ? JSON.parse(raw) : raw
            if (d) {
              console.log('dashboard event, cpu usage:', d.cpu?.usage)
              if (d.system) { mergeSystem(d.system); delete d.system }
              Object.assign(data, d)
            }
          } catch (e) { console.error('dashboard event error:', e) }
        })
      }

      document.addEventListener('keydown', handleKey)

      try {
        const info = await window.go.main.App.CheckUpdate()
        if (info && info.has_update) {
          updateInfo.has_update = info.has_update
          updateInfo.latest = info.latest
          updateInfo.current = info.current
          updateInfo.download_url = info.download_url
        }
      } catch (e) { /* ignore */ }
    })

    onUnmounted(() => {
      if (clockTimer) clearInterval(clockTimer)
      if (fpsTimer) clearInterval(fpsTimer)
      if (eventCleanup) eventCleanup()
      document.removeEventListener('keydown', handleKey)
    })

    function handleKey(e) {
      if (e.key === 'Escape') {
        toggleFs()
      }
    }

    async function toggleFs() {
      try {
        await window.go.main.App.ToggleFullscreen()
        fullscreen.value = await window.go.main.App.GetIsFullscreen()
      } catch (e) { /* ignore */ }
    }

    async function doShutdown() {
      try { await window.go.main.App.Shutdown() } catch (e) { /* ignore */ }
    }

    async function doUpdate() {
      if (updateInfo.download_url) {
        try { await window.go.main.App.OpenURL(updateInfo.download_url) } catch (e) { /* ignore */ }
      }
    }

    function onEsc() {}

    return { data, fullscreen, updateInfo, toggleFs, doShutdown, doUpdate, onEsc }
  }
}
</script>

<style scoped>
.app {
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 12px 16px;
  gap: 10px;
}

.update-bar {
  background: linear-gradient(90deg, var(--amber-dim), rgba(245, 158, 11, 0.25));
  border: 1px solid var(--amber);
  border-radius: 8px;
  color: var(--amber);
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  padding: 8px 16px;
  cursor: pointer;
  flex-shrink: 0;
  transition: opacity 0.2s;
}
.update-bar:hover {
  opacity: 0.85;
}

.main-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  flex: 1;
  min-height: 0;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
  flex-shrink: 0;
}

.overlay-btns {
  position: fixed;
  top: 6px;
  right: 6px;
  display: flex;
  gap: 4px;
  z-index: 100;
}

.overlay-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  opacity: 0.15;
  transition: opacity 0.2s;
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
  -webkit-app-region: no-drag;
}

.overlay-btn:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.2);
}
</style>
